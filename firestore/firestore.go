package firestore

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/andyday/depot"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DB struct {
	firestore *firestore.Client
}

var _ depot.Database = &DB{}

func NewDatabase(ctx context.Context, projectID, databaseID string) (c *DB, err error) {
	c = &DB{}
	if databaseID == "" {
		c.firestore, err = firestore.NewClient(ctx, projectID)
	} else {
		c.firestore, err = firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	}
	return
}

func (d *DB) Put(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		doc *firestore.DocumentRef
		m   map[string]interface{}
	)
	if doc, err = d.doc(table, entity); err != nil {
		return
	}
	if m, err = depot.EntityMap(entity, true); err != nil {
		return
	}
	_, err = doc.Set(ctx, m)
	return
}

func (d *DB) Get(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		doc *firestore.DocumentRef
		res *firestore.DocumentSnapshot
	)
	if doc, err = d.doc(table, entity); err != nil {
		return
	}
	if res, err = doc.Get(ctx); status.Code(err) == codes.NotFound {
		return depot.ErrEntityNotFound
	} else if err != nil {
		return
	}
	return depot.EntityFromMap(res.Data(), entity, true)
}

func (d *DB) Delete(ctx context.Context, table string, entity interface{}) (err error) {
	var doc *firestore.DocumentRef
	if doc, err = d.doc(table, entity); err != nil {
		return
	}
	_, err = doc.Delete(ctx)
	return
}

func (d *DB) Create(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		doc *firestore.DocumentRef
		m   map[string]interface{}
	)
	if doc, err = d.doc(table, entity); err != nil {
		return
	}
	if m, err = depot.EntityMap(entity, true); err != nil {
		return
	}
	if _, err = doc.Create(ctx, m); status.Code(err) == codes.AlreadyExists {
		return depot.ErrEntityAlreadyExists
	}
	return
}

func (d *DB) Update(ctx context.Context, table string, entity interface{}, op ...depot.UpdateOp) (err error) {
	var (
		doc          *firestore.DocumentRef
		res          *firestore.DocumentSnapshot
		depotUpdates []depot.Update
		u            depot.Update
		existing     map[string]interface{}
		updates      []firestore.Update
		v            interface{}
	)
	if doc, err = d.doc(table, entity); err != nil {
		return
	}
	if depotUpdates, err = depot.EntityUpdates(entity, op); err != nil {
		return
	}

	err = d.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) (err error) {
		if res, err = tx.Get(doc); err != nil {
			return
		}
		existing = res.Data()
		for _, u = range depotUpdates {
			v = existing[u.Name]
			switch u.Op.(type) {
			case *depot.AddUpdateOp:
				v = depot.AddValues(v, u.Value)
			case *depot.SubtractUpdateOp:
				v = depot.SubtractValues(v, u.Value)
			default:
				v = u.Value
			}
			updates = append(updates, firestore.Update{Path: u.Name, Value: v})
		}
		err = tx.Update(doc, updates)
		return
	})
	return
}

func (d *DB) Query(ctx context.Context, table, kind string, entity interface{}, entities interface{}, op ...depot.QueryOp) (page string, err error) {
	var (
		conditions []depot.Condition
		sortField  string
		offset     int
		q          firestore.Query
	)

	if sortField, conditions, err = depot.EntityConditions(kind, entity, op); err != nil {
		return
	}
	q = applyQueryConditions(d.firestore.Collection(table), conditions)
	if q, offset, err = applyQueryDirectives(q, op, sortField); err != nil {
		return
	}
	qr := newQueryRunner(table, entities, offset)
	if offset, err = qr.run(q.Documents(ctx)); err != nil {
		return
	}
	if err = qr.next(q.Offset(offset).Documents(ctx), false); errors.Is(err, iterator.Done) {
		// If this is an error it means it is done at the 1st item or there was an error
		return "", nil
	}
	if err != nil {
		return
	}

	return strconv.Itoa(offset), nil
}

func (d *DB) doc(table string, entity interface{}) (_ *firestore.DocumentRef, err error) {
	var k string
	if k, err = LoadKey(table, entity); err != nil {
		return
	}
	return d.firestore.Doc(k), nil
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

func (q *queryRunner) run(it *firestore.DocumentIterator) (offset int, err error) {
	defer it.Stop()
	for {
		if err = q.next(it, true); errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return
		}
	}
	return q.offset, nil
}

func (q *queryRunner) next(it *firestore.DocumentIterator, append bool) (err error) {
	var res *firestore.DocumentSnapshot
	if !append {
		defer it.Stop()
	}
	ev := reflect.New(q.elementType)
	if res, err = it.Next(); err != nil {
		return
	}
	if err = depot.EntityFromMap(res.Data(), ev.Interface(), true); err != nil {
		return
	}
	q.offset++
	if append {
		q.listValue.Set(reflect.Append(q.listValue, ev.Elem()))
	}
	return
}

func applyQueryConditions(in *firestore.CollectionRef, conditions []depot.Condition) (q firestore.Query) {
	q = in.Query
	for _, c := range conditions {
		switch c.Op.(type) {
		case *depot.EqualQueryCondition:
			q = q.Where(c.Name, "=", c.Value)
		case *depot.NotEqualQueryCondition:
			q = q.Where(c.Name, "!=", c.Value)
		case *depot.LTQueryCondition:
			q = q.Where(c.Name, "<", c.Value)
		case *depot.LTEQueryCondition:
			q = q.Where(c.Name, "<=", c.Value)
		case *depot.GTQueryCondition:
			q = q.Where(c.Name, ">", c.Value)
		case *depot.GTEQueryCondition:
			q = q.Where(c.Name, ">=", c.Value)
		case *depot.ExistsQueryCondition:
			q = q.Where(c.Name, "!=", "-DEAD-BEEF-")
		default:
			q = q.Where(c.Name, "==", c.Value)
		}
	}
	return
}

func applyQueryDirectives(in firestore.Query, ops []depot.QueryOp, sortField string) (q firestore.Query, offset int, err error) {
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
				q = q.OrderBy(sortField, firestore.Asc)
			case *depot.DescQueryDirective:
				if sortField == "" {
					return q, 0, depot.ErrNoSortField
				}
				q = q.OrderBy(sortField, firestore.Desc)
			}
		}
	}
	return
}

func LoadKey(table string, entity interface{}) (k string, err error) {
	var key depot.Key
	if key, err = depot.EntityKey(entity); err != nil {
		return
	}
	return fmt.Sprintf("%s/%s", table, key.String()), nil
}
