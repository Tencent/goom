package test

// unexportedGlobalIntVar 用于测试全局变量 mock
var (
	unexportedGlobalIntVar           = 1
	unexportedGlobalStrVar           = "str"
	unexportedGlobalMapVar           = map[string]int{"key": 1}
	unexportedGlobalArrVar           = []int{1, 2, 3}
	unexportedGlobalStructVar        = Struct{Field1: "1"}
	unexportedGlobalStructPointerVar = &Struct{Field1: "p1"}
)

// UnexportedGlobalIntVar 获取未导出Int全局变量
func UnexportedGlobalIntVar() int {
	return unexportedGlobalIntVar
}

// UnexportedGlobalStrVar 获取未导出Str全局变量
func UnexportedGlobalStrVar() string {
	return unexportedGlobalStrVar
}

// UnexportedGlobalMapVar 获取未导出map全局变量
func UnexportedGlobalMapVar() map[string]int {
	return unexportedGlobalMapVar
}

// UnexportedGlobalArrVar 获取未导出数组全局变量
func UnexportedGlobalArrVar() []int {
	return unexportedGlobalArrVar
}

// UnexportedGlobalStructVar 获取未导出结构体全局变量
func UnexportedGlobalStructVar() Struct {
	return unexportedGlobalStructVar
}

// UnexportedGlobalStructPointerVar 获取未导出结构体引用全局变量
func UnexportedGlobalStructPointerVar() *Struct {
	return unexportedGlobalStructPointerVar
}

type Struct struct {
	Field1 string
}
