// Package proxy 封装了给各种类型的代理(或较 patch)中间层
// 负责比如外部传如类型校验、私有函数名转换成 uintptr、trampoline 初始化、并发 proxy 等
package proxy

import (
	"errors"
	"fmt"
	"reflect"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/patch"
	"git.code.oa.com/goom/mocker/internal/unexports"
)

// StaticProxyByName 静态代理(函数或方法)
// @param genCallableFunc 函数名称
// @param proxyFunc 代理函数实现
// @param trampolineFunc 跳板函数,即代理后的原始函数定义;跳板函数的签名必须和原函数一致,值不能为空
func StaticProxyByName(funcName string, proxyFunc interface{}, trampolineFunc interface{}) (*patch.Guard, error) {
	e := checkTrampolineFunc(trampolineFunc)
	if e != nil {
		return nil, e
	}

	originFuncPtr, err := unexports.FindFuncByName(funcName)
	if err != nil {
		return nil, err
	}

	logger.LogInfo("start StaticProxyByName genCallableFunc=", funcName)

	// gomonkey 添加函数 hook
	patchGuard, err := patch.PtrTrampoline(originFuncPtr, proxyFunc, trampolineFunc)
	if err != nil {
		logger.LogError("StaticProxyByName fail genCallableFunc=", funcName, ":", err)
		return nil, err
	}

	// 构造原先方法实例值
	logger.LogDebug("OriginUintptr is:", fmt.Sprintf("0x%x", patchGuard.OriginFunc()))
	logger.LogInfo("static proxy[trampoline] ok, genCallableFunc=", funcName)

	return patchGuard, nil
}

// StaticProxyByFunc 静态代理(函数或方法)
// @param funcDef 原始函数定义
// @param proxyFunc 代理函数实现
// @param originFunc 跳板函数即代理后的原始函数定义(值为 nil 时,使用公共的跳板函数, 不为 nil 时使用指定的跳板函数)
func StaticProxyByFunc(funcDef interface{}, proxyFunc, trampolineFunc interface{}) (*patch.Guard, error) {
	e := checkTrampolineFunc(trampolineFunc)
	if e != nil {
		return nil, e
	}

	logger.LogInfo("start StaticProxyByFunc funcDef=", funcDef)

	// gomonkey 添加函数 hook
	patchGuard, err := patch.Trampoline(
		reflect.Indirect(reflect.ValueOf(funcDef)).Interface(), proxyFunc, trampolineFunc)
	if err != nil {
		logger.LogError("StaticProxyByFunc fail funcDef=", funcDef, ":", err)
		return nil, err
	}
	// 构造原先方法实例值
	logger.LogDebug("OriginUintptr is:", fmt.Sprintf("0x%x", patchGuard.OriginFunc()))

	if patch.IsPtr(trampolineFunc) {
		_, err = unexports.CreateFuncForCodePtr(trampolineFunc, patchGuard.OriginFunc())
		if err != nil {
			logger.LogError("StaticProxyByFunc fail funcDef=", funcDef, ":", err)
			patchGuard.Unpatch()

			return nil, err
		}
	}

	logger.LogDebug("static proxy ok funcDef=", funcDef)

	return patchGuard, nil
}

// StaticProxyByMethod 方法静态代理
// @param target 类型
// @param methodName 方法名
// @param proxyFunc 代理函数实现
// @param trampolineFunc 跳板函数即代理后的原始方法定义(值为 nil 时,使用公共的跳板函数, 不为 nil 时使用指定的跳板函数)
func StaticProxyByMethod(target reflect.Type, methodName string, proxyFunc,
	trampolineFunc interface{}) (*patch.Guard, error) {
	e := checkTrampolineFunc(trampolineFunc)
	if e != nil {
		return nil, e
	}

	logger.LogInfo("start StaticProxyByMethod genCallableFunc=", target, ".", methodName)

	// gomonkey 添加函数 hook
	patchGuard, err := patch.InstanceMethodTrampoline(target, methodName, proxyFunc, trampolineFunc)
	if err != nil {
		logger.LogError("StaticProxyByMethod fail type=", target, "methodName=", methodName, ":", err)
		return nil, err
	}

	// 构造原先方法实例值
	logger.LogDebug("OriginUintptr is:", fmt.Sprintf("0x%x", patchGuard.OriginFunc()))

	if patch.IsPtr(trampolineFunc) {
		_, err = unexports.CreateFuncForCodePtr(trampolineFunc, patchGuard.OriginFunc())
		if err != nil {
			logger.LogError("StaticProxyByMethod fail method=", target, ".", methodName, ":", err)
			patchGuard.Unpatch()

			return nil, err
		}
	}

	logger.LogDebug("static proxy ok genCallableFunc=", target, ".", methodName)

	return patchGuard, nil
}

// checkTrampolineFunc 检测 TrampolineFunc 类型
func checkTrampolineFunc(trampolineFunc interface{}) error {
	if trampolineFunc != nil {
		if reflect.ValueOf(trampolineFunc).Kind() != reflect.Func &&
			reflect.ValueOf(trampolineFunc).Elem().Kind() != reflect.Func {
			return errors.New("trampolineFunc has to be a exported func")
		}
	}
	return nil
}
