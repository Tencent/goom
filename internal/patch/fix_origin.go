package patch

import (
	"github.com/Jakegogo/goom_mocker/internal/logger"
)

// fixOrigin 将原函数拷贝到另外一个内存区段,并且修复
// trampoline 跳板函数地址, 不传递用0表示
// jumpDataLen jumpData 字节数组长度
func fixOrigin(origin, trampoline uintptr, jumpDataLen int) (uintptr, error) {
	logger.Infof("starting fix Origin origin=0x%x trampoline=0x%x ...", origin, trampoline)
	r, e := fixOriginFuncToTrampoline(origin, trampoline, jumpDataLen)
	if e != nil {
		logger.Errorf("fixed Origin error origin=%d trampoline=%d error:%s", origin, trampoline, e)
	}
	return r, e
}
