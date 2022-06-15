package patch

// fixOriginFuncToTrampoline 修复函数偏移量
func fixOriginFuncToTrampoline(_ uintptr, _ uintptr, _ int) (uintptr, error) {
	panic("not support yet on M1-MAC or arm CPU!")
}
