package keqing

import (
	"fmt"
	"reflect"
)

// 将 interface{} 转换为指定的切片类型；nil 切片指针返回 nil，nil 元素转换为 T 的零值。
func ToArray[T any](data interface{}) []T {
	if data == nil {
		return nil
	}

	// 情况 1: data 是 *[]T
	if ptr, ok := data.(*[]T); ok {
		if ptr == nil {
			return nil
		}
		return *ptr
	}

	// 情况 2: data 是 []*T，需要转换成 []T
	if ptrList, ok := data.([]*T); ok {
		res := make([]T, len(ptrList))
		for i := range ptrList {
			if ptrList[i] != nil {
				res[i] = *ptrList[i]
			}
		}
		return res
	}

	// 情况 3: data 是 []T
	if list, ok := data.([]T); ok {
		return list
	}

	return nil
}

// 将 interface{} 强转为指定的单个对象类型
func ToObject[T any](data interface{}) T {
	var t T
	if data == nil {
		return t
	}
	// 类型断言
	if value, ok := data.(T); ok {
		return value
	}
	return t
}

// CopyProperties 将 src 结构体的字段复制到 dst 结构体中（字段名相同），可指定忽略字段
func CopyProperties(dst, src interface{}, ignoreFields ...string) error {
	dstVal := reflect.ValueOf(dst)
	if !dstVal.IsValid() || dstVal.Kind() != reflect.Ptr || dstVal.IsNil() || dstVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst must be a non-nil pointer to a struct")
	}
	dstVal = dstVal.Elem()
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("src must be a struct")
	}
	ignoreMap := make(map[string]struct{}, len(ignoreFields))
	for _, field := range ignoreFields {
		ignoreMap[field] = struct{}{}
	}
	for i := 0; i < srcVal.NumField(); i++ {
		field := srcVal.Type().Field(i)
		source := srcVal.Field(i)
		if _, ignored := ignoreMap[field.Name]; ignored || !source.CanInterface() {
			continue
		}
		targetField, ok := dstVal.Type().FieldByName(field.Name)
		if !ok || source.Type() != targetField.Type {
			continue
		}
		target, err := dstVal.FieldByIndexErr(targetField.Index)
		if err != nil {
			return fmt.Errorf("cannot copy field %s: %w", field.Name, err)
		}
		if target.CanSet() {
			target.Set(source)
		}
	}
	return nil
}
