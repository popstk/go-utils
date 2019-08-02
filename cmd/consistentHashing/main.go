package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// https://leileiluoluo.com/posts/consistent-hashing-and-high-available-cluster-proxy.html

func main() {
	nodes := []Node{
		{"0", "http://10.10.1.10/"},
		{"1", "http://10.10.1.11/"},
		{"2", "http://10.10.1.12/"},
	}
	ch := NewConsistentHashing(nodes, 100)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sign := r.Header.Get("sign")
		node := ch.GetNode(sign)
		uri, _ := url.Parse(node.Address)
		httputil.NewSingleHostReverseProxy(uri)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}