package netio

import (
	"encoding/binary"
	"time"
)

var (
	seq          uint16
	broadcastSeq uint16
)

func init() {
	seq = icmpSeqMin
	broadcastSeq = icmpBroadMin
}

type icmpEchoType struct {
	typ      uint8
	code     uint8
	checksum uint16
	id       uint16
	seq      uint16
	payload  []byte
}

func (echo *icmpEchoType) marshal() []byte {
	p := make([]byte, 64)
	p[0] = echo.typ
	p[1] = echo.code
	binary.BigEndian.PutUint16(p[4:], echo.id)
	binary.BigEndian.PutUint16(p[6:], echo.seq)
	t := time.Now().UnixNano()
	binary.LittleEndian.PutUint64(p[8:], uint64(t))
	c := checkSum(p)
	binary.BigEndian.PutUint16(p[2:], c)
	return p
}

func buildIcmpEchoRequest() ([]byte, uint16) {
	icmp := icmpEchoType{}
	icmp.typ = 8
	icmp.seq = seq + 1
	icmp.id = uint16(pid)
	seq += 1
	if seq >= icmpSeqMax {
		seq = icmpSeqMin
	}
	return icmp.marshal(), icmp.seq
}

func buildIcmpBroadcast() ([]byte, uint16) {
	icmp := icmpEchoType{}
	icmp.typ = 8
	icmp.seq = broadcastSeq + 1
	icmp.id = uint16(pid)
	broadcastSeq += 1
	if broadcastSeq >= icmpBroadMax {
		broadcastSeq = icmpBroadMin
	}
	return icmp.marshal(), icmp.seq
}

type icmpHeadType struct {
	typ      uint8
	code     uint8
	checksum uint16
	data     []byte
}

type icmpEchoReply struct {
	typ      uint8
	code     uint8
	checksum uint16
	id       uint16
	seq      uint16
	data     []byte
}

func parseIcmpHeader(data []byte) (h *icmpHeadType) {
	h = &icmpHeadType{}
	if len(data) < 4 {
		return
	}
	h.typ = data[0]
	h.code = data[1]
	h.checksum = binary.BigEndian.Uint16(data[2:])
	return h
}

func parseIcmpEchoReply(data []byte) (reply *icmpEchoReply) {
	reply = &icmpEchoReply{}
	if len(data) < 8 {
		return
	}
	reply.typ = data[0]
	reply.code = data[1]
	reply.checksum = binary.BigEndian.Uint16(data[2:])
	reply.id = binary.BigEndian.Uint16(data[4:])
	reply.seq = binary.BigEndian.Uint16(data[6:])
	if len(data) > 8 {
		reply.data = data[8:]
	}
	return reply
}

func buildUDP(ttl int) []byte {
	b := make([]byte, 64)
	t := time.Now().UnixNano()
	binary.LittleEndian.PutUint32(b, uint32(pid))
	binary.LittleEndian.PutUint32(b[4:], uint32(ttl))
	binary.LittleEndian.PutUint64(b[8:], uint64(t))
	return b
}
