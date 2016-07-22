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
	chack := make(map[string]int)
	addrs, errs := net.InterfaceAddrs()
	if errs != nil {
		fmt.Println(errs)
		err = errs
		return
	}

	for i, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				chack[ipnet.IP.String()] = i
			}
		}
	}
	for s, _ := range chack {
		src = append(src, s)
	}
	return
}
