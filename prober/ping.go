package prober

import (
	"fmt"
	"net"
	"probec/pkt"
)

type Pinger struct {
	src  *net.IPAddr
	dest *net.IPAddr
	conn net.Conn
}

func NewPinger() *Pinger {
	p := &Pinger{}
	return p
}

func (p *Pinger) Type() int {
	return PingType
}

func (p *Pinger) Probe(laddr *net.IPAddr, raddr *net.IPAddr) (delays []int, err error) {
	p.conn, err = net.DialIP("ip:icmp", laddr, raddr)
	if err != nil {
		return delays, err
	}
	defer p.conn.Close()

	echo := pkt.NewEchoRequest(pid, 1)

	p.conn.Write(echo)

	b := make([]byte, 1024)

	size, _ := p.conn.Read(b)
	fmt.Println(size)

	return
}

func (p *Pinger) ProbeHost(src string, dest string) (delays []int, err error) {
	laddr, err := net.ResolveIPAddr("ip", src)
	if err != nil {
		return delays, err
	}

	raddr, err := net.ResolveIPAddr("ip", dest)
	if err != nil {
		return delays, err
	}
	return p.Probe(laddr, raddr)
}
