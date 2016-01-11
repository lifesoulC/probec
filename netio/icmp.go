package netio

import (
	"encoding/binary"
)

var (
	seq uint16
)

type icmpEchoType struct {
	typ      uint8
	code     uint8
	checksum uint16
	id       uint16
	seq      uint16
}

func (echo *icmpEchoType) marshal() []byte {
	p := make([]byte, 64)
	p[0] = echo.typ
	p[1] = echo.code
	binary.BigEndian.PutUint16(p[4:], echo.id)
	binary.BigEndian.PutUint16(p[6:], echo.seq)
	c := checkSum(p)
	binary.BigEndian.PutUint16(p[2:], c)
	return p
}

func buildIcmpEchoRequest(raddr string) []byte {
	icmp := icmpEchoType{}
	icmp.typ = 8
	icmp.seq = seq + 1

	seq += 1
	if seq > 30000 {
		seq = 0
	}

	return icmp.marshal()
}
