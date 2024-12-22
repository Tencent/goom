package mocker

import (
	"fmt"
	"reflect"
	"unsafe"

	"git.woa.com/goom/mocker/internal/unexports2"
)

// UnExportVarMock 未导出变量 mock
type UnExportVarMock interface {
	Mocker
	VarMock
	// Set 设置变量值, val 类型必须和变量指针指向的值的类型一致
	Set(val interface{})
}

// unExportVarMocker 未导出变量 mock 实现
type unExportVarMocker struct {
	*defaultVarMocker
	target      unsafe.Pointer
	typ         reflect.Type
	mockValue   interface{}
	originValue interface{}
	canceled    bool // canceled 是否被取消
}

// NewUnExportVarMocker 创建 UnExportVarMock
// path 变量路径, package + name 组成, 比如 "github.com/xxx/yyy.varName"
// typ 变量类型, 必须和变量指针指向的值的类型一致，否则会出现不可预测的异常行为
//  1. 可以是指针类型，
//  2. 也可以是非指针类型 relect.TypeOf(struct A{})
//  3. 基本类型 reflect.Int, reflect.Slice, reflect.Map 等
func NewUnExportVarMocker(path string, typ reflect.Type) UnExportVarMock {
	addr, err := unexports2.FindVarByName(path)
	if err != nil {
		panic(fmt.Sprintf("cannot find unexported var: %s, cause by %v", path, err))
	}
	return &unExportVarMocker{
		defaultVarMocker: newVarMocker(reflect.NewAt(typ, unsafe.Pointer(addr))),
		target:           unsafe.Pointer(addr),
		typ:              typ,
		canceled:         false,
	}
}

// String mock 的名称或描述, 方便调试和问题排查
func (m *unExportVarMocker) String() string {
	return fmt.Sprintf("var at[%d]", m.target)
}
