package demo

import (
	"errors"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"testing"
)

func TestDemo17(t *testing.T) {
	keqing.RsaLoadKey("D:\\img\\publicKey.pem", "D:\\img\\privateKey.pem", keqing.RSA_KEY_FILE_TYPE)
	str := keqing.RsaEncrypt("剑光如我,斩尽牛杂")
	fmt.Println(str)
	fmt.Println("-------------------")
	data := keqing.RsaDecrypt(str)
	fmt.Println(data)
}
func TestDemo16(t *testing.T) {
	keqing.RsaGenerateKey("D:\\img")
}

func TestDemo15(t *testing.T) {
	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuOP6JK8tT0WdIvMbK60T
x/BrbuZN/tJn+HVXwh9xdDAyygMhpI7yfX1NFzkAtoBb9dCEbmvAtImSlq+grDMQ
8qse/KI0W0viQLbIXwUdIaSx/qKYjUprYgBwlMN8woV99IB+SwZ7T5SHMb1LK6xM
9ybaeiNxOPoeGy4S6YvEpC1tjtKKq8kpBlaN3psx5H2jeR//Jj/h4V/7trzWO5Em
LgcC9jx4OazrAR85jOKgLo9L0KmKe6S4icQ8It7LN17BciqUKMrYQS0kj6R+HtoI
3f79wNBLdTna3DYsF3YDLBrlPV8GB0PSTUTqNSoX42/qlZVuH8wz8TGpQijsRc/v
3wIDAQAB
-----END PUBLIC KEY-----`
	privateKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC44/okry1PRZ0i
8xsrrRPH8Gtu5k3+0mf4dVfCH3F0MDLKAyGkjvJ9fU0XOQC2gFv10IRua8C0iZKW
r6CsMxDyqx78ojRbS+JAtshfBR0hpLH+opiNSmtiAHCUw3zChX30gH5LBntPlIcx
vUsrrEz3Jtp6I3E4+h4bLhLpi8SkLW2O0oqrySkGVo3emzHkfaN5H/8mP+HhX/u2
vNY7kSYuBwL2PHg5rOsBHzmM4qAuj0vQqYp7pLiJxDwi3ss3XsFyKpQoythBLSSP
pH4e2gjd/v3A0Et1OdrcNiwXdgMsGuU9XwYHQ9JNROo1Khfjb+qVlW4fzDPxMalC
KOxFz+/fAgMBAAECggEAHYDzt8rkdhPrwVn96fhSgcNRwX6qz5EP2kwPVwDhf+L5
F9dsFPBirbfDB4OnI3hUNGOz3lL/i0+wvq8D+raja7X22eWgaTkwv5brXo5YWbgI
V1Pm+BT6Ecd0L6kKTZgzw0KF5L8CCm7vK/bC+hMirQXcM0VYmfj/uOKfTflpxbDc
OPtzT99lVhjJg2e7GTnwpIXERI8SByjc+L1pQrN4FsoOl4VTg8o+QeLcKHghFaZ1
jnQb9B2EQ1TrsGnww1tNrNOvV9k/UBXCH+A3s6OCaHfNLoREATMqrtX11a2qzMZz
7VZtWC7uqI/FjZZEDluXYlJHqe9VY4c7cVqMIuNSMQKBgQDz72qDrdufi5z61h89
vWg0giIC6gkyy6OhwDYbv4ydI/9m3vexGbkX33YLmSxxKmgiYT6UD+LvjiHw6Waa
0iFX7L0FF1KgPRxhfs1g37YPWD1e/2tpT17VKLqjMAX4xOwZ3HfVvZBI/0g6hjzZ
sBB3Vtq0iyU+yXFZ5u2pgtqp7QKBgQDCCPec900fNF6e08Z0fxWVYqPiOw4I3iWb
emi3pGPWMLPoOwKDPfguSaHJKhi/Gq9Ipd30T2xkvuW7jMFzMOcGhoVMmFVuLfJz
ttj8J5x1YtgmNrXYwtWXQJypqmZJTYwXOBzUYRiplIZNvA1FIy0tSVyH1Zg85/43
vpElJ+wXewKBgHAOuaWIBm4CWri4CF36VpZYeXtRO6yD88VoYPLaSaQeV0NQhgRr
RqX612V4lfveeTvh5DdsHNnjNyBOd/4DLaIQdLyT/Db0G8eF0p7/5ciixn6PYy5b
cbsGHMa+Vt/yxmsS5lHf5RpDe1C3PdjakpXf5lQt34w6ScH83YyTOhP5AoGBAK5C
s28LQv4lYF0wQOlbQR0aq6h/9QjNyeSquOVFBEzXDJwicw0/WGbpxh0Oa48l/go2
vPGvat/H+jbIIOy9HJ7lrU2u+fqr1TVLH/DF+mQKU6luNT7pLD5cztYprRdkR86K
nIm4chfKxhuGKjzPbMFhQ3LSx5jbmZqi0WQXSJeFAoGARorE8VlMqDbgtSP/TApp
HAr/UuVdmf8Ag+ZUpa2o33E39L9yLnj7L60GZKVJ87weYj/CHatiCIzpStNsGkoU
81iOD/cUhLl7K9B0WeT7ECfO4Vk6xg6n0SH6YykEyudAmlrtVfxZ1EclQ2OKhxfN
Ia8OCVyYe7VFKVYFkh4PBhQ=
-----END PRIVATE KEY-----`
	str := keqing.RsaEncryptCustom(publicKey, "剑光如我,斩尽牛杂")
	fmt.Println(str)
	fmt.Println("-------------------")
	data := keqing.RsaDecryptCustom(privateKey, str)
	fmt.Println(data)
}

