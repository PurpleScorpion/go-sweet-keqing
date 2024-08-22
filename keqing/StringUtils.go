package keqing

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

/*
字符串替换 - 正则替换
mainStr: 主字符串
pattern: 正则表达式
new: 要替换成的字符串
*/
func Replace2Regex(mainStr, pattern, new string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(mainStr, new)
}

/*
字符串替换
mainStr: 主字符串
old: 要替换的字符串
new: 要替换成的字符串
*/
func Replace(mainStr, old, new string) string {
	return strings.Replace(mainStr, old, new, -1)
}

/*
查找小串在大串中的位置 - 正序
mainStr: 大串
substr: 小串
*/
func IndexOf(mainStr, substr string) int {

	return strings.Index(mainStr, substr)
}

/*
查找小串在大串中的位置 - 倒序
mainStr: 大串
substr: 小串
*/
func LastIndexOf(mainStr, substr string) int {
	return strings.LastIndex(mainStr, substr)
}

/*
将字符串解析成byte数组
*/
func GetBytes(str string) []byte {
	return []byte(str)
}

/*
字符串中大串是否包含小串
mainStr: 大串
substr: 小串
*/
func Contains(mainStr, substr string) bool {
	return strings.Contains(mainStr, substr)
}

/*
将 数组 / 结构体 转换为标准字符串打印
*/
func ToString(obj interface{}) string {
	if isSlice(obj) {
		return array2Str(obj)
	}
	if isStruct(obj) {
		return struct2Str(obj)
	}
	// 如果是map 类型
	if isMap(obj) {
		return map2Str(obj)
	}

	return ""
}

func map2Str(obj interface{}) string {
	mymap := obj.(map[string]interface{})
	var result []string
	for key, value := range mymap {
		if isStruct(value) {
			result = append(result, fmt.Sprintf("%s: %v", key, struct2Str(value)))
		} else if isMap(value) {
			result = append(result, fmt.Sprintf("%s: %v", key, map2Str(value)))
		} else if isSlice(value) {
			result = append(result, fmt.Sprintf("%s: %v", key, array2Str(value)))
		} else {
			fmt.Println("aaaaaaaaaaaaaaaa", key)
			result = append(result, fmt.Sprintf("%s: %v", key, value))
		}
	}
	return "{" + strings.Join(result, ", ") + "}"
}

func struct2Str(obj interface{}) string {
	var result []string
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result = append(result, fmt.Sprintf("%s: %v", field.Name, struct2Str(value.Interface())))
		} else if isMap(value.Interface()) {
			result = append(result, fmt.Sprintf("%s: %v", field.Name, map2Str(value.Interface())))
		} else if isSlice(value.Interface()) {
			result = append(result, fmt.Sprintf("%s: %v", field.Name, array2Str(value.Interface())))
		} else {

			result = append(result, fmt.Sprintf("%s: %v", field.Name, value))
		}
	}
	return "{" + strings.Join(result, ", ") + "}"
}

func array2Str(obj interface{}) string {
	// 获取 obj 的反射值
	objVal := reflect.ValueOf(obj)
	var slice []interface{}
	for i := 0; i < objVal.Len(); i++ {
		elem := objVal.Index(i)
		// 检查元素类型是否与 value 相同
		slice = append(slice, elem.Interface())
	}
	var result []string

	for _, item := range slice {
		switch v := item.(type) {
		case string:
			result = append(result, fmt.Sprintf(`"%s"`, v))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result = append(result, strconv.FormatInt(reflect.ValueOf(v).Int(), 10))
		case bool:
			result = append(result, strconv.FormatBool(v))
		default:
			if isStruct(item) {
				result = append(result, struct2Str(item))
			} else if isMap(item) {
				result = append(result, map2Str(item))
			} else if isSlice(item) {
				result = append(result, array2Str(item))
			} else {
				result = append(result, fmt.Sprintf("%v", v))
			}
		}
	}

	return "[" + strings.Join(result, ", ") + "]"
}

/*
字符串转字符串数组
*/
func ToStringArray(str string) []string {
	var chars []string
	for _, r := range str {
		chars = append(chars, string(r))
	}
	return chars
}

/*
字符串首字母转大写
*/
func ToFirstUpperCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if unicode.IsLetter(r[0]) {
		r[0] = unicode.ToUpper(r[0])
	}
	return string(r)
}

