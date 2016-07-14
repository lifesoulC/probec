package netio

import (
	"sync"
	"time"
)

type pktKey struct {
	srcIP   uint32
	dstIP   uint32
	icmpID  uint16
	icmpSeq uint16
	srcPort uint16
	dstPort uint16
}

type pktValue struct {
	ttl       int
	sendStamp time.Time
}

type pktsMapOptsICMP struct {
	src uint32
	dst uint32
	id  uint16
	seq uint16
}

type pktsMapOptsTTL struct {
	src     uint32
	dst     uint32
	ttl     int
	srcPort uint16
	dstPort uint16
}

type pktMap struct {
	lock sync.Mutex                 //锁
	pkts map[pktKey]pktValue        //map键值对    值存放ttl和 sendstamp time.Time
}

func newPktMap() *pktMap {
	ret := &pktMap{}
	ret.pkts = make(map[pktKey]pktValue)      //初始化map
	return ret
}

func (pkts *pktMap) clear() {
	t := time.Now()
	pkts.lock.Lock()
	defer pkts.lock.Unlock()
	for k, v := range pkts.pkts {
		if t.Sub(v.sendStamp) > time.Second*5 {
			delete(pkts.pkts, k)
		}
	}
}

func genIcmpKey(opts *pktsMapOptsICMP) (k pktKey) {
	k.dstIP = opts.dst
	k.icmpID = opts.id
	k.icmpSeq = opts.seq
	k.srcIP = opts.src
	return
}

func genTTLKey(opts *pktsMapOptsTTL) (k pktKey) {
	k.dstIP = opts.dst
	k.dstPort = opts.dstPort
	k.srcIP = opts.src
	k.srcPort = opts.srcPort
	return
}

func (pkts *pktMap) addIcmpRequest(opts *pktsMapOptsICMP) {
	k := genIcmpKey(opts)

	v := pktValue{}
	v.sendStamp = time.Now()

	pkts.lock.Lock()
	defer pkts.lock.Unlock()
	pkts.pkts[k] = v
}

func (pkts *pktMap) addTTLRequst(opts *pktsMapOptsTTL) {
	k := genTTLKey(opts)

	v := pktValue{}
	v.sendStamp = time.Now()
	v.ttl = opts.ttl

	pkts.lock.Lock()
	defer pkts.lock.Unlock()
	pkts.pkts[k] = v
}

func (pkts *pktMap) getIcmpDelay(opts *pktsMapOptsICMP) (b bool, delay int) {
	k := genIcmpKey(opts)
	t := time.Now()
	pkts.lock.Lock()
	defer pkts.lock.Unlock()
	v, ok := pkts.pkts[k]
	if !ok {
		return false, 0
	}
	delete(pkts.pkts, k)
	delay = int(t.Sub(v.sendStamp).Nanoseconds() / 1000)
	if delay < 0 || delay > 500000 {
		return false, 0
	} else {
		return true, delay
	}
}

func (pkts *pktMap) getTTLDelay(opts *pktsMapOptsTTL) (b bool, ttl int, delay int) {
	k := genTTLKey(opts)
	t := time.Now()
	pkts.lock.Lock()
	pkts.lock.Unlock()
	v, ok := pkts.pkts[k]
	if !ok {
		return false, 0, 0
	}
	delete(pkts.pkts, k)
	delay = int(t.Sub(v.sendStamp).Nanoseconds() / 1000)
	if delay < 0 || delay > 500000 {
		return false, 0, 0
	} else {
		return true, v.ttl, delay
	}
}
