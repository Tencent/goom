package patch

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Jakegogo/goom_mocker/internal/arch/x86asm"
	"github.com/Jakegogo/goom_mocker/internal/bytecode"
	"github.com/Jakegogo/goom_mocker/internal/logger"
)

// allowCopyCall 是否允许拷贝 Call 指令
const allowCopyCall = true

// fixRelativeAddr 修复函数字节码中的相对地址(如果有的话)
// from 函数起始地址
// copyOrigin 函数字节码
// trampoline 需要移动到的目标地址
// funcSize 函数字节码整体长度
// leastSize 要替换的字节长度的最小限制
func fixRelativeAddr(from uintptr, copyOrigin []byte, trampoline uintptr, funcSize int, leastSize int) (
	fixedData []byte, fixedDataSize int, err error) {

	// try to replace and get the len(endPos) to replace
	_, fixedDataSize, err = fixBlock(from, copyOrigin, trampoline, leastSize, funcSize)
	if err != nil {
		return
	}

	// check if exists jump back to [0:endPos], return not support
	if err = checkJumpBetween(from, fixedDataSize, copyOrigin, funcSize); err != nil {
		return
	}

	logger.Debugf("fix size: %d", fixedDataSize)

	// real replace
	fixedData, _, err = fixBlock(from, copyOrigin, trampoline, fixedDataSize, fixedDataSize)
	return
}

// fixBlock 替换函数字节码中的相对地址(如果有的话)
// from 起始位置(实际地址)
// block 被替换的目标字节区块
// trampoline 跳板函数起始地址
// leastSize 最少替换范围
// blockSize 目标区块范围, 用于判断地址是否超出 block 范围, 超出才需要替换
// return []byte 修复后的指令
func fixBlock(from uintptr, block []byte, trampoline uintptr,
	leastSize int, blockSize int) (fixedData []byte, fixedDataSize int, err error) {
	var (
		fixedBlock = make([]byte, 0)
	)
	logger.Debug("target fix ins >>>>>")

	for pos := 0; pos < len(block); {
		ins, _, err := bytecode.ParseIns(pos, block)
		if err != nil {
			panic("fixRelativeAddr err:" + err.Error())
		}

		if ins != nil && ins.Opcode != 0 {
			if !allowCopyCall && ins.Op.String() == bytecode.CallInsName {
				return nil, 0,
					fmt.Errorf("copy call instruction is not allowed in auto trampoline model. size: %d", leastSize)
			}
			// 拷贝一份 block,防止被修改; replace 之后的 block 是回参 fixedBlock
			copyBlock := make([]byte, len(block))
			if l := copy(copyBlock, block); l != len(block) {
				return nil, 0, errors.New("copy block array error")
			}
			fixedInsData := fixIns(ins, pos, copyBlock, blockSize, (uint64)(from), trampoline)
			fixedBlock = append(fixedBlock, fixedInsData...)

			logger.Debugf("[%d]>[%d] 0x%x:\t%s\t\t%s\t\t%s", ins.Len, len(fixedInsData),
				(uint64)(from)+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(fixedInsData))
		}

		pos = pos + ins.Len

		// for fix only first few inst, not copy all func inst
		if leastSize > 0 && pos >= leastSize {
			ins, _, err := bytecode.ParseIns(pos, block)
			if err != nil {
				panic("fixRelativeAddr err:" + err.Error())
			}
			// fix jump to RET err: signal SIGSEGV: segmentation violation
			if ins != nil && ins.String() != "RET" {
				return fixedBlock, pos, nil
			}
		}
	}

	return fixedBlock, len(fixedBlock), nil
}

// fixIns 替换单条指令的偏移地址
func fixIns(ins *x86asm.Inst, pos int, block []byte, blockSize int,
	from uint64, trampoline uintptr) []byte {
	if ins.PCRelOff <= 0 {
		// 不需要替换偏移地址
		return block[pos : pos+ins.Len]
	}
	offset := pos + ins.PCRelOff
	addr := bytecode.DecodeRelativeAddr(ins, block, offset)

	// TODO 待实现
	//if ins.PCRel <= 1 {
	//	// 1字节相对地址暂时忽略
	//	return
	//}

	logger.Debugf("ins relative [%d] need fix : ", (addr)+pos+ins.Len)

	if (addr > 0 && (addr)+pos+ins.Len >= blockSize) ||
		(addr < 0 && (addr)+pos+ins.Len < 0) {
		if ins.Op.String() == bytecode.CallInsName {
			logger.Debug((int64)(from)-(int64)(trampoline), from, trampoline, int32(addr))
		}
		if logger.LogLevel <= logger.DebugLevel {
			ins, err := x86asm.Decode(block[pos:pos+ins.Len], 64)
			if err == nil {
				logger.Infof("replaced: \t%s\t\t%s", ins.Op, ins.String())
			}
		}

		result := bytecode.EncodeAddress(block[pos:offset],
			block[offset:offset+ins.PCRel], ins.PCRel, addr, (int)(from)-(int)(trampoline))
		if len(result) > ins.PCRel {
			return result
		}
	} else {
		if ins.Op.String() == bytecode.CallInsName {
			logger.Debug((addr)+pos+ins.Len, blockSize, (addr)+pos+ins.Len)
		}
	}

	// nothing to fix
	return block[pos : pos+ins.Len]
}

// checkJumpBetween check if exists the instruction in originData of function
// that jump into address between :from and :to
// if exists, return error and not support to fix.
func checkJumpBetween(from uintptr, to int, originData []byte, funcSize int) error {
	for pos := 0; pos <= funcSize; {
		ins, code, err := bytecode.ParseIns(pos, originData)
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
		relativeAddr := bytecode.DecodeRelativeAddr(ins, originData, offset)
		if ((relativeAddr)+pos+ins.Len < to) &&
			((relativeAddr)+pos+ins.Len > 0) {
			logger.Errorf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s\t\tabs:0x%x", ins.Len,
				from+uintptr(pos), ins.Op, ins.String(), hex.EncodeToString(code),
				from+uintptr(pos)+uintptr(relativeAddr)+uintptr(ins.Len))
			return fmt.Errorf("not support of jump to inside of the first 13 bytes\n"+
				"jump address is: 0x%x", (uintptr)((relativeAddr)+pos+ins.Len)+from)
		}
		pos = pos + ins.Len
	}
	return nil
}
