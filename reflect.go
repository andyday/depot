package depot

import (
	"reflect"
	"strings"
)

type KeyPart struct {
	Name  string
	Value string
}

type Key struct {
	Partition KeyPart
	Sort      KeyPart
}

type Property struct {
	Name  string
	Value interface{}
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
			key.Partition.Value = v.Field(i).String()
		case FieldModeSort:
			key.Sort.Name = f.Name
			key.Sort.Value = v.Field(i).String()
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
		props = append(props, Property{Name: f.Name, Value: fv.Interface()})
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
			fld.Set(pv)
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
		if f.Mode == FieldModeExclude ||
			f.Mode == FieldModePartition ||
			f.Mode == FieldModeSort ||
			fv.IsZero() {
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

func EntityConditions(kind string, entity interface{}, ops []QueryOp) (conditions []Condition, err error) {
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

		if fv.IsZero() {
			value = nil
			if !(op != nil && op.Valueless()) {
				continue
			}
		}
		mode := GetMode(kind, f)
		switch mode {
		case FieldModeExclude:
			continue
		case FieldModePartition:
			kt = KeyTypePartition
		case FieldModeSort:
			kt = KeyTypeSort
		default:
			kt = KeyTypeNone
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
		err = ErrEntityNotFound
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
