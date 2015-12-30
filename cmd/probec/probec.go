package main

import (
	"fmt"
	"probec/prober"
)

func main() {
	fmt.Println("hello, probe c")
	p := prober.NewPinger()
	_, err := p.ProbeHost("0.0.0.0", "114.114.114.114")
	if err != nil {
		fmt.Println(err)
	}
}
