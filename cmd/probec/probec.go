package main

import (
	"fmt"
	"probec/prober"
)

func main() {
	srcAddrs := []string{"10.0.0.166"}
	probe, e := prober.NewProber(srcAddrs)
	if e != nil {
		fmt.Println(e)
		return
	}
	pingOpts := &prober.PingOpts{}
	pingOpts.Src = "10.0.0.166"
	pingOpts.Dest = "baidu.com"
	pingOpts.Count = 10
	pingOpts.Interval = 2
	probe.ICMPPing(pingOpts)

	c := make(chan int)
	c <- 1
}
