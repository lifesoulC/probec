package netio

import (
	"fmt"
	"syscall"
	"time"
)

func (io *NetIO) SendPing(laddr string, raddr [4]byte) {
	socket := io.getIcmpSock(laddr)
	if socket == nil {
		fmt.Println(laddr, "not a local ip")
		return
	}
	opt := &icmpOpts{}
	opt.dest = raddr
	opt.sock = socket
	io.icmpChan <- opt
}

func (io *NetIO) SendPingBroadcast(laddr string, raddr [4]byte) {
	socket := io.getIcmpSock(laddr)
	if socket == nil {
		fmt.Println(laddr, "not a local ip")
		return
	}

	for i := 1; i < 255; i++ {
		raddr[3] = byte(i)
		opt := &icmpOpts{}
		opt.sock = socket
		opt.dest = raddr
		io.icmpChan <- opt
	}

}

func (io *NetIO) SendTTL(laddr string, raddr string, data []byte, ttl int) {
	socket := io.getUdpSock(laddr)
	if socket == nil {
		fmt.Println(laddr, "not a local ip")
		return
	}
}

func (io *NetIO) sendIcmp(opts *icmpOpts) {
	opts.data = buildIcmpEchoRequest()
	e := syscall.Sendto(opts.sock.fd, opts.data, 0, &syscall.SockaddrInet4{Port: 0, Addr: opts.dest})
	if e != nil {
		fmt.Println("send to", opts.dest, e.Error())
		return
	}
	if io.handler != nil {
		req := &PingReq{}
		req.Laddr = opts.sock.laddr
		req.Raddr = opts.raddr
		io.handler.OnSendPing(req)
	}
	time.Sleep(500 * time.Microsecond)
}

func (io *NetIO) sendRoutine() {
	for {
		select {
		case icmpOpts := <-io.icmpChan:
			io.sendIcmp(icmpOpts)
		}
	}
}
