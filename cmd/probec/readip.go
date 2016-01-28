package main

import (
	"bufio"
	"os"
)

var srcIP []string

func readIPFile() (src []string, err error) {
	f, e := os.Open("ip.txt")
	if e != nil {
		err = e
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		src = append(src, line)

	}
	return
}
