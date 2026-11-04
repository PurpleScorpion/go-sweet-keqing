package keqing

import (
	"errors"
	"fmt"
	"math"
	"strconv"
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
	return int(yamlIntegerValue(key, strconv.IntSize))
}

func ValueInt64(key string) int64 {
	return int64(yamlIntegerValue(key, 64))
}

func ValueInt32(key string) int32 {
	return int32(yamlIntegerValue(key, 32))
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

// yamlIntegerValue preserves integer width and rejects conversion overflow.
// Finite fractional values keep the existing truncation-toward-zero behavior.
func yamlIntegerValue(key string, bits int) int64 {
	value := getYamlValue(key)
	if value == nil {
		return 0
	}
	var integer int64
	var floating float64
	isFloat := false
	switch v := value.(type) {
	case int:
		integer = int64(v)
	case int8:
		integer = int64(v)
	case int16:
		integer = int64(v)
	case int32:
		integer = int64(v)
	case int64:
		integer = v
	case float32:
		floating, isFloat = float64(v), true
	case float64:
		floating, isFloat = v, true
	default:
		panic(fmt.Sprintf("YAML value %s of type %T cannot be converted to integer", key, value))
	}
	if isFloat {
		floating = math.Trunc(floating)
		limit := math.Ldexp(1, bits-1)
		if math.IsNaN(floating) || math.IsInf(floating, 0) || floating < -limit || floating >= limit {
			panic(fmt.Sprintf("YAML value %s is outside int%d range", key, bits))
		}
		return int64(floating)
	}
	if bits < 64 {
		limit := int64(1) << (bits - 1)
		if integer < -limit || integer >= limit {
			panic(fmt.Sprintf("YAML value %s is outside int%d range", key, bits))
		}
	}
	return integer
}
