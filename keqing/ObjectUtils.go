package keqing

// 将 interface{} 转换为指定的切片类型
func GetList[T any](data interface{}) []T {
	if data == nil {
		return nil
	}
	// 类型断言
	if list, ok := data.([]T); ok {
		return list
	}
	return nil
}

// 将 interface{} 强转为指定的单个对象类型
func GetObject[T any](data interface{}) T {
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
