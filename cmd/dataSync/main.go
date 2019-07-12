package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-oci8"
	"log"
)

var (
	maxLines  int
	configFile string
	checkTableIndex  int
	syncTableIndex int

	config *Config
)


func init() {
	flag.IntVar(&maxLines, "lines", 1000, "max lines for each sync")
	flag.IntVar(&checkTableIndex, "test", 0, "check index table pair only")
	flag.IntVar(&syncTableIndex, "sync", 0, "sync index table pair only")
	flag.StringVar(&configFile, "f", "dataSync.json", "specify config file")
	log.SetFlags(log.Lshortfile|log.Ltime)
}

func main() {
	flag.Parse()
	config = ReadConfig(configFile)

	if isFlagPassed("test") {
		checkTable(checkTableIndex)
		return
	}

	if isFlagPassed("sync") {
		tb := config.Tables[syncTableIndex]
		Must(Sync(context.Background(), tb))
		return
	}

	for _, tb := range config.Tables {
		fmt.Printf("Sync %s -> %s\n", tb.From, tb.To)
		Must(Sync(context.Background(), tb))
	}
}
