//go:build !go1.18
// +build !go1.18

package unexports

import "git.woa.com/goom/mocker/internal/hack"

// checkOverflow 检查 hack ftab 数据的是否正确, 一般溢出则不正确
func checkOverflow(ftab hack.Functab, moduleData *hack.Moduledata) bool {
	return ftab.Funcoff >= uintptr(len(moduleData.Pclntable))
}
