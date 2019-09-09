package go_smq

import "sync"

var GSmq *Smq

func init() {
	GSmq = NewSmq()
}

type Smq struct {
	Brokers map[string]*Broker
	Mu      sync.RWMutex
}

func NewSmq() *Smq {
	return &Smq{Brokers: make(map[string]*Broker)}
}

func (a *Smq) Get(name string) *Broker {
	a.Mu.RLock()
	if _, ok := a.Brokers[name]; ok {
		defer a.Mu.RUnlock()
		return a.Brokers[name]
	} else {
		a.Mu.RUnlock()
		a.Mu.Lock()
		defer a.Mu.Unlock()
		a.Brokers[name] = NewStartedBroker(name, 1)
		return a.Brokers[name]
	}
}

func (a *Smq) Close() {
	a.Mu.RLock()
	a.Mu.RUnlock()
	for _, v := range a.Brokers {
		v.Stop()
	}
}
