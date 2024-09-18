package depot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	values      = []any{123, int8(123), int16(123), int32(123), int64(123), uint(123), uint8(123), uint16(123), uint32(123), uint64(123), float32(123), float64(123), "123"}
	greaterThan = []any{124, int8(124), int16(124), int32(124), int64(124), uint(124), uint8(124), uint16(124), uint32(124), uint64(124), float32(124), float64(124), "124"}
	lessThan    = []any{122, int8(122), int16(122), int32(122), int64(122), uint(122), uint8(122), uint16(122), uint32(122), uint64(122), float32(122), float64(122), "122"}
)

func TestAddValues(t *testing.T) {
	assert.Equal(t, 5, AddValues(nil, 5))
	assert.Equal(t, int8(5), AddValues(nil, int8(5)))
	assert.Equal(t, int16(5), AddValues(nil, int16(5)))
	assert.Equal(t, int32(5), AddValues(nil, int32(5)))
	assert.Equal(t, int64(5), AddValues(nil, int64(5)))
	assert.Equal(t, float32(5), AddValues(nil, float32(5)))
	assert.Equal(t, float64(5), AddValues(nil, float64(5)))
	assert.Equal(t, uint(5), AddValues(nil, uint(5)))
	assert.Equal(t, uint8(5), AddValues(nil, uint8(5)))
	assert.Equal(t, uint16(5), AddValues(nil, uint16(5)))
	assert.Equal(t, uint32(5), AddValues(nil, uint32(5)))
	assert.Equal(t, uint64(5), AddValues(nil, uint64(5)))

	assert.Equal(t, 15, AddValues(10, 5))
	assert.Equal(t, int8(15), AddValues(int8(10), int8(5)))
	assert.Equal(t, int16(15), AddValues(int16(10), int16(5)))
	assert.Equal(t, int32(15), AddValues(int32(10), int32(5)))
	assert.Equal(t, int64(15), AddValues(int64(10), int64(5)))
	assert.Equal(t, float32(15), AddValues(float32(10), float32(5)))
	assert.Equal(t, float64(15), AddValues(float64(10), float64(5)))
	assert.Equal(t, uint(15), AddValues(uint(10), uint(5)))
	assert.Equal(t, uint8(15), AddValues(uint8(10), uint8(5)))
	assert.Equal(t, uint16(15), AddValues(uint16(10), uint16(5)))
	assert.Equal(t, uint32(15), AddValues(uint32(10), uint32(5)))
	assert.Equal(t, uint64(15), AddValues(uint64(10), uint64(5)))

	assert.Equal(t, int8(15), AddValues(uint64(10), int8(5)))
	assert.Equal(t, int16(15), AddValues(uint32(10), int16(5)))
	assert.Equal(t, int32(15), AddValues(uint16(10), int32(5)))
	assert.Equal(t, int64(15), AddValues(uint8(10), int64(5)))
	assert.Equal(t, float32(15), AddValues(uint(10), float32(5)))
	assert.Equal(t, float64(15), AddValues(float64(10), float64(5)))
	assert.Equal(t, uint(15), AddValues(float32(10), uint(5)))
	assert.Equal(t, uint8(15), AddValues(int64(10), uint8(5)))
	assert.Equal(t, uint16(15), AddValues(int32(10), uint16(5)))
	assert.Equal(t, uint32(15), AddValues(int16(10), uint32(5)))
	assert.Equal(t, uint64(15), AddValues(int8(10), uint64(5)))
}

func TestSubtractValues(t *testing.T) {
	assert.Equal(t, -5, SubtractValues(nil, 5))
	assert.Equal(t, int8(-5), SubtractValues(nil, int8(5)))
	assert.Equal(t, int16(-5), SubtractValues(nil, int16(5)))
	assert.Equal(t, int32(-5), SubtractValues(nil, int32(5)))
	assert.Equal(t, int64(-5), SubtractValues(nil, int64(5)))
	assert.Equal(t, float32(-5), SubtractValues(nil, float32(5)))
	assert.Equal(t, float64(-5), SubtractValues(nil, float64(5)))

	assert.Equal(t, 15, SubtractValues(20, 5))
	assert.Equal(t, int8(15), SubtractValues(int8(20), int8(5)))
	assert.Equal(t, int16(15), SubtractValues(int16(20), int16(5)))
	assert.Equal(t, int32(15), SubtractValues(int32(20), int32(5)))
	assert.Equal(t, int64(15), SubtractValues(int64(20), int64(5)))
	assert.Equal(t, float32(15), SubtractValues(float32(20), float32(5)))
	assert.Equal(t, float64(15), SubtractValues(float64(20), float64(5)))
	assert.Equal(t, uint(15), SubtractValues(uint(20), uint(5)))
	assert.Equal(t, uint8(15), SubtractValues(uint8(20), uint8(5)))
	assert.Equal(t, uint16(15), SubtractValues(uint16(20), uint16(5)))
	assert.Equal(t, uint32(15), SubtractValues(uint32(20), uint32(5)))
	assert.Equal(t, uint64(15), SubtractValues(uint64(20), uint64(5)))

	assert.Equal(t, int8(15), SubtractValues(uint64(20), int8(5)))
	assert.Equal(t, int16(15), SubtractValues(uint32(20), int16(5)))
	assert.Equal(t, int32(15), SubtractValues(uint16(20), int32(5)))
	assert.Equal(t, int64(15), SubtractValues(uint8(20), int64(5)))
	assert.Equal(t, float32(15), SubtractValues(uint(20), float32(5)))
	assert.Equal(t, float64(15), SubtractValues(float64(20), float64(5)))
	assert.Equal(t, uint(15), SubtractValues(float32(20), uint(5)))
	assert.Equal(t, uint8(15), SubtractValues(int64(20), uint8(5)))
	assert.Equal(t, uint16(15), SubtractValues(int32(20), uint16(5)))
	assert.Equal(t, uint32(15), SubtractValues(int16(20), uint32(5)))
	assert.Equal(t, uint64(15), SubtractValues(int8(20), uint64(5)))
}

