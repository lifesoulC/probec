package prober

import (
	"fmt"
	"probec/internal/addr"
	"probec/netio"
	"time"
)

type Prober struct {
	src              []string
	io               *netio.NetIO
	icmpResults      *icmpResultsType
	icmpBroadResults *icmpBroadResultsType
}

func NewProber(src []string) (p *Prober, e error) {
	p = &Prober{}
	p.src = src
	p.io, e = netio.NewNetIO(src)
	if e != nil {
		return
	}
	p.icmpResults = newIcmpResults()
	p.icmpBroadResults = newIcmpBroadResults()
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
	fmt.Printf("ping %s from %s \n", opts.dest.String, opts.src.String)
	p.icmpResults.beginWait(opts.src, opts.dest)

	for i := 0; i < opts.Count; i++ {
		p.io.SendPing(opts.src, opts.dest)
		time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
	}

	delays = p.icmpResults.endWait(opts.src, opts.dest, 2)
	return
}

func (p *Prober) BroadCastPing(opts *IcmpBroadcastOpts) (ret []*DestDelays, e error) {
	opts.src, e = addr.FromString(opts.Src)
	if e != nil {
		return
	}
	opts.dest, e = addr.FromString(opts.Dest)
	if e != nil {
		return
	}
	fmt.Printf("broadcast %s from %s \n", opts.dest.String, opts.src.String)
	p.icmpBroadResults.beginWait(opts.src, opts.dest)
	for i := 0; i < opts.Count; i++ {
		p.io.SendPingBroadcast(opts.src, opts.dest)
		time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
	}

	delays := p.icmpBroadResults.endWait(opts.src, opts.dest, 1)
	ret = searchBroadcastDelays(opts.dest, delays)
	return
}

func searchBroadcastDelays(src *addr.IPAddr, delays []*icmpBroadResultType) []*DestDelays {
	less := &DestDelays{}
	equal := &DestDelays{}
	greater := &DestDelays{}
	for _, result := range delays {
		if result == nil {
			continue
		}
		if src.Equal(result.dest) {
			equal.Dest = result.dest
			equal.Delays = append(equal.Delays, result.delays...)
			continue
		}

		if result.dest.Less(src) {
			if less.Dest == nil {
				less.Dest = result.dest
				less.Delays = make([]int, 0)
				less.Delays = append(less.Delays, result.delays...)
				continue
			}
			if result.dest.Great(less.Dest) {
				less.Dest = result.dest
				less.Delays = make([]int, 0)
				less.Delays = append(less.Delays, result.delays...)
				continue
			}
		}

		if result.dest.Great(src) {
			if greater.Dest == nil {
				greater.Dest = result.dest
				greater.Delays = make([]int, 0)
				greater.Delays = append(greater.Delays, result.delays...)
				continue
			}

			if result.dest.Less(greater.Dest) {
				greater.Dest = result.dest
				greater.Delays = make([]int, 0)
				greater.Delays = append(greater.Delays, result.delays...)
				continue
			}
		}
	}

	result := []*DestDelays{less, equal, greater}
	return result
}
