//go:build test

package appmock

import (
	"fmt"
	"reflect"

	"github.com/stretchr/testify/mock"
)

func ReturnArgs(value ...any) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		if len(args) != len(value) {
			panic(fmt.Sprintf("expected %d arguments, got %d", len(value), len(args)))
		}

		for i, v := range value {
			arg := args.Get(i)
			if arg == nil {
				panic(fmt.Sprintf("argument %d is nil", i))
			}

			argVal := reflect.ValueOf(arg)
			if argVal.Kind() != reflect.Ptr {
				panic(fmt.Sprintf("argument %d is not a pointer", i))
			}

			argElem := argVal.Elem()
			val := reflect.ValueOf(v)
			if !val.Type().AssignableTo(argElem.Type()) {
				panic(fmt.Sprintf("type mismatch at argument %d: cannot assign %v to %v", i, val.Kind(), argElem.Kind()))
			}

			argElem.Set(val)
		}
	}
}
