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

type DestDelays struct {
	Dest   *addr.IPAddr
	Delays []int
}

type broadcastMap map[uint64][]*icmpBroadResultType

type icmpBroadResultsType struct {
	results broadcastMap
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
	r, ok := results.results[id]
	if !ok {
		return
	}
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
	r := make([]*icmpBroadResultType, 0)
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
	time.Sleep(time.Duration(t) * time.Second)
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
	r, ok := results.results[id]
	if !ok {
		return
	}
	newResult := true
	for _, v := range r {
		// fmt.Println(v)
		if v != nil && v.dest.Equal(dest) {
			v.delays = append(v.delays, delay)
			newResult = false
		}
	}

	if newResult {
		n := &icmpBroadResultType{}
		n.dest = dest
		n.delays = append(n.delays, delay)
		r = append(r, n)
	}

	results.results[id] = r
}

type TraceResultType struct {
	TTL    int
	Host   *addr.IPAddr
	Delays []int
}

type traceResultsType struct {
	results map[uint64][]*TraceResultType
	lock    sync.Mutex
	cond    *sync.Cond
}

func newTraceResults() *traceResultsType {
	results := &traceResultsType{}
	results.cond = sync.NewCond(&results.lock)
	results.results = make(map[uint64][]*TraceResultType)
	return results
}

func (r *traceResultsType) addResult(src *addr.IPAddr, dest *addr.IPAddr, host *addr.IPAddr, ttl int, delay int) {
	id := addr.AddrPair(src, dest)
	r.lock.Lock()
	defer r.lock.Unlock()
	result, ok := r.results[id]
	if !ok {
		return
	}

	for k, v := range result {
		if v.TTL == ttl && v.Host.Equal(host) {
			result[k].Delays = append(result[k].Delays, delay)
			return
		}
	}

	t := &TraceResultType{}
	t.Host = host
	t.Delays = make([]int, 1)
	t.Delays[0] = delay
	t.TTL = ttl
	result = append(result, t)
	r.results[id] = result
}

func (r *traceResultsType) beginWait(src *addr.IPAddr, dest *addr.IPAddr) {
	id := addr.AddrPair(src, dest)
	r.lock.Lock()
	defer r.lock.Unlock()
	for {
		_, ok := r.results[id]
		if !ok {
			d := make([]*TraceResultType, 0)
			r.results[id] = d
			return
		}
		r.cond.Wait()
	}
}

func (r *traceResultsType) endWait(src *addr.IPAddr, dest *addr.IPAddr, t int) (ret []*TraceResultType) {
	id := addr.AddrPair(src, dest)
	time.Sleep(time.Second * time.Duration(t))
	r.lock.Lock()
	defer r.lock.Unlock()
	d, ok := r.results[id]
	if ok {
		ret = d
		delete(r.results, id)
	}
	r.cond.Signal()
	return
}
