package go_smq

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
)

// broker 生产者 消费者
// 简单模拟消息队列

type Broker struct {
	event      chan Message
	handlers   []func(message Message)
	handlersMu sync.RWMutex
	Type       reflect.Type
	wait       *sync.WaitGroup
	onceStart  sync.Once
	onceStop   sync.Once
}

// NewStartedBroker 创建broker,并开始
func NewStartedBroker(msgType reflect.Type, chanBuf int) *Broker {
	if !msgType.Implements(reflect.TypeOf((*Message)(nil)).Elem()) {
		panic("broker type is not a message")
	}
	b := &Broker{
		event:    make(chan Message, chanBuf),
		handlers: make([]func(Message), 0),
		Type:     msgType,
		wait:     &sync.WaitGroup{},
	}
	b.Start()
	return b
}

// Send 注册事件
func (b *Broker) Send(ctx context.Context, msg Message) (err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = fmt.Errorf("消息处理异常 name:%v msg:%v err:%v", b.Type, msg, errs)
			log.Println(err)
		}
	}()
	b.event <- msg
	return
}

// Register 注册事件
func (b *Broker) Register(ctx context.Context, f interface{}) (err error) {
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

	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()
	b.handlers = append(b.handlers, func(o Message) {
		defer func() {
			if err := recover(); err != nil {
				err = fmt.Errorf("panic on broker handler msg name:%v err:%v msg:%v", b.Type, err, o)
				log.Println(err)
			}
		}()
		f2Value := reflect.ValueOf(f2.In(0))
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
