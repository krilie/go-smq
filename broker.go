package go_smq

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// broker 生产者 消费者
// 简单模拟消息队列

type Broker struct {
	event      chan interface{}
	handlers   []func(interface{})
	handlersMu sync.RWMutex
	Name       string
	wait       *sync.WaitGroup
	onceStart  sync.Once
	onceStop   sync.Once
}

// NewStartedBroker 创建broker,并开始
func NewStartedBroker(name string, chanBuf int) *Broker {
	b := &Broker{
		event:    make(chan interface{}, chanBuf),
		handlers: make([]func(interface{}), 0),
		Name:     name,
		wait:     &sync.WaitGroup{},
	}
	b.Start()
	return b
}

// Send 注册事件
func (b *Broker) Send(ctx context.Context, o interface{}) (err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = fmt.Errorf("消息处理异常 name:%v msg:%v err:%v", b.Name, o, errs)
			log.Println(err)
		}
	}()
	b.event <- o
	return
}

// Register 注册事件
func (b *Broker) Register(ctx context.Context, f func(interface{})) (err error) {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()
	b.handlers = append(b.handlers, func(o interface{}) {
		defer func() {
			if err := recover(); err != nil {
				err = fmt.Errorf("panic on broker handler msg name:%v err:%v msg:%v", b.Name, err, o)
				log.Println(err)
			}
		}()
		f(o)
	})
	return nil
}

func (b *Broker) Clear() {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()
	b.handlers = b.handlers[0:0]
}

// Stop 调用stop之前确保写入方都已经退出了,不然要panic
func (b *Broker) Stop() {
	b.onceStop.Do(func() {
		close(b.event)
		b.wait.Wait()
	})
}

func (b *Broker) Start() {
	b.onceStart.Do(func() {
		b.wait.Add(1)
		go func() {
			for {
				event, ok := <-b.event
				if ok {
					// 事件分发
					b.handlersMu.RLock()
					for _, v := range b.handlers {
						v(event) // 有recover
					}
					b.handlersMu.RUnlock()
				} else {
					// 通道已经关闭
					b.wait.Done()
					return
				}
			}
		}()
	})
}
