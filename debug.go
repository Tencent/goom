package mocker

import (
	"reflect"

	"git.woa.com/goom/mocker/arg"
	"git.woa.com/goom/mocker/internal/hack"
	"git.woa.com/goom/mocker/internal/iface"
	"git.woa.com/goom/mocker/internal/logger"
)

// excludeFunc 对 excludeFunc 不进行拦截
const (
	excludeFunc = "time.Now"
)

// interceptDebugInfo 添加对 apply 的拦截代理，截取函数调用信息用于 debug
func interceptDebugInfo(imp interface{}, pFunc iface.PFunc, mocker Mocker) (interface{}, iface.PFunc) {
	if !logger.IsDebugOpen() {
		return imp, pFunc
	}

	// 因为当使用了 when 时候,imp 代理会被覆盖,pFunc 会生效; 所以优先拦截有 pFunc 代理的 mock 回调
	if pFunc != nil {
		originPFunc := pFunc
		pFunc = func(params []reflect.Value) []reflect.Value {
			results := originPFunc(params)
			// 日志打印用到了 time.Now,避免递归死循环
			if mocker.String() == excludeFunc {
				return results
			}
			logger.Consolefc(logger.DebugLevel, "mocker [%s] called, args [%s], results [%s]",
				logger.Caller(hack.InterceptCallerSkip), mocker.String(), arg.SprintV(params), arg.SprintV(results))
			return results
		}
		return imp, pFunc
	}

	if imp != nil {
		originImp := imp
		impType := reflect.TypeOf(imp)
		imp = reflect.MakeFunc(impType, func(params []reflect.Value) []reflect.Value {
			var results []reflect.Value
			if impType.IsVariadic() {
				results = reflect.ValueOf(originImp).CallSlice(params)
			} else {
				results = reflect.ValueOf(originImp).Call(params)
			}
			// 日志打印用到了 time.Now,避免递归死循环
			if mocker.String() == excludeFunc {
				return results
			}
			logger.Consolefc(logger.DebugLevel, "mocker [%s] called, args [%s], results [%s]",
				logger.Caller(hack.InterceptCallerSkip), mocker.String(), arg.SprintV(params), arg.SprintV(results))
			return results
		}).Interface()
		return imp, pFunc
	}

	return imp, pFunc
}
