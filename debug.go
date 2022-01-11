package mocker

import (
	"reflect"

	"git.code.oa.com/goom/mocker/arg"
	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// excludeFunc 对excludeFunc不进行拦截
const excludeFunc = "time.Now"

// interceptDebugInfo 添加对apply的拦截代理，截取函数调用信息用于debug
func interceptDebugInfo(imp interface{}, pFunc proxy.PFunc, mocker Mocker) (interface{}, proxy.PFunc) {
	if !logger.IsDebugOpen() {
		return imp, nil
	}

	// 因为当使用了when时候,imp代理会被覆盖,pFunc会生效; 所以优先拦截有pFunc代理的mock回调
	if pFunc != nil {
		originPFunc := pFunc
		pFunc = func(args []reflect.Value) []reflect.Value {
			results := originPFunc(args)
			// 日志打印用到了time.Now,避免递归死循环
			if mocker.String() == excludeFunc {
				return results
			}
			logger.Log2Consolef(logger.DebugLevel, "mocker [%s] called, args [%s], results [%s]",
				mocker.String(), arg.SprintV(args), arg.SprintV(results))
			return results
		}
		return imp, pFunc
	}

	if imp != nil {
		originImp := imp
		imp = reflect.MakeFunc(reflect.TypeOf(imp), func(args []reflect.Value) []reflect.Value {
			results := reflect.ValueOf(originImp).Call(args)
			// 日志打印用到了time.Now,避免递归死循环
			if mocker.String() == excludeFunc {
				return results
			}

			logger.Log2Consolef(logger.DebugLevel, "mocker [%s] called, args [%s], results [%s]",
				mocker.String(), arg.SprintV(args), arg.SprintV(results))
			return results
		}).Interface()
		return imp, pFunc
	}

	return nil, nil
}
