package fnstore

import (
	"fmt"
	"reflect"
)

type FunctionData struct {
	Func   reflect.Value
	Input  []reflect.Type
}

type FunctionStore interface {
	AddFunction(name string, fn interface{}) error
	CallFunction(name string, args ...interface{}) ([]interface{}, error)
}

type FunctionStoreImpl[K comparable] struct {
	Functions map[K]FunctionData
}

func NewFunctionStore[K comparable]() *FunctionStoreImpl[K] {
	return &FunctionStoreImpl[K]{Functions: make(map[K]FunctionData)}
}

func (fs *FunctionStoreImpl[K]) AddFunction(name K, fn interface{}) error {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("no function passed")
	}
	if existingFn, exists := fs.Functions[name]; exists && existingFn.Func.Pointer() != fnValue.Pointer() {
		return fmt.Errorf("key %v is already used by a different function", name)
	}
	inputTypes := make([]reflect.Type, fnValue.Type().NumIn())
	for i := 0; i < fnValue.Type().NumIn(); i++ {
		inputTypes[i] = fnValue.Type().In(i)
	}
	fs.Functions[name] = FunctionData{
		Func:   fnValue,
		Input:  inputTypes,
	}
	return nil
}

func (fs *FunctionStoreImpl[K]) CallFunction(name K, args []interface{}) ([]interface{}, error) {
	fnMeta, exists := fs.Functions[name]
	if !exists {
		return nil, fmt.Errorf("function %v not found", name)
	}

	if len(args) != len(fnMeta.Input) {
		return nil, fmt.Errorf("incorrect number of arguments")
	}

	inputs := make([]reflect.Value, len(args))
	for i, arg := range args {
		if reflect.TypeOf(arg) != fnMeta.Input[i] {
			return nil, fmt.Errorf("argument %d does not match function parameter type", i)
		}
		inputs[i] = reflect.ValueOf(arg)
	}

	results := fnMeta.Func.Call(inputs)
	outputs := make([]interface{}, len(results))
	for i, result := range results {
		outputs[i] = result.Interface()
	}
	return outputs, nil
}
