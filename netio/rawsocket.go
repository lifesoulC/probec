package netio

import (
	"errors"
	"syscall"
)

type udpSocket struct {
	laddr string
	fd    int
}

type icmpSocket struct {
	laddr string
	fd    int
}

type recvSock struct {
	fd int
}

func newUDPSocket(laddr string) (*udpSocket, error) {
	a, err := addrIpArray(laddr)
	if err != nil {
		return nil, err
	}
	s := &udpSocket{}
	s.laddr = laddr
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, errors.New(laddr + " create udp socket:" + err.Error())
	}

	err = syscall.Bind(s.fd, &syscall.SockaddrInet4{Port: localUDPPort, Addr: a})
	if err != nil {
		return nil, errors.New(laddr + " bind udp socket:" + err.Error())
	}
	return s, nil
}

func newIcmpSocket(laddr string) (*icmpSocket, error) {
	a, err := addrIpArray(laddr)
	if err != nil {
		return nil, err
	}
	s := &icmpSocket{}
	s.laddr = laddr
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, errors.New(laddr + " create icmp socket:" + err.Error())
	}

	err = syscall.Bind(s.fd, &syscall.SockaddrInet4{Port: 0, Addr: a})
	if err != nil {
		return nil, errors.New(laddr + " bind icmp socket:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, 1024*1024)
	if err != nil {
		return nil, errors.New(laddr + " set icmp socket send buff:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 1024*1024)
	if err != nil {
		return nil, errors.New(laddr + " set icmp socket recv buff:" + err.Error())
	}

	return s, nil
}

func newRecvSocket() (s *icmpSocket, err error) {
	s = &icmpSocket{}
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, errors.New("create icmp recieve socket:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 1024*1024)
	if err != nil {
		return nil, errors.New("set icmp  recv buff:" + err.Error())
	}
	return s, nil
}
