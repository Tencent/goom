package iface

import "unsafe"

const (
	_0b1      = 1  // _0b1
	_0b10     = 2  // 0b10
	_0b11     = 3  // 0b11
	_0b100101 = 37 // 0b100101
)

// jmpWithRdx Assembles a jump to a clourse function value
// dx DX 寄存器
func jmpWithRdx(dx uintptr) (value []byte) {
	res := make([]byte, 0, 24)
	d0d1 := dx & 0xFFFF
	d2d3 := dx >> 16 & 0xFFFF
	d4d5 := dx >> 32 & 0xFFFF
	d6d7 := dx >> 48 & 0xFFFF

	res = append(res, movImm(_0b10, 0, d0d1)...)         // MOVZ x26, double[16:0]
	res = append(res, movImm(_0b11, 1, d2d3)...)         // MOVK x26, double[32:16]
	res = append(res, movImm(_0b11, 2, d4d5)...)         // MOVK x26, double[48:32]
	res = append(res, movImm(_0b11, 3, d6d7)...)         // MOVK x26, double[64:48]
	res = append(res, []byte{0x5B, 0x03, 0x40, 0xF9}...) // LDR x27, [x26]
	res = append(res, []byte{0x60, 0x03, 0x1F, 0xD6}...) // BR x27
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

// jmpWithRdxAndCtx Assembles a jump to a function value
// ctx context 地址
// to 跳转目标地址
// from 跳转来源地址
func jmpWithRdxAndCtx(ctx, _, _ uintptr) (value []byte) {
	res := make([]byte, 0, 40)
	d0d1 := ctx & 0xFFFF
	d2d3 := ctx >> 16 & 0xFFFF
	d4d5 := ctx >> 32 & 0xFFFF
	d6d7 := ctx >> 48 & 0xFFFF

	res = append(res, movImm(_0b10, 0, d0d1)...)         // MOVZ x26, double[16:0]
	res = append(res, movImm(_0b11, 1, d2d3)...)         // MOVK x26, double[32:16]
	res = append(res, movImm(_0b11, 2, d4d5)...)         // MOVK x26, double[48:32]
	res = append(res, movImm(_0b11, 3, d6d7)...)         // MOVK x26, double[64:48]
	res = append(res, []byte{0x5B, 0x03, 0x40, 0xF9}...) // LDR x27, [x26]
	res = append(res, []byte{0x60, 0x03, 0x1F, 0xD6}...) // BR x27
	return res
}
