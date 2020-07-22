package mocker

import "reflect"

// If 条件表达式 TODO
type If struct {
	ExportedMocker

	// funcTyp reflect.Type
}

// Arg 获取第0个参数
func (i *If) Arg() *If {
	return i
}

// Args 获取第j个参数
func (i *If) Args(j int) *If {
	return i
}

// Gt 大于某个值
func (i *If) Gt(j int) *If {
	return i
}

// Lt 小于某个值
func (i *If) Lt(j int) *If {
	return i
}

// Ge 大于等于某个值
func (i *If) Ge(j int) *If {
	return i
}

// Le 小于等于某个值
func (i *If) Le(j int) *If {
	return i
}

// And 与表达式
func (i *If) And() *If {
	return i
}

// Or 或表达式
func (i *If) Or() *If {
	return i
}

// Between 在一个范围内[j(包含),k(不包含)]
func (i *If) Between(j, k interface{}) *If {
	return i
}

// NotEqual 不等于
func (i *If) NotEqual(v interface{}) *If {
	return i
}

// Arg0 获取第0个参数
func (i *If) Arg0() *If {
	return i
}

// Arg1 获取第1个参数
func (i *If) Arg1() *If {
	return i
}

// Arg2 获取第2个参数
func (i *If) Arg2() *If {
	return i
}

// Arg3 获取第3个参数
func (i *If) Arg3() *If {
	return i
}

// Arg4 获取第4个参数
func (i *If) Arg4() *If {
	return i
}

// Arg5 获取第5个参数
func (i *If) Arg5() *If {
	return i
}

func (i *If) invoke(args1 []reflect.Value) (results []reflect.Value) {
	return nil
}
