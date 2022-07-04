package memory

// WriteTo this function is super unsafe
// 因为M1芯片不支持MProtect同时拥有写Write和执行Exec两个权限, 因此只能设置读和写来绕过系统检查
// 但是当前函数所在内存区段也需要保证没有被修改权限， 否则设置非执行的权限之后，执行后续代码会抛异常
// 因此采用较为hack的方式, 在当前函数前后填充空函数(BeforeSpace、Space)来避免被修改到权限
// 注意: BeforeSpace、Space长度均需要大于pageSize, 且因go中同一个包的函数一般会连续编译到附近
func WriteTo(addr uintptr, data []byte) error {
	memoryAccessLock.Lock()
	defer memoryAccessLock.Unlock()
	if err := writeTo(addr, data); err != nil {
		return err
	}
	ClearICache(addr)
	return nil
}

// WriteToNoFlush 写入 .text 区, 不刷新 icache
func WriteToNoFlush(addr uintptr, data []byte) error {
	memoryAccessLock.Lock()
	defer memoryAccessLock.Unlock()
	return writeTo(addr, data)
}

// WriteToNoFlushNoLock 写入 .text 区, 不刷新 icache, 不加锁
func WriteToNoFlushNoLock(addr uintptr, data []byte) error {
	return writeTo(addr, data)
}
