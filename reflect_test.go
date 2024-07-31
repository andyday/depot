package depot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealSlice(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, RealSlice([]interface{}{"a", "b"}))
	assert.Equal(t, []int{1, 2}, RealSlice([]interface{}{1, 2}))
	assert.Equal(t, []int64{1, 2}, RealSlice([]interface{}{int64(1), int64(2)}))
}
