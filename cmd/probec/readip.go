package main

import (
	"bufio"
	"os"
)

var srcIP []string

func readIPFile() (src []string, err error) {   //从ip.txt中读取信息
	f, e := os.Open("ip.txt")
	if e != nil {
		err = e
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()      //读出其中的文本段
		src = append(src, line)

	}
	return
}
