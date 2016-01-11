package netio

import (
	"syscall"
)

type ipSocket struct {
	laddr string
	fd    int
}

type icmpSocket struct {
	fd int
}

func newIpSocket(laddr string) (*ipSocket, error) {
	a, err := addrIpArray(laddr)
	if err != nil {
		return nil, err
	}
	s := &ipSocket{}
	s.laddr = laddr
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, err
	}

	err = syscall.Bind(s.fd, &syscall.SockaddrInet4{Port: 0, Addr: a})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func newIcmpSocket() (s *icmpSocket, err error) {
	s = &icmpSocket{}
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *icmpSocket) readFrom() (p []byte, r syscall.Sockaddr, e error) {
	p = make([]byte, 512)
	_, r, e = syscall.Recvfrom(s.fd, p, 0)
	return
}
