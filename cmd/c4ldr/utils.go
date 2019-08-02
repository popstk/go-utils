package main

import (
	"fmt"
	"strings"
)

// {time} -> m["time"]
func stringFormat(s string, m map[string]string) string {
	for k, v := range m {
		s = strings.Replace(s, fmt.Sprintf("{%s}", k), v, -1)
	}
	return s
}

// 20190730 -> [2019073000, ..., 2019073023]
// 2019073001 -> [2019073001]
func TimeExpand(raw string) []string {
	if len(raw) != 8 && len(raw) != 10 {
		panic(fmt.Sprintf("invalid time: %s", raw))
	}

	var points []string
	if len(raw) == 8 {
		for i:=0; i < 24; i++ {
			points = append(points, fmt.Sprintf("%s%02d", raw, i))
		}
	} else {
		points = append(points, raw)
	}

	return points
}

