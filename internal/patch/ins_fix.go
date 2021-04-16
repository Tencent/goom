package patch

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/x86asm"
)

// callInsName call指令名称
const callInsName = "CALL"

// opAddrExpand 短地址指令 -> 长地址指令
// 原始函数内部的短地址跳转无法满足长距离跳转时候,需要修改为长地址跳转, 因此同时需要将指令修改为对应的长地址指令
var opAddrExpand = map[uint32][]byte{
	0x74: {0x0F, 0x84}, // JE 74->0F
	0x7F: {0x0F, 0x8F},
}

// replaceRelativeAddr 替换函数字节码中的相对地址(如果有的话)
// from 函数起始地址
// copyOrigin 函数字节码
// placeholder 需要移动到的目标地址
// funcSize 函数字节码整体长度
// leastSize 要替换的字节长度的最小限制
// allowCopyCall 是否允许拷贝Call指令
func replaceRelativeAddr(from uintptr, copyOrigin []byte, placeholder uintptr, funcSize int, leastSize int,
	allowCopyCall bool) ([]byte, int, error) {
	replaceOrigin, err := doReplaceRelativeAddr(from, copyOrigin, placeholder, funcSize, leastSize, allowCopyCall)
	if err != nil {
		return nil, 0, err
	}

	var replaceNew = replaceOrigin

	if leastSize > 0 {
		replaceNew, err = doReplaceRelativeAddr(from, replaceOrigin, placeholder, len(replaceOrigin), leastSize,
			allowCopyCall)
	}

	return replaceNew, len(replaceOrigin), err
}

//doReplaceRelativeAddr 替换函数字节码中的相对地址(如果有的话)
func doReplaceRelativeAddr(from uintptr, copyOrigin []byte, placeholder uintptr, funcSize int, leastSize int,
	allowCopyCall bool) ([]byte, error) {
	startAddr := (uint64)(from)
	result := make([]byte, 0)

	logger.LogDebug("target fix ins >>>>>")

	for pos := 0; pos < len(copyOrigin); {
		ins, err := nextIns(pos, copyOrigin)
		if err != nil {
			panic("replaceRelativeAddr err:" + err.Error())
		}

		if ins != nil && ins.Opcode != 0 {
			if !allowCopyCall && ins.Op.String() == callInsName {
				return nil, fmt.Errorf("copy call instruction is not allowed in auto trampoline model. size: %d", leastSize)
			}

			replaced := replaceIns(ins, pos, copyOrigin, funcSize, startAddr, placeholder)
			result = append(result, replaced...)

			logger.LogDebugf("[%d]>[%d] 0x%x:\t%s\t\t%s\t\t%s", ins.Len, len(replaced),
				startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(replaced))
		}

		pos = pos + ins.Len

		// for fix only first few inst, not copy all func inst
		if leastSize > 0 && pos >= leastSize {
			ins, err := nextIns(pos, copyOrigin)
			if err != nil {
				panic("replaceRelativeAddr err:" + err.Error())
			}
			// fix jump to RET err: signal SIGSEGV: segmentation violation
			if ins != nil && ins.String() != "RET" {
				return result, nil
			}
		}
	}

	return result, nil
}

// nextIns nextIns
func nextIns(pos int, copyOrigin []byte) (*x86asm.Inst, error) {
	if pos >= len(copyOrigin) {
		return nil, nil
	}
	// read 16 bytes at most each time
	endPos := pos + 16
	if endPos > len(copyOrigin) {
		endPos = len(copyOrigin)
	}

	code := copyOrigin[pos:endPos]

	ins, err := x86asm.Decode(code, 64)
	if err != nil {
		logger.LogError("decode assembly code err:", err)
	}

	return &ins, err
}

