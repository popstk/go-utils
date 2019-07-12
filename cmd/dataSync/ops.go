package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

var SupportedTestDB = map[string]string{
	"mysql": "select * from %s limit 1",
	"oci8": "select * from %s where rownum <= 1",
}

func OutputRowInfo(cols []*sql.ColumnType, values []interface{}) {
	for i, col := range cols {
		var v interface{}
		val := values[i]
		b, ok := val.([]byte)// sql.rawBytes
		byteArray := false
		if ok {
			v = string(b)
			byteArray = true
		} else {
			v = val
		}
		fmt.Printf("%s [%s][%s] -> %v ", col.Name(), col.DatabaseTypeName(), col.ScanType(), v)
		if byteArray {
			fmt.Printf("string -> %d, raw bytes -> %d", utf8.RuneCountInString(string(b)), len(b))
		}
		fmt.Println()
	}
}

func checkStat(db *sql.DB, table, querySQL string) {
	query := fmt.Sprintf("select count(*) from %s", table)
	log.Println("exec sql => ", query)
	row := db.QueryRow(query)
	{
		var dec sql.NullFloat64
		Must(row.Scan(&dec))
		if dec.Valid {
			count := int64(dec.Float64)
			fmt.Printf("table %s total %d lines\n", table, count)
		}
	}

	querySQL = fmt.Sprintf(querySQL, table)
	log.Println("exec sql => ", querySQL)
	rows, err := db.Query(querySQL)
	Must(err)
	colTypes, err := rows.ColumnTypes()
	Must(err)

	count := len(colTypes)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range colTypes {
			valuePtrs[i] = &values[i]
		}
		Must(rows.Scan(valuePtrs...))
		OutputRowInfo(colTypes, values)
	}
	Must(rows.Close())
}


func checkTable(index int) {
	tb := config.Tables[index]

	src, err := sql.Open(config.Src.Driver, config.Src.URL)
	Must(err)
	querySQL, ok := SupportedTestDB[config.Src.Driver]
	if ok {
		checkStat(src, tb.From, querySQL)
	} else {
		fmt.Println("Unsupported driver: ", config.Src.Driver)
	}

	fmt.Println(strings.Repeat("=", 30))

	dest, err := sql.Open(config.Dest.Driver, config.Dest.URL)
	Must(err)
	querySQL, ok = SupportedTestDB[config.Dest.Driver]
	if ok {
		checkStat(dest, tb.To, querySQL)
	} else {
		fmt.Println("Unsupported driver: ", config.Dest.Driver)
	}
}

func Sync(ctx context.Context, tb TableInfo) error {
	fromDB, err := sql.Open(config.Src.Driver, config.Src.URL)
	Must(err)
	log.Println("check source table: ", tb.From)

	var srcID int64
	query := fmt.Sprintf("select max(id) from %s", tb.From)
	row := fromDB.QueryRow(query)
	Must(row.Scan(&srcID))
	log.Println("max id on source table is ", srcID)

	toDB, err := sql.Open(config.Dest.Driver, config.Dest.URL)
	Must(err)
	log.Println("check destination table: ", tb.To)


	query = fmt.Sprintf("select max(id) from %s", tb.To)
	row = toDB.QueryRow(query)
	var destID int64
	{
		var v sql.NullFloat64
		Must(row.Scan(&v))
		if v.Valid {
			destID = int64(v.Float64)
		} else {
			log.Println("destination table is empty...")
			destID = 0
		}
	}
	log.Println("max id on destination table is ", destID)

	if destID == srcID {
		fmt.Println("No need to sync..")
		return nil
	}

	if srcID < destID {
		return errors.New("destination id bigger than source id")
	}

	querySQL := fmt.Sprintf("select * from %s where id > ? and id <= ?", tb.From)
	for ; srcID > destID ; {
		end := destID +int64(maxLines)
		log.Printf("sync id (%d, %d]\n", destID, end)
		rows, err := fromDB.Query(querySQL, destID, end)
		Must(err)

		cols, err := rows.Columns()
		Must(err)

		placeHold := ""
		for i := range cols {
			if i != 0 {
				placeHold += ", "
			}
			placeHold += fmt.Sprintf(":%d", i)
		}

		insertSQL := fmt.Sprintf("insert into %s(%s) values (%s)", tb.To,
			strings.Join(cols, ", "), placeHold)

		//log.Println("insert sql is ", insertSQL)

		colType, err := rows.ColumnTypes()
		Must(err)
		count := len(cols)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		tx, err := toDB.BeginTx(ctx, nil)
		Must(err)
		stmt, err := tx.Prepare(insertSQL)
		Must(err)

		fmt.Println("start tx insert...")

		var rowsAffected int64
		for rows.Next() {
			Must(rows.Scan(valuePtrs...))

			for i, col := range colType {
				// vchar -> sql.rawBytes 转为 string
				t := col.DatabaseTypeName()
				if t == "VARCHAR" || t == "TEXT" {
					val := values[i]
					if val != nil {
						b := values[i].([]uint8) // interface{} 转换不能使用type别名， sql.rawBytes等价于[]uint8
						values[i] = string(b)
					}
				} else if t == "DECIMAL" {
					val := values[i]
					if val != nil {
						b := values[i].([]uint8)
						values[i] = string(b)
					}
				}
			}

			timeoutCtx, _ := context.WithTimeout(ctx, 5 * time.Second)
			_, err := stmt.ExecContext(timeoutCtx, values...)
			if err != nil {
				fmt.Println(strings.Repeat("-", 79))
				OutputRowInfo(colType, values)
				panic(err)
			}
			rowsAffected += 1
		}
		destID = end

		Must(tx.Commit())
		Must(stmt.Close())
		Must(rows.Close())
		log.Printf("insert %d lines\n", rowsAffected)
	}

	return nil
}