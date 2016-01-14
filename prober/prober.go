package prober

import (
	"fmt"
	"probec/internal/addr"
	"probec/netio"
	"time"
)

type Prober struct {
	src         []string
	io          *netio.NetIO
	icmpResults *icmpResultsType
}

type PingOpts struct {
	Src      string
	Dest     string
	src      *addr.IPAddr
	dest     *addr.IPAddr
	Count    int
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
	p.icmpResults = newIcmpResults()
	p.io.SetHandler(p)
	return
}

func (p *Prober) ICMPPing(opts *PingOpts) (delays []int, e error) {

	opts.src, e = addr.FromString(opts.Src)
	if e != nil {
		return
	}
	opts.dest, e = addr.FromString(opts.Dest)
	if e != nil {
		return
	}

	p.icmpResults.beginWait(opts.src, opts.dest)

	for i := 0; i < opts.Count; i++ {
		p.io.SendPing(opts.src, opts.dest)
		time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
	}

	delays = p.icmpResults.endWait(opts.src, opts.dest, 2)
	fmt.Println("delays:", delays)
	return
}

func ipArray(ip []byte) [4]byte {
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}
