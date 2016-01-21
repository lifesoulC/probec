package prober

import (
	"probec/netio"
)

func (prober *Prober) OnRecvPing(pkt *netio.PingResp) {
	// fmt.Println(pkt.Src.String, "->", pkt.Dest.String, pkt.Delay)
	prober.icmpResults.addResult(pkt.Src, pkt.Dest, pkt.Delay)
}

func (prober *Prober) OnRecvTTL(pkt *netio.TTLResp) {
	// fmt.Println(pkt.Src.String, "->", pkt.Dest.String, pkt.Delay)

}

func (prober *Prober) OnRecvICMPBroadcast(pkt *netio.PingResp) {
	// fmt.Println(pkt.Src.String, "->", pkt.Dest.String, pkt.Delay)
	prober.icmpBroadResults.addResult(pkt.Src, pkt.Dest, pkt.Delay)
}
