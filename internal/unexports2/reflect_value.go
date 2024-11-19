package unexports2

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

func getRVFlagPtr(v *reflect.Value) *uintptr {
	return (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + rvFlagOffset))
}

var (
	rvFlagAddr uintptr
	rvFlagRO   uintptr

	rvFlagOffset uintptr
	rvFlagsFound bool
	rvFlagsError = fmt.Errorf("This function is disabled because the internal " +
		"flags structure has changed with this go release. Please open " +
		"an issue at https://github.com/kstenerud/go-subvert/issues/new")
)

type rvFlagTester struct {
	A   int // reflect/value.go: flagAddr
	a   int // reflect/value.go: flagStickyRO
	int     // reflect/value.go: flagEmbedRO
	// Note: flagRO = flagStickyRO | flagEmbedRO as of go 1.5
}

func initReflectValue() {
	initReflectValueFlags()
}

func initReflectValueFlags() {
	fail := func(reason string) {
		rvFlagsFound = false
		log.Println(fmt.Sprintf("reflect.Value flags could not be determined because %v."+
			"Please open an issue at https://github.com/kstenerud/go-subvert/issues", reason))
	}
	getFlag := func(v reflect.Value) uintptr {
		return uintptr(reflect.ValueOf(v).FieldByName("flag").Uint())
	}
	getFldFlag := func(v reflect.Value, fieldName string) uintptr {
		return getFlag(v.FieldByName(fieldName))
	}

	if field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag"); ok {
		rvFlagOffset = field.Offset
	} else {
		fail("reflect.Value no longer has a flag field")
		return
	}

	v := rvFlagTester{}
	rv := reflect.ValueOf(&v).Elem()
	rvFlagRO = (getFldFlag(rv, "a") | getFldFlag(rv, "int")) ^ getFldFlag(rv, "A")
	if rvFlagRO == 0 {
		fail("reflect.Value.flag no longer has flagEmbedRO or flagStickyRO bit")
		return
	}

	rvFlagAddr = getFlag(reflect.ValueOf(int(1))) ^ getFldFlag(rv, "A")
	if rvFlagAddr == 0 {
		fail("reflect.Value.flag no longer has a flagAddr bit")
		return
	}
	rvFlagsFound = true
}
