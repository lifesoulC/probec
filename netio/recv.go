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
		p := make([]byte, 128)
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
