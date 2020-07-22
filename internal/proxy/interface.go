package proxy

import _ "unsafe"

func InterfaceCallStub()

//go:linkname MakeFuncStub reflect.makeFuncStub
func MakeFuncStub()
