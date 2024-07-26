package datastore

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/andyday/depot"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DB struct {
	datastore *datastore.Client
}

var _ depot.Database = &DB{}

func NewDatabase(ctx context.Context, projectID, databaseID string) (c *DB, err error) {
	c = &DB{}
	c.datastore, err = datastore.NewClientWithDatabase(ctx, projectID, databaseID)
	return
}

func (d *DB) Put(ctx context.Context, table string, entity interface{}) (err error) {
	var k *datastore.Key
	if k, err = LoadKey(table, entity); err != nil {
		return
	}
	_, err = d.datastore.Put(ctx, k, &datastoreEntity{entity: entity})
	return
}

func (d *DB) Get(ctx context.Context, table string, entity interface{}) (err error) {
	var k *datastore.Key
	if k, err = LoadKey(table, entity); err != nil {
		return
	}
	if err = d.datastore.Get(ctx, k, &datastoreEntity{entity: entity}); errors.Is(err, datastore.ErrNoSuchEntity) {
		return depot.ErrEntityNotFound
	}
	return
}

func (d *DB) Delete(ctx context.Context, table string, entity interface{}) (err error) {
	var k *datastore.Key
	if k, err = LoadKey(table, entity); err != nil {
		return
	}
	return d.datastore.Delete(ctx, k)
}

func (d *DB) Create(ctx context.Context, table string, entity interface{}) (err error) {
	var k *datastore.Key
	if k, err = LoadKey(table, entity); err != nil {
		return
	}
	if _, err = d.datastore.Mutate(ctx, datastore.NewInsert(k, &datastoreEntity{entity: entity})); status.Code(err) == codes.AlreadyExists {
		return depot.ErrEntityAlreadyExists
	}
	return
}

func (d *DB) Update(ctx context.Context, table string, entity interface{}, op ...depot.UpdateOp) (err error) {
	var (
		k       *datastore.Key
		updates []depot.Update
		u       depot.Update
		prop    depot.Property
		propMap = make(datastoreMap)
		ok      bool
	)
	if k, err = LoadKey(table, entity); err != nil {
		return
	}
	if updates, err = depot.EntityUpdates(entity, op); err != nil {
		return
	}
	_, err = d.datastore.RunInTransaction(ctx, func(tx *datastore.Transaction) (err error) {
		if err = d.datastore.Get(ctx, k, propMap); errors.Is(err, datastore.ErrNoSuchEntity) {
			return depot.ErrEntityNotFound
		} else if err != nil {
			return
		}
		for _, u = range updates {
			if prop, ok = propMap[u.Name]; !ok {
				prop = depot.Property{Name: u.Name}
			}
			switch u.Op.(type) {
			case *depot.AddUpdateOp:
				prop.Value = depot.AddValues(prop.Value, u.Value)
			case *depot.SubtractUpdateOp:
				prop.Value = depot.SubtractValues(prop.Value, u.Value)
			default:
				prop.Value = u.Value
			}
			propMap[u.Name] = prop
		}

		if err = depot.EntityFromPropertyMap(propMap, entity); err != nil {
			return
		}
		_, err = d.datastore.Put(ctx, k, &datastoreEntity{entity: entity})
		return
	})
	return
}

func (d *DB) Query(ctx context.Context, table, kind string, entity interface{}, entities interface{}, op ...depot.QueryOp) (page string, err error) {
	var (
		conditions []depot.Condition
		offset     int
		q          = datastore.NewQuery(table)
	)

	if conditions, err = depot.EntityConditions(kind, entity, op); err != nil {
		return
	}
	sortField := getSortField(conditions)
	q = applyQueryConditions(q, conditions)
	if q, offset, err = applyQueryDirectives(q, op, sortField); err != nil {
		return
	}
	qr := newQueryRunner(table, entities, offset)
	if offset, err = qr.run(d.datastore.Run(ctx, q)); err != nil {
		return
	}
	if err = qr.next(d.datastore.Run(ctx, q.Offset(offset))); errors.Is(err, iterator.Done) {
		// If this is an error it means it is done at the 1st item or there was an error
		return "", nil
	}
	if err != nil {
		return
	}

	return strconv.Itoa(offset), nil
}

type queryRunner struct {
	table       string
	list        interface{}
	listValue   reflect.Value
	elementType reflect.Type
	offset      int
}

func newQueryRunner(table string, list interface{}, offset int) *queryRunner {
	q := &queryRunner{table: table, list: list, offset: offset}
	q.listValue = reflect.ValueOf(q.list)
	if q.listValue.Kind() == reflect.Ptr {
		q.listValue = q.listValue.Elem()
	}
	q.elementType = q.listValue.Type().Elem()
	return q
}

func (q *queryRunner) run(it *datastore.Iterator) (offset int, err error) {
	for {
		if err = q.next(it); errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return
		}
	}
	return q.offset, nil
}

func (q *queryRunner) next(it *datastore.Iterator) (err error) {
	ev := reflect.New(q.elementType)
	if _, err = it.Next(&datastoreEntity{entity: ev.Interface()}); err != nil {
		return
	}
	q.listValue.Set(reflect.Append(q.listValue, ev.Elem()))
	q.offset++
	return
}

