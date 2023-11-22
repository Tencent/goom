package bytecode

import (
	"reflect"

	"github.com/Jakegogo/goom_mocker/internal/arch/x86asm"
	"github.com/Jakegogo/goom_mocker/internal/bytecode/memory"
	"github.com/Jakegogo/goom_mocker/internal/logger"
)

func init() {
	// call fake
	_ = callFakeFunc()
	// check call asm code
	checkInlineDisable()
}

// checkInlineDisable 检测是否关闭 inline
func checkInlineDisable() {
	addr := reflect.ValueOf(callFakeFunc).Pointer()
	bytes := memory.RawRead(addr, 100)

	existsCallIns := false
	for pos := 0; pos < len(bytes); {
		ins, _, err := ParseIns(pos, bytes)
		if err != nil {
			logger.Warningf("goom resolve inline err: %v", err)
			break
		}

		if ins.Op == x86asm.CALL {
			existsCallIns = true
			break
		}
		if ins.Op == x86asm.INT && ins.Args[0].String() == "0x3" {
			break
		}
		pos = pos + ins.Len
	}

	if !existsCallIns {
		logger.Warningf("go inline is not disable, please use the build param: -gcflags=all=-l")
		logger.Consolef(logger.WarningLevel, "go inline is not disable, please use the build param: -gcflags=all=-l")
	}
}

// callFakeFunc 内联测试函数
func callFakeFunc() int {
	return checkTarget(1)
}

// checkTarget 测试目标 mock 函数
func checkTarget(i int) int {
	// example short code
	return i + 1
}
