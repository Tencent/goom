// Package mocker 定义了mock的外层用户使用API定义, 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件实现了mock守卫能力，
// 不同类型的Mocker具备不一样的Apply、Cancel具体行为，本MockGuard抽象了统一各类Mocker的Guard，
// 以供BaseMocker使用其接口类MockGuard的Apply、Cancel方法。
package mocker

import (
	"git.code.oa.com/goom/mocker/internal/patch"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// MockGuard Mock守卫
// 负责Mock应用和取消
type MockGuard interface {
	// Apply 应用Mock
	Apply()
	// Cancel 取消Mock
	Cancel()
}

// IFaceMockGuard 接口Mock守卫
type IFaceMockGuard struct {
	ctx *proxy.IContext
}

// NewIFaceMockGuard 创建IFaceMockGuard
func NewIFaceMockGuard(ctx *proxy.IContext) *IFaceMockGuard {
	return &IFaceMockGuard{ctx: ctx}
}

// Apply 应用mock
func (i *IFaceMockGuard) Apply() {
	// 无需操作
}

// Cancel 取消mock
func (i *IFaceMockGuard) Cancel() {
	i.ctx.Cancel()
}

// PatchMockGuard Patch类型的Mock守卫
type PatchMockGuard struct {
	patchGuard *patch.Guard
}

// NewPatchMockGuard 创建PatchMockGuard
func NewPatchMockGuard(patchGuard *patch.Guard) *PatchMockGuard {
	return &PatchMockGuard{patchGuard: patchGuard}
}

// Apply 应用mock
func (p *PatchMockGuard) Apply() {
	p.patchGuard.Apply()
}

// Cancel 取消mock
func (p *PatchMockGuard) Cancel() {
	p.patchGuard.UnpatchWithLock()
}
