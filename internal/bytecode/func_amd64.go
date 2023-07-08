package bytecode

import (
	"bytes"
	"encoding/hex"
	"strings"

	"git.woa.com/goom/mocker/internal/arch/x86asm"
	"git.woa.com/goom/mocker/internal/bytecode/memory"
	"git.woa.com/goom/mocker/internal/logger"
)

// defaultInsLen 默认一次解析指令的长度
const defaultInsLen = 16

// funcPrologue 函数的开头指纹,用于不同 OS 获取不同的默认值
var funcPrologue = defaultFuncPrologue64

// CallInsName call 指令名称
const CallInsName = "CALL"

// GetFuncSize get func binary size
// not absolutely safe
func GetFuncSize(mode int, start uintptr, minimal bool) (length int, err error) {
	funcSizeReadLock.Lock()
	defer func() {
		funcSizeCache[start] = length
		funcSizeReadLock.Unlock()
	}()

	if l, ok := funcSizeCache[start]; ok {
		return l, nil
	}

	prologueLen := len(funcPrologue)
	code := memory.RawRead(start, defaultInsLen)

	var (
		int3Found = false
		curLen    = 0
	)
	for {
		inst, err := x86asm.Decode(code, mode)
		if err != nil || (inst.Opcode == 0 && inst.Len == 1 && inst.Prefix[0] == x86asm.Prefix(code[0])) {
			return curLen, nil
		}

		if inst.Len == 1 && code[0] == 0xcc {
			// 0xcc -> int3, trap to debugger, padding to function end
			if minimal {
				return curLen, nil
			}
			int3Found = true
		} else if int3Found {
			return curLen, nil
		}

		curLen = curLen + inst.Len
		code = memory.RawRead(start+uintptr(curLen), defaultInsLen)
		if bytes.Equal(funcPrologue, code[:prologueLen]) {
			return curLen, nil
		}
	}
}

// PrintInstf 调试内存指令替换,对原指令、替换之后的指令进行输出对比
func PrintInstf(title string, from uintptr, copyOrigin []byte, level int) {
	if logger.LogLevel < level {
		return
	}
	logger.Important(title)

	startAddr := (uint64)(from)
	for pos := 0; pos < len(copyOrigin); {
		to := pos + defaultInsLen
		if to > len(copyOrigin) {
			to = len(copyOrigin)
		}

		code := copyOrigin[pos:to]
		ins, err := x86asm.Decode(code, 64)

		if err != nil {
			logger.Importantf("[0] 0x%x: inst decode error:%s", startAddr+(uint64)(pos), err)
			if ins.Len == 0 {
				pos = pos + 1
			} else {
				pos = pos + ins.Len
			}
			continue
		}

		if ins.Opcode == 0 {
			if ins.Len == 0 {
				pos = pos + 1
			} else {
				pos = pos + ins.Len
			}
			continue
		}

		if ins.PCRelOff <= 0 {
			logger.Importantf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s", ins.Len,
				startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(code[:ins.Len]))
			pos = pos + ins.Len
			continue
		}

		offset := pos + ins.PCRelOff
		relativeAddr := DecodeAddress(copyOrigin[offset:offset+ins.PCRel], ins.PCRel)
		if !isRelativeAdd(ins) && relativeAddr > 0 {
			relativeAddr = -relativeAddr
		}

		logger.Importantf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s\t\tabs:0x%x", ins.Len,
			startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(code[:ins.Len]),
			from+uintptr(pos)+uintptr(relativeAddr)+uintptr(ins.Len))

		pos = pos + ins.Len
	}
}

func isRelativeAdd(ins x86asm.Inst) bool {
	isAdd := true
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
	return isAdd
}

// GetInnerFunc Get the first real func location from wrapper
// not absolutely safe
func GetInnerFunc(mode int, start uintptr) (uintptr, error) {
	prologueLen := len(funcPrologue)
	code := memory.RawRead(start, defaultInsLen)

	var (
		int3Found = false
		curLen    = 0
	)
	for {
		inst, err := x86asm.Decode(code, mode)
		if err != nil || (inst.Opcode == 0 && inst.Len == 1 && inst.Prefix[0] == x86asm.Prefix(code[0])) {
			return 0, nil
		}

		if inst.Len == 1 && code[0] == 0xcc {
			int3Found = true
		} else if int3Found {
			return 0, nil
		}

		if inst.Op.String() == CallInsName {
			relativeAddr := DecodeRelativeAddr(&inst, code, inst.PCRelOff)
			if relativeAddr >= 0 {
				return start + uintptr(curLen) + uintptr(relativeAddr) + uintptr(inst.Len), nil
			}
			if curLen+int(relativeAddr) < 0 {
				return start + uintptr(curLen) - uintptr(-relativeAddr) + uintptr(inst.Len), nil
			}
		}

		curLen = curLen + inst.Len
		code = memory.RawRead(start+uintptr(curLen), defaultInsLen)
		if bytes.Equal(funcPrologue, code[:prologueLen]) {
			return 0, nil
		}
	}
}
