package stub

// nolint
import _ "unsafe"

//go:linkname MakeFuncStub reflect.makeFuncStub
func MakeFuncStub()
