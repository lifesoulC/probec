package pkt

import (
	"errors"
	"fmt"
	"net"
)

func ipErrorFromError(ip string, err error) error {
	msg := fmt.Sprintf("%s: %s", ip, err.Error())
	return errors.New(msg)
}

func ipErrorFromMsg(ip string, msg string) error {
	msg = fmt.Sprintf("%s: %s", ip, msg)
	return errors.New(msg)
}

func validIPAddr(src []string) (addrs []*net.IPAddr, err error) {
	for _, addr := range src {
		ipaddr, err := net.ResolveIPAddr("ip4", addr)
		if err != nil {
			return addrs, ipErrorFromError(addr, err)
		}

		ip := ipaddr.IP.To4()
		if ip.IsLoopback() {
			err = ipErrorFromMsg(addr, "loopback addr or broadcast addr")
			return addrs, err
		}

		addrs = append(addrs, ipaddr)
	}
	return
}

func ipSlice(ip [4]byte) []byte {
	return []byte{ip[0], ip[1], ip[2], ip[3]}
}

func ipArrary(ip []byte) [4]byte {
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}
