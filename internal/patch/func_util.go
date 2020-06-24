package patch

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/unexports"
	"git.code.oa.com/goom/mocker/internal/x86asm"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)


func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

var funcSizeCache = make(map[uintptr]int)
var funcSizeGetLock sync.Mutex

// get func binary size
// not absolutly safe
func GetFuncSize(mode int, start uintptr, minimal bool) (lenth int, err error) {
	funcSizeGetLock.Lock()
	defer func() {
		funcSizeCache[start] = lenth
		funcSizeGetLock.Unlock()
	}()

	if lenth, ok := funcSizeCache[start]; ok {
		return lenth, nil
	}

	prologueLen := len(funcPrologue)
	code := rawMemoryRead(start, 16) // instruction takes at most 16 bytes

	/* prologue is not required
	if !bytes.Equal(funcPrologue, code[:prologueLen]) { // not valid function start or invalid prologue
		return 0, errors.New(fmt.Sprintf("no func prologue, addr:0x%x", start))
	}
	*/

	int3_found := false
	curLen := 0

	for {
		inst, err := x86asm.Decode(code, mode)
		if err != nil || (inst.Opcode == 0 && inst.Len == 1 && inst.Prefix[0] == x86asm.Prefix(code[0])) {
			return curLen, nil
		}

		if inst.Len == 1 && code[0] == 0xcc {
			// 0xcc -> int3, trap to debugger, padding to function end
			if minimal {
				return curLen, nil
			}
			int3_found = true
		} else if int3_found {
			return curLen, nil
		}

		curLen = curLen + inst.Len
		code = rawMemoryRead(start+uintptr(curLen), 16) // instruction takes at most 16 bytes

		if bytes.Equal(funcPrologue, code[:prologueLen]) {
			return curLen, nil
		}
	}

	return 0, nil
}


type value struct {
	_   uintptr
	ptr unsafe.Pointer
}


// getPtr 获取函数的调用地址(和函数的指令地址不一样)
func getPtr(v reflect.Value) unsafe.Pointer {
	return (*value)(unsafe.Pointer(&v)).ptr
}

// isNil 判断interface{}是否为空
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).Elem().IsNil()
	}
	return reflect.ValueOf(i).IsNil()
}

// getTrampolinePtr 获取跳板函数的地址
func getTrampolinePtr(trampoline interface{}) (uintptr, error) {
	trampolinePtr := uintptr(0)
	if !isNil(trampoline) {
		trampolineType := reflect.TypeOf(trampoline)
		if trampolineType.Kind() == reflect.Ptr {
			trampolinePtr = reflect.ValueOf(trampoline).Elem().Pointer()
			//return  0, errors.New("param trampoline must be a pointer")
		} else if trampolineType.Kind() == reflect.Func {
			trampolinePtr = reflect.ValueOf(trampoline).Pointer()
		}

		firsPtr := unsafe.Pointer(&trampoline)
		secondPtr := ((*uintptr)(unsafe.Pointer(uintptr(firsPtr))))
		thirdPtr := ((*uintptr)(unsafe.Pointer(*secondPtr)))

		logger.LogDebugf("trampoline caller: 0x%x 0x%x 0x%x", uintptr(firsPtr), *secondPtr, *thirdPtr)
		logger.LogDebugf("trampoline value: 0x%x 0x%x", getPtr(reflect.ValueOf(trampoline)), trampolinePtr)
	}

	return trampolinePtr, nil
}

// IsPtr 判断interface{}是否为指针类型
func IsPtr(value interface{}) bool {
	if value == nil {
		return false
	}
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		return true
	}
	return false
}

func LoadUnit(s int64) string {
	suffix := ""
	b := s
	if s > (1 << 40) {
		suffix = "G"
		b = s / (1 << 30)
	} else if s > (1 << 30) {
		suffix = "M"
		b = s / (1 << 20)
	} else if s > (1 << 20) {
		suffix = "K"
		b = s / (1 << 10)
	}
	return fmt.Sprintf("%d%s", b, suffix)
}


func ShowInst(name string, from uintptr, size int, level int) {
	_, funcName, _ := unexports.FindFuncByPtr(from)
	instBytes := rawMemoryRead(from, size)
	showInst(fmt.Sprintf("show [%s = %s] inst>>: ", name, funcName), from, instBytes, level)
}


func minSize(showSize int, fixOrigin []byte) int {
	if showSize > len(fixOrigin) {
		showSize = len(fixOrigin)
	}
	return showSize
}

func showInst(title string, from uintptr, copyOrigin []byte, level int) {
	if logger.LogLevel < level {
		return
	}

	logger.LogImportant(title)
	startAddr := (uint64)(from)
	for pos := 0; pos < len(copyOrigin); {
		// read 16 bytes atmost each time
		endPos := pos + 16
		if endPos > len(copyOrigin) {
			endPos = len(copyOrigin)
		}

		code := copyOrigin[pos:endPos]

		ins, err := x86asm.Decode(code, 64)

		if err != nil {
			logger.LogImportantf("[0] 0x%x: inst decode error:%s", startAddr+(uint64)(pos) , err)
			if ins.Len == 0 {
				pos = pos + 1
			} else {
				pos = pos + ins.Len
			}
			continue
		}
		if ins.Opcode == 0 {
			if ins.Len == 0 {
				pos = pos + 1
			} else {
				pos = pos + ins.Len
			}
			continue
		}

		if ins.PCRelOff > 0 {

			var isAdd = true
			for i := 0; i < len(ins.Args); i++ {
				arg := ins.Args[i]
				if arg == nil {
					break
				}
				addrArgs := arg.String()
				if strings.HasPrefix(addrArgs, ".-") || strings.Contains(addrArgs, "RIP-") {
					isAdd = false
				}
			}

			offset := pos + ins.PCRelOff
			relativeAddr := decodeAddress(copyOrigin[offset:offset+ins.PCRel], ins.PCRel)
			if !isAdd && relativeAddr > 0 {
				relativeAddr = -relativeAddr
			}

			logger.LogImportantf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s\t\tabs:0x%x", ins.Len,
				startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(code[:ins.Len]),
				from + uintptr(pos) + uintptr(relativeAddr) + uintptr(ins.Len))

		} else {
			logger.LogImportantf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s", ins.Len,
				startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(code[:ins.Len]))
		}

		pos = pos + ins.Len
	}
}
