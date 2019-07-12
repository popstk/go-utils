package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	user     string
	password string
)

func init() {
	flag.StringVar(&user, "u", "", "set user name")
	flag.StringVar(&password, "p", "", "set password")
}

func Must(err error) {
	if err !=  nil {
		panic(err)
	}
}

func test() error {
	r, err := http.Get("http://www.baidu.com")
	if err != nil {
		return err
	}

	if strings.Contains(r.Request.URL.String(), "http://10.0.0.8/redirect") {
		fmt.Println("Try Login...")
		login()
	}

	return nil
}

func login() {
	defer func() {
		if err := recover(); err!= nil {
			fmt.Println(err)
		}
	}()

	req, err := http.NewRequest("GET", "http://10.0.0.8/login", nil)
	Must(err)
	req.SetBasicAuth(user, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	Must(err)

	fmt.Println("status code is ", resp.StatusCode )
}

func main() {
	flag.Parse()

	for {
		if err := test(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Minute)
	}
}
