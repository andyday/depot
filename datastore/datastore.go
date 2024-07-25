package datastore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/andyday/depot"
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
		return depot.ErrEntityNotFound
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
		return depot.ErrEntityAlreadyExists
	}
	return
}

func (d *DB) Update(ctx context.Context, table string, entity interface{}, op ...depot.UpdateOp) (err error) {
	var (
		updates []depot.Update
		u       depot.Update
		prop    depot.Property
		propMap = make(datastoreMap)
		ok      bool
	)
	e := &datastoreEntity{kind: table, entity: entity}
	k := datastore.Key{}
	if err = e.LoadKey(&k); err != nil {
		return
	}
	if updates, err = depot.EntityUpdates(entity, op); err != nil {
		return
	}
	_, err = d.datastore.RunInTransaction(ctx, func(tx *datastore.Transaction) (err error) {
		if err = d.datastore.Get(ctx, &k, propMap); errors.Is(err, datastore.ErrNoSuchEntity) {
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
		_, err = d.datastore.Put(ctx, &k, e)
		return
	})
	return
}

func (d *DB) Query(ctx context.Context, table, kind string, entity interface{}, entities interface{}, op ...depot.QueryOp) (string, error) {
	// TODO implement me
	panic("implement me")
}

type datastoreEntity struct {
	kind   string
	entity interface{}
}

func (d *datastoreEntity) LoadKey(k *datastore.Key) (err error) {
	var key depot.Key
	if key, err = depot.EntityKey(d.entity); err != nil {
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
			NoIndex: true,
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
