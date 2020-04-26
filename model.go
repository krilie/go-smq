package go_smq

import (
	"context"
	"reflect"
)

type Message interface {
	GetName() string
	GetCtx() context.Context
	GetType() reflect.Type
}

type MessageName struct{}

func (m *MessageName) GetName() string {
	return "MessageName"
}

func (m *MessageName) GetCtx() context.Context {
	return context.Background()
}

func (m *MessageName) GetType() reflect.Type {
	return reflect.TypeOf(m)
}
