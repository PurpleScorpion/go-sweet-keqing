package keqing

import (
	"reflect"
)

func IsEmpty(obj interface{}) bool {
	if obj == nil {
		return true
	}

	switch v := obj.(type) {
	case string:
		return len(v) == 0 || ToUpperCase(v) == "NULL"
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case bool:
		return v
	}

	if isSlice(obj) {
		v := reflect.ValueOf(obj)
		return v.Len() == 0
	}

	if isMap(obj) {
		v := reflect.ValueOf(obj)
		return v.Len() == 0
	}

	if isStruct(obj) {
		panic("not support struct")
	}

	return false
}

func IsNotEmpty(obj interface{}) bool {
	return !IsEmpty(obj)
}

func isSlice(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	kind := value.Kind()

	return kind == reflect.Slice
}

func isStruct(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	kind := value.Kind()

	return kind == reflect.Struct
}

func isMap(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	kind := value.Kind()

	return kind == reflect.Map
}
