package mocker

import "reflect"

// VarMock 变量mock
// 支持全局变量, 任意类型包括不限于基本类型，结构体，函数变量，指针与非指针类型
type VarMock interface {
	Mocker
	// Return 执行返回值
	Return(ret ...interface{}) *When
}

// defaultVarMocker 默认变量mock实现
type defaultVarMocker struct {
	target      interface{}
	mockValue   interface{}
	originValue reflect.Value
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
// 注意: Apply会覆盖之前设定的Return
func (m *defaultVarMocker) Apply(valueCallback interface{}) {
	f := reflect.ValueOf(valueCallback)
	if f.Kind() != reflect.Func {
		panic("VarMock Apply argument(valueCallback) must be a func.")
	}
	ret := f.Call([]reflect.Value{})
	if ret == nil || len(ret) != 1 {
		panic("VarMock Apply valueCallback's returns length must be 1.")
	}

	m.Return(ret[0].Interface())
}

// Cancel 取消mock
func (m *defaultVarMocker) Cancel() {
	t := reflect.ValueOf(m.target)
	t.Elem().Set(m.originValue)
	m.canceled = true
}

// Canceled 是否取消了mock
func (m *defaultVarMocker) Canceled() bool {
	return m.canceled
}

// Return 设置变量值
// 注意: Return会覆盖之前设定的Apply; 参数ret个数必须为1, 如果大于1则去第0个
func (m *defaultVarMocker) Return(ret ...interface{}) *When {
	if len(ret) == 0 {
		panic("VarMock return value must not be empty.")
	}

	t := reflect.ValueOf(m.target)
	m.originValue = reflect.ValueOf(t.Elem().Interface())
	d := reflect.ValueOf(ret[0])
	t.Elem().Set(d)
	m.mockValue = ret[0]
	return nil
}
