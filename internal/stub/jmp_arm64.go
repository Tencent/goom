package stub

// jmpWithRdx Assembles a jump to a clourse function value
// dx DX 寄存器
func jmpWithRdx(dx uintptr) (value []byte) {
	panic("not support yet")
}

// jmpWithRdxAndCtx Assembles a jump to a function value
// ctx context 地址
// to 跳转目标地址
// from 跳转来源地址
func jmpWithRdxAndCtx(ctx, to, from uintptr) (value []byte) {
	panic("not support yet")
}
