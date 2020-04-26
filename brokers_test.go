package go_smq

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type TestMsg struct{ C int }

func (t *TestMsg) GetName() string {
	return "TestMsg" + strconv.Itoa(t.C)
}

func (t *TestMsg) GetCtx() context.Context {
	return context.Background()
}

func (t *TestMsg) GetType() reflect.Type {
	return reflect.TypeOf(t)
}

func TestSmq_Close(t *testing.T) {
	_ = GSmq.Register(context.Background(), func(msg *TestMsg) {
		fmt.Println(msg.GetCtx(), msg.GetName(), msg.GetType())
	})
	_ = GSmq.Send(context.Background(), &TestMsg{})
	GSmq.Close()
}

func BenchmarkSmq_Register(b *testing.B) {
	var smq = NewSmq()
	defer smq.Close()
	_ = smq.Register(context.Background(), func(msg *TestMsg) {
		fmt.Println(msg.GetCtx(), msg.GetName(), msg.GetType())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = smq.Send(context.Background(), &TestMsg{C: i})
	}
}

func TestSmq_Close2(t *testing.T) {
	msg := &TestMsg{}
	msg2 := (Message)(msg)
	fmt.Println(reflect.TypeOf(msg))
	fmt.Println(reflect.TypeOf(msg2))
}
