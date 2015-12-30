package pkt

import (
	"fmt"
	"sync"
)

type traceID int64

func (id traceID) String() string {
	return fmt.Sprintf("%ld", int64(id))
}

type pktTracer struct {
	pkts map[traceID]bool
	lock sync.Mutex
}

func newPktTracer() *pktTracer {
	t := &pktTracer{}

}
