package depot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateOp(t *testing.T) {
	updateOps := []UpdateOp{Add("a"), Subtract("a"), Force("a")}
	for _, o := range updateOps {
		o.isUpdateOp()
		assert.Equal(t, "a", o.Field())
	}

	conditions := []Condition{
		Equal("a"), NotEqual("a"), LessThan("a"), LessThanOrEqual("a"),
		GreaterThan("a"), GreaterThanOrEqual("a"), Exists("a"),
		In("a", 1), NotIn("a", 1),
	}

	for _, o := range conditions {
		o.isCondition()
		u, ok := o.(UpdateOp)
		u.isUpdateOp()
		assert.True(t, ok)
		q, ok := o.(QueryOp)
		q.isQueryOp()
		assert.True(t, ok)
		assert.Equal(t, "a", o.Field())
		switch o.(type) {
		case *ExistsCondition, *InCondition, *NotInCondition:
			assert.True(t, o.Valueless())
		default:
			assert.False(t, o.Valueless())
		}
	}

	directives := []QueryDirective{Asc(), Desc(), Limit(10), Page("p")}
	for _, o := range directives {
		o.isQueryDirective()
		q, ok := o.(QueryOp)
		q.isQueryOp()
		assert.True(t, ok)
		switch v := o.(type) {
		case *LimitQueryDirective:
			assert.Equal(t, 10, v.Limit)
		case *PageQueryDirective:
			assert.Equal(t, "p", v.Page)
		}
	}
}
