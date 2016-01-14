package netio

import (
	"fmt"
	"probec/internal/addr"
	"syscall"
	"time"
)

func (io *NetIO) SendPing(src *addr.IPAddr, dest *addr.IPAddr) {
	socket := io.getIcmpSock(src.String)
	if socket == nil {
		fmt.Println(src.String, "not a local ip")
		return
	}
	opt := &icmpOpts{}
	opt.dest = dest.Array
	opt.sock = socket
	opt.broad = false
	io.icmpChan <- opt
}

func (io *NetIO) SendPingBroadcast(src *addr.IPAddr, dest *addr.IPAddr) {
	socket := io.getIcmpSock(src.String)
	if socket == nil {
		fmt.Println(src.String, "not a local ip")
		return
	}

	for i := 1; i < 255; i++ {
		raddr := dest.Array
		raddr[3] = byte(i)
		opt := &icmpOpts{}
		opt.sock = socket
		opt.dest = dest.Array
		opt.broad = true
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
	if opts.broad {
		opts.data = buildIcmpBroadcast()
	} else {
		opts.data = buildIcmpEchoRequest()
	}

	e := syscall.Sendto(opts.sock.fd, opts.data, 0, &syscall.SockaddrInet4{Port: 0, Addr: opts.dest})
	if e != nil {
		fmt.Println("send to", opts.dest, e.Error())
		return
	}

	time.Sleep(10 * time.Microsecond)
}

func (io *NetIO) sendRoutine() {
	for {
		select {
		case icmpOpts := <-io.icmpChan:
			io.sendIcmp(icmpOpts)
		}
	}
}
