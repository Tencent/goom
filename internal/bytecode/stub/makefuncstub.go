// Package stub 负责生成和应用桩函数
package stub

// unsafe 包使用方式
import _ "unsafe" // 匿名引入

// MakeFuncStub 用于调用 makeFunc 创建出来的函数
// keep sync with reflect.makeFuncStub
//
//go:linkname MakeFuncStub reflect.makeFuncStub
func MakeFuncStub()
