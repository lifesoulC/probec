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
	pingOpts.Src = "192.168.199.138"
	pingOpts.Dest = "sina.com.cn"
	pingOpts.Count = 100
	pingOpts.Interval = 200
	probe.ICMPPing(pingOpts)

	c := make(chan int)
	c <- 1
}