/*
字符串首字母转小写
*/
func ToFirstLowerCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if unicode.IsLetter(r[0]) {
		r[0] = unicode.ToLower(r[0])
	}
	return string(r)
}

/*
字符串转大写
*/
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

/*
字符串转小写
*/
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

/*
移除字符串两端空白
*/
func Trim(str string) string {
	return strings.TrimSpace(str)
}

/*
去除字符串最后一个指定字符(如果有的话)
*/
func TrimSuffix(s string, suffix string) string {
	if len(s) == 0 {
		return s
	}
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

/*
下划线转小驼峰
*/
func Snake2Camel(s string) string {
	words := strings.Split(s, "_")
	var camelCase strings.Builder

	for i, word := range words {
		if i > 0 {
			camelCase.WriteString(strings.Title(word))
		} else {
			camelCase.WriteString(word)
		}
	}
	return camelCase.String()
}

/*
下划线转大驼峰
*/
func Snake2BigCamel(s string) string {
	words := strings.Split(s, "_")
	var pascalCase strings.Builder

	for _, word := range words {
		pascalCase.WriteString(strings.Title(word))
	}
	return pascalCase.String()
}

/*
驼峰转下划线
*/
func Camel2Snake(s string) string {
	if IsEmpty(s) {
		return s
	}

	if unicode.IsUpper(rune(s[0])) {
		return bigCamelToSnake(s)
	}
	return smallCamelToSnake(s)
}

func bigCamelToSnake(s string) string {
	var snakeCase strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			snakeCase.WriteRune('_')
		}
		snakeCase.WriteRune(unicode.ToLower(r))
	}

	return snakeCase.String()
}

func smallCamelToSnake(s string) string {
	var snakeCase strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				snakeCase.WriteRune('_')
			}
			snakeCase.WriteRune(unicode.ToLower(r))
		} else {
			snakeCase.WriteRune(r)
		}
	}

	return snakeCase.String()
}

/*
数字类型转换为字符串
包含的类型有
int int8 int16 int32 int64
uint uint8 uint16 uint32 uint64
float32 float64
*/
func Num2Str(obj interface{}) string {
	if obj == nil {
		return ""
	}
	switch num := obj.(type) {
	case int:
		return int64ToStr(int64(num))
	case int8:
		return int64ToStr(int64(num))
	case int16:
		return int64ToStr(int64(num))
	case int32:
		return int64ToStr(int64(num))
	case int64:
		return int64ToStr(num)
	case uint:
		return uint64ToStr(uint64(num))
	case uint8:
		return uint64ToStr(uint64(num))
	case uint16:
		return uint64ToStr(uint64(num))
	case uint32:
		return uint64ToStr(uint64(num))
	case uint64:
		return uint64ToStr(num)
	case float32:
		return strconv.FormatFloat(float64(num), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(num, 'f', -1, 64)
	default:
		return ""
	}
}

/*
浮点数转字符串 (可选择保留的小数位数)
若不需要精度，则直接调用Num2Str即可
*/
func Float2Str(obj interface{}, precision int) string {
	if obj == nil {
		return ""
	}

	switch num := obj.(type) {
	case float32:
		if precision < 0 {
			return Num2Str(obj)
		}
		return strconv.FormatFloat(float64(num), 'f', precision, 32)
	case float64:
		if precision < 0 {
			return Num2Str(obj)
		}
		return strconv.FormatFloat(num, 'f', precision, 64)
	default:
		return ""
	}
}

func ParseInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(str + " Cannot be converted to int64 type")
	}
	return num
}

func ParseInt32(str string) int32 {
	return int32(ParseInt64(str))
}

func ParseInt16(str string) int16 {
	return int16(ParseInt64(str))
}

func ParseInt8(str string) int8 {
	return int8(ParseInt64(str))
}

func ParseInt(str string) int {
	return int(ParseInt64(str))
}

func ParseFloat32(dataStr string) float32 {
	f, _ := strconv.ParseFloat(dataStr, 32)
	return float32(f)
}

func ParseFloat64(dataStr string) float64 {
	f, _ := strconv.ParseFloat(dataStr, 64)
	return f
}

func Split(str string, sep string) []string {
	if IsEmpty(str) {
		return []string{str}
	}
	if !strings.Contains(str, sep) {
		return []string{str}
	}
	return strings.Split(str, sep)
}

func int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func uint64ToStr(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func float64ToStr(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
