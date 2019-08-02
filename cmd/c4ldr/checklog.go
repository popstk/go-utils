package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	oraReg = `ORA-\d+:.+`
)

func checkLog(item Item, args ...string) {
	if len(args) <= 0 {
		panic("time: 20190730 or 2019073000")
	}

	points := TimeExpand(args[0])
	data := map[string]string{}

	for _, p := range points {
		prefix := strings.Builder{}
		data["time"] = p
		path := stringFormat(item.Path, data)

		prefix.WriteString(fmt.Sprintf("> %s \n", p))

		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Println(prefix.String())
			fmt.Println(err)
			continue
		}

		for _, m := range matches {
			prefix.WriteString(fmt.Sprintf(">> %s \n", m))

			content, err := ioutil.ReadFile(m)
			if err != nil {
				if prefix.Len() > 0 {
					fmt.Print(prefix.String())
					prefix.Reset()
				}
				fmt.Println(">>> ", err)
				fmt.Println("")
				continue
			}

			r := regexp.MustCompile(oraReg)
			result := r.FindAllString(string(content), -1)
			if len(result) > 0 {
				if prefix.Len() > 0 {
					fmt.Print(prefix.String())
					prefix.Reset()
				}
				for _, match := range result {
					fmt.Println(">>> ", match)
				}

				fmt.Println("")
			}
		}
	}
}