package main

import (
	"fmt"
	"os"
)

func reset(item Item, args ...string) {
	if item.Cwd != "" {
		if err := os.Chdir(item.Cwd); err != nil {
			fmt.Println(err)
		}
	}

	points := TimeExpand(args[0])
	data := map[string]string{}
	for _, p := range points {
		data["time"] = p
		fmt.Println(">> "+ p)

		for _, cmd := range item.Reset {
			cmd = stringFormat(cmd, data)
			fmt.Println(">>> "+ cmd)
		}
	}
}

