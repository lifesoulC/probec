package netio

import (
	"encoding/binary"
	"fmt"
	"probec/internal/addr"
	"syscall"
	"time"
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
			io.onTraceReply(p)
			continue
		}
	}
}

func (io *NetIO) onIcmpReply(data []byte) {
	t := time.Now()
	dest := addr.FromSlice(data[12:])
	src := addr.FromSlice(data[16:])
	reply := parseIcmpEchoReply(data[20:])
	if reply.id != uint16(pid) {
		return
	}
	resp := &PingResp{}
	resp.Src = src
	resp.Dest = dest
	t1 := int64(binary.LittleEndian.Uint64(data[28:]))
	if t1 < 0 {
		return
	}
	delay := (t.UnixNano() - t1) / 1000
	if delay < 0 || delay > 100000000 {
		return
	}
	resp.Delay = int(delay)

	if io.handler != nil {
		if reply.seq < icmpSeqMax {
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
	if len(data) < 72 {
		return
	}

	hostIP := binary.BigEndian.Uint32(data[12:])
	srcIP := binary.BigEndian.Uint32(data[16:])
	destIP := binary.BigEndian.Uint32(data[44:])
	srcPort := binary.BigEndian.Uint16(data[48:])
	destPort := binary.BigEndian.Uint16(data[50:])
	// pid := binary.LittleEndian.Uint32(data[56:])
	ttl := binary.LittleEndian.Uint32(data[60:])
	sendStamp := binary.LittleEndian.Uint64(data[64:])

	src := addr.FromInt(srcIP)
	dest := addr.FromInt(destIP)
	host := addr.FromInt(hostIP)

	if srcPort != localUDPPort || destPort != remoteUDPPort {
		return
	}

	nowStamp := time.Now().UnixNano()
	delay := (nowStamp - int64(sendStamp)) / 1000
	if delay < 0 || delay > 10000000 {
		return
	}

	resp := &TTLResp{}
	resp.Src = src
	resp.Dest = dest
	resp.Host = host
	resp.TTL = int(ttl)
	resp.Delay = int(delay)

	io.handler.OnRecvTTL(resp)
	// fmt.Println(src.String, dest.String, host.String, pid, ttl, sendStamp, srcPort, destPort, delay)
}
