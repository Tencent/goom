package patch

// nopOpcode 空指令插入到原函数开头第一个字节, 用于判断原函数是否已经被Patch过
const nopOpcode = 0x90

// funcPrologue 函数的开头指纹,用于不同OS获取不同的默认值
var funcPrologue = defaultFuncPrologue32

// jmpToFunctionValue Assembles a jump to a function value
func jmpToFunctionValue(to uintptr) []byte {
	return []byte{
		0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24), // mov edx,to
		0xFF, 0x22,     // jmp DWORD PTR [edx]
	}
}
