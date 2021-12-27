//go:build go1.18
// +build go1.18

// Package unexports 实现了对未导出函数的获取
// 基于github.com/alangpierce/go-forceexport进行了修改和扩展。
package unexports

import "git.code.oa.com/goom/mocker/internal/hack"

func checkOverflow(ftab hack.Functab, moduleData *hack.Moduledata) bool {
	return ftab.Funcoff >= uint32(len(moduleData.Pclntable))
}
