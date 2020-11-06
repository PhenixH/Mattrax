package main

import (
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout:   time.Second * 10,
	Transport: &customTransport{http.DefaultTransport},
}

type customTransport struct {
	T http.RoundTripper
}

func (ct *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "MattraxAgent "+Version)
	return ct.T.RoundTrip(req)
}
