package netio

import (
	"container/list"
	"sync"
	"time"
)

type record struct {
	typ    int
	icmpId int
	ttl    int
	src    string
	dest   string
	t      time.Time
}

type recMap struct {
	records map[int64][]*list.List
	lock    sync.Mutex
}

func newRecMap() *recMap {
	rec := &recMap{}
	rec.records = make(map[int64][]*list.List)
	return rec
}
