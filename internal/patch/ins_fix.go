package patch

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/x86asm"
)

// callInsName call指令名称
const callInsName = "CALL"

// opExpand 短地址指令 -> 长地址指令
// 原始函数内部的短地址跳转无法满足长距离跳转时候,需要修改为长地址跳转, 因此同时需要将操作符指令修改为对应的长地址操作符指令
var opExpand = map[uint32][]byte{
	0x74: {0x0F, 0x84}, // JE: 74->0F
	0x76: {0x0F, 0x86}, // JBE: 76->0F
	0x7F: {0x0F, 0x8F},
	0xEB: {0xE9}, // JMP: EB->E9 (Jump near, relative, displacement relative to next instruction.)
}

// replaceRelativeAddr 替换函数字节码中的相对地址(如果有的话)
// from 函数起始地址
// copyOrigin 函数字节码
// trampoline 需要移动到的目标地址
// funcSize 函数字节码整体长度
// leastSize 要替换的字节长度的最小限制
// allowCopyCall 是否允许拷贝Call指令
func replaceRelativeAddr(from uintptr, copyOrigin []byte, trampoline uintptr, funcSize int, leastSize int,
	allowCopyCall bool) ([]byte, int, error) {

	// try replace and get the len(applyEndPos) to replace
	_, applyEndPos, err := replaceBlock(from, copyOrigin, trampoline, leastSize, funcSize, allowCopyCall)
	if err != nil {
		return nil, 0, err
	}

	// check if exists jump back to [0:applyEndPos], return not support
	if err = checkJumpBetween(from, applyEndPos, copyOrigin, funcSize); err != nil {
		return nil, 0, err
	}

	logger.LogDebugf("fix size: %d", applyEndPos)

	// real replace
	replacedBlock, _, err := replaceBlock(from, copyOrigin, trampoline, applyEndPos, applyEndPos, allowCopyCall)
	return replacedBlock, applyEndPos, err
}

// replaceBlock 替换函数字节码中的相对地址(如果有的话)
// from 起始位置(实际地址)
// block 被替换的目标字节区块
// trampoline 跳板函数起始地址
// leastSize 最少替换范围
// blockSize 目标区块范围, 用于判断地址是否超出block范围, 超出才需要替换
// return []byte 修复后的指令
func replaceBlock(from uintptr, block []byte, trampoline uintptr,
	leastSize int, blockSize int, allowCopyCall bool) ([]byte, int, error) {
	startAddr := (uint64)(from)
	replacedBlock := make([]byte, 0)

	logger.LogDebug("target fix ins >>>>>")

	for pos := 0; pos < len(block); {
		ins, _, err := nextIns(pos, block)
		if err != nil {
			panic("replaceRelativeAddr err:" + err.Error())
		}

		if ins != nil && ins.Opcode != 0 {
			if !allowCopyCall && ins.Op.String() == callInsName {
				return nil, 0, fmt.Errorf("copy call instruction is not allowed in auto trampoline model. size: %d", leastSize)
			}

			// 拷贝一份block,防止被修改; replace之后的block是回参replacedBlock
			copyBlock := make([]byte, len(block))
			if l := copy(copyBlock, block); l != len(block) {
				return nil, 0, errors.New("copy block array error")
			}
			replaced := replaceIns(ins, pos, copyBlock, blockSize, startAddr, trampoline)
			replacedBlock = append(replacedBlock, replaced...)

			logger.LogDebugf("[%d]>[%d] 0x%x:\t%s\t\t%s\t\t%s", ins.Len, len(replaced),
				startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(replaced))
		}

		pos = pos + ins.Len

		// for fix only first few inst, not copy all func inst
		if leastSize > 0 && pos >= leastSize {
			ins, _, err := nextIns(pos, block)
			if err != nil {
				panic("replaceRelativeAddr err:" + err.Error())
			}
			// fix jump to RET err: signal SIGSEGV: segmentation violation
			if ins != nil && ins.String() != "RET" {
				return replacedBlock, pos, nil
			}
		}
	}

	return replacedBlock, 0, nil
}

