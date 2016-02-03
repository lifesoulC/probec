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
		opt.dest = raddr
		opt.broad = true
		io.icmpChan <- opt
		time.Sleep(5 * time.Millisecond)
	}

}

func (io *NetIO) SendTTL(src *addr.IPAddr, dest *addr.IPAddr, ttl int) {
	socket := io.getUdpSock(src.String)
	if socket == nil {
		fmt.Println(src.String, "not a local ip")
		return
	}

	for i := 1; i <= ttl; i++ {
		opts := &ttlOpts{}
		opts.sock = socket
		opts.dest = dest.Array
		opts.ttl = i
		io.ttlChan <- opts
		time.Sleep(5 * time.Millisecond)
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
}

func (io *NetIO) sendTTLUDP(opts *ttlOpts) {
	e := syscall.SetsockoptInt(opts.sock.fd, 0, syscall.IP_TTL, opts.ttl)
	if e != nil {
		fmt.Println(e)
		return
	}

	b := buildUDP(opts.ttl)

	e = syscall.Sendto(opts.sock.fd, b, 0, &syscall.SockaddrInet4{Port: remoteUDPPort, Addr: opts.dest})
	if e != nil {
		fmt.Println(e)
		return
	}
}

func (io *NetIO) sendRoutine() {
	for {
		select {
		case icmpOpts := <-io.icmpChan:
			io.sendIcmp(icmpOpts)
		case ttlOpts := <-io.ttlChan:
			io.sendTTLUDP(ttlOpts)
		}
	}
}
