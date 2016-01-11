package netio

import (
	"net"
	"strconv"
	"strings"
)

func ipSlice(ip [4]byte) []byte {
	return []byte{ip[0], ip[1], ip[2], ip[3]}
}

func ipArrary(ip []byte) [4]byte {
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}

func addrIpArray(addr string) (ip [4]byte, err error) {
	ip = [4]byte{0, 0, 0, 0}

	a, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return ip, err
	}

	return ipArrary(a.IP.To4()), nil
}

func splitAddr(addr string) [4]byte {
	ret := [4]byte{0, 0, 0, 0}
	ips := strings.Split(addr, ".")
	if len(ips) != 4 {
		return ret
	}
	ip0, _ := strconv.Atoi(ips[0])
	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ret[0] = byte(ip0)
	ret[1] = byte(ip1)
	ret[2] = byte(ip2)
	ret[3] = byte(ip3)
	return ret
}
func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}
