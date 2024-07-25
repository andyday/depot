package parse

import (
	"go/ast"
	"go/parser"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	e, err := parser.ParseExpr("(a == 10 || b == \"bar\") && c > h && strings.StartsWith(d, \"foo\")")
	assert.NoError(t, err)
	err = ast.Print(nil, e)
	e.Pos()
	assert.NoError(t, err)
	slices.Contains([]int{1, 2, 3}, 2)
}
