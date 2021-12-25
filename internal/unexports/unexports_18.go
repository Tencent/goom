// +build go1.18

// Package unexports 实现了对未导出函数的获取
// 基于github.com/alangpierce/go-forceexport进行了修改和扩展。
package unexports

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"

	"git.code.oa.com/goom/mocker/erro"
	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/logger"
)

const ptrMax uintptr = (1<<31 - 1) * 100

// FindFuncByName searches through the moduledata table created by the linker
// and returns the function's code pointer. If the function was not found, it
// returns an error. Since the data structures here are not exported, we copy
// them below (and they need to stay in sync or else things will fail
// catastrophically).
func FindFuncByName(name string) (uintptr, error) {
	if len(name) == 0 {
		return 0, errors.New("FindFuncByName error: func name is empty")
	}

	suggester := newSuggester(name)
	for moduleData := &hack.Firstmoduledata; moduleData != nil; moduleData = moduleData.Next {
		for _, ftab := range moduleData.Ftab {
			if ftab.Funcoff >= uint32(len(moduleData.Pclntable)) {
				break
			}
			f := (*runtime.Func)(unsafe.Pointer(&moduleData.Pclntable[ftab.Funcoff]))
			if f == nil {
				continue
			}

			if f.Entry() > ptrMax {
				continue
			}

			fName := funcName(f)
			if fName == name {
				return f.Entry(), nil
			}

			suggester.AddItem(fName)
		}
	}
	logger.LogDebugf("FindFuncByName not found %s", name)

	return 0, erro.NewFuncNotFoundErrorWithSuggestion(name, suggester.Suggestions())
}

// funcName 获取函数名字
func funcName(f *runtime.Func) string {
	defer func() {
		if err := recover(); err != nil {
			var buf = make([]byte, 1024)

			runtime.Stack(buf, true)
			logger.LogErrorf("get funcName error:[%+v]\n%s", err, buf)
		}
	}()

	return f.Name()
}

// FindFuncByPtr 根据地址找到函数信息
// 返回: *runtime.Func 函数信息结构体
// string 函数名
func FindFuncByPtr(ptr uintptr) (*runtime.Func, string, error) {
	for moduleData := &hack.Firstmoduledata; moduleData != nil; moduleData = moduleData.Next {
		for _, ftab := range moduleData.Ftab {
			if ftab.Funcoff >= uint32(len(moduleData.Pclntable)) {
				break
			}
			f := (*runtime.Func)(unsafe.Pointer(&moduleData.Pclntable[ftab.Funcoff]))
			if f == nil {
				continue
			}

			if f.Entry() > ptrMax {
				continue
			}

			fName := funcName(f)

			if f.Entry() == ptr {
				return f, fName, nil
			}
		}
	}

	return nil, "", fmt.Errorf("invalid function ptr: %d", ptr)
}

// CreateFuncForCodePtr is given a code pointer and creates a function value
// that uses that pointer. The outFun argument should be a pointer to a function
// of the proper type (e.g. the address of a local variable), and will be set to
// the result function value.
func CreateFuncForCodePtr(outFuncPtr interface{}, codePtr uintptr) (*hack.Func, error) {
	outFunc := reflect.ValueOf(outFuncPtr)
	if outFunc.Kind() != reflect.Ptr {
		return nil, errors.New("func param must be ptr")
	}

	outFuncVal := outFunc.Elem()
	// Use reflect.MakeGlobalFunc to create a well-formed function value that's
	// guaranteed to be of the right type and guaranteed to be on the heap
	// (so that we can modify it). We give a nil delegate function because
	// it will never actually be called.
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	// Use reflection on the reflect.Value (yep!) to grab the underling
	// function value pointer. Trying to call newFuncVal.Pointer() wouldn't
	// work because it gives the code pointer rather than the function value
	// pointer. The function value is a struct that starts with its code
	// pointer, so we can swap out the code pointer with our desired value.
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	// nolint hack用法
	funcPtr := (*hack.Func)(unsafe.Pointer(funcValuePtr))
	funcPtr.CodePtr = codePtr

	outFuncVal.Set(newFuncVal)

	return funcPtr, nil
}

// NewFuncWithCodePtr 根据类型和函数地址进行构造reflect.Value
func NewFuncWithCodePtr(typ reflect.Type, codePtr uintptr) reflect.Value {
	var ptr2Ptr = &codePtr
	pointer := unsafe.Pointer(ptr2Ptr)
	funcVal := reflect.NewAt(typ, pointer).Elem()
	(*hack.Value)(unsafe.Pointer(&funcVal)).Flag = uintptr(reflect.Func)

	return funcVal
}
