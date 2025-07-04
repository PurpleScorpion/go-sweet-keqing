package keqing

import (
	"fmt"
	"reflect"
)

// 将 interface{} 转换为指定的切片类型
func ToArray[T any](data interface{}) []T {
	if data == nil {
		return nil
	}

	// 情况 1: data 是 *[]T
	if ptr, ok := data.(*[]T); ok {
		return *ptr
	}

	// 情况 2: data 是 []*T，需要转换成 []T
	if ptrList, ok := data.([]*T); ok {
		res := make([]T, len(ptrList))
		for i := range ptrList {
			res[i] = *ptrList[i]
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
	// 构建忽略字段的 map，提升查找效率
	ignoreMap := make(map[string]struct{})
	for _, field := range ignoreFields {
		ignoreMap[field] = struct{}{}
	}

	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src)

	// 解除指针指向
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	// 确保 src 是结构体类型
	if srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("src must be a struct")
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcType := srcVal.Type()
		field := srcType.Field(i)
		fieldName := field.Name

		// 跳过忽略字段
		if _, ignored := ignoreMap[fieldName]; ignored {
			continue
		}

		dstField, ok := dstVal.Type().FieldByName(fieldName)
		if !ok || !dstVal.FieldByName(fieldName).CanSet() {
			continue
		}

		srcFieldType := srcVal.Field(i).Type()
		dstFieldType := dstField.Type
		if srcFieldType == dstFieldType {
			dstVal.FieldByName(fieldName).Set(srcVal.Field(i))
		}
	}

	return nil
}
