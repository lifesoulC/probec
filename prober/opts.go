package prober

import (
	"probec/internal/addr"
)

type PingOpts struct {
	Src      string
	Dest     string
	src      *addr.IPAddr
	dest     *addr.IPAddr
	Count    int
	Interval int
	Delays   []int
}

type BroadcastDelays struct {
	LessAddr    *addr.IPAddr
	LessDelays  []int
	EqualAddr   *addr.IPAddr
	EqualDelays []int
	GreatAddr   *addr.IPAddr
	GreadDelays []int
}

type IcmpBroadcastOpts struct {
	Src      string
	Dest     string
	src      *addr.IPAddr
	dest     *addr.IPAddr
	Count    int
	Interval int
}

type TraceOpts struct {
	Src      string
	Dest     string
	src      *addr.IPAddr
	dest     *addr.IPAddr
	Count    int
	Interval int
}
