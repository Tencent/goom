package mocker_test

import (
	"fmt"
	"reflect"
	"testing"
)

// TestMakeFuncData 测试构建函数的调用方式
func TestMakeFuncData(t *testing.T) {
	fun := genFunc()

	funcIml := toFuncIml(fun)
	funcIml(nil)
}

//go:noinline
func toFuncIml(fun interface{}) func(data *Impl2) int {
	return fun.(func(data *Impl2) int)
}

//go:noinline
func genFunc() interface{} {
	methodTyp := reflect.TypeOf(func(data *Impl2) int {
		fmt.Println("proxy")
		return 3
	})
	return reflect.MakeFunc(methodTyp, func(args []reflect.Value) (results []reflect.Value) {
		fmt.Println("called")
		return []reflect.Value{reflect.ValueOf(1)}
	}).Interface()
}
