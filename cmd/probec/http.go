package main

import (
	"net/http"
)

func StartHTTP(addr string) error {
	http.HandleFunc("/probe/ping", icmpPing)
	http.HandleFunc("/probe/broadcast", icmpBroadcast)
	return http.ListenAndServe(addr, nil)
}
