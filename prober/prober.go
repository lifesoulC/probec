package prober

import (
	"net"
	"probec/netio"
	"time"
)

type Prober struct {
	src []string
	io  *netio.NetIO
}

type PingOpts struct {
	Src      string
	Dest     string
	DestIP   string
	Count    int
	ip       [4]byte
	Interval int
	Delays   []int
}

func NewProber(src []string) (p *Prober, e error) {
	p = &Prober{}
	p.src = src
	p.io, e = netio.NewNetIO(src)
	if e != nil {
		return
	}
	return
}

// func (p *Prober) Ping(laddr string, raddr string) {
// 	p.io.SendPing(laddr, raddr)

// }

func (p *Prober) ICMPPing(opts *PingOpts) {
	radd, e := net.ResolveIPAddr("ip4", opts.Dest)
	if e != nil {
		return
	}

	opts.DestIP = radd.String()
	opts.ip = ipArray(radd.IP.To4())
	for i := 0; i < opts.Count; i++ {
		p.io.SendPingBroadcast(opts.Src, opts.ip)
		time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
	}
}

func ipArray(ip []byte) [4]byte {
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}
