package depot

type numbers interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func AddValues(a, b interface{}) interface{} {
	switch v := b.(type) {
	case int8:
		return addValue(a, v)
	case int16:
		return addValue(a, v)
	case int32:
		return addValue(a, v)
	case int64:
		return addValue(a, v)
	case int:
		return addValue(a, v)
	case float32:
		return addValue(a, v)
	case float64:
		return addValue(a, v)
	case uint:
		return addValue(a, v)
	case uint8:
		return addValue(a, v)
	case uint16:
		return addValue(a, v)
	case uint32:
		return addValue(a, v)
	case uint64:
		return addValue(a, v)
	default:
		return v
	}
}

func addValue[T numbers](a interface{}, b T) T {
	switch v := a.(type) {
	case int8:
		return T(v) + b
	case int16:
		return T(v) + b
	case int32:
		return T(v) + b
	case int64:
		return T(v) + b
	case int:
		return T(v) + b
	case uint:
		return T(v) + b
	case uint8:
		return T(v) + b
	case uint16:
		return T(v) + b
	case uint32:
		return T(v) + b
	case uint64:
		return T(v) + b
	case float32:
		return T(v) + b
	case float64:
		return T(v) + b
	default:
		return b
	}
}

func SubtractValues(a, b interface{}) interface{} {
	switch v := b.(type) {
	case int8:
		return subtractValue(a, v)
	case int16:
		return subtractValue(a, v)
	case int32:
		return subtractValue(a, v)
	case int64:
		return subtractValue(a, v)
	case int:
		return subtractValue(a, v)
	case float32:
		return subtractValue(a, v)
	case float64:
		return subtractValue(a, v)
	case uint:
		return subtractValue(a, v)
	case uint8:
		return subtractValue(a, v)
	case uint16:
		return subtractValue(a, v)
	case uint32:
		return subtractValue(a, v)
	case uint64:
		return subtractValue(a, v)
	default:
		return v
	}
}

func subtractValue[T numbers](a interface{}, b T) T {
	switch v := a.(type) {
	case int8:
		return T(v) - b
	case int16:
		return T(v) - b
	case int32:
		return T(v) - b
	case int64:
		return T(v) - b
	case int:
		return T(v) - b
	case uint:
		return T(v) - b
	case uint8:
		return T(v) - b
	case uint16:
		return T(v) - b
	case uint32:
		return T(v) - b
	case uint64:
		return T(v) - b
	case float32:
		return T(v) - b
	case float64:
		return T(v) - b
	default:
		return -b
	}
}

func NegateValue(in interface{}) interface{} {
	switch v := in.(type) {
	case int8:
		return -v
	case int16:
		return -v
	case int32:
		return -v
	case int64:
		return -v
	case int:
		return -v
	case float32:
		return -v
	case float64:
		return -v
	default:
		return v
	}
}
