package main

import "flag"

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}


func Must(err error) {
	if err != nil {
		panic(err)
	}
}