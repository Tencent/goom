package patch

const defaultLength = 4

// GetFuncSize get func binary size
// not absolutely safe
func GetFuncSize(_ int, start uintptr, minimal bool) (length int, err error) {
	panic("not support arm64 yet")
}
