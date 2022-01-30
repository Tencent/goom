package mocker

import (
	"reflect"

	"git.code.oa.com/goom/mocker/arg"
	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// excludeFunc 对 excludeFunc 不进行拦截
const (
	excludeFunc = "time.Now"
)

// interceptDebugInfo 添加对 apply 的拦截代理，截取函数调用信息用于 debug
func interceptDebugInfo(imp interface{}, pFunc proxy.PFunc, mocker Mocker) (interface{}, proxy.PFunc) {
	if !logger.IsDebugOpen() {
		return imp, pFunc
	}

	// 因为当使用了 when 时候,imp 代理会被覆盖,pFunc 会生效; 所以优先拦截有 pFunc 代理的 mock 回调
	if pFunc != nil {
		originPFunc := pFunc
		pFunc = func(args []reflect.Value) []reflect.Value {
			results := originPFunc(args)
			// 日志打印用到了 time.Now,避免递归死循环
			if mocker.String() == excludeFunc {
				return results
			}
			logger.Log2Consolefc(logger.DebugLevel, "mocker [%s] called, args [%s], results [%s]",
				logger.Caller(hack.InterceptCallerSkip), mocker.String(), arg.SprintV(args), arg.SprintV(results))
			return results
		}
		return imp, pFunc
	}

	if imp != nil {
		originImp := imp
		imp = reflect.MakeFunc(reflect.TypeOf(imp), func(args []reflect.Value) []reflect.Value {
			results := reflect.ValueOf(originImp).Call(args)
			// 日志打印用到了 time.Now,避免递归死循环
			if mocker.String() == excludeFunc {
				return results
			}

			logger.Log2Consolefc(logger.DebugLevel, "mocker [%s] called, args [%s], results [%s]",
				logger.Caller(hack.InterceptCallerSkip), mocker.String(), arg.SprintV(args), arg.SprintV(results))
			return results
		}).Interface()
		return imp, pFunc
	}

	return imp, pFunc
}
