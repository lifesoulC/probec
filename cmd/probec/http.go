package main

import (
	"net/http"
)

func StartHTTP(addr string) error {
	http.HandleFunc("/probe/ping", icmpPing)
	return http.ListenAndServe(addr, nil)
}
