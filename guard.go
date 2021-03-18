// Package mocker定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件实现了mock守卫能力。
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

// IfaceMockGuard 接口Mock守卫
type IfaceMockGuard struct {
	ctx *proxy.IContext
}

// NewIfaceMockGuard 创建IfaceMockGuard
func NewIfaceMockGuard(ctx *proxy.IContext) *IfaceMockGuard {
	return &IfaceMockGuard{ctx: ctx}
}

//noLint
func (i *IfaceMockGuard) Apply() {
	// do nothing
}

// Cancel() Cancel()
func (i *IfaceMockGuard) Cancel() {
	i.ctx.Cancel()
}

// PatchMockGuard Patch类型的Mock守卫
type PatchMockGuard struct {
	patchGuard *patch.PatchGuard
}

// NewPatchMockGuard 创建PatchMockGuard
func NewPatchMockGuard(patchGuard *patch.PatchGuard) *PatchMockGuard {
	return &PatchMockGuard{patchGuard: patchGuard}
}

//noLint
func (p *PatchMockGuard) Apply() {
	p.patchGuard.Apply()
}

// Cancel() Cancel()
func (p *PatchMockGuard) Cancel() {
	p.patchGuard.UnpatchWithLock()
}
