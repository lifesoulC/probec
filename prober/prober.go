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
	traceResults     *traceResultsType
}

func NewProber(src []string) (p *Prober, e error) {
	p = &Prober{}
	p.src = src
	p.io, e = netio.NewNetIO(src) //在netio中实现  将netio初始化
	if e != nil {
		return
	}
	p.icmpResults = newIcmpResults() //在results中实现
	p.icmpBroadResults = newIcmpBroadResults()
	p.traceResults = newTraceResults()
	p.io.SetHandler(p)
	return
}

func (p *Prober) ICMPPing(opts *PingOpts) (delays []int, e error) { //返回延迟

	opts.src, e = addr.FromString(opts.Src) // 源ip字符序列转换和 转换IPv4 定义在 ipaddr.go 中
	if e != nil {
		return
	}
	opts.dest, e = addr.FromString(opts.Dest) //目的ip字符序列转换 ...
	if e != nil {
		return
	}
	fmt.Printf("ping %s from %s \n", opts.dest.String, opts.src.String)
	p.icmpResults.beginWait(opts.src, opts.dest) //将源和目的地址整合到一起 map[uint64][]int 中

	for i := 0; i < opts.Count; i++ { //向通讯管道中发送要测的源和目的ip 发送次数为count次
		p.io.SendPing(opts.src, opts.dest)
		time.Sleep(time.Duration(opts.Interval) * time.Millisecond) //每次发送延迟为 interval
	}

	delays = p.icmpResults.endWait(opts.src, opts.dest, 500) //在results.go中定义
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

	delays := p.icmpBroadResults.endWait(opts.src, opts.dest, 500)
	ret = searchBroadcastDelays(opts.dest, delays)
	return
}

func (p *Prober) Trace(opts *TraceOpts) (delays []*TraceResultType, e error) {
	opts.src, e = addr.FromString(opts.Src)
	if e != nil {
		return
	}
	opts.dest, e = addr.FromString(opts.Dest)
	if e != nil {
		return
	}
	fmt.Printf("trace %s from %s \n", opts.dest.String, opts.src.String)

	p.traceResults.beginWait(opts.src, opts.dest)
	for i := 0; i < opts.Count; i++ {
		p.io.SendTTL(opts.src, opts.dest, 64)
		time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
	}
	delays = p.traceResults.endWait(opts.src, opts.dest, 500) //在results.go中定义
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
