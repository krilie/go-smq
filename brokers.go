package go_smq

import "sync"

type Smq struct {
	Brokers map[string]*Broker
	Mu      sync.RWMutex
}

func NewSmq() *Smq {
	return &Smq{Brokers: make(map[string]*Broker)}
}

func (a *Smq) Get(name string) *Broker {
	a.Mu.RLock()
	defer a.Mu.RUnlock()
	if _, ok := a.Brokers[name]; ok {
		return a.Brokers[name]
	} else {
		a.Mu.Lock()
		defer a.Mu.Unlock()
		a.Brokers[name] = NewStartedBroker(name, 1)
		return a.Brokers[name]
	}
}
