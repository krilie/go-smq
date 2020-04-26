package go_smq

import (
	"context"
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

func (a *Smq) Register(ctx context.Context, f interface{}) (err error) {
	f2 := reflect.TypeOf(f)
	if f2.Kind() != reflect.Func {
		panic("message func not a function")
	}
	if f2.NumOut() != 0 {
		panic("message func with not void out param")
	}
	if f2.NumIn() != 1 {
		panic("message func with not one in param")
	}
	inParam := f2.In(0)
	if !inParam.Implements(reflect.TypeOf((*Message)(nil)).Elem()) {
		panic("message func with param is not message")
	}
	return a.Get(inParam).Register(ctx, f)
}

func (a *Smq) Send(ctx context.Context, msg Message) (err error) {
	return a.Get(reflect.TypeOf(msg)).Send(ctx, msg)
}
