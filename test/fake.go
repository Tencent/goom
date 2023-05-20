package test

// Fake 导出结构体
type Fake struct{}

// Call 普通方法
//
//go:noinline
func (f *Fake) Call(i int) int {
	if i < -10000 {
		dummy()
	}
	return f.call(i)
}

// call 未导出方法
//
//go:noinline
func (f *Fake) call(i int) int {
	if i < -10000 {
		dummy()
	}
	return i
}

// Invokecall 测试调用未导出函数
//
//go:noinline
func (f *Fake) Invokecall(i int) int {
	return f.call(i)
}

// Call2 普通方法
//
//go:noinline
func (f *Fake) Call2(i int) int {
	if i < -10000 {
		dummy()
	}
	return i
}
