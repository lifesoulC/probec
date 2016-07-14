package main

import (
	"fmt"
	"os"
	"probec/prober"
	"strconv"
)

var prob *prober.Prober

func main() {

	src, e := readIPFile() //从ip.txt中得到源IP
	if e != nil {
		fmt.Println(e)
		return
	}
	srcIP = src
	for _, s := range srcIP { //打印出内容
		//	addr //??
		fmt.Println(s)
	}
	// var e error
	prob, e = prober.NewProber(src) //初始化Prober结构体 在prober.go中

	if e != nil {
		fmt.Println(e)
		return
	}
	if len(os.Args) < 2 { //判断是否开启端口号
		fmt.Println("usage: probec [port]")
		return
	}

	port, e := strconv.Atoi(os.Args[1]) //将端口号赋值于 port
	if e != nil {
		fmt.Println("invalid port number")
		return
	}

	listenPort := fmt.Sprintf(":%d", port) //设置server监听端口号
	fmt.Println("listen port", port)
	e = StartHTTP(listenPort) //开启http服务器
	if e != nil {
		fmt.Println(e)
	}
}
