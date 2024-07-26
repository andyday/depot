package dynamo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andyday/depot"
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
		return depot.ErrEntityNotFound
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
		k    depot.Key
	)
	if item, err = attributevalue.MarshalMapWithOptions(entity, encoderOptions); err != nil {
		return
	}
	if k, err = depot.EntityKey(entity); err != nil {
		return
	}

	if _, err = d.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                aws.String(table),
		Item:                     item,
		ExpressionAttributeNames: map[string]string{"#pk": k.Partition.Name},
		ConditionExpression:      aws.String("attribute_not_exists(#pk)"),
	}); errorIsConditionCheckFailure(err) {
		return depot.ErrEntityAlreadyExists
	}
	return
}

func (d *DB) Update(ctx context.Context, table string, entity interface{}, op ...depot.UpdateOp) (err error) {
	var (
		key     map[string]types.AttributeValue
		updates []depot.Update
		names   = make(map[string]string)
		values  = make(map[string]types.AttributeValue)
		mv      types.AttributeValue
		exp     strings.Builder
	)
	if key, err = keyFromEntity(entity); err != nil {
		return
	}
	if updates, err = depot.EntityUpdates(entity, op); err != nil {
		return
	}

	for _, u := range updates {
		names["#"+u.Name] = u.Name
		if mv, err = updateValue(u); err != nil {
			return err
		}
		values[":"+u.Name] = mv
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

func (d *DB) Query(ctx context.Context, table, kind string, entity interface{}, entities interface{}, op ...depot.QueryOp) (nextPage string, err error) {
	var (
		idx        *string
		names      = make(map[string]string)
		values     = make(map[string]types.AttributeValue)
		conditions []depot.Condition
		cv         types.AttributeValue
		res        *dynamodb.QueryOutput
		keyExp     *string
		filterExp  *string
		limit      *int32
		page       map[string]types.AttributeValue
		asc        *bool
	)
	if kind != "" {
		idx = &kind
	}

	if conditions, err = depot.EntityConditions(kind, entity, op); err != nil {
		return
	}

	for _, c := range conditions {
		names["#"+c.Name] = c.Name
		if c.Value != nil {
			if cv, err = conditionValue(c); err != nil {
				return
			}
			values[":"+c.Name] = cv
		}
	}

	keyExp, filterExp = queryExpressionParts(conditions)
	if limit, page, asc, err = queryDirectives(op); err != nil {
		return
	}

	if res, err = d.dynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(table),
		IndexName:                 idx,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
		KeyConditionExpression:    keyExp,
		FilterExpression:          filterExp,
		Limit:                     limit,
		ExclusiveStartKey:         page,
		ScanIndexForward:          asc,
	}); err != nil {
		return
	}

	if err = unmarshalEntities(res.Items, entities); err != nil {
		return
	}
	return EncodePage(res.LastEvaluatedKey)
}

func keyMap(k depot.Key) (m map[string]interface{}) {
	m = make(map[string]interface{})
	m[k.Partition.Name] = k.Partition.Value
	if k.Sort.Name != "" {
		m[k.Sort.Name] = k.Sort.Value
	}
	return
}

