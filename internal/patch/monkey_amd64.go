package patch

import "unsafe"

var funcPrologue = defaultFuncPrologue64

// NOP_OPCODE 空指令插入到原函数开头第一个字节, 用于判断原函数是否已经被Patch过
const NopOpcode = 0x90

// Assembles a jump to a function value
func jmpToFunctionValue(from, to uintptr) (value []byte) {
	return []byte{
		0x90, // NOP
		0x48, 0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}
}

// Assembles a jump to a function value
func jmpToOriginFunctionValue(from, to uintptr) (value []byte) {

	if relative(from, to) {
		var dis uint32
		if to > from {
			dis = uint32(int32(to-from) - 5)
		} else {
			dis = uint32(-int32(from-to) - 5)
		}
		return []byte{
			0xe9,
			byte(dis),
			byte(dis >> 8),
			byte(dis >> 16),
			byte(dis >> 24),
		}
	}

	return []byte{
		0x48, 0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}
}

func relative(from uintptr, to uintptr) bool {
	delta := int64(from - to)
	if unsafe.Sizeof(uintptr(0)) == unsafe.Sizeof(int32(0)) {
		delta = int64(int32(from - to))
	}

	// 跨度大于2G时
	relative := (delta <= 0x7fffffff)

	if delta < 0 {
		delta = -delta
		relative = (delta <= 0x80000000)
	}
	return relative
}