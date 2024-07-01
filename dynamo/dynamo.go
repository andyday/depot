package dynamo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andyday/depot"
	"github.com/andyday/depot/internal"
	"github.com/andyday/depot/transform"
	types2 "github.com/andyday/depot/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DB struct {
	dynamo  *dynamodb.Client
	encoder *attributevalue.Encoder
	decoder *attributevalue.Decoder
}

var _ depot.Database = &DB{}

func encoderOptions(opts *attributevalue.EncoderOptions) {
	opts.TagKey = "depot"
}

func decoderOptions(opts *attributevalue.DecoderOptions) {
	opts.TagKey = "depot"
}

func NewDatabase(cfg aws.Config) (c *DB, err error) {
	c = &DB{
		dynamo: dynamodb.NewFromConfig(cfg),
		encoder: attributevalue.NewEncoder(func(opts *attributevalue.EncoderOptions) {
			opts.TagKey = "depot"
		}),
		decoder: attributevalue.NewDecoder(func(opts *attributevalue.DecoderOptions) {
			opts.TagKey = "depot"
		}),
	}
	return
}

func (d *DB) Table(name string) *depot.Table {
	return depot.NewTable(d, name)
}

func (d *DB) Get(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		out *dynamodb.GetItemOutput
		in  = &dynamodb.GetItemInput{TableName: aws.String(table)}
	)
	if in.Key, err = keyFromEntity(entity); err != nil {
		return
	}

	if out, err = d.dynamo.GetItem(ctx, in); err != nil {
		return
	}
	if len(out.Item) <= 0 {
		return types2.ErrEntityNotFound
	}
	err = attributevalue.UnmarshalMapWithOptions(out.Item, entity, decoderOptions)
	return
}

func (d *DB) Put(ctx context.Context, table string, entity interface{}) (err error) {
	var in = &dynamodb.PutItemInput{TableName: aws.String(table)}
	if in.Item, err = marshalEntity(entity); err != nil {
		return
	}
	_, err = d.dynamo.PutItem(ctx, in)
	return
}

func (d *DB) Delete(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		inp = &dynamodb.DeleteItemInput{TableName: aws.String(table), ReturnValues: types.ReturnValueAllOld}
		out *dynamodb.DeleteItemOutput
	)
	if inp.Key, err = keyFromEntity(entity); err != nil {
		return
	}

	if out, err = d.dynamo.DeleteItem(ctx, inp); err != nil {
		return
	}
	return unmarshalEntity(out.Attributes, entity)
}

func (d *DB) Create(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		item map[string]types.AttributeValue
		k    internal.Key
	)
	if item, err = attributevalue.MarshalMapWithOptions(entity, encoderOptions); err != nil {
		return
	}
	if k, err = internal.EntityKey(entity); err != nil {
		return
	}

	if _, err = d.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                aws.String(table),
		Item:                     item,
		ExpressionAttributeNames: map[string]string{"#pk": k.Partition.Name},
		ConditionExpression:      aws.String("attribute_not_exists(#pk)"),
	}); errorIsConditionCheckFailure(err) {
		return types2.ErrEntityAlreadyExists
	}
	return
}

func (d *DB) Update(ctx context.Context, table string, entity interface{}) (err error) {
	var (
		key     map[string]types.AttributeValue
		updates map[string]interface{}
		names   = make(map[string]string)
		values  = make(map[string]types.AttributeValue)
		mv      types.AttributeValue
		exp     strings.Builder
	)
	if key, err = keyFromEntity(entity); err != nil {
		return
	}
	if updates, err = internal.EntityUpdates(entity); err != nil {
		return
	}

	for k, v := range updates {
		names["#"+k] = k
		if mv, err = updateValue(v); err != nil {
			return err
		}
		values[":"+k] = mv
	}

	set, add := updateExpressionParts(updates)
	if len(set) > 0 {
		exp.WriteString("SET ")
		exp.WriteString(strings.Join(set, ", "))
	}
	if len(set) > 0 && len(add) > 0 {
		exp.WriteRune(' ')
	}
	if len(add) > 0 {
		exp.WriteString("ADD ")
		exp.WriteString(strings.Join(add, ", "))
	}

	_, err = d.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       key,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
		UpdateExpression:          aws.String(strings.TrimSpace(exp.String())),
		ReturnValues:              types.ReturnValueAllNew,
	})
	return
}

func keyMap(k internal.Key) (m map[string]interface{}) {
	m = make(map[string]interface{})
	m[k.Partition.Name] = k.Partition.Value
	if k.Sort.Name != "" {
		m[k.Sort.Name] = k.Sort.Value
	}
	return
}

func keyFromEntity(entity interface{}) (key map[string]types.AttributeValue, err error) {
	var k internal.Key
	if k, err = internal.EntityKey(entity); err != nil {
		return
	}
	key, err = attributevalue.MarshalMap(keyMap(k))
	return
}

func marshalEntity(entity interface{}) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMapWithOptions(entity, encoderOptions)
}

func unmarshalEntity(item map[string]types.AttributeValue, entity interface{}) (err error) {
	v := reflect.ValueOf(entity)
	if len(item) > 0 && v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		return attributevalue.UnmarshalMapWithOptions(item, entity, decoderOptions)
	}
	return
}

func updateExpressionParts(updates map[string]interface{}) (set, add []string) {
	for k, v := range updates {
		if tf, ok := v.(*transform.Transform); ok {
			switch tf.Type {
			case transform.TypeAdd:
				add = append(add, fmt.Sprintf("#%s :%s", k, k))
			case transform.TypeSubtract:
				add = append(add, fmt.Sprintf("#%s :%s", k, k))
			default:
				continue
			}
		} else {
			set = append(set, fmt.Sprintf("#%s = :%s", k, k))
		}
	}
	return
}

func updateValue(v interface{}) (av types.AttributeValue, err error) {
	if tf, ok := v.(*transform.Transform); ok {
		switch tf.Type {
		case transform.TypeSubtract:
			v = internal.NegateValue(tf.Value)
		default:
			v = tf.Value
		}
	}
	av, err = attributevalue.Marshal(v)
	return
}

func errorIsConditionCheckFailure(err error) bool {
	var conditionCheckFailure *types.ConditionalCheckFailedException
	return errors.As(err, &conditionCheckFailure)
}
