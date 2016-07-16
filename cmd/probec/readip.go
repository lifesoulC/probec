package main

import (
//	"bufio"
	"net"
//	"os"
  "fmt"
)

var srcIP []string

func readIPFile() (src []string, err error) { //从ip.txt中读取信息
	//	f, e := os.Open("ip.txt")
	//	if e != nil {
	//		err = e
	//		return
	//	}
	//	scanner := bufio.NewScanner(f)
	//	for scanner.Scan() {
	//		line := scanner.Text() //读出其中的文本段
	//		src = append(src, line)

	//	}
	//	defer f.Close()
	addrs, errs := net.InterfaceAddrs()

	if errs != nil {
		fmt.Println(errs)
		err = errs
		return
	}
//	fmt.Println(time.Now())
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				src = append(src, ipnet.IP.String())
			}

		}
	}
	return
}
