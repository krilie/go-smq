package go_smq

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewStartedBroker(t *testing.T) {
	brokerTest := NewStartedBroker("test", 3) // 0 相当于同步调用
	brokerTest.Register(context.Background(), func(i interface{}) {
		fmt.Println(i)
		j := 0
		m := 6 / j
		fmt.Print(m)
	})
	brokerTest.Register(context.Background(), func(i interface{}) {
		fmt.Println("seconed", i)
	})
	brokerTest.Register(context.Background(), func(i interface{}) {
		fmt.Println("seconedseconedseconedseconedseconed", i)
		time.Sleep(time.Second * 2)
	})
	brokerTest.Send(context.Background(), "assssss")
	brokerTest.Stop()
	err := brokerTest.Send(context.Background(), "123")
	if err != nil {
		t.Log(err)
	}
}

func TestNewSmq(t *testing.T) {
	mu := sync.RWMutex{}
	mu.RLock()
	mu.RUnlock()
}
