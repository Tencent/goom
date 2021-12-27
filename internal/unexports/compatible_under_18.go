// +build !go1.18

package unexports

import "git.code.oa.com/goom/mocker/internal/hack"

func checkOverflow(ftab hack.Functab, moduleData *hack.Moduledata) bool {
	return ftab.Funcoff >= uintptr(len(moduleData.Pclntable))
}
