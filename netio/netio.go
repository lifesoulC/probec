package netio

import (
	"errors"
	"os"
	"probec/internal/addr"
)

const (
	pktTypeICMPEcho = 1
	pktTypeUDP      = 2
	localUDPPort    = 33333
	remoteUDPPort   = 33486
	icmpSeqMin      = 1
	icmpSeqMax      = 5000
	icmpBroadMin    = 5001
	icmpBroadMax    = 10000
)

type PingResp struct {
	Src   *addr.IPAddr
	Dest  *addr.IPAddr
	Delay int
}

type TTLResp struct {
	Src   *addr.IPAddr
	Dest  *addr.IPAddr
	Host  *addr.IPAddr
	TTL   int
	Delay int
}

type ICMPOpts struct {
	Src      *addr.IPAddr
	Dest     *addr.IPAddr
	Interval int
	Count    int
	ip       [4]byte
}

type NetIOHandler interface {
	OnRecvPing(*PingResp)
	OnRecvTTL(*TTLResp)
	OnRecvICMPBroadcast(*PingResp)
}

var (
	ErrNotLocalAddr = errors.New("not local addr")
)

var (
	pid = 0
)

func init() {
	pid = os.Getpid() & 0xffff

}

type icmpOpts struct {
	sock  *icmpSocket
	broad bool
	dest  [4]byte
	data  []byte
}

type ttlOpts struct {
	sock *udpSocket
	dest [4]byte
	ttl  int
	data []byte
}

type NetIO struct {
	udpSocks   []*udpSocket
	icmpSocks  []*icmpSocket
	recvSocket *icmpSocket
	icmpChan   chan *icmpOpts
	ttlChan    chan *ttlOpts
	handler    NetIOHandler
}

func NewNetIO(srcAddrs []string) (*NetIO, error) {
	io := &NetIO{}
	for _, addr := range srcAddrs {
		udp, e := newUDPSocket(addr)
		if e != nil {
			return nil, e

		} else {
			io.udpSocks = append(io.udpSocks, udp)
		}

		icmp, e := newIcmpSocket(addr)
		if e != nil {
			return nil, e

		} else {
			io.icmpSocks = append(io.icmpSocks, icmp)
		}
	}

	recv, e := newRecvSocket()
	if e != nil {
		return nil, e
	}
	io.recvSocket = recv

	io.icmpChan = make(chan *icmpOpts, 1024)
	io.ttlChan = make(chan *ttlOpts, 1024)

	go io.sendRoutine()
	go io.recvRoutine()

	return io, nil
}

func (io *NetIO) SetHandler(h NetIOHandler) {
	io.handler = h
}
func (io *NetIO) getIcmpSock(addr string) *icmpSocket {
	for _, s := range io.icmpSocks {
		if s.laddr == addr {
			return s
		}
	}
	return nil
}

func (io *NetIO) getUdpSock(addr string) *udpSocket {
	for _, s := range io.udpSocks {
		if s.laddr == addr {
			return s
		}
	}
	return nil
}
