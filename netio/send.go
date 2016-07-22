package netio

import (
	"encoding/binary"
	"fmt"
	"probec/internal/addr"
	"syscall"
	"time"
)

func (io *NetIO) SendPing(src *addr.IPAddr, dest *addr.IPAddr) (err error) {
	socket, e := io.getIcmpSock(src.String) //获得已经绑定好的socket
	if socket == nil {
		fmt.Println(src.String, "not a local ip")
		err = e
		return
	}
	opt := &icmpOpts{}
	opt.dest = dest.Array
	opt.sock = socket
	opt.srcInt = src.Int()
	opt.dstInt = dest.Int()
	opt.broad = false
	io.icmpChan <- opt //送到管道
}

func (io *NetIO) SendPingBroadcast(src *addr.IPAddr, dest *addr.IPAddr) (err error) {
	socket, e := io.getIcmpSock(src.String)
	if socket == nil {
		fmt.Println(src.String, "not a local ip")
		err = e
		return
	}

	for i := 1; i < 255; i++ {
		raddr := dest.Array
		raddr[3] = byte(i)
		opt := &icmpOpts{}
		opt.sock = socket
		tmp := make([]byte, 4)
		tmp = []byte{raddr[0], raddr[1], raddr[2], raddr[3]}
		opt.dstInt = binary.BigEndian.Uint32(tmp)
		opt.srcInt = src.Int()
		opt.dest = raddr
		opt.broad = true
		io.icmpChan <- opt
		time.Sleep(1000 * time.Microsecond)
	}

}

func (io *NetIO) SendTTL(src *addr.IPAddr, dest *addr.IPAddr, ttl int) (err error) {
	socket, e := io.getUdpSock(src.String)
	if socket == nil {
		fmt.Println(src.String, "not a local ip")
		err = e
		return
	}

	for i := 1; i <= ttl; i++ {
		opts := &ttlOpts{}
		opts.sock = socket
		opts.dest = dest.Array
		opts.srcInt = src.Int()
		opts.dstInt = dest.Int()
		opts.ttl = i
		io.ttlChan <- opts
		time.Sleep(5 * time.Millisecond)
	}
}

func (io *NetIO) sendIcmp(opts *icmpOpts) { //发送icmp包
	var s uint16
	if opts.broad {
		opts.data, s = buildIcmpBroadcast()
	} else {
		opts.data, s = buildIcmpEchoRequest()
	}

	e := syscall.Sendto(opts.sock.fd, opts.data, 0, &syscall.SockaddrInet4{Port: 0, Addr: opts.dest}) //从零号端口发出
	if e != nil {
		fmt.Println("send to", opts.dest, e.Error())
		return
	}
	keyOpts := &pktsMapOptsICMP{}
	keyOpts.dst = opts.dstInt
	keyOpts.id = uint16(pid)
	keyOpts.seq = s
	keyOpts.src = opts.srcInt
	io.pkts.addIcmpRequest(keyOpts)
}

func (io *NetIO) sendTTLUDP(opts *ttlOpts) {
	e := syscall.SetsockoptInt(opts.sock.fd, 0, syscall.IP_TTL, opts.ttl)
	if e != nil {
		fmt.Println(e)
		return
	}

	b := buildUDP(opts.ttl)

	e = syscall.Sendto(opts.sock.fd, b, 0, &syscall.SockaddrInet4{Port: int(remoteUDPPort), Addr: opts.dest})

	if e != nil {
		fmt.Println(e)
		return
	}

	keyOpts := &pktsMapOptsTTL{}
	keyOpts.dst = opts.dstInt
	keyOpts.ttl = opts.ttl
	keyOpts.dstPort = remoteUDPPort
	keyOpts.src = opts.srcInt
	keyOpts.srcPort = localUDPPort

	remoteUDPPort = remoteUDPPort + 1
	if remoteUDPPort >= remoteUDPPortMax {
		remoteUDPPort = remoteUDPPortMin
	}
	io.pkts.addTTLRequst(keyOpts)
}

func (io *NetIO) sendRoutine() {
	for {
		select {
		case icmpOpts := <-io.icmpChan: //从icmpchan队列中读出一个放入icmp发送出去
			io.sendIcmp(icmpOpts)
		case ttlOpts := <-io.ttlChan:
			io.sendTTLUDP(ttlOpts)
		}
		t := time.Now()
		if t.Sub(io.lastClear) > 5*time.Second {
			io.pkts.clear()
			io.lastClear = t
		}
	}
}
