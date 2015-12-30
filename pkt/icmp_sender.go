package pkt

import (
	"net"
	"syscall"
)

type sendOpt struct {
	ttl  int
	dest *net.IPAddr
	data []byte
}

type IPSender struct {
	icmpFD []int
	ch     chan *sendOpt
}

func NewIPSender(src []string) (sender *IPSender, err error) {
	sender = &IPSender{}
	addrs, err := validIPAddr(src)
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		fd, err := createSocket(addr)
		if err != nil {
			return nil, ipErrorFromError(addr.String(), err)
		}
		sender.icmpFD = append(sender.icmpFD, fd)
	}

	sender.ch = make(chan *sendOpt, 1024)
	return
}

func (s *IPSender) WriteTo(data []byte, to *net.IPAddr) {
	s.TTLWriteTo(data, to, 64)
}

func (s *IPSender) TTLWriteTo(data []byte, to *net.IPAddr, ttl int) {
	opts := &sendOpt{}
	opts.data = data
	opts.dest = to
	opts.ttl = ttl
	s.ch <- opts
}

func (s *IPSender) sendRoutine() {

}

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
