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

type icmpBroadResultType struct {
	dest   *addr.IPAddr
	delays []int
}
type icmpBroadResultsType struct {
	results map[uint64][]*icmpBroadResultType
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

func newIcmpBroadResults() *icmpBroadResultsType {
	results := &icmpBroadResultsType{}
	results.results = make(map[uint64][]*icmpBroadResultType)
	results.cond = sync.NewCond(&results.lock)
	return results
}

func (results *icmpBroadResultsType) beginWait(src *addr.IPAddr, dest *addr.IPAddr) {
	id := addr.AddrSectionPair(src, dest)
	r := make([]*icmpBroadResultType, 255)
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

func (results *icmpBroadResultsType) endWait(src *addr.IPAddr, dest *addr.IPAddr, t int) []*icmpBroadResultType {
	time.Sleep(time.Duration(t) * time.Millisecond)
	id := addr.AddrSectionPair(src, dest)
	results.lock.Lock()
	r := results.results[id]
	delete(results.results, id)
	results.lock.Unlock()
	results.cond.Signal()
	return r
}

func (results *icmpBroadResultsType) addResult(src *addr.IPAddr, dest *addr.IPAddr, delay int) {
	id := addr.AddrSectionPair(src, dest)
	results.lock.Lock()
	defer results.lock.Unlock()
	r := results.results[id]
	for _, v := range r {
		if v.dest.Equal(dest) {
			v.delays = append(v.delays, delay)
		}
	}
	results.results[id] = r
}
