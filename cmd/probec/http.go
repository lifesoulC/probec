package main

import (
	"net/http"
)

func StartHTTP(addr string) error {    //每次有http request的时候都会匹配“/"进行注册的函数。
	http.HandleFunc("/probe/ping", icmpPing)  //定义于request.go中
	http.HandleFunc("/probe/broadcast", icmpBroadcast)
	http.HandleFunc("/probe/trace", udpTrace)
	return http.ListenAndServe(addr, nil)
}
