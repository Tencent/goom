package memory

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/unexports"
)

// 指令生成相关
const (
	_0b1      = 1  // 0b1
	_0b10     = 2  // 0b10
	_0b11     = 3  // 0b11
	_0b100101 = 37 // 0b100101
	_X0       = 0  // X0 寄存器
	_X1       = 1  // X1 寄存器
)

// cacheSize 需要失效的缓存长度, 一般情况不会超过这个数
const cacheSize = 2 * 1024

// icache 清理指令内容
var (
	insPadding = []byte{
		0xEB, 0x07, 0x40, 0xF9, // # 	ldr   	x11, [sp, #8]
		0x2B, 0x75, 0x0b, 0xd5, // #	ic		ivau, x11
		0x6B, 0x01, 0x01, 0x91, // # 	add		x11, x11, #1<<6
		0x2B, 0x75, 0x0b, 0xd5, // #	ic		ivau, x11
		0x6B, 0x01, 0x01, 0x91, // # 	add		x11, x11, #1<<6
		0x2B, 0x75, 0x0b, 0xd5, // #	ic		ivau, x11
		0x6B, 0x01, 0x01, 0x91, // # 	add		x11, x11, #1<<6
		0x2B, 0x75, 0x0b, 0xd5, // #	ic		ivau, x11
	}

	// clearICacheIns 清除指令缓存指令
	// 新版mac使用 AARCH64 (arm64v8), 指令缓存命中率较高, 会出现patch成功但是使用缓存中原来的指令来执行，导致 mock 失败
	clearICacheIns = []byte{
		//0x60, 0x00, 0xa0, 0xd2,//# 	mov 	x1,  [size] // 8个指令 $ 1-4
		0xE0, 0x07, 0x40, 0xF9, // # 	ldr   	x0, [sp, #8] $ 5
		0x09, 0xe4, 0x7a, 0x92, // # 	and		x9, x0, #~((1<<6)-1) $ 6 cacheline align address
		0x0a, 0x14, 0x40, 0x92, // # 	and		x10, x0, #((1<<6)-1) $ 7 extend length by alignment
		0x2a, 0x00, 0x0a, 0x8b, // # 	add		x10, x1, x10 $ 8
		0x4a, 0x05, 0x00, 0xd1, // # 	sub		x10, x10, #1 $ 9
		0x0b, 0x00, 0x80, 0x92, // # 	mov		x11, #-1 $ 10
		0x6a, 0x19, 0x4a, 0xca, // # 	eor		x10, x11, x10, lsr #6 $ 11 compute cacheline counter
		0x9f, 0x3b, 0x03, 0xd5, // # 	dsb		ish $ 12
	}

	clearICacheIns1 = []byte{
		// 循环
		0x29, 0x75, 0x0b, 0xd5, // #	ic		ivau, x9 $ 13
		0x29, 0x01, 0x01, 0x91, // # 	add		x9, x9, #1<<6 $ 14
		0x4a, 0x05, 0x00, 0xb1, // # 	adds	x10, x10, #1 $ 15
		0xa1, 0xff, 0xff, 0x54, // # 	b.ne 	#0xfffffffffffffff4 $ 16
	}

	clearICacheIns2 = []byte{
		0x9f, 0x3b, 0x03, 0xd5, // # 	dsb		ish $ 20
		0xdf, 0x3f, 0x03, 0xd5, // # 	isb # 21
		0xc0, 0x03, 0x5f, 0xd6, // #	ret # 22
	}
)

// clearFuncCache 缓存已经生成的函数
var clearFuncCache func(uintptr)

// ClearICache 清除指令缓存
func ClearICache(addr uintptr) {
	if clearFuncCache != nil {
		clearFuncCache(addr)
		return
	}
	if ok := initClearICacheFunc(); !ok {
		return
	}
	clearFuncCache(addr)
}

// initClearICacheFunc 初始化 clearFuncCache
func initClearICacheFunc() bool {
	f, err := makeFunc()
	if err != nil {
		clearFuncCache = func(uintptr) {}
		logger.Warningf("ClearICache failed: %v", err)
		return false
	}
	clearFuncCache = f
	return true
}

// makeFunc 构造 clearICacheIns 的函数实例
func makeFunc() (func(uintptr), error) {
	code := make([]byte, 0, 256)
	code = append(code, insPadding...)
	code = append(code, movAddr(_X1, uintptr(cacheSize))...)
	code = append(code, clearICacheIns...)
	code = append(code, clearICacheIns1...)
	code = append(code, clearICacheIns2...)

	var (
		addr uintptr
		err  error
	)
	if addr, err = WriteICacheFn(code); err != nil {
		return nil, err
	}
	var f func(uintptr)
	fn := unexports.NewFuncWithCodePtr(reflect.TypeOf(f), addr).Interface()
	return (fn).(func(uintptr)), nil
}

// WriteICacheFn 写入 icache clear 函数数据
//go:linkname WriteICacheFn git.code.oa.com/goom/mocker/internal/bytecode/stub.WriteICacheFn
func WriteICacheFn([]byte) (uintptr, error)

// movAddr 生成 mov x[?] [addr] 四个指令
func movAddr(r uint32, addr uintptr) (value []byte) {
	res := make([]byte, 0, 24)
	d0d1 := addr & 0xFFFF
	d2d3 := addr >> 16 & 0xFFFF
	d4d5 := addr >> 32 & 0xFFFF
	d6d7 := addr >> 48 & 0xFFFF

	res = append(res, movImm(r, _0b10, 0, d0d1)...) // MOVZ x0, double[16:0]
	res = append(res, movImm(r, _0b11, 1, d2d3)...) // MOVK x0, double[32:16]
	res = append(res, movImm(r, _0b11, 2, d4d5)...) // MOVK x0, double[48:32]
	res = append(res, movImm(r, _0b11, 3, d6d7)...) // MOVK x0, double[64:48]
	return res
}

// movImm 动态生成 mov 指令
// r 寄存器 X0~x27
// opc 操作数, MOVZ/MOVK
// shift 寄存器高低位段
// val 值数据
func movImm(r uint32, opc, shift int, val uintptr) []byte {
	var m uint32 = r           // rd
	m |= uint32(val) << 5      // imm16
	m |= uint32(shift&3) << 21 // hw
	m |= _0b100101 << 23       // const
	m |= uint32(opc&0x3) << 29 // opc
	m |= _0b1 << 31            // sf
	res := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&res[0])) = m
	return res
}
