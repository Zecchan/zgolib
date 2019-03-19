package reflection

/*
	Reflection package
	ver 1.1 - 2019-02-21
	Copyright (c) 2019 - Zecchan Silverlake

	This package contains useful function to map values between structs
*/

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Map(from interface{}, to interface{}) error {
	frVal, frTyp, frOK := GetType(from)
	if !frOK {
		return errors.New("mapper.Map() - from must be a pointer")
	}
	toVal, toTyp, toOK := GetType(to)
	if !toOK {
		return errors.New("mapper.Map() - to must be a pointer")
	}
	if toTyp.Kind() != frTyp.Kind() {
		return errors.New("mapper.Map() - from and to must be the same kind")
	}
	kind := frTyp.Kind()
	if kind == reflect.Struct {
		for i := 0; i < toTyp.NumField(); i++ {
			fldName := toTyp.Field(i).Name
			_, found := frTyp.FieldByName(fldName)
			if found && toVal.FieldByName(fldName).CanSet() {
				if frVal.FieldByName(fldName).CanAddr() && toVal.FieldByName(fldName).CanAddr() {
					err := Map(frVal.FieldByName(fldName).Addr().Interface(), toVal.FieldByName(fldName).Addr().Interface())
					if err != nil {
						return err
					}
				} else {
					err := Map(frVal.FieldByName(fldName).Interface(), toVal.FieldByName(fldName).Interface())
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
	if toVal.CanSet() {
		if kind == reflect.Int || kind == reflect.String || kind == reflect.Uint || kind == reflect.Float64 || kind == reflect.Float32 || kind == reflect.Bool {
			toVal.Set(frVal)
			return nil
		}
		if kind == reflect.Slice {
			return MapSlice(from, to)
		}
	}

	return errors.New("mapper.Map() - mapping is not supported for this type")
}

func MapSlice(from interface{}, to interface{}) error {
	frVal, frTyp, frOK := GetType(from)
	if !frOK {
		return errors.New("mapper.SliceMapper() - from must be a pointer")
	}
	toVal, toTyp, toOK := GetType(to)
	if !toOK {
		return errors.New("mapper.SliceMapper() - to must be a pointer")
	}
	if toTyp.Kind() != frTyp.Kind() {
		return errors.New("mapper.SliceMapper() - from and to must be the same kind")
	}
	if frTyp.Kind() != reflect.Slice {
		return errors.New("mapper.SliceMapper() - from and to must be a slice")
	}

	l := frVal.Len()
	c := frVal.Cap()

	newSlice := reflect.MakeSlice(toTyp, l, c)

	for i := 0; i < l; i++ {
		if frVal.Index(i).CanAddr() && newSlice.Index(i).CanAddr() {
			err := Map(frVal.Index(i).Addr().Interface(), newSlice.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		} else {
			err := Map(frVal.Index(i).Interface(), newSlice.Index(i).Interface())
			if err != nil {
				return err
			}
		}
	}

	toVal.Set(newSlice)

	return nil
}

func GetType(obj interface{}) (val reflect.Value, typ reflect.Type, ok bool) {
	otyp := reflect.TypeOf(obj)
	kind := otyp.Kind()
	if kind != reflect.Ptr {
		return reflect.ValueOf(nil), reflect.TypeOf(nil), false
	}
	oval := reflect.ValueOf(obj).Elem()
	kind = oval.Kind()
	otyp = oval.Type()
	return oval, otyp, true
}

func ToString(value interface{}) string {
	if value == nil {
		return "null"
	}
	val, typ, ok := GetType(&value)

	if !ok {
		return "unknown{}"
	}

	var kind = val.Kind()
	if kind == reflect.Int || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 || kind == reflect.Int8 {
		return strconv.FormatInt(val.Int(), 10)
	}
	if kind == reflect.Uint || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 || kind == reflect.Uint8 {
		return strconv.FormatUint(val.Uint(), 10)
	}
	if kind == reflect.Float32 || kind == reflect.Float64 {
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)
	}
	if kind == reflect.String {
		return val.String()
	}
	if kind == reflect.Array {
		return "array[" + typ.Name() + "]"
	}
	if kind == reflect.Slice {
		return "slice[" + typ.Name() + "]"
	}
	if kind == reflect.Struct {
		return "struct{" + typ.Name() + "}"
	}
	if kind == reflect.Bool {
		if val.Bool() {
			return "true"
		}
		return "false"
	}

	return "object{" + typ.Name() + "}"
}

type GoType struct {
	Kind reflect.Kind
	Type reflect.Type
}

// Create creates a new instance of the current GoType initialized to specified obj. If obj is nil, zero values is returned. References are passed as is.
func (typ *GoType) Create(obj interface{}) (interface{}, error) {
	if typ.IsPtrToPtr() {
		return nil, errors.New("Pointer of a pointer is not supported")
	}
	if obj == nil {
		var ptrRes reflect.Value
		if typ.Kind == reflect.Ptr {
			ptrRes = reflect.New(typ.Type.Elem())
		} else {
			ptrRes = reflect.New(typ.Type)
		}
		return ptrRes.Elem().Interface(), nil
	}
	t := GetGoType(obj)
	if t.IsPtrToPtr() {
		return nil, errors.New("Pointer of a pointer is not supported")
	}
	v := reflect.ValueOf(obj)
	if t.IsPtr() {
		v = v.Elem()
	}
	var ptrRes reflect.Value
	if typ.Kind == reflect.Ptr {
		ptrRes = reflect.New(typ.Type.Elem())
	} else {
		ptrRes = reflect.New(typ.Type)
	}
	e := Assign(v, ptrRes)
	if e != nil {
		return nil, e
	}
	return ptrRes.Elem().Interface(), nil
}

// Instantiate is the same as Create, but it does not return an error. If an error occured, nil is returned.
func (typ *GoType) Instantiate(obj interface{}) interface{} {
	val, err := typ.Create(obj)
	if err != nil {
		return nil
	}
	return val
}

// IsPtr checks whether the current GoType is a pointer
func (typ *GoType) IsPtr() bool {
	if typ.Type == nil {
		return false
	}
	return typ.Kind == reflect.Ptr
}

// IsPtrToPtr checks whether the current GoType is a pointer to pointer
func (typ *GoType) IsPtrToPtr() bool {
	if typ.Type == nil {
		return false
	}
	if typ.Type.Kind() != reflect.Ptr {
		return false
	}
	return typ.Type.Elem().Kind() == reflect.Ptr
}

// GetGoType will get a GoType of specified object
func GetGoType(obj interface{}) *GoType {
	typ := GoType{}
	typ.Type = reflect.TypeOf(obj)
	if typ.Type == nil {
		typ.Kind = reflect.Invalid
		return &typ
	}
	typ.Kind = typ.Type.Kind()
	return &typ
}

// GoTypeOf converts reflect.Type to a GoType
func GoTypeOf(rtype reflect.Type) *GoType {
	typ := GoType{}
	typ.Type = rtype
	if typ.Type == nil {
		typ.Kind = reflect.Invalid
		return &typ
	}
	typ.Kind = rtype.Kind()
	return &typ
}

// Assign will assigns a value to a pointer of a value
func Assign(valFrom reflect.Value, ptrTo reflect.Value) error {
	if valFrom.Kind() != ptrTo.Elem().Kind() {
		return errors.New("Cannot assign a " + KindToString(valFrom.Kind()) + " to a " + KindToString(ptrTo.Elem().Kind()))
	}
	fmt.Println(KindToString(valFrom.Kind()) + " to a " + KindToString(ptrTo.Elem().Kind()))
	var kind = valFrom.Kind()
	valTo := ptrTo.Elem()

	if kind == reflect.Struct {
		toType := valTo.Type()
		frType := valFrom.Type()
		for i := 0; i < toType.NumField(); i++ {
			fldName := toType.Field(i).Name
			_, found := frType.FieldByName(fldName)
			if found && valTo.FieldByName(fldName).CanSet() {
				valTo.FieldByName(fldName).Set(valFrom.FieldByName(fldName))
			}
		}
	}

	return nil
}

// KindToString gets a string representation of reflect.Kind
func KindToString(kind reflect.Kind) string {
	switch kind {
	case reflect.Array:
		return "Array"
	case reflect.Bool:
		return "Bool"
	case reflect.Chan:
		return "Chan"
	case reflect.Complex128:
		return "Complex128"
	case reflect.Complex64:
		return "Complex64"
	case reflect.Float32:
		return "Float32"
	case reflect.Float64:
		return "Float64"
	case reflect.Func:
		return "Func"
	case reflect.Int:
		return "Int"
	case reflect.Int16:
		return "Int16"
	case reflect.Int32:
		return "Int32"
	case reflect.Int64:
		return "Int64"
	case reflect.Int8:
		return "Int8"
	case reflect.Interface:
		return "Interface"
	case reflect.Map:
		return "Map"
	case reflect.Ptr:
		return "Pointer"
	case reflect.Slice:
		return "Slice"
	case reflect.String:
		return "String"
	case reflect.Struct:
		return "Struct"
	case reflect.Uint:
		return "Uint"
	case reflect.Uint16:
		return "Uint16"
	case reflect.Uint32:
		return "Uint32"
	case reflect.Uint64:
		return "Uint64"
	case reflect.Uint8:
		return "Uint8"
	case reflect.Uintptr:
		return "Uintptr"
	case reflect.UnsafePointer:
		return "UnsafePointer"
	default:
		return "Invalid"
	}
}
