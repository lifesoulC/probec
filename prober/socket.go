package prober

import (
	"net"
	"syscall"
)

func createSocket(addr *net.IPAddr) (fd int, err error) {
	fd = -1
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
	if err != nil {
		return
	}

	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 0, Addr: ipArrary(addr.IP.To4())})
	if err != nil {
		return
	}
	return
}

func createICMPSocket(addr *net.IPAddr) (fd int, err error) {
	fd = -1
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return
	}

	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 0, Addr: ipArrary(addr.IP.To4())})
	if err != nil {
		return
	}
	return
}

func closeSocket(fd int) {
	syscall.Close(fd)
}

func ipSlice(ip [4]byte) []byte {
	return []byte{ip[0], ip[1], ip[2], ip[3]}
}

func ipArrary(ip []byte) [4]byte {
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}
