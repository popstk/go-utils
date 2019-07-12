package main

import (
	"encoding/json"
	"io/ioutil"
)

type DataBaseInfo struct {
	Driver string `json:"driver"`
	URL string `json:"url"`
}

type TableInfo struct {
	From string `json:"from"`
	To string  `json:"To"`
}

type Config struct {
	Src DataBaseInfo `json:"src"`
	Dest DataBaseInfo `json:"dest"`
	Tables []TableInfo `json:"tables"`
}

func ReadConfig(p string) *Config {
	data, err := ioutil.ReadFile(p)
	Must(err)

	var conf Config
	Must(json.Unmarshal(data, &conf))

	return &conf
}



