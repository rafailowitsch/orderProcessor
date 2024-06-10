package validate

import (
	"fmt"
	"reflect"
)

func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			if err := ValidateStruct(field.Interface()); err != nil {
				return err
			}
		} else if field.Kind() == reflect.Slice {
			if field.Len() == 0 {
				return fmt.Errorf("missing or empty field: %s", v.Type().Field(i).Name)
			}
			for j := 0; j < field.Len(); j++ {
				if field.Index(j).Kind() == reflect.Struct {
					if err := ValidateStruct(field.Index(j).Interface()); err != nil {
						return err
					}
				}
			}
		} else {
			if IsEmptyValue(field) {
				return fmt.Errorf("missing or empty field: %s", v.Type().Field(i).Name)
			}
		}
	}
	return nil
}

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}
