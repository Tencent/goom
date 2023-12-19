package mocker

import (
	"fmt"
	"reflect"

	"github.com/tencent/goom/internal/logger"
)

// VarMock 变量 mock
// 支持全局变量, 任意类型包括不限于基本类型，结构体，函数变量，指针与非指针类型
// 主要提供方便变量 mock/reset 场景的能力支持
type VarMock interface {
	Mocker
	// Set 设置变量值, val 类型必须和变量指针指向的值的类型一致
	Set(val interface{})
}

// defaultVarMocker 默认变量 mock 实现
type defaultVarMocker struct {
	target      interface{}
	mockValue   interface{}
	originValue interface{}
	canceled    bool // canceled 是否被取消
}

// String mock 的名称或描述, 方便调试和问题排查
func (m *defaultVarMocker) String() string {
	return fmt.Sprintf("var at[%d]", reflect.ValueOf(m.target).Pointer())
}

// NewVarMocker 创建 VarMock
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
// 注意: Apply 会覆盖之前设定 Set 的值
func (m *defaultVarMocker) Apply(valueCallback interface{}) {
	f := reflect.ValueOf(valueCallback)
	if f.Kind() != reflect.Func {
		panic("VarMock Apply argument(valueCallback) must be a func.")
	}
	ret := f.Call([]reflect.Value{})
	if ret == nil || len(ret) != 1 {
		panic("VarMock Apply valueCallback's returns length must be 1.")
	}

	m.doSet(ret[0].Interface())
	logger.Consolefc(logger.DebugLevel, "mocker [%s] apply.", logger.Caller(5), m.String())
}

// Cancel 取消 mock
func (m *defaultVarMocker) Cancel() {
	t := reflect.ValueOf(m.target)
	t.Elem().Set(reflect.ValueOf(m.originValue))
	m.canceled = true
}

// Canceled 是否取消了 mock
func (m *defaultVarMocker) Canceled() bool {
	return m.canceled
}

// Set 设置变量值
// 注意: Set 会覆盖之前设定 Apply 的值
func (m *defaultVarMocker) Set(val interface{}) {
	m.doSet(val)
	logger.Consolefc(logger.DebugLevel, "mocker [%s] apply.", logger.Caller(5), m.String())
}

func (m *defaultVarMocker) doSet(val interface{}) {
	t := reflect.ValueOf(m.target)
	m.originValue = t.Elem().Interface()
	d := reflect.ValueOf(val)
	t.Elem().Set(d)
	m.mockValue = val
}