// it := d.datastore.Run(ctx, q)
// v := reflect.ValueOf(entities)
// v = v.Elem()
// elemType := v.Type().Elem()
//
// for {
// 	ev := reflect.New(elemType)
// 	item := datastoreEntity{kind: table, entity: ev.Interface()}
// 	if _, err = it.Next(&item); errors.Is(err, iterator.Done) {
// 		if cursor, err = it.Cursor(); err != nil {
// 			return
// 		}
// 		return cursor.String(), nil
// 	}
// 	if err != nil {
// 		return
// 	}
// dv.Set(reflect.Append(dv, ev))
// 	reflect.Append(v, ev.Elem())
// }

func getSortField(conditions []depot.Condition) string {
	for _, c := range conditions {
		if c.KeyType == depot.KeyTypeSort {
			return c.Name
		}
	}
	return ""
}

func applyQueryConditions(in *datastore.Query, conditions []depot.Condition) (q *datastore.Query) {
	q = in
	for _, c := range conditions {
		switch c.Op.(type) {
		case *depot.EqualQueryCondition:
			q = q.FilterField(c.Name, "=", c.Value)
		case *depot.NotEqualQueryCondition:
			q = q.FilterField(c.Name, "!=", c.Value)
		case *depot.LTQueryCondition:
			q = q.FilterField(c.Name, "<", c.Value)
		case *depot.LTEQueryCondition:
			q = q.FilterField(c.Name, "<=", c.Value)
		case *depot.GTQueryCondition:
			q = q.FilterField(c.Name, ">", c.Value)
		case *depot.GTEQueryCondition:
			q = q.FilterField(c.Name, ">=", c.Value)
		case *depot.ExistsQueryCondition:
			q = q.FilterField(c.Name, "!=", "-DEAD-BEEF-")
		default:
			q = q.FilterField(c.Name, "=", c.Value)
		}
	}
	return
}

func applyQueryDirectives(in *datastore.Query, ops []depot.QueryOp, sortField string) (q *datastore.Query, offset int, err error) {
	q = in
	for _, op := range ops {
		if d, ok := op.(depot.QueryDirective); ok {
			switch v := d.(type) {
			case *depot.LimitQueryDirective:
				q = q.Limit(v.Limit)
			case *depot.PageQueryDirective:
				// if cursor, err = datastore.DecodeCursor(v.Page); err != nil {
				// 	return
				// }
				// q = q.Start(cursor)
				if offset, err = strconv.Atoi(v.Page); err != nil {
					return
				}
				q = q.Offset(offset)
			case *depot.AscQueryDirective:
				if sortField == "" {
					return q, 0, depot.ErrNoSortField
				}
				q = q.Order(sortField)
			case *depot.DescQueryDirective:
				if sortField == "" {
					return q, 0, depot.ErrNoSortField
				}
				q = q.Order("-" + sortField)
			}
		}
	}
	return
}

// type datastoreKey struct {
// 	kind   string
// 	entity interface{}
// }
//
// func (d *datastoreKey) LoadKey(k *datastore.Key) (err error) {
// 	var key depot.Key
// 	if key, err = depot.EntityKey(d.entity); err != nil {
// 		return
// 	}
// 	k.Kind = d.kind
// 	if key.Sort.Value != "" {
// 		k.Name = fmt.Sprintf("%s:%s", key.Partition.Value, key.Sort.Value)
// 	} else {
// 		k.Name = key.Partition.Value
// 	}
// 	return
// }

func LoadKey(kind string, entity interface{}) (k *datastore.Key, err error) {
	var key depot.Key
	if key, err = depot.EntityKey(entity); err != nil {
		return
	}
	k = &datastore.Key{Kind: kind}
	if key.Sort.Value != "" {
		k.Name = fmt.Sprintf("%s:%s", key.Partition.Value, key.Sort.Value)
	} else {
		k.Name = key.Partition.Value
	}
	return
}

type datastoreEntity struct {
	entity interface{}
}

func (d *datastoreEntity) Load(properties []datastore.Property) (err error) {
	return depot.EntityFromProperties(fromDatastoreProps(properties), d.entity)
}

func (d *datastoreEntity) Save() (datastoreProps []datastore.Property, err error) {
	var props []depot.Property
	if props, err = depot.EntityProperties(d.entity); err != nil {
		return
	}
	return toDatastoreProps(props), nil
}

type datastoreMap map[string]depot.Property

func (d datastoreMap) Load(properties []datastore.Property) (err error) {
	for _, prop := range properties {
		d[prop.Name] = depot.Property{Name: prop.Name, Value: prop.Value}
	}
	return
}

func (d datastoreMap) Save() (datastoreProps []datastore.Property, err error) {
	panic("not implemented")
}

func toDatastoreProps(in []depot.Property) (out []datastore.Property) {
	for _, prop := range in {
		out = append(out, datastore.Property{
			Name:    prop.Name,
			Value:   prop.Value,
			NoIndex: !prop.Index,
		})
	}
	return
}

func fromDatastoreProps(in []datastore.Property) (out []depot.Property) {
	for _, prop := range in {
		out = append(out, depot.Property{
			Name:  prop.Name,
			Value: prop.Value,
		})
	}
	return
}
