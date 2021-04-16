// Package stub 负责生成和应用桩函数
package stub

// unsafe包使用方式
import _ "unsafe"

// MakeFuncStub keep sync with reflect.makeFuncStub
//go:linkname MakeFuncStub reflect.makeFuncStub
func MakeFuncStub()
