package depot

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestMaxInt64(t *testing.T) {
	// 9223372036854775807

	digits := 0
	for i := math.MaxInt64; i > 0; digits++ {
		fmt.Printf("%019d\n", i)
		i = i / 10
	}
	assert.Equal(t, 19, digits)

	digits = 0
	for i := math.MaxInt64; i > 0; digits++ {
		i = i / 16
	}
	assert.Equal(t, 16, digits)

	bi := big.NewInt(math.MaxInt64)
	assert.Equal(t, "foo", base64.RawURLEncoding.EncodeToString(bi.Bytes()))
}
