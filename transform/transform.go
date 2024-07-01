package transform

type Type int8

const (
	TypeNone Type = iota
	TypeAdd
	TypeSubtract
)

type Transform struct {
	Type  Type
	Value interface{}
}

func Add(value interface{}) *Transform {
	return &Transform{Type: TypeAdd, Value: value}
}

func Subtract(value interface{}) *Transform {
	return &Transform{Type: TypeSubtract, Value: value}
}
