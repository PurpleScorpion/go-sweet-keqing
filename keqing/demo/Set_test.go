package demo

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"hash/fnv"
	"testing"
)

type User struct {
	UserName string
	Age      int
}

// 基本数据类型添加
func TestSet1(t *testing.T) {
	set := keqing.NewSet[string](nil)

	set.Add("bbb")
	set.Add("aaa")
	set.Add("ccc")
	set.Add("ddd")

	for obj := range set.Iterator() {
		fmt.Println("遍历结果: ", obj)
	}

	set.Add("aaa")

	fmt.Println(set.ToString())
}

// 基本数据类型移除
func TestSet2(t *testing.T) {
	set := keqing.NewSet[string](nil)

	set.Add("bbb")
	set.Add("aaa")
	set.Add("ccc")
	set.Add("ddd")

	set.Remove("aaa")

	fmt.Println(set.ToString())
}

// 对象类型添加
func TestSet3(t *testing.T) {
	set := keqing.NewSet[User](nil)

	set.Add(User{UserName: "bbb"})
	set.Add(User{UserName: "aaa"})
	set.Add(User{UserName: "ccc"})
	set.Add(User{UserName: "ddd"})

	set.Add(User{UserName: "bbb"})

	fmt.Println(set.ToString())
}

// 对象类型添加 - 自定义Hash函数
func TestSet4(t *testing.T) {
	set := keqing.NewSet[User](func(u User) uint64 {
		h := fnv.New64()
		_, _ = h.Write([]byte(u.UserName))
		return h.Sum64()
	})

	set.Add(User{UserName: "bbb"})
	set.Add(User{UserName: "aaa"})
	set.Add(User{UserName: "ccc"})
	set.Add(User{UserName: "ddd"})

	set.Add(User{UserName: "bbb"})

	fmt.Println(set.ToString())
}

func TestSet5(t *testing.T) {
	set := keqing.NewSet[int](nil)

	// 添加
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	fmt.Println("初始 Set:", set.ToString())

	// 删除单个
	set.Remove(2)
	fmt.Println("删除 2 后:", set.ToString())

	// 批量删除
	set.RemoveAll(1, 3)
	fmt.Println("批量删除 1 和 3 后:", set.ToString())

	// 再次添加
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(5)
	fmt.Println("再次添加后:", set.ToString())

	// 保留指定
	set.RetainAll(3, 5)
	fmt.Println("保留 3 和 5 后:", set.ToString())

	// 条件过滤
	newSet := set.Filter(func(i int) bool { return i > 3 })
	fmt.Println("原始元素:", set.ToString())
	fmt.Println("过滤大于 3 的元素后:", newSet.ToString())
}

func TestSet6(t *testing.T) {
	set := keqing.NewSet[string](nil)

	set.Add("bbb")
	set.Add("aaa")
	set.Add("ccc")
	set.Add("ddd")

	fmt.Println("aaa是否存在", set.Contains("aaa"))
	fmt.Println("vvv是否存在", set.Contains("vvv"))

	fmt.Println("集合大小:", set.Size())
}

// 求差集
func TestSet_Difference(t *testing.T) {

	set1 := keqing.NewSet[int](nil)
	set2 := keqing.NewSet[int](nil)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2.Add(2)
	set2.Add(3)
	fmt.Println("set1 - set2:", set1.Difference(set2).ToString())

}

// 求交集
func TestSet_Intersect(t *testing.T) {

	set1 := keqing.NewSet[int](nil)
	set2 := keqing.NewSet[int](nil)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	fmt.Println("set1 & set2:", set1.Intersect(set2).ToString())
}

// 求并集
func TestSet_Union(t *testing.T) {
	set1 := keqing.NewSet[int](nil)
	set2 := keqing.NewSet[int](nil)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	fmt.Println("set1 & set2:", set1.Union(set2).ToString())
}

// 测试 Sort 方法：按 age 排序
func TestSet_ToList(t *testing.T) {
	set := keqing.NewSet[Cat](nil)

	set.Add(Cat{Name: "Alice", Age: 30})
	set.Add(Cat{Name: "Bob", Age: 25})
	set.Add(Cat{Name: "Charlie", Age: 35})
	set.Add(Cat{Name: "Charlie", Age: 35})

	list := set.ToList()

	fmt.Println(list.ToString())
	fmt.Println(set.ToString())
}

// 测试 Sort 方法：按 age 排序
func TestSet_Set2Map(t *testing.T) {
	set := keqing.NewSet[Cat](nil)

	set.Add(Cat{Name: "Alice", Age: 30})
	set.Add(Cat{Name: "Bob", Age: 25})
	set.Add(Cat{Name: "Charlie", Age: 35})
	set.Add(Cat{Name: "Charlie", Age: 35})

	res := keqing.Set2Map(set, func(cat Cat) string {
		return cat.Name
	}, keqing.OverwriteFirst)

	fmt.Println(set.ToString())
	fmt.Println(keqing.ToJsonString(res))
}

func TestSet_AsSet(t *testing.T) {
	// 写一个数组
	arr := []int{1, 2, 3, 3, 4, 5}
	set := keqing.AsSet(arr, nil)
	fmt.Println(set.ToString())
}
