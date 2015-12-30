package misc

import (
	"net"
	"syscall"
)

func CreateSocket(addr *net.IPAddr) (fd int, err error) {
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

func CreateICMPSocket(addr *net.IPAddr) (fd int, err error) {
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

func CloseSocket(fd int) {
	syscall.Close(fd)
}
