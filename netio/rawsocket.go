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

func newUDPSocket(laddr string) (*udpSocket, error) { //创建udp套接字
	a, err := addrIpArray(laddr) //将字符Ip转换为标准ip
	if err != nil {
		return nil, err
	}
	s := &udpSocket{}
	s.laddr = laddr
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP) //CreateSocket
	if err != nil {
		return nil, errors.New(laddr + " create udp socket:" + err.Error())
	}

	err = syscall.Bind(s.fd, &syscall.SockaddrInet4{Port: localUDPPort, Addr: a}) //BindSocket
	if err != nil {
		return nil, errors.New(laddr + " bind udp socket:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, 1024*1024) //setbuf
	if err != nil {
		return nil, errors.New(laddr + " set udp socket send buff:" + err.Error())
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
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP) //CreateSocket
	if err != nil {
		return nil, errors.New(laddr + " create icmp socket:" + err.Error())
	}

	err = syscall.Bind(s.fd, &syscall.SockaddrInet4{Port: 0, Addr: a}) //BindSocket
	if err != nil {
		return nil, errors.New(laddr + " bind icmp socket:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, 1024*1024) //setsendbuf
	if err != nil {
		return nil, errors.New(laddr + " set icmp socket send buff:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 1024*1024) //SetRcvbuf
	if err != nil {
		return nil, errors.New(laddr + " set icmp socket recv buff:" + err.Error())
	}

	return s, nil
}

func newRecvSocket() (s *icmpSocket, err error) {
	s = &icmpSocket{}
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP) //CreateSocket
	if err != nil {
		return nil, errors.New("create icmp recieve socket:" + err.Error())
	}

	err = syscall.SetsockoptInt(s.fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 1024*1024) //SetRcvbuf
	if err != nil {
		return nil, errors.New("set icmp  recv buff:" + err.Error())
	}
	return s, nil
}
