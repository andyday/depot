package datastore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/andyday/depot"
	"github.com/andyday/depot/internal"
	"github.com/andyday/depot/transform"
	"github.com/andyday/depot/types"
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

func (d *DB) Table(name string) *depot.Table {
	return depot.NewTable(d, name)
}

func (d *DB) Put(ctx context.Context, table string, entity interface{}) (err error) {
	e := &datastoreEntity{kind: table, entity: entity}
	k := datastore.Key{}
	if err = e.LoadKey(&k); err != nil {
		return
	}
	_, err = d.datastore.Put(ctx, &k, e)
	return
}

func (d *DB) Get(ctx context.Context, table string, entity interface{}) (err error) {
	e := &datastoreEntity{kind: table, entity: entity}
	k := datastore.Key{}
	if err = e.LoadKey(&k); err != nil {
		return
	}
	if err = d.datastore.Get(ctx, &k, e); errors.Is(err, datastore.ErrNoSuchEntity) {
		return types.ErrEntityNotFound
	}
	return
}

func (d *DB) Delete(ctx context.Context, table string, entity interface{}) (err error) {
	e := &datastoreEntity{kind: table, entity: entity}
	k := datastore.Key{}
	if err = e.LoadKey(&k); err != nil {
		return
	}
	return d.datastore.Delete(ctx, &k)
}

func (d *DB) Create(ctx context.Context, table string, entity interface{}) (err error) {
	e := &datastoreEntity{kind: table, entity: entity}
	k := datastore.Key{}
	if err = e.LoadKey(&k); err != nil {
		return
	}
	if _, err = d.datastore.Mutate(ctx, datastore.NewInsert(&k, e)); status.Code(err) == codes.AlreadyExists {
		return types.ErrEntityAlreadyExists
	}
	return
}

func (d *DB) Update(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		updates map[string]interface{}
		v       interface{}
		ok      bool
		tf      *transform.Transform
	)
	e := &datastoreEntity{kind: table, entity: entity}
	k := datastore.Key{}
	if err = e.LoadKey(&k); err != nil {
		return
	}
	if updates, err = internal.EntityUpdates(entity); err != nil {
		return
	}
	_, err = d.datastore.RunInTransaction(ctx, func(tx *datastore.Transaction) (err error) {
		var props []datastore.Property
		if err = d.datastore.Get(ctx, &k, e); errors.Is(err, datastore.ErrNoSuchEntity) {
			return types.ErrEntityNotFound
		} else if err != nil {
			return
		}
		if props, err = e.Save(); err != nil {
			return
		}
		for i, prop := range props {
			if v, ok = updates[prop.Name]; !ok {
				continue
			}
			if tf, ok = v.(*transform.Transform); ok {
				switch tf.Type {
				case transform.TypeAdd:
					prop.Value = internal.AddValues(prop.Value, tf.Value)
				case transform.TypeSubtract:
					prop.Value = internal.SubtractValues(prop.Value, tf.Value)
				default:
					continue
				}
			} else {
				prop.Value = v
			}
			props[i] = prop
		}
		if err = e.Load(props); err != nil {
			return
		}
		_, err = d.datastore.Put(ctx, &k, e)
		return
	})
	return
}

type datastoreEntity struct {
	kind   string
	entity interface{}
}

func (d *datastoreEntity) LoadKey(k *datastore.Key) (err error) {
	var key internal.Key
	if key, err = internal.EntityKey(d.entity); err != nil {
		return
	}
	k.Kind = d.kind
	if key.Sort.Value != "" {
		k.Name = fmt.Sprintf("%s:%s", key.Partition.Value, key.Sort.Value)
	} else {
		k.Name = key.Partition.Value
	}
	return
}

func (d *datastoreEntity) Load(properties []datastore.Property) (err error) {
	return internal.EntityFromProperties(fromDatastoreProps(properties), d.entity)
}

func (d *datastoreEntity) Save() (datastoreProps []datastore.Property, err error) {
	var props []internal.Property
	if props, err = internal.EntityProperties(d.entity); err != nil {
		return
	}
	return toDatastoreProps(props), nil
}

func toDatastoreProps(in []internal.Property) (out []datastore.Property) {
	for _, prop := range in {
		out = append(out, datastore.Property{
			Name:    prop.Name,
			Value:   prop.Value,
			NoIndex: true,
		})
	}
	return
}

func fromDatastoreProps(in []datastore.Property) (out []internal.Property) {
	for _, prop := range in {
		out = append(out, internal.Property{
			Name:  prop.Name,
			Value: prop.Value,
		})
	}
	return
}
