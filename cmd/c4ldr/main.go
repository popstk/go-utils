package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var (
	config Config
	item string
	verbose bool

	cmdMap = map[string]func(item Item,args ...string) {
		"check": checkLog,
		"reset": reset,
	}
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	data, err := ioutil.ReadFile("c4ldr.json")
	Must(err)
	Must(json.Unmarshal(data, &config))

	flag.StringVar(&item, "i", "pica", "specify item in c4ldr.json")
	flag.BoolVar(&verbose, "v", false, "verbose")
}




func main() {
	flag.Parse()

	i, ok := config.Items[item]
	if !ok {
		panic("invalid item: %s" + item)
	}

	cmd := flag.Arg(0)
	f, ok := cmdMap[cmd]
	if !ok {
		panic("invalid command: " + cmd)
	}

	f(i, flag.Args()[1:]...)
}
