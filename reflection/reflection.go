package reflection

import (
	"errors"
	"reflect"
	"strconv"
)

/*
	Reflection package
	ver 1.0 - 2019-02-21
	by Zecchan Silverlake

	This package contains useful function to map values between structs
*/

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

	newSlice := reflect.MakeSlice(frTyp, l, c)

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