func demo14() {
	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 15, 20, 215, 30, 35, 27, 40, 55, 77}

	arr := keqing.Difference(intSlice1, intSlice2)
	fmt.Println(keqing.ToString(arr))

	stringSlice := []string{"banana", "apple1", "orange1", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	arr2 := keqing.Difference(stringSlice, expectedStringSlice)
	fmt.Println(keqing.ToString(arr2))
}

func demo13() {
	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 15, 20, 215, 30, 35, 27, 40, 55, 77}

	arr := keqing.Union(intSlice1, intSlice2)
	fmt.Println(keqing.ToString(arr))

	stringSlice := []string{"banana", "apple1", "orange1", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	arr2 := keqing.Union(stringSlice, expectedStringSlice)
	fmt.Println(keqing.ToString(arr2))
}

func demo12() {
	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 15, 20, 215, 30, 35, 27, 40, 55, 77}

	arr := keqing.Intersection(intSlice1, intSlice2)
	fmt.Println(keqing.ToString(arr))

	stringSlice := []string{"banana", "apple1", "orange1", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	arr2 := keqing.Intersection(stringSlice, expectedStringSlice)
	fmt.Println(keqing.ToString(arr2))
}

func demo11() {

	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	expectedIntSlice1 := []int32{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}
	expectedIntSlice2 := []int64{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}

	keqing.SortArrayDesc(&intSlice1)
	keqing.SortArrayDesc(&intSlice2)
	keqing.SortArrayDesc(&expectedIntSlice1)
	keqing.SortArrayDesc(&expectedIntSlice2)
	fmt.Println(keqing.ToString(intSlice1))
	fmt.Println(keqing.ToString(intSlice2))
	fmt.Println(keqing.ToString(expectedIntSlice1))
	fmt.Println(keqing.ToString(expectedIntSlice2))

	stringSlice := []string{"banana", "apple", "orange", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	keqing.SortArrayDesc(&stringSlice)
	keqing.SortArrayDesc(&expectedStringSlice)
	fmt.Println(keqing.ToString(stringSlice))
	fmt.Println(keqing.ToString(expectedStringSlice))

}

func demo10() {

	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	expectedIntSlice1 := []int32{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}
	expectedIntSlice2 := []int64{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}

	keqing.SortArray(&intSlice1)
	keqing.SortArray(&intSlice2)
	keqing.SortArray(&expectedIntSlice1)
	keqing.SortArray(&expectedIntSlice2)
	fmt.Println(keqing.ToString(intSlice1))
	fmt.Println(keqing.ToString(intSlice2))
	fmt.Println(keqing.ToString(expectedIntSlice1))
	fmt.Println(keqing.ToString(expectedIntSlice2))

	stringSlice := []string{"banana", "apple", "orange", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	keqing.SortArray(&stringSlice)
	keqing.SortArray(&expectedStringSlice)
	fmt.Println(keqing.ToString(stringSlice))
	fmt.Println(keqing.ToString(expectedStringSlice))

}

func demo9() {
	fmt.Println("NowDate:  ", keqing.NowDate())
	fmt.Println("NowUTCDate:  ", keqing.NowUTCDate())
	fmt.Println("FormatNowDate:  ", keqing.NowDateStr())
	fmt.Println("FormatNowUTCDate:  ", keqing.NowUTCDateStr())
	fmt.Println("CurrentTimeMillis:  ", keqing.CurrentTimeMillis())

	fmt.Println("NowLocalDate:  ", keqing.ToString(keqing.DateInfo(keqing.NowDate())))
	fmt.Println("NowUTCDate:  ", keqing.ToString(keqing.DateInfo(keqing.NowUTCDate())))
}

func demo8() {

	var arr = make([]interface{}, 0)
	arr = append(arr, "a")
	arr = append(arr, "2")
	arr = append(arr, 3)
	arr = append(arr, true)
	arr = append(arr, false)

	fmt.Println(keqing.ToString(arr))

	//var obj = dog{
	//	Age:  1,
	//	Name: "dog",
	//}
	//fmt.Println(keqing.ToString(obj))
}

func demo7() {
	str := "  123456  "
	fmt.Println("AAA" + keqing.Trim(str) + "BBB")

	fmt.Println("TrimSuffix: ", keqing.TrimSuffix("123456789", "789"))
	fmt.Println("TrimSuffix: ", keqing.TrimSuffix("123456789", "9"))
	fmt.Println("TrimSuffix: ", keqing.TrimSuffix("123456789", "asd"))

	fmt.Println("Split: ", keqing.Split("123,456,789", ","))
	fmt.Println("ToFirstUpperCase: ", keqing.ToFirstUpperCase("a123456"))
	fmt.Println("ToFirstUpperCase: ", keqing.ToFirstUpperCase("123456"))
	fmt.Println("ToFirstLowerCase: ", keqing.ToFirstLowerCase("B123456"))
	fmt.Println("ToFirstLowerCase: ", keqing.ToFirstLowerCase("123456"))
	fmt.Println("ToUpperCase: ", keqing.ToUpperCase("123aWSsda456"))
	fmt.Println("ToLowerCase: ", keqing.ToLowerCase("12ADWsaSA3456"))

	fmt.Println("ToStringArray: ", keqing.ToStringArray("剑光如网,斩尽芜杂"))

}

func demo6() {
	for i := 0; i < 20; i++ {
		fmt.Println("UUID: ", keqing.UUID())
	}

}

func demo5() {
	a := "simple_test"
	b := "http_request"

	c := "JsonResponse"
	d := "xmlHttpRequest"

	fmt.Println("下划线转驼峰: ", keqing.Snake2Camel(a))
	fmt.Println("下划线转大驼峰: ", keqing.Snake2BigCamel(b))

	fmt.Println("大驼峰转下划线", keqing.Camel2Snake(c))
	fmt.Println("小驼峰转下划线", keqing.Camel2Snake(d))

}

func demo4() {
	fmt.Println("MD5: ", keqing.MD5("123456"))
	fmt.Println("MD5Salt: ", keqing.MD5Salt("123456", "asd"))
}

func demo3() {
	var a int = 10
	var b int8 = 20
	var c int16 = 30
	var d float32 = 40.1234
	var e float64 = 50.123456789

	fmt.Println("int-string", keqing.Num2Str(a))
	fmt.Println("int8-string", keqing.Num2Str(b))
	fmt.Println("int16-string", keqing.Num2Str(c))
	fmt.Println("float32-string", keqing.Num2Str(d))
	fmt.Println("float64-string", keqing.Num2Str(e))

	fmt.Println("float32-string-2", keqing.Float2Str(d, 2))
	fmt.Println("float64-string-4", keqing.Float2Str(e, 4))

}

// emptyData
func demo1() {
	var a int = 1
	var a1 int = 0
	fmt.Println("int-非空....", keqing.IsEmpty(a))
	fmt.Println("int-空....", keqing.IsEmpty(a1))

	var b int8 = 1
	var b1 int8 = 0
	fmt.Println("int8-非空....", keqing.IsEmpty(b))
	fmt.Println("int8-空....", keqing.IsEmpty(b1))

	var c int16 = 1
	var c1 int16 = 0
	fmt.Println("int8-非空....", keqing.IsEmpty(c))
	fmt.Println("int8-空....", keqing.IsEmpty(c1))

	d := make(map[string]int)
	d["one"] = 1
	d["two"] = 2
	d1 := make(map[string]int)
	fmt.Println("map-非空....", keqing.IsEmpty(d))
	fmt.Println("map-空....", keqing.IsEmpty(d1))

	var e []string
	var e1 []string
	e = append(e, "1")
	e = append(e, "2")
	fmt.Println("切片-非空....", keqing.IsEmpty(e))
	fmt.Println("切片-空....", keqing.IsEmpty(e1))

	var f = "a"
	var f1 = ""
	fmt.Println("string-非空....", keqing.IsEmpty(f))
	fmt.Println("string-空....", keqing.IsEmpty(f1))

	var err = errors.New("aaaaa")

	fmt.Println("nil值的变量-非空....", keqing.IsEmpty(err))
	fmt.Println("nil值的变量-空....", keqing.IsEmpty(getNilError()))

	//var obj = dog{
	//	Age:  1,
	//	Name: "dog",
	//}
	//fmt.Println("结构体....", keqing.IsEmpty(obj))
}

func demo2() {

	arr := make([]string, 0)
	arr = append(arr, "1")
	arr = append(arr, "2")
	arr = append(arr, "3")

	a := keqing.ArrayContains(arr, "2")
	a1 := keqing.ArrayContains(arr, "4")

	fmt.Println("包含: ", a)
	fmt.Println("不包含: ", a1)
}
func getNilError() error {
	return nil
}
