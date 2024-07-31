package depot

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type KeyPart struct {
	Name  string
	Value interface{}
}

type Key struct {
	Partition KeyPart
	Sort      KeyPart
}

func (k Key) String() string {
	if k.Sort.Value != nil {
		return fmt.Sprintf("%v:%v", k.Partition.Value, k.Sort.Value)
	} else {
		return fmt.Sprintf("%v", k.Partition.Value)
	}
}

type Property struct {
	Name  string
	Value interface{}
	Index bool
}

func EntityKey(entity interface{}) (key Key, err error) {
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		switch f.Mode {
		case FieldModePartition:
			key.Partition.Name = f.Name
			key.Partition.Value = v.Field(i).Interface()
		case FieldModeSort:
			key.Sort.Name = f.Name
			key.Sort.Value = v.Field(i).Interface()
		default:
		}
	}
	return
}

func EntityProperties(entity interface{}) (props []Property, err error) {
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		fv := v.Field(i)
		if f.Mode == FieldModeExclude || (f.Mode == FieldModeOmitEmpty && fv.IsZero()) {
			continue
		}
		props = append(props, Property{Name: f.Name, Value: fv.Interface(), Index: NeedsIndex(f)})
	}
	return
}

func EntityMap(entity interface{}, convertTTL bool) (m map[string]interface{}, err error) {
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	m = make(map[string]interface{})
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		fv := v.Field(i)
		if f.Mode == FieldModeExclude || (f.Mode == FieldModeOmitEmpty && fv.IsZero()) {
			continue
		}
		fvi := fv.Interface()
		if f.TTL && convertTTL {
			fvi = time.Unix(fvi.(int64), 0)
		}
		m[f.Name] = fvi
	}
	return
}

func EntityFromMap(m map[string]interface{}, entity interface{}, convertTTL bool) (err error) {
	var (
		s  Struct
		ev = reflect.ValueOf(entity)
	)
	if s, ev, err = GetStruct(ev); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		fld := ev.Field(i)
		if v, ok := m[f.Name]; ok {
			v = RealSlice(v)
			if f.TTL && convertTTL {
				v = (v.(time.Time)).Unix()
			}
			pv := reflect.ValueOf(v)
			if fld.Kind() == reflect.Ptr && pv.Kind() != reflect.Ptr {
				if fld.IsNil() {
					fld.Set(reflect.New(fld.Type().Elem()))
				}
				fld.Elem().Set(pv)
			} else {
				fld.Set(pv)
			}
		}
	}
	return
}

func RealSlice(v interface{}) interface{} {
	if s, ok := v.([]interface{}); !ok {
		return v
	} else if len(s) <= 0 {
		return v
	} else {
		switch s[0].(type) {
		case string:
			return ConvertSlice[string](s)
		case int:
			return ConvertSlice[int](s)
		case int64:
			return ConvertSlice[int64](s)
		default:
			return v
		}
	}
}

func ConvertSlice[T any](in []interface{}) (out []T) {
	for _, e := range in {
		out = append(out, e.(T))
	}
	return
}

func EntityFromProperties(props []Property, entity interface{}) (err error) {
	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Name] = prop
	}
	return EntityFromPropertyMap(propMap, entity)
}

func EntityFromPropertyMap(props map[string]Property, entity interface{}) (err error) {
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		fld := v.Field(i)
		if p, ok := props[f.Name]; ok {
			pv := reflect.ValueOf(p.Value)
			if fld.Kind() == reflect.Ptr && pv.Kind() != reflect.Ptr {
				if fld.IsNil() {
					fld.Set(reflect.New(fld.Type().Elem()))
				}
				fld.Elem().Set(pv)
			} else {
				fld.Set(pv)
			}
		}
	}
	return
}

type Update struct {
	Name  string
	Value interface{}
	Op    UpdateOp
}

func EntityUpdates(entity interface{}, ops []UpdateOp) (updates []Update, err error) {
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		fv := v.Field(i)
		op := GetUpdateOp(ops, f.Name)
		_, force := op.(*ForceUpdateOp)
		if f.Mode == FieldModeExclude ||
			f.Mode == FieldModePartition ||
			f.Mode == FieldModeSort ||
			(fv.IsZero() && !force) {
			continue
		}
		updates = append(updates, Update{
			Name:  f.Name,
			Value: fv.Interface(),
			Op:    GetUpdateOp(ops, f.Name),
		})
	}
	return
}

