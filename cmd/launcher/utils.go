package main

import (
	"fmt"
	"strings"
)

func FmtNamedVariable(format string, p map[string]string) string {
	args := make([]string, 0, len(p)*2)
	for k, v := range p {
		args = append(args, fmt.Sprintf("{%s}",k), v)
	}

	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

