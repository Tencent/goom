// Package hack 对 go 系统包的 hack, 包含一些系统结构体的 copy，需要和不同的 go 版本保持同步
package hack

import "unsafe"

// Func convenience struct for modifying the underlying code pointer of a function
// value. The actual struct has other values, but always starts with a code
// pointer.
// keep async with runtime.Func
type Func struct {
	CodePtr uintptr
}

// Value reflect.Value
// keep async with runtime.Value
type Value struct {
	Typ  *uintptr
	Ptr  unsafe.Pointer
	Flag uintptr
}
