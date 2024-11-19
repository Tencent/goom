package unexports2

var (
	pageSize           int
	pageBeginMask      uintptr
	processBaseAddress uintptr
	processCopy        []byte
)

func initProcess() {
	osInitProcess()
	pageSize = osGetPageSize()
	pageBeginMask = ^uintptr(pageSize - 1)
	//processBaseAddress = osGetProcessBaseAddress()
}
