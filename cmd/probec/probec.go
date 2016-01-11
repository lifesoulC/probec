package main

import (
	"fmt"
	"probec/netio"
)

func main() {
	fmt.Println("hello, probe c")
	srcAddrs := []string{"192.168.2.6"}
	netio.NewNetIO(srcAddrs)
	c := make(chan int)
	c <- 1
}
