package memory

import (
	"reflect"
	"unsafe"

	"github.com/tencent/goom/internal/logger"
	"github.com/tencent/goom/internal/unexports"
)

// 指令生成相关
// nolint
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
//
//go:linkname WriteICacheFn github.com/tencent/goom/internal/bytecode/stub.WriteICacheFn
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
