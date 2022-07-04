package iface

// jmpWithRdx Assembles a jump to a clourse function value
// dx DX 寄存器
func jmpWithRdx(dx uintptr) (value []byte) {
	return []byte{
		0x48, 0xBA,
		byte(dx),
		byte(dx >> 8),
		byte(dx >> 16),
		byte(dx >> 24),
		byte(dx >> 32),
		byte(dx >> 40),
		byte(dx >> 48),
		byte(dx >> 56), // movabs rdx,dx
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}
}
