package demo

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"reflect"
	"testing"
)

// 测试辅助函数：断言是否 panic
func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but did not happen")
		}
	}()
	f()
}

// 测试基础 Add / Get / Size
func TestList_Add_Get_Size(t *testing.T) {
	list := keqing.NewList[string](nil)
	list.Add("a")
	list.Add("b")

	fmt.Println(list.ToString())
	fmt.Println(list.Size())
	fmt.Println(list.Get(0))
	//fmt.Println(list.Get(10))
}

// 测试 GetFirst / GetLast 空列表时 panic
func TestList_GetFirst_Empty(t *testing.T) {
	list := keqing.NewList[int](nil)
	fmt.Println(list.GetFirst())
}

func TestList_GetLast_Empty(t *testing.T) {
	list := keqing.NewList[int](nil)
	fmt.Println(list.GetLast())
}

// 测试 Contains / IndexOf / LastIndexOf
func TestList_Contains_IndexOf_LastIndexOf(t *testing.T) {
	list := keqing.NewList[string](nil)
	list.Add("apple")
	list.Add("banana")
	list.Add("apple")

	fmt.Println(list.Contains("apple"))
	fmt.Println(list.IndexOf("apple"))
	fmt.Println(list.LastIndexOf("apple"))

}

// 测试 Set 方法
func TestList_Set(t *testing.T) {
	list := keqing.NewList[string](nil)
	list.Add("a")
	list.Add("b")
	fmt.Println(list.ToString())
	list.Set(0, "x")
	fmt.Println(list.ToString())

}

// 测试 Remove / RemoveAt
func TestList_Remove_RemoveAt(t *testing.T) {
	list := keqing.NewList[string](nil)
	list.Add("a")
	list.Add("b")
	list.Add("c")
	fmt.Println(list.ToString())
	list.Remove("b")
	fmt.Println(list.ToString())
	list.RemoveAt(0)
	fmt.Println(list.ToString())
}

// 测试 RemoveAll / RetainAll
func TestList_RemoveAll_RetainAll(t *testing.T) {
	list := keqing.NewList[string](nil)
	list.Add("a")
	list.Add("b")
	list.Add("c")
	fmt.Println(list.ToString())
	list.RemoveAll("a", "c")
	fmt.Println(list.ToString())

	list.Add("a")
	list.Add("c")
	list.RetainAll("b", "c")
	fmt.Println(list.ToString())
}

// 测试 Union / Intersect / Difference
func TestList_Union_Intersect_Difference(t *testing.T) {
	a := keqing.NewList[string](nil)
	b := keqing.NewList[string](nil)

	a.Add("a")
	a.Add("b")
	b.Add("b")
	b.Add("c")

	union := a.Union(b)
	fmt.Println(union.ToString())

	intersect := a.Intersect(b)
	fmt.Println(intersect.ToString())

	diff := a.Difference(b)
	fmt.Println(diff.ToString())
}

// 测试 Filter / SubList
func TestList_Filter_SubList(t *testing.T) {
	list := keqing.NewList[int](nil)
	for i := 0; i < 5; i++ {
		list.Add(i)
	}

	even := list.Filter(func(x int) bool { return x%2 == 0 })
	fmt.Println(even.ToString())

	sub := list.SubList(1, 4)
	fmt.Println(sub.ToString())
}

// 测试 SortASC / SortDESC
func TestList_Sort(t *testing.T) {
	list := keqing.NewList[int](nil)
	list.Add(3)
	list.Add(1)
	list.Add(2)

	list.SortASC()
	fmt.Println(list.ToString())

	list.SortDESC()
	fmt.Println(list.ToString())
}

// 测试 Copy / DeepCopy
func TestList_Copy_DeepCopy(t *testing.T) {
	list := keqing.NewList[string](nil)
	list.Add("a")
	list.Add("b")

	copy1 := list.Copy()
	fmt.Println(copy1.ToString())

	deep := list.DeepCopy(func(s string) string { return s + "-copy" })
	expected := []string{"a-copy", "b-copy"}
	if !reflect.DeepEqual(deep.GetData(), expected) {
		t.Errorf("DeepCopy failed")
	}
}

// 测试 GetFirst / GetLast
func TestList_GetFirst_GetLast(t *testing.T) {
	list := keqing.NewList[int](nil)
	list.Add(1)
	list.Add(2)

	fmt.Println(list.GetFirst())
	fmt.Println(list.GetLast())

}

type Cat struct {
	Name string
	Age  int
}

type Dog struct {
	Name string
	Age  int
	Hah  []string
}

// 测试 Sort 方法：按 age 排序
func TestList_Sort_CustomStruct(t *testing.T) {
	// 方案1
	//list := keqing.NewList[*Dog](nil)

	// 方案2 使用指针类型，指针是可比较的
	list := keqing.NewList[*Dog](func(cat *Dog) uint64 {
		if cat == nil {
			return 0
		}
		return uint64(cat.Age)
	})

	list.Add(&Dog{Name: "Alice", Age: 30})
	list.Add(&Dog{Name: "Bob", Age: 25})
	list.Add(&Dog{Name: "Charlie", Age: 35})

	list.Sort(func(a, b *Dog) bool {
		return a.Age > b.Age
	})

	fmt.Println(list.ToString())
	fmt.Println(list.GetFirst())
}

// 测试 Sort 方法：按 age 排序
func TestList_ToSet(t *testing.T) {
	list := keqing.NewList[Cat](nil)

	list.Add(Cat{Name: "Alice", Age: 30})
	list.Add(Cat{Name: "Bob", Age: 25})
	list.Add(Cat{Name: "Charlie", Age: 35})
	list.Add(Cat{Name: "Charlie", Age: 35})

	set := list.ToSet()

	fmt.Println(list.ToString())
	fmt.Println(set.ToString())
}

// 测试 Sort 方法：按 age 排序
func TestList_List2Map(t *testing.T) {
	list := keqing.NewList[Cat](nil)

	list.Add(Cat{Name: "Alice", Age: 30})
	list.Add(Cat{Name: "Bob", Age: 25})
	list.Add(Cat{Name: "Charlie", Age: 35})
	list.Add(Cat{Name: "Charlie", Age: 36})

	res := keqing.List2Map(list, func(cat Cat) string {
		return cat.Name
	}, keqing.OverwriteFirst)

	fmt.Println(list.ToString())
	fmt.Println(keqing.ToJsonString(res))
}

func TestList_AsList(t *testing.T) {
	// 写一个数组
	arr := []int{1, 2, 3, 4, 5}
	list := keqing.AsList(arr, nil)
	fmt.Println(list.ToString())
}
