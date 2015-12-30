package prober

import (
	"net"
	"os"
)

const (
	PingType    = 1
	SectionType = 2
	TraceType   = 3
)

var pid int

func init() {
	pid = os.Getpid() & 0xffff
}

type Prober interface {
	Type() int
	Probe(*net.IPAddr, *net.IPAddr) ([]int, error)
}
