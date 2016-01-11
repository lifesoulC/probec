package netio

import (
	"encoding/binary"
	"fmt"
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
	destIP := fmt.Sprintf("%d.%d.%d.%d", data[12], data[13], data[14], data[15])
	srcIP := fmt.Sprintf("%d.%d.%d.%d", data[16], data[17], data[18], data[19])
	reply := parseIcmpEchoReply(data[20:])
	if reply.id != uint16(pid) {
		return
	}
	resp := &PingResp{}
	resp.Data = reply.data
	resp.Stamp = t
	resp.Laddr = srcIP
	resp.Raddr = destIP
	t1 := int64(binary.LittleEndian.Uint64(data[28:]))
	if t1 < 0 {
		return
	}
	delay := (t.UnixNano() - t1) / 1000
	if delay < 0 || delay > 100000000 {
		return
	}
	resp.Delay = delay
	fmt.Println(resp.Laddr, "->", resp.Raddr, ":", resp.Delay)

	if io.handler != nil {
		io.handler.OnRecvPing(resp)
	}
}
