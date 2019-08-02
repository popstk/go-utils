package main

type Item struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Queue string `json:"queue"`
	Cwd string `json:"cwd"`
	Cron string `json:"cron"`
	Reset []string `json:"reset"`
}

type Config struct {
	Items map[string]Item `json:"items"`
	Redis string `json:"redis"`
}

