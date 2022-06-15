// Package test 兼容性测试、跨包结构测试工具类
package test

import "fmt"

// GlobalVar 用于测试全局变量 mock
var GlobalVar = 1

// Foo foo 测试函数
//go:noinline
func Foo(i int) int {
	// check 对 defer 的支持
	defer func() { fmt.Printf("defer\n") }()
	return i * 1
}

//go:noinline
// foo foo 测试未导出函数
func foo(i int) int {
	// check 对 defer 的支持
	defer func() { fmt.Printf("defer\n") }()
	return i * 1
}

// Invokefoo foo 测试调用未导出函数
//go:noinline
func Invokefoo(i int) int {
	return foo(i)
}

// fake 未导出结构体
type fake struct {
	field1 string
	field2 int
}

// NewUnexportedFake 构建未导出fake
//nolint
func NewUnexportedFake() *fake {
	return &fake{
		field1: "field1",
		field2: 2,
	}
}

// Call 普通方法
//go:noinline
func (f *fake) Call(i int) int {
	if i < -10000 {
		dummy()
	}
	return i
}

// Call2 普通方法
//go:noinline
func (f *fake) Call2(i int) int {
	if i < -10000 {
		dummy()
	}
	return i
}

// call 未导出方法
//go:noinline
func (f *fake) call(i int) int {
	if i < -10000 {
		dummy()
	}
	return i
}

// Invokecall 测试调用未导出函数
//go:noinline
func (f *fake) Invokecall(i int) int {
	return f.call(i)
}

func dummy() {
	fmt.Println("never print, just from fill function length")
}

// S 测试返回复杂类型
type S struct {
	Field1 string
	Field2 int
}

// S1 测试返回同等结构(不同Type)的值
type S1 struct {
	Field1 string
	Field2 int
}

// Foo1 foo1 测试函数返回复杂类型
//go:noinline
func Foo1() *S {
	return &S{
		Field1: "ok",
		Field2: 2,
	}
}

// GetS 测试返回多参返回值
func GetS() ([]byte, error) {
	return []byte("hello"), nil
}
