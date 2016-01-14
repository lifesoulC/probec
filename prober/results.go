package prober

import (
	"probec/internal/addr"
	"sync"
	"time"
)

type icmpResultsType struct {
	results map[uint64][]int
	lock    sync.Mutex
	cond    *sync.Cond
}

func newIcmpResults() *icmpResultsType {
	r := &icmpResultsType{}
	r.results = make(map[uint64][]int)

	r.cond = sync.NewCond(&r.lock)
	return r
}

func (results *icmpResultsType) beginWait(src *addr.IPAddr, dest *addr.IPAddr) {
	id := addr.AddrPair(src, dest)
	r := make([]int, 0, 64)
	results.lock.Lock()

	for {
		_, ok := results.results[id]
		if ok {
			results.cond.Wait()
			continue
		} else {
			break
		}
	}
	results.results[id] = r
	results.lock.Unlock()
}

func (results *icmpResultsType) endWait(src *addr.IPAddr, dest *addr.IPAddr, t int) []int {
	time.Sleep(time.Duration(t) * time.Millisecond)
	id := addr.AddrPair(src, dest)
	results.lock.Lock()
	r := results.results[id]
	delete(results.results, id)
	results.lock.Unlock()
	results.cond.Signal()
	return r
}

func (results *icmpResultsType) addResult(src *addr.IPAddr, dest *addr.IPAddr, delay int) {
	id := addr.AddrPair(src, dest)
	results.lock.Lock()
	defer results.lock.Unlock()
	r := results.results[id]
	r = append(r, delay)
	results.results[id] = r
}

func buildIdFromAddr(src *addr.IPAddr, dest *addr.IPAddr) uint64 {
	return addr.AddrPair(src, dest)
}
