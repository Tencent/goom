// Package stub负责生成和应用桩函数
package stub

// nolint unsafe包使用方式
import _ "unsafe"

//go:linkname MakeFuncStub reflect.makeFuncStub
func MakeFuncStub()
