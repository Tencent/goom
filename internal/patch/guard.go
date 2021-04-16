package patch

import (
	"fmt"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// PatchGuard 代理执行控制句柄, 可通过此对象进行代理还原
type PatchGuard struct {
	target        uintptr // 被patch的函数
	originFuncPtr uintptr // 修复的函数指针
	jumpBytes     []byte  // 跳转指令字节
	originalBytes []byte  // 原始字节码
	applied       bool    // 是否已经被应用
}

// Apply 执行
func (g *PatchGuard) Apply() {
	Lock()
	defer Unlock()

	g.applied = true
	// 执行函数调用地址替换(延迟执行)
	if err := CopyToLocation(g.target, g.jumpBytes); err != nil {
		logger.LogWarningf("Apply to 0x%x error: %s", g.target, err)
	}

	Debug(fmt.Sprintf("apply copy to 0x%x", g.target), g.target, 20, logger.DebugLevel)
}

// Unpatch 取消代理,还原指令码
// 外部调用请使用PatchGuard.UnpatchWithLock()
func (g *PatchGuard) Unpatch() {
	if g != nil && g.applied {
		_ = CopyToLocation(g.target, g.originalBytes)
		Debug(fmt.Sprintf("unpatch copy to 0x%x", g.target), g.target, 20, logger.DebugLevel)
	}
}

// UnpatchWithLock 外部调用需要加锁
func (g *PatchGuard) UnpatchWithLock() {
	Lock()
	defer Unlock()

	g.Unpatch()
}

// Restore 重新应用代理
func (g *PatchGuard) Restore() {
	Lock()
	defer Unlock()

	if g != nil && g.applied {
		_ = CopyToLocation(g.target, g.jumpBytes)
		Debug(fmt.Sprintf("unpatch copy to 0x%x", g.target), g.target, 20, logger.DebugLevel)
	}
}

// OriginFunc 获取应用代理后的原函数地址(和代理前的原函数地址不一样)
func (g *PatchGuard) OriginFunc() uintptr {
	return g.originFuncPtr
}
