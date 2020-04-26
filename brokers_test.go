package go_smq

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type TestMsg struct{}

func (t *TestMsg) GetName() string {
	return "TestMsg"
}

func (t *TestMsg) GetCtx() context.Context {
	return context.Background()
}

func (t *TestMsg) GetType() reflect.Type {
	return reflect.TypeOf(t)
}

func TestSmq_Close(t *testing.T) {
	GSmq.Register(context.Background(), func(msg *TestMsg) {
		fmt.Println(msg.GetCtx(), msg.GetName(), msg.GetType())
	})
	GSmq.Send(context.Background(), &TestMsg{})
	GSmq.Close()
}

func TestSmq_Close2(t *testing.T) {
	msg := &TestMsg{}
	msg2 := (Message)(msg)
	fmt.Println(reflect.TypeOf(msg))
	fmt.Println(reflect.TypeOf(msg2))
}
