package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-oci8"
	"github.com/pkg/errors"
	"log"
	"strings"
)

var (
	mysqlURL  string
	oracleURL string
	srcTable  string
	destTable string
	maxLines  int
)

const (
	defaultMysql  = ""
	defaultOracle = ""

	DBOracle int = iota
	DBMySQL
)

func init() {
	flag.StringVar(&mysqlURL, "mysql", defaultMysql, "mysql connect source")
	flag.StringVar(&oracleURL, "oracle", defaultOracle, "oracle connect source")
	flag.StringVar(&srcTable, "src", "tb_app_user", "source table name")
	flag.StringVar(&destTable, "dest", "tb_app_user", "destination table name")
	flag.IntVar(&maxLines, "lines", 1000, "max lines for each sync")
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}



// mysql -> oracle
func FullSync(ctx context.Context, fromDB, toDB *sql.DB) error {
	log.Println("check source table: ", srcTable)

	// 检查mysql最大id
	var srcID int64
	query := fmt.Sprintf("select max(id) from %s", srcTable)
	row := fromDB.QueryRow(query)
	Must(row.Scan(&srcID))
	log.Println("max id on source table is ", srcID)

	// 检查oracle最大id
	log.Println("check destination table: ", destTable)

	var destID int64
	var v interface{}
	query = fmt.Sprintf("select max(id) from %s", destTable)
	row = toDB.QueryRow(query)
	Must(row.Scan(&v))
	if v == nil {
		log.Println("destination table is empty...")
		destID = 0
	} else {
		Must(row.Scan(&destID))
	}
	log.Println("max id on destination table is ", destID)

	if destID == srcID {
		fmt.Println("No need to sync..")
		return nil
	}

	if srcID < destID {
		return errors.New("destination id bigger than source id")
	}

	querlSQL := fmt.Sprintf("select * from %s where id > ? and id <= ?", srcTable)

	for ; srcID > destID ; {
		end := destID +int64(maxLines)
		log.Printf("sync id (%d, %d]\n", destID, end)
		rows, err := fromDB.Query(querlSQL, destID, end)
		Must(err)

		cols, err := rows.Columns()
		Must(err)

		//placeHold := strings.Repeat("?,", len(cols))
		//placeHold = placeHold[:len(placeHold)-1]
		placeHold := ""
		for i := range cols {
			if i != 0 {
				placeHold += ", "
			}
			placeHold += fmt.Sprintf(":%d", i)
		}

		inertSQL := fmt.Sprintf("insert into %s(%s) values (%s)", destTable,
			strings.Join(cols, ", "), placeHold)

		log.Println("insert sql is ", inertSQL)

		Must(err)
		for rows.Next() {
			count := len(cols)
			values := make([]interface{}, count)
			valuePtrs := make([]interface{}, count)
			for i := range cols {
				valuePtrs[i] = &values[i]
			}
			Must(rows.Scan(valuePtrs...))

			log.Println("values len is ", len(values))

			r, err := toDB.Exec(inertSQL, values...)
			Must(err)
			log.Println("result is ", r)
		}

		destID = end
	}

	return nil
}

func checkTable(db *sql.DB, t int) {
	table := srcTable
	if t == DBOracle {
		table = destTable
	}

	// 获取总行数
	query := fmt.Sprintf("select count(*) from %s", table)
	log.Println("exec sql => ", query)
	row := db.QueryRow(query)
	{
		var count int64
		Must(row.Scan(&count))
		fmt.Printf("table %s total %d lines\n", table, count)
	}

	if t == DBOracle {
		query = fmt.Sprintf("select * from %s where rownum <= 5", table)
	} else if t == DBMySQL {
		query = fmt.Sprintf("select * from %s limit 5", table)
	} else {
		panic(t)
	}

	log.Println("exec sql => ", query)
	rows, err := db.Query(query)
	Must(err)
	colTypes, err := rows.ColumnTypes()
	Must(err)

	count := len(colTypes)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	defer rows.Close()
	for _, col := range colTypes {
		fmt.Printf("%s-> %s\n", col.Name(), col.DatabaseTypeName())
	}
	fmt.Println("")

	for rows.Next() {
		for i := range colTypes {
			valuePtrs[i] = &values[i]
		}
		Must(rows.Scan(valuePtrs...))

		for i := range colTypes {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			fmt.Printf("%v ", v)
		}
		fmt.Println("")
	}
}

func main() {
	flag.Parse()

	log.Println("mysql url is ", mysqlURL)
	srcDB, err := sql.Open("mysql", mysqlURL)
	Must(err)

	log.Println("oracle url is ", oracleURL)
	destDB, err := sql.Open("oci8", oracleURL)
	Must(err)

	Must(FullSync(context.Background(),srcDB, destDB))
}
