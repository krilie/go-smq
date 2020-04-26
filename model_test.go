package go_smq

import (
	"reflect"
	"testing"
)

func TestTypeOf(t *testing.T) {
	println(reflect.TypeOf(&MessageName{}).Implements(reflect.TypeOf((*Message)(nil)).Elem()))
}

func TestTypeOf2(t *testing.T) {
	var msg Message = nil
	println(reflect.TypeOf(&msg).Elem().Kind())
}
