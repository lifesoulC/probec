package main

import (
	"errors"
	"net"
	"fmt"
)

const (
	errSuccess = 1
	errSrcIP   = 2
	errJson    = 3
	errUnkown  = 99
)

func checkSrcIP(ip string) error {
	//	for _, v := range srcIP {
	//		if v == ip {
	//			return nil
	//		}
	//	}
	//	return errors.New("invalid src ip")
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		return errors.New("invalid src ip")
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ip == ipnet.IP.String() {
				return nil
			}
		}
	}
	return errors.New("invalid src ip")

}