// replaceIns 替换单条指令
func replaceIns(ins *x86asm.Inst, pos int, copyOrigin []byte, funcSize int,
	startAddr uint64, placeholder uintptr) []byte {
	// 需要替换偏移地址
	if ins.PCRelOff <= 0 {
		return copyOrigin[pos : pos+ins.Len]
	}

	offset := pos + ins.PCRelOff

	var isAdd = true

	for i := 0; i < len(ins.Args); i++ {
		arg := ins.Args[i]
		if arg == nil {
			break
		}

		addrArgs := arg.String()
		if strings.HasPrefix(addrArgs, ".-") || strings.Contains(addrArgs, "RIP-") {
			isAdd = false
		}
	}

	relativeAddr := decodeAddress(copyOrigin[offset:offset+ins.PCRel], ins.PCRel)
	if !isAdd && relativeAddr > 0 {
		relativeAddr = -relativeAddr
	}

	// TODO 待实现
	//if ins.PCRel <= 1 {
	//	// 1字节相对地址暂时忽略, 不太可能跳出当前函数地址范围
	//	return
	//}

	logger.LogDebug("ins relative:", (relativeAddr)+pos+ins.Len)

	if (isAdd && (relativeAddr)+pos+ins.Len >= funcSize) ||
		(!isAdd && (relativeAddr)+pos+ins.Len < 0) {
		if ins.Op.String() == callInsName {
			logger.LogDebug((int64)(startAddr)-(int64)(placeholder), startAddr, placeholder, int32(relativeAddr))
		}

		var encoded = encodeAddress(copyOrigin[pos:offset],
			copyOrigin[offset:offset+ins.PCRel], ins.PCRel, relativeAddr, (int)(startAddr)-(int)(placeholder))

		ins, err := x86asm.Decode(copyOrigin[pos:pos+ins.Len], 64)
		if err == nil {
			//d := color.New(color.FgGreen, color.BgGray)
			logger.LogInfof("replaced: \t%s\t\t%s", ins.Op, ins.String())
		}

		if len(encoded) > ins.PCRel {
			return encoded
		}
	} else {
		if ins.Op.String() == callInsName {
			logger.LogDebug((relativeAddr)+pos+ins.Len, funcSize, (relativeAddr)+pos+ins.Len)
			logger.LogDebug("called")
		}
	}

	return copyOrigin[pos : pos+ins.Len]
}

// encodeAddress 写入地址参数到函数字节码
// len 地址值的位数
// val 地址值
// add 偏移量, 可为负数
func encodeAddress(ops []byte, addr []byte, addrLen int, val int, add int) []byte {
	result := make([]byte, 0)

	if addrLen == 1 {
		if isByteOverflow((int32)(int8(val)) + (int32)(add)) {
			if opsNew, ok := opAddrExpand[uint32(ops[0])]; ok {
				addr = make([]byte, 4)
				LittleEndian.PutInt32(addr, (int32)(int8(val))+int32(add)-
					int32(len(addr)-addrLen)-int32(len(opsNew)-len(ops))) // 新增了4个字节,需要减去

				ops = opsNew
			} else {
				panic("address overflow:" + hex.EncodeToString(ops) + ", addr:" + hex.EncodeToString(addr[:addrLen]))
			}
		} else {
			addr[0] = byte((int)(int8(val)) + add)
		}
	} else if addrLen == 2 {
		if isInt16Overflow((int32)(int8(val)) + (int32)(add)) {
			if opsNew, ok := opAddrExpand[uint32(ops[0]<<16+ops[1])]; ok {
				addr = make([]byte, 4)
				LittleEndian.PutInt32(addr, (int32)(int8(val))+int32(add)-
					int32(len(addr)-addrLen)-int32(len(opsNew)-len(ops))) // 新增了4个字节,需要减去
				ops = opsNew
			} else {
				panic("address overflow:" + hex.EncodeToString(ops) + ", addr:" + hex.EncodeToString(addr[:addrLen]))
			}
		}
		LittleEndian.PutInt16(addr, int16(val)+int16(add))
	} else if addrLen == 4 {
		LittleEndian.PutInt32(addr, int32(val)+int32(add))
	} else if addrLen == 8 {
		LittleEndian.PutInt64(addr, int64(val)+int64(add))
	}

	result = append(result, ops...)

	return append(result, addr...)
}

// isByteOverflow 字节是否溢出
func isByteOverflow(v int32) bool {
	if v > 0 {
		if v > math.MaxInt8 {
			return true
		}
	} else {
		if v < math.MinInt8 {
			return true
		}
	}

	return false
}

// isInt16Overflow  init16是否溢出
func isInt16Overflow(v int32) bool {
	if v > 0 {
		if v > math.MaxInt16 {
			return true
		}
	} else {
		if v < math.MinInt16 {
			return true
		}
	}

	return false
}

// nolint
func isInt32Overflow(v int64) bool {
	if v > 0 {
		if v > math.MaxInt32 {
			return true
		}
	} else {
		if v < math.MinInt32 {
			return true
		}
	}
	return false
}

// decodeAddress 从函数字节码中解析地址数值
// len 地址值的位数
func decodeAddress(bytes []byte, len int) int {
	if len == 1 {
		return int(int8(bytes[0]))
	} else if len == 2 {
		return int(LittleEndian.Int16(bytes))
	} else if len == 4 {
		return int(LittleEndian.Int32(bytes))
	} else if len == 8 {
		return int(LittleEndian.Int64(bytes))
	}

	return 0
}
