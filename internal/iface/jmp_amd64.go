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

// jmpWithRdxAndCtx Assembles a jump to a function value
// ctx context 地址
// to 跳转目标地址
// from 跳转来源地址
func jmpWithRdxAndCtx(ctx, to, from uintptr) (value []byte) {
	var dis uint32
	if to > from {
		dis = uint32(int32(to-from) - 5)
		dis = dis + 10
	} else {
		dis = uint32(-int32(from-to) - 5)
		dis = dis - 10
	}

	return []byte{
		0x48, 0xBA,
		byte(ctx),
		byte(ctx >> 8),
		byte(ctx >> 16),
		byte(ctx >> 24),
		byte(ctx >> 32),
		byte(ctx >> 40),
		byte(ctx >> 48),
		byte(ctx >> 56), // movabs rdx,ctx
		0xe9,
		byte(dis),
		byte(dis >> 8),
		byte(dis >> 16),
		byte(dis >> 24), // jmp dis
	}
}
