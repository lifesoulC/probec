package main

import (
	"fmt"
	"probec/prober"
)

var prob *prober.Prober

func main() {

	src, e := readIPFile()
	if e != nil {
		fmt.Println(e)
		return
	}
	for _, s := range src {
		fmt.Println(s)
	}

	prob, e = prober.NewProber(src)
	if e != nil {
		fmt.Println(e)
		return
	}

	e = StartHTTP(":8088")
	if e != nil {
		fmt.Println(e)
	}
}
