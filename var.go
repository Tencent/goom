package mocker

import "reflect"

// VarMock 变量mock
// 支持全局变量, 任意类型包括不限于基本类型，结构体，函数变量，指针与非指针类型
// 主要提供方便变量mock/reset场景的能力支持
type VarMock interface {
	Mocker
	// Set 设置返回值, val类型必须和变量指针指向的值的类型一致
	Set(val interface{})
}

// defaultVarMocker 默认变量mock实现
type defaultVarMocker struct {
	target      interface{}
	mockValue   interface{}
	originValue interface{}
	canceled    bool // canceled 是否被取消
}

// NewVarMocker 创建VarMock
func NewVarMocker(target interface{}) VarMock {
	t := reflect.ValueOf(target)
	if t.Type().Kind() != reflect.Ptr {
		panic("VarMock target must be a pointer.")
	}
	return &defaultVarMocker{
		target:   target,
		canceled: false,
	}
}

// Apply 变量取值回调函数, 只会执行一次
// 注意: Apply会覆盖之前设定的Set
func (m *defaultVarMocker) Apply(valueCallback interface{}) {
	f := reflect.ValueOf(valueCallback)
	if f.Kind() != reflect.Func {
		panic("VarMock Apply argument(valueCallback) must be a func.")
	}
	ret := f.Call([]reflect.Value{})
	if ret == nil || len(ret) != 1 {
		panic("VarMock Apply valueCallback's returns length must be 1.")
	}

	m.Set(ret[0].Interface())
}

// Cancel 取消mock
func (m *defaultVarMocker) Cancel() {
	t := reflect.ValueOf(m.target)
	t.Elem().Set(reflect.ValueOf(m.originValue))
	m.canceled = true
}

// Canceled 是否取消了mock
func (m *defaultVarMocker) Canceled() bool {
	return m.canceled
}

// Set 设置变量值
// 注意: Set会覆盖之前设定的Apply
func (m *defaultVarMocker) Set(val interface{}) {
	t := reflect.ValueOf(m.target)
	m.originValue = t.Elem().Interface()
	d := reflect.ValueOf(val)
	t.Elem().Set(d)
	m.mockValue = val
}
