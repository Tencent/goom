package patch

import (
	"reflect"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/x86asm"
)

func init() {
	// call fake
	callFakeFunc()
	// check call asm code
	checkInlineDisbale()
}

// checkInlineDisbale 检测是否关闭inline
func checkInlineDisbale() {
	addr := reflect.ValueOf(callFakeFunc).Pointer()
	bytes := rawMemoryRead(addr, 100)

	hasCallIns := false
	for pos := 0; pos < len(bytes); {
		ins, err := nextIns(pos, bytes)
		if err != nil {
			logger.LogWarningf("goom resolve inline err: %v", err)
			break
		}

		if ins.Op == x86asm.CALL {
			hasCallIns = true
			break
		}
		if ins.Op == x86asm.INT && ins.Args[0].String() == "0x3" {
			break
		}
		pos = pos + ins.Len
	}

	if !hasCallIns {
		logger.LogWarningf("go inline is not disable, please use the build param: -gcflags=all=-l")
		logger.Log2Consolef("go inline is not disable, please use the build param: -gcflags=all=-l")
	}
}

// callFakeFunc 内联测试函数
func callFakeFunc() int {
	return target(1)
}

// target 测试目标mock函数
func target(i int) int {
	// example short code
	return i + 1
}
