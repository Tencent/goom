package errortype

// TypeNotFound 类型没有找到
type TypeNotFound struct {
	typName string
}

func (t *TypeNotFound) Error() string {
	return "type not found:" + t.typName
}

// NewTypeNotFoundError 创建类型未找到异常
// typName 类型名称
func NewTypeNotFoundError(typName string) error {
	return &TypeNotFound{typName:typName}
}
