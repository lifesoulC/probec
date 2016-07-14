package netio

import (
	"errors"
	"os"
	"probec/internal/addr"
	"time"
)

const (
	pktTypeICMPEcho  = 1
	pktTypeUDP       = 2
	localUDPPort     = 33333
	remoteUDPPortMin = 33000
	remoteUDPPortMax = 48000
	icmpSeqMin       = 1
	icmpSeqMax       = 5000
	icmpBroadMin     = 5001
	icmpBroadMax     = 10000
)

var (
	// remoteUDPPort
	remoteUDPPort uint16 = remoteUDPPortMin
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
	sock   *icmpSocket
	broad  bool
	srcInt uint32
	dstInt uint32
	dest   [4]byte
	data   []byte
}

type ttlOpts struct {
	sock   *udpSocket
	srcInt uint32
	dstInt uint32
	dest   [4]byte
	ttl    int
	data   []byte
}

type NetIO struct {
	udpSocks   []*udpSocket
	icmpSocks  []*icmpSocket
	icmpMap    map[string]*icmpSocket //存放icmp socket
	udpMap     map[string]*udpSocket  //存放udp  socket
	recvSocket *icmpSocket
	icmpChan   chan *icmpOpts
	ttlChan    chan *ttlOpts
	pkts       *pktMap
	lastClear  time.Time
	handler    NetIOHandler
}

func NewNetIO() (*NetIO, error) { //一个IP绑定连个socket 一个UDP 一个ICMP
	io := &NetIO{}
	//	for _, addr := range srcAddrs {               //依次绑定本地IP
	//		udp, e := newUDPSocket(addr)                 //将源地址绑定socket放入udpSocket中
	//		if e != nil {
	//			return nil, e

	//		} else {
	//			io.udpSocks = append(io.udpSocks, udp)    //添加进队列
	//		}

	//		icmp, e := newIcmpSocket(addr)             //将原地址依次绑定socket放入udpsocket队列中
	//		if e != nil {
	//			return nil, e

	//		} else {
	//			io.icmpSocks = append(io.icmpSocks, icmp)       //添加进队列
	//		}
	//	}
	io.icmpMap = make(map[string]*icmpSocket)
	io.udpMap = make(map[string]*udpSocket)

	io.pkts = newPktMap() //在 stamp.go中定义
	io.lastClear = time.Now()

	recv, e := newRecvSocket() //创建接收socket icmp接口
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

	if v, ok := io.icmpMap[addr]; ok {
		return v
	} else {
		icmp, e := newIcmpSocket(addr) //将原地址依次绑定socket放入udpsocket队列中
		if e != nil {
			return nil

		} else {
			io.icmpMap[addr] = icmp //添加Map
			return icmp
		}
	}

	//	for _, s := range io.icmpSocks {
	//		if s.laddr == addr {
	//			return s
	//		}
	//	}
	//	return nil
}

func (io *NetIO) getUdpSock(addr string) *udpSocket {

	if v, ok := io.udpMap[addr]; ok {
		return v
	} else {
		udp, e := newUDPSocket(addr) //将原地址依次绑定socket放入udpsocket队列中
		if e != nil {
			return nil

		} else {
			io.udpMap[addr] = udp //添加Map
			return udp
		}
	}
	//	for _, s := range io.udpSocks {
	//		if s.laddr == addr {
	//			return s
	//		}
	//	}
	//	return nil
}
