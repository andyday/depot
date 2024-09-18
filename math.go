package depot

import "fmt"

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

func ValuesEqual(a, b interface{}) bool {
	switch v := b.(type) {
	case int8:
		return numbersEqual(a, v)
	case int16:
		return numbersEqual(a, v)
	case int32:
		return numbersEqual(a, v)
	case int64:
		return numbersEqual(a, v)
	case int:
		return numbersEqual(a, v)
	case float32:
		return numbersEqual(a, v)
	case float64:
		return numbersEqual(a, v)
	case uint:
		return numbersEqual(a, v)
	case uint8:
		return numbersEqual(a, v)
	case uint16:
		return numbersEqual(a, v)
	case uint32:
		return numbersEqual(a, v)
	case uint64:
		return numbersEqual(a, v)
	case string:
		return fmt.Sprintf("%v", a) == v
	default:
		return a == b
	}
}

func numbersEqual[T numbers](a interface{}, b T) bool {
	switch v := a.(type) {
	case int8:
		return T(v) == b
	case int16:
		return T(v) == b
	case int32:
		return T(v) == b
	case int64:
		return T(v) == b
	case int:
		return T(v) == b
	case uint:
		return T(v) == b
	case uint8:
		return T(v) == b
	case uint16:
		return T(v) == b
	case uint32:
		return T(v) == b
	case uint64:
		return T(v) == b
	case float32:
		return T(v) == b
	case float64:
		return T(v) == b
	case string:
		return fmt.Sprintf("%v", b) == v

	default:
		return false
	}
}

func ValuesNotEqual(a, b interface{}) bool {
	switch v := b.(type) {
	case int8:
		return numbersNotEqual(a, v)
	case int16:
		return numbersNotEqual(a, v)
	case int32:
		return numbersNotEqual(a, v)
	case int64:
		return numbersNotEqual(a, v)
	case int:
		return numbersNotEqual(a, v)
	case float32:
		return numbersNotEqual(a, v)
	case float64:
		return numbersNotEqual(a, v)
	case uint:
		return numbersNotEqual(a, v)
	case uint8:
		return numbersNotEqual(a, v)
	case uint16:
		return numbersNotEqual(a, v)
	case uint32:
		return numbersNotEqual(a, v)
	case uint64:
		return numbersNotEqual(a, v)
	case string:
		return fmt.Sprintf("%v", a) != v
	default:
		return a != b
	}
}

func numbersNotEqual[T numbers](a interface{}, b T) bool {
	switch v := a.(type) {
	case int8:
		return T(v) != b
	case int16:
		return T(v) != b
	case int32:
		return T(v) != b
	case int64:
		return T(v) != b
	case int:
		return T(v) != b
	case uint:
		return T(v) != b
	case uint8:
		return T(v) != b
	case uint16:
		return T(v) != b
	case uint32:
		return T(v) != b
	case uint64:
		return T(v) != b
	case float32:
		return T(v) != b
	case float64:
		return T(v) != b
	case string:
		return fmt.Sprintf("%v", b) != v
	default:
		return true
	}
}

func ValuesGreaterThan(a, b interface{}) bool {
	switch v := b.(type) {
	case int8:
		return numbersGreaterThan(a, v)
	case int16:
		return numbersGreaterThan(a, v)
	case int32:
		return numbersGreaterThan(a, v)
	case int64:
		return numbersGreaterThan(a, v)
	case int:
		return numbersGreaterThan(a, v)
	case float32:
		return numbersGreaterThan(a, v)
	case float64:
		return numbersGreaterThan(a, v)
	case uint:
		return numbersGreaterThan(a, v)
	case uint8:
		return numbersGreaterThan(a, v)
	case uint16:
		return numbersGreaterThan(a, v)
	case uint32:
		return numbersGreaterThan(a, v)
	case uint64:
		return numbersGreaterThan(a, v)
	case string:
		if a == nil {
			return false
		}
		return fmt.Sprintf("%v", a) > v
	default:
		return false
	}
}

func numbersGreaterThan[T numbers](a interface{}, b T) bool {
	switch v := a.(type) {
	case int8:
		return T(v) > b
	case int16:
		return T(v) > b
	case int32:
		return T(v) > b
	case int64:
		return T(v) > b
	case int:
		return T(v) > b
	case uint:
		return T(v) > b
	case uint8:
		return T(v) > b
	case uint16:
		return T(v) > b
	case uint32:
		return T(v) > b
	case uint64:
		return T(v) > b
	case float32:
		return T(v) > b
	case float64:
		return T(v) > b
	case string:
		return v > fmt.Sprintf("%v", b)
	default:
		return false
	}
}

