package errortype

// FieldNotFound 类型没有找到
type FieldNotFound struct {
	typName   string
	fieldName string
}

func (t *FieldNotFound) Error() string {
	return "field not found:" + t.typName + "." + t.fieldName
}

// NewFieldNotFoundError 创建类型未找到异常
// typName 类型名称
// fieldName 属性名称
func NewFieldNotFoundError(typName string, fieldName string) error {
	return &FieldNotFound{typName: typName, fieldName: fieldName}
}
