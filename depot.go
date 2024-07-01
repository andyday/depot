package depot

import "context"

type Database interface {
	Table(name string) *Table
	Get(ctx context.Context, table string, entity interface{}) error
	Put(ctx context.Context, table string, entity interface{}) error
	Delete(ctx context.Context, table string, entity interface{}) error
	Create(ctx context.Context, table string, entity interface{}) error
	Update(ctx context.Context, table string, entity interface{}) error
}

type Table struct {
	db    Database
	table string
}

func NewTable(db Database, table string) *Table {
	return &Table{db: db, table: table}
}

func (t *Table) Put(ctx context.Context, entity interface{}) error {
	return t.db.Put(ctx, t.table, entity)
}

func (t *Table) Get(ctx context.Context, entity interface{}) error {
	return t.db.Get(ctx, t.table, entity)
}

func (t *Table) Delete(ctx context.Context, entity interface{}) error {
	return t.db.Delete(ctx, t.table, entity)
}

func (t *Table) Create(ctx context.Context, entity interface{}) error {
	return t.db.Create(ctx, t.table, entity)
}

func (t *Table) Update(ctx context.Context, entity interface{}) error {
	return t.db.Update(ctx, t.table, entity)
}
