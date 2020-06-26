package mocker

import (
	"reflect"
	"runtime"
	"strings"
)

// CurrentPackage 获取当前调用的包路径
func CurrentPackage() string {
	return currentPackage(2)
}

// currentPackage 获取调用者的包路径
func currentPackage(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	callerName := runtime.FuncForPC(pc).Name()
	if i := strings.Index(callerName, ".("); i > -1 {
		return callerName[:i]
	}
	if i := strings.LastIndex(callerName, "."); i > -1 {
		return callerName[:i]
	}
	return callerName
}

// getFunctionName 获取函数名称
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
