// Package mocker 定义了 mock 的外层用户使用 API 定义, 包括函数、方法、接口、未导出函数(或方法的)的 Mocker 的实现。
// 当前文件实现了 mock 守卫能力，
// 不同类型的 Mocker 具备不一样的 Apply、Cancel 具体行为，本 MockGuard 抽象了统一各类 Mocker 的 Guard，
// 以供 BaseMocker 使用其接口类 MockGuard 的 Apply、Cancel 方法。
package mocker

import (
	"git.code.oa.com/goom/mocker/internal/iface"
	"git.code.oa.com/goom/mocker/internal/patch"
)

// MockGuard Mock 守卫
// 负责 Mock 应用和取消
type MockGuard interface {
	// Apply 应用 Mock
	Apply()
	// Cancel 取消 Mock
	Cancel()
}

// iFaceMockGuard 接口 Mock 守卫
type iFaceMockGuard struct {
	ctx *iface.IContext
}

// newIFaceMockGuard 创建 iFaceMockGuard
func newIFaceMockGuard(ctx *iface.IContext) *iFaceMockGuard {
	return &iFaceMockGuard{ctx: ctx}
}

// Apply 应用 mock
func (i *iFaceMockGuard) Apply() {
	// 无需操作
}

// Cancel 取消 mock
func (i *iFaceMockGuard) Cancel() {
	i.ctx.Cancel()
}

// patchMockGuard Patch 类型的 Mock 守卫
type patchMockGuard struct {
	patchGuard *patch.Guard
}

// newPatchMockGuard 创建 patchMockGuard
func newPatchMockGuard(patchGuard *patch.Guard) *patchMockGuard {
	return &patchMockGuard{patchGuard: patchGuard}
}

// Apply 应用 mock
func (p *patchMockGuard) Apply() {
	p.patchGuard.Apply()
}

// Cancel 取消 mock
func (p *patchMockGuard) Cancel() {
	p.patchGuard.UnpatchWithLock()
}
