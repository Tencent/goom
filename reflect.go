// Package mocker 定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件汇聚了公共的工具类，比如类型转换。
package mocker

import (
	"reflect"
	"runtime"
)

// functionName 获取函数名称
func functionName(fnc interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()
}

// typeName 获取类型名称
func typeName(fnc interface{}) string {
	t := reflect.TypeOf(fnc)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}

// inTypes 获取类型
func inTypes(isMethod bool, funTyp reflect.Type) []reflect.Type {
	skip := 0
	if isMethod {
		skip = 1
	}

	numIn := funTyp.NumIn()
	typeList := make([]reflect.Type, numIn-skip)

	for i := 0; i < numIn-skip; i++ {
		typeList[i] = funTyp.In(i + skip)
	}

	return typeList
}

// outTypes 获取类型
func outTypes(funTyp reflect.Type) []reflect.Type {
	numOut := funTyp.NumOut()
	typeList := make([]reflect.Type, numOut)

	for i := 0; i < numOut; i++ {
		typeList[i] = funTyp.Out(i)
	}

	return typeList
}
