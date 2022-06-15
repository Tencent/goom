package patch

import (
	"unsafe"
)

const (
	_0b1      = 1  // _0b1
	_0b10     = 2  // 0b10
	_0b11     = 3  // 0b11
	_0b100101 = 37 // 0b100101
)

// nopOpcode 空指令插入到原函数开头第一个字节, 用于判断原函数是否已经被Patch过
var nopOpcode = []byte{0xD5, 0x03, 0x20, 0x1F}

func jmpToFunctionValue(_, double uintptr) []byte {
	//func buildJmpDirective(double uintptr) []byte {
	res := make([]byte, 0, 24)
	d0d1 := double & 0xFFFF
	d2d3 := double >> 16 & 0xFFFF
	d4d5 := double >> 32 & 0xFFFF
	d6d7 := double >> 48 & 0xFFFF

	res = append(res, movImm(_0b10, 0, d0d1)...)         // MOVZ x26, double[16:0]
	res = append(res, movImm(_0b11, 1, d2d3)...)         // MOVK x26, double[32:16]
	res = append(res, movImm(_0b11, 2, d4d5)...)         // MOVK x26, double[48:32]
	res = append(res, movImm(_0b11, 3, d6d7)...)         // MOVK x26, double[64:48]
	res = append(res, []byte{0x4A, 0x03, 0x40, 0xF9}...) // LDR x10, [x26]
	res = append(res, []byte{0x40, 0x01, 0x1F, 0xD6}...) // BR x10

	return res
}

func movImm(opc, shift int, val uintptr) []byte {
	var m uint32 = 26          // rd
	m |= uint32(val) << 5      // imm16
	m |= uint32(shift&3) << 21 // hw
	m |= _0b100101 << 23       // const
	m |= uint32(opc&0x3) << 29 // opc
	m |= _0b1 << 31            // sf

	res := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&res[0])) = m

	return res
}

// jmpToOriginFunctionValue Assembles a jump to a function value
func jmpToOriginFunctionValue(_, _ uintptr) (value []byte) {
	panic("not support yet")
}

// checkAlreadyPatch 检测是否已经patch
func checkAlreadyPatch(origin []byte) bool {
	for i := 0; i < len(nopOpcode); i++ {
		if origin[i] != nopOpcode[i] {
			return false
		}
	}
	return true
}
