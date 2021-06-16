package patch

import "unsafe"

// nopOpcode 空指令插入到原函数开头第一个字节, 用于判断原函数是否已经被Patch过
var nopOpcode []byte = []byte{0xD5, 0x03, 0x20, 0x1F}

// funcPrologue 函数的开头指纹,用于不同OS获取不同的默认值
var funcPrologue = armFuncPrologue64

// Reference instruction set:
//
// MOVZ x27, 0xffff, LSL #16
// MOVK x27, 0x1234, LSL #32
// LDR  x28, [x27]
// BR   x28
//
// Copy test data into so.s and then
//
//   clang -c so.s -o so.o
//   objdump-d so.o
//
// You will see something like
//
// 0000000000000000 <ltmp0>:
//     0: fb ff bf d2    mov  x27, #4294901760
//     4: 9b 46 c2 f2    movk x27, #4660, lsl #32
//     8: 7c 03 40 f9    ldr  x28, [x27]
//     c: 80 03 1f d6    br   x28

// Assembles a jump to a function value
func jmpToFunctionValue(_, to uintptr) []byte {
	var res []byte

	// as you probably know there's no generic way to set a direct
	// ("immediate") 64-bit value in one instruction. The easiest
	// method is to split qword into 4 words. In our case, we can
	// just map a number into [4]uint and then iterate over nthem
	var buf [4]uint16
	*(*uint64)(unsafe.Pointer(&buf[0])) = uint64(to)

	// now assemble a sequence of `MOVx x0, <vi>, LSL i*16 instructions,
	// where the first of the sequence is
	//   MOVZ x27, vN, LSL #16*N
	// and the rest are
	//   MOVK x27, vi, LSL #18*i
	// where vX denotes a word that is not zero and X is its index in
	// their sequence
	//
	// Example: if we are to set value 0xffff_0000_00ff_0000
	// then the sequence of instructions is
	//    MOVZ x27, 0x00ff, LSL #16
	//    MOVK x27, 0xffff, LSL #48
	var thereWerentNonZeroes bool
	for i := range buf {
		if buf[i] == 0 {
			continue
		}
		if !thereWerentNonZeroes {
			res = append(res, armv8ImmediateMovz(27, 0, buf[i])...)
			thereWerentNonZeroes = true
		} else {
			res = append(res, arm8ImmediateMovk(27, i, buf[i])...)
		}
	}

	// assemble the rest:
	//  LDR x28, [x27]
	//  BR	x28
	res = append(res, 0x7c, 0x03, 0x40, 0xf9)
	res = append(res, 0x80, 0x03, 0x1f, 0xd6)

	return res
}

func armv8ImmediateMovz(regN, shift int, value uint16) []byte {
	return movx(0b10, regN, shift, value)
}

func arm8ImmediateMovk(regN, shift int, value uint16) []byte {
	return movx(0b11, regN, shift, value)
}

func movx(movX, regN, shift int, value uint16) []byte {
	var val uint32

	val = 0b10010010 << 24
	val |= 0b10000000 << 16
	val |= uint32(movX&0x3) << 29
	val |= uint32(shift&3) << 21
	val |= uint32(value) << 5
	val |= uint32(regN)

	var res [4]byte
	*(*uint32)(unsafe.Pointer(&res[0])) = val

	return res[:]
}

// jmpToOriginFunctionValue Assembles a jump to a function value
func jmpToOriginFunctionValue(from, to uintptr) (value []byte) {
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
