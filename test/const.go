package test

// unexportedGlobalIntVar 用于测试全局常量 mock
var (
	unexportedGlobalIntConst = 1
	unexportedGlobalStrConst = "str"
	unexportedGlobalMapConst = map[string]int{"key": 1}
	unexportedGlobalArrConst = []int{1, 2, 3}
)

// UnexportedGlobalIntConst 获取未导出Int全局常量
func UnexportedGlobalIntConst() int {
	return unexportedGlobalIntConst
}

// UnexportedGlobalStrConst 获取未导出Str全局常量
func UnexportedGlobalStrConst() string {
	return unexportedGlobalStrConst
}

// UnexportedGlobalMapConst 获取未导出map全局常量
func UnexportedGlobalMapConst() map[string]int {
	return unexportedGlobalMapConst
}

// UnexportedGlobalArrConst 获取未导出数组全局常量
func UnexportedGlobalArrConst() []int {
	return unexportedGlobalArrConst
}
