package main

import (
	"fmt"
	"os"
	"probec/prober"
	"strconv"
)

var prob *prober.Prober

func main() {

	src, e := readIPFile()
	if e != nil {
		fmt.Println(e)
		return
	}
	srcIP = src
	for _, s := range srcIP {
		fmt.Println(s)
	}

	prob, e = prober.NewProber(src)
	if e != nil {
		fmt.Println(e)
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("usage: probec [port]")
		return
	}

	port, e := strconv.Atoi(os.Args[1])
	if e != nil {
		fmt.Println("invalid port number")
		return
	}

	listenPort := fmt.Sprintf(":%d", port)
	fmt.Println("listen port", port)
	e = StartHTTP(listenPort)
	if e != nil {
		fmt.Println(e)
	}
}
