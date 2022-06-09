// Package testdata 测试数据包, 包含被测函数层
package testdata

// Fake fake 结构体
type Fake struct{}

// Call call测试方法
//go:noinline
func (f *Fake) Call(i int) int {
	return f.call(i)
}

//go:noinline
func (f *Fake) call(i int) int {
	return i
}
