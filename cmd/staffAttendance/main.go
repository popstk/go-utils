package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ResultAttendanceDay struct {
	DutyDate string `json:"dutyDate"`
	LateMinutes string `json:"lateMinutes"`
	OnDutyDescr string `json:"onDutyDescr"`
	OnDutyTime string `json:"onDutyTime"`
	OffDutyTime string `json:"offDutyTime"`
	StaffCode string `json:"staffCode"`
	StaffName string `json:"staffName"`
	StatusInfo string `json:"statusInfo"`
	DateType string `json:"dateType"` // 0工作日 1双休日
}

type ResultData struct {
	List []ResultAttendanceDay `json:"list"`
}

type Result struct {
	Data ResultData `json:"data"`
}

type QueryCondition struct {
	StaffCode string `json:"staffCode"`
	Month string `json:"month"`
	Year string `json:"year"`
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	staffCode string
	queryTime string
)

func init() {
	flag.StringVar(&staffCode, "s", "", "staff Code")
	flag.StringVar(&queryTime, "t", "", "year and month")
}

func TimeOnly(s string) string {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	Must(err)

	return t.Format("15:04:05")
}

func main() {
	flag.Parse()

	if staffCode == "" {
		fmt.Println("staff Code is empty")
		return
	}

	if !strings.HasPrefix(staffCode,"CX") {
		staffCode = "CX"+staffCode
	}

	t := time.Now()
	if queryTime != "" {
		var err error
		t, err = time.Parse("20060102", queryTime)
		Must(err)
	}
	var c QueryCondition
	c.StaffCode = staffCode
	c.Year = strconv.Itoa(t.Year())
	c.Month = fmt.Sprintf("%02d", t.Month())
	condition, err := json.Marshal(c)
	Must(err)

	query := url.Values{}
	query.Set("page", "1")
	query.Set("rows", "30")
	query.Set("condition", string(condition))

	client := NewHTTPClient()
	resp, err := client.PostForm("http://hr.richinfo.cn/queryAttendances4NoSesssion.action", query)
	Must(err)

	data, err := ioutil.ReadAll(resp.Body)
	Must(err)
	Must(resp.Body.Close())
	var rspResult map[string]string
	Must(json.Unmarshal(data, &rspResult))

	var result Result
	Must(json.Unmarshal([]byte(rspResult["result"]), &result))

	if len(result.Data.List) <= 0 {
		fmt.Println("empty....")
		fmt.Println("result is ", rspResult["result"])
		return
	}

	fmt.Println("-> ", result.Data.List[0].StaffName)
	for _, i := range result.Data.List {
		if i.DateType == "0" {
			if i.OnDutyTime != "" &&  i.OffDutyTime != "" {
				fmt.Printf("%s	%s	->	%s\n", i.DutyDate ,TimeOnly(i.OnDutyTime), TimeOnly(i.OffDutyTime))
			} else {
				fmt.Printf("%s	%s\n",i.DutyDate, i.StatusInfo)
			}
		}
	}
}
