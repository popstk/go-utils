package main

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.90 Safari/537.36"

type CustomTransport struct {
	T http.RoundTripper
}

func (adt *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	return adt.T.RoundTrip(req)
}

func NewCustomTransport(T http.RoundTripper) *CustomTransport {
	if T == nil {
		T = http.DefaultTransport
	}
	return &CustomTransport{T}
}

func NewHTTPClient() *http.Client {
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	transport := NewCustomTransport(httpTransport)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:       cookieJar,
		Transport: transport,
	}

	return client
}