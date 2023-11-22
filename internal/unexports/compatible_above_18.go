//go:build go1.18
// +build go1.18

// Package unexports 实现了对未导出函数的获取
// 基于 github.com/alangpierce/go-forceexport 进行了修改和扩展。
package unexports

import "github.com/Jakegogo/goom_mocker/internal/hack"

// checkOverflow 检查 hack ftab 数据的是否正确, 一般溢出则不正确
func checkOverflow(ftab hack.Functab, moduleData *hack.Moduledata) bool {
	return ftab.Funcoff >= uint32(len(moduleData.Pclntable))
}