func TestNegateValue(t *testing.T) {
	assert.Equal(t, -5, NegateValue(5))
	assert.Equal(t, int8(-5), NegateValue(int8(5)))
	assert.Equal(t, int16(-5), NegateValue(int16(5)))
	assert.Equal(t, int32(-5), NegateValue(int32(5)))
	assert.Equal(t, int64(-5), NegateValue(int64(5)))
}

func TestValuesEqual(t *testing.T) {
	var notEqual []any
	copy(notEqual, greaterThan)
	notEqual = append(notEqual, nil)
	for _, a := range values {
		for _, b := range values {
			assert.True(t, ValuesEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range values {
		for _, b := range notEqual {
			assert.False(t, ValuesEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
			assert.False(t, ValuesEqual(b, a), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(b), b, reflect.TypeOf(a), a))
		}
	}
}

func TestValuesNotEqual(t *testing.T) {
	var notEqual []any
	copy(notEqual, greaterThan)
	notEqual = append(notEqual, nil)
	for _, a := range values {
		for _, b := range values {
			assert.False(t, ValuesNotEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}
	for _, a := range values {
		for _, b := range notEqual {
			assert.True(t, ValuesNotEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
			assert.True(t, ValuesNotEqual(b, a), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(b), b, reflect.TypeOf(a), a))
		}
	}
}

func TestValuesGreaterThan(t *testing.T) {
	var notGreaterThan []any
	copy(notGreaterThan, values)
	notGreaterThan = append(notGreaterThan, nil)
	for _, a := range values {
		for _, b := range lessThan {
			assert.True(t, ValuesGreaterThan(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range values {
		for _, b := range notGreaterThan {
			assert.False(t, ValuesGreaterThan(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
			assert.False(t, ValuesGreaterThan(b, a), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(b), b, reflect.TypeOf(a), a))
		}
	}
}

func TestValuesGreaterThanOrEqual(t *testing.T) {
	var valuesPlus, lessThanPlus []any
	copy(valuesPlus, values)
	copy(lessThanPlus, lessThan)
	valuesPlus = append(valuesPlus, nil)
	lessThanPlus = append(lessThanPlus, nil)
	for _, a := range values {
		for _, b := range lessThan {
			assert.True(t, ValuesGreaterThanOrEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range values {
		for _, b := range values {
			assert.True(t, ValuesGreaterThanOrEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range lessThanPlus {
		for _, b := range valuesPlus {
			assert.False(t, ValuesGreaterThanOrEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}
}

func TestValuesLessThan(t *testing.T) {
	var notLessThan []any
	copy(notLessThan, values)
	notLessThan = append(notLessThan, nil)
	for _, a := range values {
		for _, b := range greaterThan {
			assert.True(t, ValuesLessThan(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range values {
		for _, b := range notLessThan {
			assert.False(t, ValuesLessThan(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
			assert.False(t, ValuesLessThan(b, a), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(b), b, reflect.TypeOf(a), a))
		}
	}
}

func TestValuesLessThanOrEqual(t *testing.T) {
	var valuesPlus, greaterThanPlus []any
	copy(valuesPlus, values)
	copy(greaterThanPlus, lessThan)
	valuesPlus = append(valuesPlus, nil)
	greaterThanPlus = append(greaterThanPlus, nil)
	for _, a := range values {
		for _, b := range greaterThan {
			assert.True(t, ValuesLessThanOrEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range values {
		for _, b := range values {
			assert.True(t, ValuesLessThanOrEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}

	for _, a := range greaterThanPlus {
		for _, b := range valuesPlus {
			assert.False(t, ValuesLessThanOrEqual(a, b), fmt.Sprintf("%v(%v) vs %v(%v)", reflect.TypeOf(a), a, reflect.TypeOf(b), b))
		}
	}
}
