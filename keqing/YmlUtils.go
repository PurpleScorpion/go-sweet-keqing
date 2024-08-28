package keqing

import (
	"errors"
	"strings"
)

var (
	applicationYml      map[string]interface{} = nil
	applicationChildYml map[string]interface{} = nil
)

func LoadYml(parent map[string]interface{}, child map[string]interface{}) error {
	if parent == nil {
		return errors.New("Parent Yml is nil")
	}
	if child == nil {
		return errors.New("Child Yml is nil")
	}
	applicationYml = parent
	applicationChildYml = child
	return nil
}

func ValueObject(key string) interface{} {
	return getYamlValue(key)
}

func ValueInt(key string) int {
	val := getYamlValue(key)
	if val == nil {
		return 0
	}

	switch v := val.(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	}

	return val.(int)
}

func ValueInt64(key string) int64 {
	return int64(ValueInt(key))
}

func ValueInt32(key string) int32 {
	return int32(ValueInt(key))
}

func ValueFloat32(key string) float32 {
	return float32(ValueFloat64(key))
}

func ValueFloat64(key string) float64 {
	val := getYamlValue(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	}
	return val.(float64)
}

func ValueBool(key string) bool {
	val := getYamlValue(key)
	if val == nil {
		return false
	}
	return val.(bool)
}

func ValueString(key string) string {
	val := getYamlValue(key)
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case float64:
		return Num2Str(v)
	case float32:
		return Num2Str(v)
	case int:
		return Num2Str(v)
	case int8:
		return Num2Str(v)
	case int16:
		return Num2Str(v)
	case int32:
		return Num2Str(v)
	case int64:
		return Num2Str(v)
	}
	return val.(string)
}

func ValueStringArr(key string) []string {
	val := getYamlValue(key)
	if val == nil {
		return []string{}
	}
	val2 := val.([]interface{})
	var arr []string
	for i := 0; i < len(val2); i++ {
		arr = append(arr, val2[i].(string))
	}
	return arr
}

func getYamlValue(key string) interface{} {
	if !(strings.HasPrefix(key, "${") && strings.HasSuffix(key, "}")) {
		panic("key must start with ${ and end with }")
	}
	key = key[2 : len(key)-1]
	arr := strings.Split(key, ".")
	val := applicationYml[arr[0]]
	val2 := applicationChildYml[arr[0]]
	if len(arr) == 1 {
		if val2 == nil {
			return val
		}
		return val2
	}

	for i := 1; i < len(arr); i++ {
		tmp := arr[i]
		if val != nil {
			v := val.(map[string]interface{})
			val = v[tmp]
		}
		if val2 != nil {
			v := val2.(map[string]interface{})
			val2 = v[tmp]
		}
	}
	if val2 == nil {
		return val
	}
	return val2
}
