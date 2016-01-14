package prober

import (
	"probec/internal/addr"
	"sync"
	"time"
)

type icmpResultsType struct {
	results map[uint64][]int
	cond    sync.Cond
}

func newIcmpResults() *icmpResultsType {
	r := &icmpResultsType{}
	r.results = make(map[uint64][]int)
	return r
}

func (results *icmpResultsType) waitResult(src string, dest string, t int) []int {
	id := buildIdFromAddr(src, dest)
	results.cond.L.Lock()
	for {
		_, ok := results.results[id]
		if ok {
			results.cond.Wait()
			continue
		} else {
			break
		}
	}

	r := make([]int, 0, 64)
	results.results[id] = r
	results.cond.L.Unlock()

	time.Sleep(time.Duration(t) * time.Millisecond)

	results.cond.L.Lock()
	r = results.results[id]
	delete(results.results, id)
	results.cond.L.Unlock()
	results.cond.Signal()
	return r
}

func buildIdFromAddr(src string, dest string) uint64 {
	return addr.AddPair(src, dest)
}
