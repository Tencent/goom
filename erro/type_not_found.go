package erro

// TypeNotFound 类型没有找到
type TypeNotFound struct {
	typName string
}

// TypeNotFound 类型没有找到
func (t *TypeNotFound) Error() string {
	return "type not found: " + t.typName
}

// NewTypeNotFoundError 创建类型未找到异常
// typName 类型名称
func NewTypeNotFoundError(typName string) error {
	return &TypeNotFound{
		typName: typName,
	}
}
