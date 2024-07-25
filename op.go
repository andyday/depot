package depot

type UpdateOp interface {
	isUpdateOp()
	Field() string
}

type AddUpdateOp struct{ field string }
type SubtractUpdateOp struct{ field string }

func (*AddUpdateOp) isUpdateOp()      {}
func (*SubtractUpdateOp) isUpdateOp() {}

func (o *AddUpdateOp) Field() string      { return o.field }
func (o *SubtractUpdateOp) Field() string { return o.field }

func Add(field string) *AddUpdateOp           { return &AddUpdateOp{field: field} }
func Subtract(field string) *SubtractUpdateOp { return &SubtractUpdateOp{field: field} }

type QueryOp interface{ isQueryOp() }
type QueryCondition interface {
	isQueryCondition()
	Field() string
	Valueless() bool
}
type QueryDirective interface{ isQueryDirective() }

type EqualQueryCondition struct{ field string }
type NotEqualQueryCondition struct{ field string }
type LTQueryCondition struct{ field string }
type LTEQueryCondition struct{ field string }
type GTQueryCondition struct{ field string }
type GTEQueryCondition struct{ field string }
type ExistsQueryCondition struct{ field string }
type NotExistsQueryCondition struct{ field string }
type PrefixQueryCondition struct{ field string }
type ContainsQueryCondition struct{ field string }

func (*EqualQueryCondition) isQueryOp()     {}
func (*NotEqualQueryCondition) isQueryOp()  {}
func (*LTQueryCondition) isQueryOp()        {}
func (*LTEQueryCondition) isQueryOp()       {}
func (*GTQueryCondition) isQueryOp()        {}
func (*GTEQueryCondition) isQueryOp()       {}
func (*ExistsQueryCondition) isQueryOp()    {}
func (*NotExistsQueryCondition) isQueryOp() {}
func (*PrefixQueryCondition) isQueryOp()    {}
func (*ContainsQueryCondition) isQueryOp()  {}

func (q *EqualQueryCondition) Field() string     { return q.field }
func (q *NotEqualQueryCondition) Field() string  { return q.field }
func (q *LTQueryCondition) Field() string        { return q.field }
func (q *LTEQueryCondition) Field() string       { return q.field }
func (q *GTQueryCondition) Field() string        { return q.field }
func (q *GTEQueryCondition) Field() string       { return q.field }
func (q *ExistsQueryCondition) Field() string    { return q.field }
func (q *NotExistsQueryCondition) Field() string { return q.field }
func (q *PrefixQueryCondition) Field() string    { return q.field }
func (q *ContainsQueryCondition) Field() string  { return q.field }

func (q *EqualQueryCondition) Valueless() bool     { return false }
func (q *NotEqualQueryCondition) Valueless() bool  { return false }
func (q *LTQueryCondition) Valueless() bool        { return false }
func (q *LTEQueryCondition) Valueless() bool       { return false }
func (q *GTQueryCondition) Valueless() bool        { return false }
func (q *GTEQueryCondition) Valueless() bool       { return false }
func (q *ExistsQueryCondition) Valueless() bool    { return true }
func (q *NotExistsQueryCondition) Valueless() bool { return true }
func (q *PrefixQueryCondition) Valueless() bool    { return false }
func (q *ContainsQueryCondition) Valueless() bool  { return false }

func (*EqualQueryCondition) isQueryCondition()     {}
func (*NotEqualQueryCondition) isQueryCondition()  {}
func (*LTQueryCondition) isQueryCondition()        {}
func (*LTEQueryCondition) isQueryCondition()       {}
func (*GTQueryCondition) isQueryCondition()        {}
func (*GTEQueryCondition) isQueryCondition()       {}
func (*ExistsQueryCondition) isQueryCondition()    {}
func (*NotExistsQueryCondition) isQueryCondition() {}
func (*PrefixQueryCondition) isQueryCondition()    {}
func (*ContainsQueryCondition) isQueryCondition()  {}

func Equal(field string) *EqualQueryCondition            { return &EqualQueryCondition{field: field} }
func NotEqual(field string) *NotEqualQueryCondition      { return &NotEqualQueryCondition{field: field} }
func LessThan(field string) *LTQueryCondition            { return &LTQueryCondition{field: field} }
func LessThanOrEqual(field string) *LTEQueryCondition    { return &LTEQueryCondition{field: field} }
func GreaterThan(field string) *GTQueryCondition         { return &GTQueryCondition{field: field} }
func GreaterThanOrEqual(field string) *GTEQueryCondition { return &GTEQueryCondition{field: field} }
func Exists(field string) *ExistsQueryCondition          { return &ExistsQueryCondition{field: field} }
func NotExists(field string) *NotExistsQueryCondition    { return &NotExistsQueryCondition{field: field} }
func Prefix(field string) *PrefixQueryCondition          { return &PrefixQueryCondition{field: field} }
func Contains(field string) *ContainsQueryCondition      { return &ContainsQueryCondition{field: field} }

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
