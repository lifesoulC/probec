package addr

import (
	"encoding/binary"
	"fmt"
	"net"
	// "strconv"
	// "strings"
)

type IPAddr struct {
	Domain string
	Array  [4]byte
	Slice  []byte
	String string
	addr   uint32
}

func FromString(addr string) (ipAddr *IPAddr, e error) {
	a, e := net.ResolveIPAddr("ip4", addr)
	if e != nil {
		return nil, e
	}
	ipAddr = &IPAddr{}
	ipAddr.Slice = make([]byte, 4)

	ipAddr.Domain = addr
	ipAddr.Slice = a.IP.To4()
	ipAddr.Array[0], ipAddr.Array[1], ipAddr.Array[2], ipAddr.Array[3] = ipAddr.Slice[0], ipAddr.Slice[1], ipAddr.Slice[2], ipAddr.Slice[3]
	ipAddr.addr = binary.BigEndian.Uint32(ipAddr.Slice)
	ipAddr.String = fmt.Sprintf("%d.%d.%d.%d", ipAddr.Array[0], ipAddr.Array[1], ipAddr.Array[2], ipAddr.Array[3])
	return
}

func FromSlice(addr []byte) (ipAddr *IPAddr) {
	ipAddr = &IPAddr{}
	ipAddr.Slice = make([]byte, 4)
	ipAddr.Slice[0], ipAddr.Slice[1], ipAddr.Slice[2], ipAddr.Slice[3] = addr[0], addr[1], addr[2], addr[3]
	ipAddr.Array[0], ipAddr.Array[1], ipAddr.Array[2], ipAddr.Array[3] = ipAddr.Slice[0], ipAddr.Slice[1], ipAddr.Slice[2], ipAddr.Slice[3]
	ipAddr.String = fmt.Sprintf("%d.%d.%d.%d", ipAddr.Array[0], ipAddr.Array[1], ipAddr.Array[2], ipAddr.Array[3])
	ipAddr.Domain = ipAddr.String
	ipAddr.addr = binary.BigEndian.Uint32(ipAddr.Slice)
	return ipAddr
}

func AddrPair(addr1 *IPAddr, addr2 *IPAddr) uint64 {
	return (uint64(addr1.addr) << 32) | uint64(addr2.addr)
}

// func (addr IPAddr) String() string {
// 	ip := addr.Array()
// 	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
// }

// func (addr IPAddr) Array() [4]byte {
// 	ip0 := (uint32(addr) & 0xff000000) >> 24
// 	ip1 := (uint32(addr) & 0xff0000) >> 16
// 	ip2 := (uint32(addr) & 0xff00) >> 8
// 	ip3 := (uint32(addr) & 0xff)
// 	r := [4]byte{byte(ip0), byte(ip1), byte(ip2), byte(ip3)}
// 	return r
// }

// func (addr IPAddr) Slice() []byte {
// 	a := addr.Array()
// 	r := []byte{a[0], a[1], a[2], a[3]}
// 	return r
// }

// func (addr IPAddr) Int() int {
// 	return int(addr)
// }

// func (addr IPAddr) Uint32() uint32 {
// 	return uint32(addr)
// }

// func FromString(s string) IPAddr {
// 	ips := strings.Split(s, ".")
// 	if len(ips) != 4 {
// 		return IPAddr(0)
// 	}
// 	var ip0, ip1, ip2, ip3 int

// 	ip0, _ = strconv.Atoi(ips[0])
// 	ip1, _ = strconv.Atoi(ips[1])
// 	ip2, _ = strconv.Atoi(ips[2])
// 	ip3, _ = strconv.Atoi(ips[3])

// 	ip := (ip0 << 24)
// 	ip = ip | (ip1 << 16)
// 	ip = ip | (ip2 << 8)
// 	ip = ip | ip3
// 	return IPAddr(ip)
// }

// func AddrPair(s1 string, s2 string) uint64 {
// 	ip1 := FromString(s1)
// 	ip2 := FromString(s2)
// 	return (uint64(ip1.Uint32()) << 32) | uint64(ip2.Uint32())
// }

// func IPPair(ip1 IPAddr, ip2 IPAddr) uint64 {
// 	return (uint64(ip1.Uint32()) << 32) | uint64(ip2.Uint32())
// }
