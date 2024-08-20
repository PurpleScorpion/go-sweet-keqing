package demo

import (
	"errors"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func Test() {
	demo14()
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

	var dog = Dog{
		Age:  1,
		Name: "dog",
	}
	fmt.Println(keqing.ToString(dog))
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

	var dog = Dog{
		Age:  1,
		Name: "dog",
	}
	fmt.Println("结构体....", keqing.IsEmpty(dog))
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
