package patch

import (
	"bytes"

	"git.code.oa.com/goom/mocker/internal/x86asm"
)

// GetFuncSize get func binary size
// not absolutely safe
func GetFuncSize(mode int, start uintptr, minimal bool) (lenth int, err error) {
	funcSizeReadLock.Lock()
	defer func() {
		funcSizeCache[start] = lenth
		funcSizeReadLock.Unlock()
	}()

	if lenth, ok := funcSizeCache[start]; ok {
		return lenth, nil
	}

	prologueLen := len(funcPrologue)
	code := rawMemoryRead(start, 16) // instruction takes at most 16 bytes

	int3Found := false
	curLen := 0

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
		code = rawMemoryRead(start+uintptr(curLen), 16) // instruction takes at most 16 bytes

		if bytes.Equal(funcPrologue, code[:prologueLen]) {
			return curLen, nil
		}
	}
}
