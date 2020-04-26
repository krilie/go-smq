package go_smq

import (
	"reflect"
	"sync"
)

var GSmq *Smq

func init() {
	GSmq = NewSmq()
}

type Smq struct {
	Brokers map[reflect.Type]*Broker
	Mu      sync.RWMutex
}

func NewSmq() *Smq {
	return &Smq{Brokers: make(map[reflect.Type]*Broker)}
}

func (a *Smq) Get(msgType reflect.Type) *Broker {
	a.Mu.RLock()
	if _, ok := a.Brokers[msgType]; ok {
		defer a.Mu.RUnlock()
		return a.Brokers[msgType]
	} else {
		a.Mu.RUnlock()
		a.Mu.Lock()
		defer a.Mu.Unlock()
		a.Brokers[msgType] = NewStartedBroker(msgType, 1)
		return a.Brokers[msgType]
	}
}

func (a *Smq) Close() {
	a.Mu.RLock()
	a.Mu.RUnlock()
	for _, v := range a.Brokers {
		v.Stop()
	}
}
