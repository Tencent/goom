package erro

import (
	"reflect"
	"strconv"
)

// ReturnsNotMatch 返回参数不匹配异常
type ReturnsNotMatch struct {
	funcDef   interface{}
	argLen    int
	expectLen int
}

// Error 返回错误字符串
func (i *ReturnsNotMatch) Error() string {
	if i.funcDef != nil {
		return "returns lenth not match of func " + reflect.ValueOf(i.funcDef).String() +
			": " + strconv.Itoa(i.argLen) + ", expect: " + strconv.Itoa(i.expectLen)
	}
	return "returns lenth not match: " + strconv.Itoa(i.argLen) + ", expect: " + strconv.Itoa(i.expectLen)
}

// NewReturnsNotMatchError 创建参数异常
// funcDef 函数定义
// argLen 参数长度
// expectLen 期望长度
func NewReturnsNotMatchError(funcDef interface{}, argLen int, expectLen int) error {
	return &ArgsNotMatch{funcDef: funcDef, argLen: argLen, expectLen: expectLen}
}
