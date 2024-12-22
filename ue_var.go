package mocker

import (
	"fmt"
	"git.woa.com/goom/mocker/internal/logger"
	"reflect"
	"unsafe"

	"git.woa.com/goom/mocker/internal/unexports2"
)

// UnExportedVarMock 未导出变量 mock
type UnExportedVarMock interface {
	Mocker
	VarMock
	// Set 设置变量值, val 类型必须和变量指针指向的值的类型一致
	Set(val interface{})
}

// unExportedVarMocker 未导出变量 mock 实现
type unExportedVarMocker struct {
	*defaultVarMocker
	target      unsafe.Pointer
	typ         reflect.Type
	mockValue   interface{}
	originValue interface{}
	canceled    bool // canceled 是否被取消
}

// NewUnExportedVarMocker 创建 UnExportedVarMock
// path 变量路径, package + name 组成, 比如 "github.com/xxx/yyy.varName"
func NewUnExportedVarMocker(path string) UnExportedVarMock {
	addr, err := unexports2.FindVarByName(path)
	if err != nil {
		panic(fmt.Sprintf("cannot find unexported var: %s, cause by %v", path, err))
	}
	return &unExportedVarMocker{
		defaultVarMocker: newVarMocker(reflect.Value{}),
		target:           unsafe.Pointer(addr),
		canceled:         false,
	}
}

// String mock 的名称或描述, 方便调试和问题排查
func (m *unExportedVarMocker) String() string {
	return fmt.Sprintf("var at[%d]", m.target)
}

// Set 设置变量值
// value 变量值, 必须和变量原值的类型一致，否则会出现不可预测的异常行为
//  1. 可以是指针类型，
//  2. 也可以是非指针类型 比如struct A{}
//  3. 基本类型 reflect.Int, reflect.Slice, reflect.Map 等
// 注意: Set 会覆盖之前设定 Apply 的值
func (m *unExportedVarMocker) Set(value interface{}) {
	m.typ = reflect.TypeOf(value)
	m.targetValue = reflect.NewAt(m.typ, m.target)
	// TODO 检查类型是否匹配
	m.defaultVarMocker.doSet(value)
	logger.Consolefc(logger.DebugLevel, "mocker [%s] apply.", logger.Caller(5), m.String())
}
