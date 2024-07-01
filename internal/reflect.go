package internal

import (
	"reflect"
	"strings"

	"github.com/andyday/depot/types"
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
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	if s, v, err = GetStruct(v); err != nil {
		return
	}
	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Name] = prop
	}
	ln := len(s)
	for i := 0; i < ln; i++ {
		f := s[i]
		if p, ok := propMap[f.Name]; ok {
			pv := reflect.ValueOf(p.Value)
			v.Field(i).Set(pv)
		}
	}
	return
}

func EntityUpdates(entity interface{}) (updates map[string]interface{}, err error) {
	var (
		s Struct
		v = reflect.ValueOf(entity)
	)
	updates = make(map[string]interface{})
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
		updates[f.Name] = fv.Interface()
	}
	return
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
		err = types.ErrEntityNotFound
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
			}
		}
	}
	return
}
