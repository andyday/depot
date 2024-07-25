package depot

import (
	"context"
)

type Database interface {
	Get(ctx context.Context, table string, entity interface{}) error
	Put(ctx context.Context, table string, entity interface{}) error
	Delete(ctx context.Context, table string, entity interface{}) error
	Create(ctx context.Context, table string, entity interface{}) error
	Update(ctx context.Context, table string, entity interface{}, op ...UpdateOp) error
	Query(ctx context.Context, table, kind string, entity interface{}, entities interface{}, op ...QueryOp) (string, error)
}

type Table[T any] interface {
	Put(ctx context.Context, entity T) (T, error)
	Get(ctx context.Context, entity T) (T, error)
	Delete(ctx context.Context, entity T) (T, error)
	Create(ctx context.Context, entity T) (T, error)
	Update(ctx context.Context, entity T, op ...UpdateOp) (T, error)
	Query(ctx context.Context, kind string, entity T, op ...QueryOp) ([]T, string, error)
}

type table[T any] struct {
	db    Database
	table string
}

func NewTable[T any](db Database, tbl string) Table[T] {
	return &table[T]{db: db, table: tbl}
}

func (t *table[T]) Put(ctx context.Context, entity T) (out T, err error) {
	if err = t.db.Put(ctx, t.table, &entity); err != nil {
		return
	}
	return entity, nil
}

func (t *table[T]) Get(ctx context.Context, entity T) (out T, err error) {
	if err = t.db.Get(ctx, t.table, &entity); err != nil {
		return
	}
	return entity, nil
}

func (t *table[T]) Delete(ctx context.Context, entity T) (out T, err error) {
	if err = t.db.Delete(ctx, t.table, &entity); err != nil {
		return
	}
	return entity, nil
}

func (t *table[T]) Create(ctx context.Context, entity T) (out T, err error) {
	if err = t.db.Create(ctx, t.table, &entity); err != nil {
		return
	}
	return entity, nil
}

func (t *table[T]) Update(ctx context.Context, entity T, op ...UpdateOp) (out T, err error) {
	if err = t.db.Update(ctx, t.table, &entity, op...); err != nil {
		return
	}
	return entity, nil
}
func (t *table[T]) Query(ctx context.Context, kind string, entityFilter T, op ...QueryOp) (entities []T, nextPage string, err error) {
	nextPage, err = t.db.Query(ctx, t.table, kind, &entityFilter, &entities, op...)
	return
}