func keyFromEntity(entity interface{}) (key map[string]types.AttributeValue, err error) {
	var k depot.Key
	if k, err = depot.EntityKey(entity); err != nil {
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

func unmarshalEntities(items []map[string]types.AttributeValue, entities interface{}) (err error) {
	return attributevalue.UnmarshalListOfMapsWithOptions(items, entities, decoderOptions)
}

func updateExpressionParts(updates []depot.Update) (set, add []string) {
	for _, u := range updates {
		switch u.Op.(type) {
		case *depot.AddUpdateOp:
			add = append(add, fmt.Sprintf("#%s :%s", u.Name, u.Name))
		case *depot.SubtractUpdateOp:
			add = append(add, fmt.Sprintf("#%s :%s", u.Name, u.Name))
		default:
			set = append(set, fmt.Sprintf("#%s = :%s", u.Name, u.Name))
		}
	}
	return
}

func queryExpressionParts(conditions []depot.Condition) (keyExp, filterExp *string) {
	var (
		keyParts    []string
		filterParts []string
	)
	for _, c := range conditions {
		var exp string
		switch c.Op.(type) {
		case *depot.EqualQueryCondition:
			exp = fmt.Sprintf("#%s = :%s", c.Name, c.Name)
		case *depot.NotEqualQueryCondition:
			exp = fmt.Sprintf("#%s <> :%s", c.Name, c.Name)
		case *depot.LTQueryCondition:
			exp = fmt.Sprintf("#%s < :%s", c.Name, c.Name)
		case *depot.LTEQueryCondition:
			exp = fmt.Sprintf("#%s <= :%s", c.Name, c.Name)
		case *depot.GTQueryCondition:
			exp = fmt.Sprintf("#%s > :%s", c.Name, c.Name)
		case *depot.GTEQueryCondition:
			exp = fmt.Sprintf("#%s >= :%s", c.Name, c.Name)
		case *depot.ExistsQueryCondition:
			exp = fmt.Sprintf("attribute_exists(#%s)", c.Name)
		// case *depot.NotExistsQueryCondition:
		// 	exp = fmt.Sprintf("attribute_not_exists(#%s)", c.Name)
		// case *depot.PrefixQueryCondition:
		// 	exp = fmt.Sprintf("begins_with(#%s, :%s)", c.Name, c.Name)
		// case *depot.ContainsQueryCondition:
		// 	exp = fmt.Sprintf("contains(#%s, :%s)", c.Name, c.Name)
		default:
			exp = fmt.Sprintf("#%s = :%s", c.Name, c.Name)
		}

		if c.KeyType == depot.KeyTypeNone {
			filterParts = append(filterParts, exp)
		} else {
			keyParts = append(keyParts, exp)
		}
	}
	if len(keyParts) > 0 {
		keyExp = aws.String(strings.Join(keyParts, " AND "))
	}
	if len(filterParts) > 0 {
		filterExp = aws.String(strings.Join(filterParts, " AND "))
	}
	return
}

func queryDirectives(ops []depot.QueryOp) (limit *int32, page map[string]types.AttributeValue, asc *bool, err error) {
	for _, op := range ops {
		if d, ok := op.(depot.QueryDirective); ok {
			switch v := d.(type) {
			case *depot.LimitQueryDirective:
				limit = aws.Int32(int32(v.Limit))
			case *depot.PageQueryDirective:
				if page, err = DecodePage(v.Page); err != nil {
					return
				}
			case *depot.AscQueryDirective:
				asc = aws.Bool(true)
			case *depot.DescQueryDirective:
				asc = aws.Bool(false)
			}
		}
	}
	return
}

func DecodePage(encoded string) (decoded map[string]types.AttributeValue, err error) {
	var (
		bytes []byte
		m     = make(map[string]map[string]string)
	)
	decoded = make(map[string]types.AttributeValue)
	if bytes, err = base64.RawURLEncoding.DecodeString(encoded); err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &m); err != nil {
		return
	}
	for k, v := range m {
		for kk, vv := range v {
			switch kk {
			case "S":
				decoded[k] = &types.AttributeValueMemberS{Value: vv}
			case "N":
				decoded[k] = &types.AttributeValueMemberN{Value: vv}
			}
		}
	}

	return

}

func EncodePage(decoded map[string]types.AttributeValue) (encoded string, err error) {
	if decoded == nil {
		return "", nil
	}
	var (
		bytes []byte
		m     = make(map[string]map[string]string)
	)
	for k, v := range decoded {
		switch vv := v.(type) {
		case *types.AttributeValueMemberS:
			m[k] = map[string]string{"S": vv.Value}
		case *types.AttributeValueMemberN:
			m[k] = map[string]string{"N": vv.Value}
		}
	}
	if bytes, err = json.Marshal(m); err != nil {
		return
	}
	encoded = base64.RawURLEncoding.EncodeToString(bytes)
	return
}

func updateValue(u depot.Update) (av types.AttributeValue, err error) {
	var v interface{}
	switch u.Op.(type) {
	case *depot.SubtractUpdateOp:
		v = depot.NegateValue(u.Value)
	case *depot.AddUpdateOp:
		v = u.Value
	default:
		v = u.Value
	}
	av, err = attributevalue.Marshal(v)
	return
}

func conditionValue(c depot.Condition) (av types.AttributeValue, err error) {
	// var v interface{}
	// switch c.Op.Type {
	// case depot.TypeSubtract:
	// 	v = depot.NegateValue(c.Value)
	// default:
	// 	v = c.Value
	// }
	av, err = attributevalue.Marshal(c.Value)
	return
}

func errorIsConditionCheckFailure(err error) bool {
	var conditionCheckFailure *types.ConditionalCheckFailedException
	return errors.As(err, &conditionCheckFailure)
}
