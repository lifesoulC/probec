package main

import (
	"fmt"
	"probec/prober"
)

func main() {

	src, e := readIPFile()
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(src)

	probe, e := prober.NewProber(src)
	if e != nil {
		fmt.Println(e)
		return
	}
	pingOpts := &prober.PingOpts{}
	pingOpts.Src = "10.0.0.166"
	pingOpts.Dest = "baidu.com"
	pingOpts.Count = 5
	pingOpts.Interval = 200
	probe.ICMPPing(pingOpts)

	probe.ICMPPing(pingOpts)
	c := make(chan int)
	c <- 1
}
