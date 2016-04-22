package addr

import (
	"encoding/binary"
	"fmt"
	"net"
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

func (addr *IPAddr) Int() uint32 {
	return addr.addr
}

func FromInt(ip uint32) (ipAddr *IPAddr) {
	ipAddr = &IPAddr{}
	ipAddr.addr = ip
	ipAddr.Slice = make([]byte, 4)
	binary.BigEndian.PutUint32(ipAddr.Slice, ip)
	ipAddr.Array[0], ipAddr.Array[1], ipAddr.Array[2], ipAddr.Array[3] = ipAddr.Slice[0], ipAddr.Slice[1], ipAddr.Slice[2], ipAddr.Slice[3]
	ipAddr.String = fmt.Sprintf("%d.%d.%d.%d", ipAddr.Array[0], ipAddr.Array[1], ipAddr.Array[2], ipAddr.Array[3])
	ipAddr.Domain = ipAddr.String
	return ipAddr
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

func (addr *IPAddr) Section() uint32 {
	return (addr.addr & 0xffffff00)
}

func AddrSectionPair(addr1 *IPAddr, addr2 *IPAddr) uint64 {
	return (uint64(addr1.addr) << 32) | uint64(addr2.Section())
}

func (addr *IPAddr) Equal(rhs *IPAddr) bool {
	return addr.addr == rhs.addr
}

func (addr *IPAddr) Less(rhs *IPAddr) bool {
	return addr.addr < rhs.addr
}

func (addr *IPAddr) Great(rhs *IPAddr) bool {
	return addr.addr > rhs.addr
}
