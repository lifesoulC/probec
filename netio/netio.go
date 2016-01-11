package netio

import (
	"errors"
	"fmt"
	"os"
)

const (
	pktTypeICMPEcho = 1
	pktTypeUDP      = 2
)

type WriteOpt struct {
	Src  string
	Dest string
	TTL  int
	Data []byte
	typ  int
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

type NetIO struct {
	writeChan   chan *WriteOpt
	writeSocket []*ipSocket
	readSocket  *icmpSocket
}

func NewNetIO(srcAddrs []string) *NetIO {
	io := &NetIO{}
	io.writeChan = make(chan *WriteOpt, 1024)

	for _, addr := range srcAddrs {
		s, e := newIpSocket(addr)
		if e == nil {
			io.writeSocket = append(io.writeSocket, s)
		} else {
			fmt.Println(e)
		}
	}
	io.readSocket, _ = newIcmpSocket()
	go io.readRoutine()
	return io
}

func (io *NetIO) writeRoutine() {

}

func (io *NetIO) readRoutine() {
	for {
		p, r, e := io.readSocket.readFrom()
		if e != nil {
			fmt.Println(e)
			continue
		}
		fmt.Println(p, r)

	}
}
