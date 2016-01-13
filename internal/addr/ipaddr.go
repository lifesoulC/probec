package addr

import (
	"fmt"
	"strconv"
	"strings"
)

type IPAddr uint32

func (addr IPAddr) String() string {
	ip := addr.Array()
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (addr IPAddr) Array() [4]byte {
	ip0 := (uint32(addr) & 0xff000000) >> 24
	ip1 := (uint32(addr) & 0xff0000) >> 16
	ip2 := (uint32(addr) & 0xff00) >> 8
	ip3 := (uint32(addr) & 0xff)
	r := [4]byte{byte(ip0), byte(ip1), byte(ip2), byte(ip3)}
	return r
}

func (addr IPAddr) Slice() []byte {
	a := addr.Array()
	r := []byte{a[0], a[1], a[2], a[3]}
	return r
}

func (addr IPAddr) Int() int {
	return int(addr)
}

func (addr IPAddr) Uint32() uint32 {
	return uint32(addr)
}

func FromString(s string) IPAddr {
	ips := strings.Split(s, ".")
	if len(ips) != 4 {
		return IPAddr(0)
	}
	var ip0, ip1, ip2, ip3 int

	ip0, _ = strconv.Atoi(ips[0])
	ip1, _ = strconv.Atoi(ips[1])
	ip2, _ = strconv.Atoi(ips[2])
	ip3, _ = strconv.Atoi(ips[3])

	ip := (ip0 << 24)
	ip = ip | (ip1 << 16)
	ip = ip | (ip2 << 8)
	ip = ip | ip3
	return IPAddr(ip)
}

func AddPair(s1 string, s2 string) uint64 {
	ip1 := FromString(s1)
	ip2 := FromString(s2)
	return (uint64(ip1.Uint32()) << 32) | uint64(ip2.Uint32())
}
