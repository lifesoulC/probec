package netio

import (
	"encoding/binary"
	"fmt"
	"probec/internal/addr"
	"syscall"
)

func (io *NetIO) recvRoutine() {
	for {
		p := make([]byte, 256)
		size, _, e := syscall.Recvfrom(io.recvSocket.fd, p, 0)
		if size < 36 {
			continue
		}

		if e != nil {
			fmt.Print("recv data error", e)
			continue
		}
		p = p[:size]

		icmp := parseIcmpHeader(p[20:])
		if icmp.typ == 0 {
			io.onIcmpReply(p)
			continue
		}

		if icmp.typ == 11 || icmp.typ == 3 {
			fmt.Println("trace pkt", len(p))
			io.onTraceReply(p)
			continue
		}
	}
}

func (io *NetIO) onIcmpReply(data []byte) {
	dest := addr.FromSlice(data[12:])
	src := addr.FromSlice(data[16:])
	reply := parseIcmpEchoReply(data[20:])
	if reply.id != uint16(pid) {
		return
	}

	keyOpts := &pktsMapOptsICMP{}
	keyOpts.dst = dest.Int()
	keyOpts.src = src.Int()
	keyOpts.id = uint16(pid)
	keyOpts.seq = reply.seq

	_, d := io.pkts.getIcmpDelay(keyOpts)

	resp := &PingResp{}
	resp.Src = src
	resp.Dest = dest

	resp.Delay = d

	if io.handler != nil {
		if reply.seq <= icmpSeqMax {
			io.handler.OnRecvPing(resp)
			return
		}
		if reply.seq < icmpBroadMax {
			io.handler.OnRecvICMPBroadcast(resp)
			return
		}
	}
}

func (io *NetIO) onTraceReply(data []byte) {
	if len(data) < 56 {
		return
	}

	hostIP := binary.BigEndian.Uint32(data[12:])
	srcIP := binary.BigEndian.Uint32(data[16:])
	destIP := binary.BigEndian.Uint32(data[44:])
	srcPort := binary.BigEndian.Uint16(data[48:])
	destPort := binary.BigEndian.Uint16(data[50:])

	src := addr.FromInt(srcIP)
	dest := addr.FromInt(destIP)
	host := addr.FromInt(hostIP)

	keyOpts := &pktsMapOptsTTL{}
	keyOpts.dst = destIP
	keyOpts.dstPort = destPort
	keyOpts.src = srcIP
	keyOpts.srcPort = srcPort

	_, ttl, delay := io.pkts.getTTLDelay(keyOpts)

	resp := &TTLResp{}
	resp.Src = src
	resp.Dest = dest
	resp.Host = host
	resp.TTL = int(ttl)
	resp.Delay = int(delay)

	io.handler.OnRecvTTL(resp)
	// fmt.Println(src.String, dest.String, host.String, pid, ttl, sendStamp, srcPort, destPort, delay)
}