type KeyType uint8

const (
	KeyTypeNone KeyType = iota
	KeyTypePartition
	KeyTypeSort
)

type Condition struct {
	Name    string
	Value   interface{}
	KeyType KeyType
	Op      QueryCondition
}

func EntityConditions(kind string, entity interface{}, ops []QueryOp) (sortField string, conditions []Condition, err error) {
	var (
		s  Struct
		v  = reflect.ValueOf(entity)
		kt KeyType
	)
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		fv := v.Field(i)
		op := GetQueryCondition(ops, f.Name)
		value := fv.Interface()

		mode := GetMode(kind, f)
		switch mode {
		case FieldModeExclude:
			continue
		case FieldModePartition:
			kt = KeyTypePartition
		case FieldModeSort:
			kt = KeyTypeSort
			sortField = f.Name
		default:
			kt = KeyTypeNone
		}
		if fv.IsZero() {
			value = nil
			if !(op != nil && op.Valueless()) {
				continue
			}
		}
		conditions = append(conditions, Condition{
			Name:    f.Name,
			Value:   value,
			KeyType: kt,
			Op:      GetQueryCondition(ops, f.Name),
		})
	}
	return
}

func GetMode(kind string, f Field) FieldMode {
	if kind != "" {
		for _, index := range f.Indexes {
			if index.Name == kind {
				return index.Mode
			}
		}
	}
	return f.Mode
}

func NeedsIndex(f Field) bool {
	if f.Mode == FieldModePartition || f.Mode == FieldModeSort {
		return true
	}
	for _, index := range f.Indexes {
		if index.Mode == FieldModeExclude || index.Mode == FieldModeSort {
			return true
		}
	}
	return false
}

func GetUpdateOp(ops []UpdateOp, field string) UpdateOp {
	for _, op := range ops {
		if op.Field() == field {
			return op
		}
	}
	return nil
}

func GetQueryCondition(ops []QueryOp, field string) QueryCondition {
	for _, op := range ops {
		if qc, ok := op.(QueryCondition); ok && qc.Field() == field {
			return qc
		}
	}
	return nil
}

type FieldMode int8

const (
	FieldModeExclude FieldMode = iota
	FieldModeOmitEmpty
	FieldModeInclude
	FieldModePartition
	FieldModeSort
)

type Field struct {
	Name    string
	Mode    FieldMode
	Indexes []Index
	TTL     bool
}

type Index struct {
	Name string
	Mode FieldMode
}

type Struct []Field

var structs = make(map[reflect.Type]Struct)

func GetStruct(v reflect.Value) (s Struct, sv reflect.Value, err error) {
	var ok bool
	sv = v
	if sv.Kind() == reflect.Ptr {
		sv = v.Elem()
	}
	if sv.Kind() == reflect.Interface {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		err = ErrInvalidEntityType
		return
	}

	t := sv.Type()
	if s, ok = structs[t]; ok {
		return
	}

	s = make(Struct, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		s[i] = field(t, i)
	}
	structs[t] = s
	return
}

func field(typ reflect.Type, i int) (fld Field) {
	f := typ.Field(i)
	t := f.Tag
	s := t.Get("depot")
	fld.Name = f.Name
	if s == "-" {
		fld.Mode = FieldModeExclude
		return
	}
	fld.Mode = FieldModeInclude
	parts := strings.Split(s, ",")
	if parts[0] != "" {
		fld.Name = parts[0]
	}
	if len(parts) > 1 {
		for _, p := range parts[1:] {
			switch p {
			case "pk":
				fld.Mode = FieldModePartition
			case "sk":
				fld.Mode = FieldModeSort
			case "omitempty":
				fld.Mode = FieldModeOmitEmpty
			case "ttl":
				fld.TTL = true
			default:
				if strings.HasPrefix(p, "index:") {
					indexParts := strings.Split(p, ":")
					if len(indexParts) != 3 {
						panic("invalid index tag " + p)
					}
					switch indexParts[2] {
					case "pk":
						fld.Indexes = append(fld.Indexes, Index{Name: indexParts[1], Mode: FieldModePartition})
					case "sk":
						fld.Indexes = append(fld.Indexes, Index{Name: indexParts[1], Mode: FieldModeSort})
					default:
						panic("invalid index tag key identifier" + p)
					}

				}
			}
		}
	}
	return
}