func ValuesGreaterThanOrEqual(a, b interface{}) bool {
	switch v := b.(type) {
	case int8:
		return numbersGreaterThanOrEqual(a, v)
	case int16:
		return numbersGreaterThanOrEqual(a, v)
	case int32:
		return numbersGreaterThanOrEqual(a, v)
	case int64:
		return numbersGreaterThanOrEqual(a, v)
	case int:
		return numbersGreaterThanOrEqual(a, v)
	case float32:
		return numbersGreaterThanOrEqual(a, v)
	case float64:
		return numbersGreaterThanOrEqual(a, v)
	case uint:
		return numbersGreaterThanOrEqual(a, v)
	case uint8:
		return numbersGreaterThanOrEqual(a, v)
	case uint16:
		return numbersGreaterThanOrEqual(a, v)
	case uint32:
		return numbersGreaterThanOrEqual(a, v)
	case uint64:
		return numbersGreaterThanOrEqual(a, v)
	case string:
		if a == nil {
			return false
		}
		return fmt.Sprintf("%v", a) >= v
	default:
		return false
	}
}

func numbersGreaterThanOrEqual[T numbers](a interface{}, b T) bool {
	switch v := a.(type) {
	case int8:
		return T(v) >= b
	case int16:
		return T(v) >= b
	case int32:
		return T(v) >= b
	case int64:
		return T(v) >= b
	case int:
		return T(v) >= b
	case uint:
		return T(v) >= b
	case uint8:
		return T(v) >= b
	case uint16:
		return T(v) >= b
	case uint32:
		return T(v) >= b
	case uint64:
		return T(v) >= b
	case float32:
		return T(v) >= b
	case float64:
		return T(v) >= b
	case string:
		return v >= fmt.Sprintf("%v", b)
	default:
		return false
	}
}

func ValuesLessThan(a, b interface{}) bool {
	switch v := b.(type) {
	case int8:
		return numbersLessThan(a, v)
	case int16:
		return numbersLessThan(a, v)
	case int32:
		return numbersLessThan(a, v)
	case int64:
		return numbersLessThan(a, v)
	case int:
		return numbersLessThan(a, v)
	case float32:
		return numbersLessThan(a, v)
	case float64:
		return numbersLessThan(a, v)
	case uint:
		return numbersLessThan(a, v)
	case uint8:
		return numbersLessThan(a, v)
	case uint16:
		return numbersLessThan(a, v)
	case uint32:
		return numbersLessThan(a, v)
	case uint64:
		return numbersLessThan(a, v)
	case string:
		if a == nil {
			return false
		}
		return fmt.Sprintf("%v", a) < v
	default:
		return false
	}
}

func numbersLessThan[T numbers](a interface{}, b T) bool {
	switch v := a.(type) {
	case int8:
		return T(v) < b
	case int16:
		return T(v) < b
	case int32:
		return T(v) < b
	case int64:
		return T(v) < b
	case int:
		return T(v) < b
	case uint:
		return T(v) < b
	case uint8:
		return T(v) < b
	case uint16:
		return T(v) < b
	case uint32:
		return T(v) < b
	case uint64:
		return T(v) < b
	case float32:
		return T(v) < b
	case float64:
		return T(v) < b
	case string:
		return v < fmt.Sprintf("%v", b)
	default:
		return false
	}
}

func ValuesLessThanOrEqual(a, b interface{}) bool {
	switch v := b.(type) {
	case int8:
		return numbersLessThanOrEqual(a, v)
	case int16:
		return numbersLessThanOrEqual(a, v)
	case int32:
		return numbersLessThanOrEqual(a, v)
	case int64:
		return numbersLessThanOrEqual(a, v)
	case int:
		return numbersLessThanOrEqual(a, v)
	case float32:
		return numbersLessThanOrEqual(a, v)
	case float64:
		return numbersLessThanOrEqual(a, v)
	case uint:
		return numbersLessThanOrEqual(a, v)
	case uint8:
		return numbersLessThanOrEqual(a, v)
	case uint16:
		return numbersLessThanOrEqual(a, v)
	case uint32:
		return numbersLessThanOrEqual(a, v)
	case uint64:
		return numbersLessThanOrEqual(a, v)
	case string:
		if a == nil {
			return false
		}
		return fmt.Sprintf("%v", a) <= v
	default:
		return false
	}
}

func numbersLessThanOrEqual[T numbers](a interface{}, b T) bool {
	switch v := a.(type) {
	case int8:
		return T(v) <= b
	case int16:
		return T(v) <= b
	case int32:
		return T(v) <= b
	case int64:
		return T(v) <= b
	case int:
		return T(v) <= b
	case uint:
		return T(v) <= b
	case uint8:
		return T(v) <= b
	case uint16:
		return T(v) <= b
	case uint32:
		return T(v) <= b
	case uint64:
		return T(v) <= b
	case float32:
		return T(v) <= b
	case float64:
		return T(v) <= b
	case string:
		return v <= fmt.Sprintf("%v", b)
	default:
		return false
	}
}
