package depot

type UpdateOp interface {
	isUpdateOp()
	Field() string
}

type AddUpdateOp struct{ field string }
type SubtractUpdateOp struct{ field string }
type ForceUpdateOp struct{ field string }

func (*AddUpdateOp) isUpdateOp()      {}
func (*SubtractUpdateOp) isUpdateOp() {}
func (*ForceUpdateOp) isUpdateOp()    {}

func (o *AddUpdateOp) Field() string      { return o.field }
func (o *SubtractUpdateOp) Field() string { return o.field }
func (o *ForceUpdateOp) Field() string    { return o.field }

func Add(field string) *AddUpdateOp           { return &AddUpdateOp{field: field} }
func Subtract(field string) *SubtractUpdateOp { return &SubtractUpdateOp{field: field} }
func Force(field string) *ForceUpdateOp       { return &ForceUpdateOp{field: field} }

type QueryOp interface{ isQueryOp() }
type Condition interface {
	isCondition()
	Field() string
	Valueless() bool
}
type QueryDirective interface{ isQueryDirective() }

type EqualCondition struct{ field string }
type NotEqualCondition struct{ field string }
type LTCondition struct{ field string }
type LTECondition struct{ field string }
type GTCondition struct{ field string }
type GTECondition struct{ field string }
type ExistsCondition struct{ field string }
type InCondition struct {
	field string
	list  []interface{}
}
type NotInCondition struct {
	field string
	list  []interface{}
}

func (*EqualCondition) isQueryOp()    {}
func (*NotEqualCondition) isQueryOp() {}
func (*LTCondition) isQueryOp()       {}
func (*LTECondition) isQueryOp()      {}
func (*GTCondition) isQueryOp()       {}
func (*GTECondition) isQueryOp()      {}
func (*ExistsCondition) isQueryOp()   {}
func (*InCondition) isQueryOp()       {}
func (*NotInCondition) isQueryOp()    {}

func (*EqualCondition) isUpdateOp()    {}
func (*NotEqualCondition) isUpdateOp() {}
func (*LTCondition) isUpdateOp()       {}
func (*LTECondition) isUpdateOp()      {}
func (*GTCondition) isUpdateOp()       {}
func (*GTECondition) isUpdateOp()      {}
func (*ExistsCondition) isUpdateOp()   {}
func (*InCondition) isUpdateOp()       {}
func (*NotInCondition) isUpdateOp()    {}

func (q *EqualCondition) Field() string    { return q.field }
func (q *NotEqualCondition) Field() string { return q.field }
func (q *LTCondition) Field() string       { return q.field }
func (q *LTECondition) Field() string      { return q.field }
func (q *GTCondition) Field() string       { return q.field }
func (q *GTECondition) Field() string      { return q.field }
func (q *ExistsCondition) Field() string   { return q.field }
func (q *InCondition) Field() string       { return q.field }
func (q *NotInCondition) Field() string    { return q.field }

func (*EqualCondition) Valueless() bool    { return false }
func (*NotEqualCondition) Valueless() bool { return false }
func (*LTCondition) Valueless() bool       { return false }
func (*LTECondition) Valueless() bool      { return false }
func (*GTCondition) Valueless() bool       { return false }
func (*GTECondition) Valueless() bool      { return false }
func (*ExistsCondition) Valueless() bool   { return true }
func (*InCondition) Valueless() bool       { return true }
func (*NotInCondition) Valueless() bool    { return true }

func (*EqualCondition) isCondition()    {}
func (*NotEqualCondition) isCondition() {}
func (*LTCondition) isCondition()       {}
func (*LTECondition) isCondition()      {}
func (*GTCondition) isCondition()       {}
func (*GTECondition) isCondition()      {}
func (*ExistsCondition) isCondition()   {}
func (*InCondition) isCondition()       {}
func (*NotInCondition) isCondition()    {}

func Equal(field string) *EqualCondition            { return &EqualCondition{field: field} }
func NotEqual(field string) *NotEqualCondition      { return &NotEqualCondition{field: field} }
func LessThan(field string) *LTCondition            { return &LTCondition{field: field} }
func LessThanOrEqual(field string) *LTECondition    { return &LTECondition{field: field} }
func GreaterThan(field string) *GTCondition         { return &GTCondition{field: field} }
func GreaterThanOrEqual(field string) *GTECondition { return &GTECondition{field: field} }
func Exists(field string) *ExistsCondition          { return &ExistsCondition{field: field} }

func In(field string, list ...interface{}) *InCondition {
	return &InCondition{field: field, list: list}
}
func NotIn(field string, list ...interface{}) *NotInCondition {
	return &NotInCondition{field: field, list: list}
}

type AscQueryDirective struct{}
type DescQueryDirective struct{}
type LimitQueryDirective struct{ Limit int }
type PageQueryDirective struct{ Page string }

func (*AscQueryDirective) isQueryOp()   {}
func (*DescQueryDirective) isQueryOp()  {}
func (*LimitQueryDirective) isQueryOp() {}
func (*PageQueryDirective) isQueryOp()  {}

func (*AscQueryDirective) isQueryDirective()   {}
func (*DescQueryDirective) isQueryDirective()  {}
func (*LimitQueryDirective) isQueryDirective() {}
func (*PageQueryDirective) isQueryDirective()  {}

func Asc() *AscQueryDirective              { return &AscQueryDirective{} }
func Desc() *DescQueryDirective            { return &DescQueryDirective{} }
func Limit(limit int) *LimitQueryDirective { return &LimitQueryDirective{Limit: limit} }
func Page(page string) *PageQueryDirective { return &PageQueryDirective{Page: page} }
