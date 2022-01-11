package patch

import (
	"git.code.oa.com/goom/mocker/internal/logger"
)

func init() {
	logger.Log2Consolef(logger.ErrorLevel, "not support arm cpu yet! "+
		"please use go-amd64(open with rosetta) instead on MACOS.")
}