// replaceIns 替换单条指令
func replaceIns(ins *x86asm.Inst, pos int, block []byte, blockSize int,
	startAddr uint64, trampoline uintptr) []byte {
	// 需要替换偏移地址
	if ins.PCRelOff <= 0 {
		return block[pos : pos+ins.Len]
	}
	offset := pos + ins.PCRelOff

	relativeAddr := decodeRelativeAddr(ins, block, offset)

	// TODO 待实现
	//if ins.PCRel <= 1 {
	//	// 1字节相对地址暂时忽略
	//	return
	//}

	logger.LogDebugf("ins relative [%d] need fix : ", (relativeAddr)+pos+ins.Len)

	if (relativeAddr > 0 && (relativeAddr)+pos+ins.Len >= blockSize) ||
		(relativeAddr < 0 && (relativeAddr)+pos+ins.Len < 0) {
		if ins.Op.String() == callInsName {
			logger.LogDebug((int64)(startAddr)-(int64)(trampoline), startAddr, trampoline, int32(relativeAddr))
		}

		var encoded = encodeAddress(block[pos:offset],
			block[offset:offset+ins.PCRel], ins.PCRel, relativeAddr, (int)(startAddr)-(int)(trampoline))

		if logger.LogLevel <= logger.DebugLevel {
			// 打印替换之后的指令
			ins, err := x86asm.Decode(block[pos:pos+ins.Len], 64)
			if err == nil {
				logger.LogInfof("replaced: \t%s\t\t%s", ins.Op, ins.String())
			}
		}

		if len(encoded) > ins.PCRel {
			return encoded
		}
	} else {
		if ins.Op.String() == callInsName {
			logger.LogDebug((relativeAddr)+pos+ins.Len, blockSize, (relativeAddr)+pos+ins.Len)
			logger.LogDebug("called")
		}
	}

	return block[pos : pos+ins.Len]
}

// nextIns nextIns
func nextIns(pos int, copyOrigin []byte) (*x86asm.Inst, []byte, error) {
	if pos >= len(copyOrigin) {
		return nil, nil, nil
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

	return &ins, code, err
}

// checkJumpBetween check if exists the instruct of function that jump into address between :from and :to
// if exists, return error
func checkJumpBetween(from uintptr, to int, copyOrigin []byte, funcSize int) error {
	for pos := 0; pos <= funcSize; {
		ins, code, err := nextIns(pos, copyOrigin)
		if err != nil {
			panic("checkJumpBetween err:" + err.Error())
		}
		if ins == nil {
			break
		}
		if ins.PCRelOff <= 0 {
			pos = pos + ins.Len
			continue
		}
		offset := pos + ins.PCRelOff
		relativeAddr := decodeRelativeAddr(ins, copyOrigin, offset)
		if ((relativeAddr)+pos+ins.Len < to) &&
			((relativeAddr)+pos+ins.Len > 0) {
			logger.LogErrorf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s\t\tabs:0x%x", ins.Len,
				from+uintptr(pos), ins.Op, ins.String(), hex.EncodeToString(code),
				from+uintptr(pos)+uintptr(relativeAddr)+uintptr(ins.Len))
			return fmt.Errorf("not support of jump to inside of the first 13 bytes\n"+
				"jump address is: 0x%x", (uintptr)((relativeAddr)+pos+ins.Len)+from)
		}
		pos = pos + ins.Len
	}
	return nil
}

// decodeRelativeAddr decode relative address, if jump to the front of current pos, return negative values
func decodeRelativeAddr(ins *x86asm.Inst, block []byte, offset int) int {
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

	relativeAddr := decodeAddress(block[offset:offset+ins.PCRel], ins.PCRel)
	if !isAdd && relativeAddr > 0 {
		relativeAddr = -relativeAddr
	}
	return relativeAddr
}

// encodeAddress 写入地址参数到函数字节码
// len 地址值的位数
// val 地址值
// add 偏移量, 可为负数
func encodeAddress(ops []byte, addr []byte, addrLen int, val int, add int) []byte {
	result := make([]byte, 0)

	if addrLen == 1 {
		if isByteOverflow((int32)(int8(val)) + (int32)(add)) {
			if opsNew, ok := opExpand[uint32(ops[0])]; ok {
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
			if opsNew, ok := opExpand[uint32(ops[0])<<16+uint32(ops[1])]; ok {
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
