package patch

// OpAddrExpand 短地址指令 -> 长地址指令
// 原始函数内部的短地址跳转无法满足长距离跳转时候,需要修改为长地址跳转, 因此同时需要将指令修改为对应的长地址指令
var OpAddrExpand = map[uint32][]byte{
	0x74: []byte{0x0F, 0x84}, // JE 74->0F
	0x7F: []byte{0x0F, 0x8F},
}
