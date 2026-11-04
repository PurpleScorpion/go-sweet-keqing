# go-sweet-keqing

适用于 Go 的刻晴版工具库，提供集合、切片、日期、字符串、金额、配置取值、文件和 RSA 等工具。要求 **Go 1.20 及以上**，仅依赖 Go 标准库。

> 剑光如我，斩尽牛杂！

## 安装与快速开始

```sh
go get github.com/PurpleScorpion/go-sweet-keqing
```

包路径为 `github.com/PurpleScorpion/go-sweet-keqing/keqing`。本文对应当前分支源码；使用已发布版本时，请查看对应 tag 的 README。

下文的 Go 代码块都是独立程序，可分别保存为 `main.go`，在已引入本库的 Go 模块中运行 `go run .`。

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	list := keqing.AsList([]int{3, 1}, nil) // nil 使用默认 hash
	list.Add(2)                             // 无返回值
	list.Set(0, 4)                          // 无返回值
	copied := list.DeepCopy(func(v int) int { return v })
	values := copied.ToArray()
	keqing.Shuffle(values) // 原地打乱切片，无返回值
	keqing.SortArrayDesc(values)
	keqing.SortArray(values)
	fmt.Println(values)

	set := keqing.AsSet([]int{1, 1, 2}, nil)
	fmt.Println(set.DeepCopy(func(v int) int { return v }).ToArray())
	fmt.Println(keqing.ToArray[int]([]int{1, 2}))
	fmt.Println(keqing.ToObject[string]("keqing"))
	date := keqing.ParseDate("2024-02-29", "2006-01-02", nil)
	fmt.Println(date.Format("2006-01-02 MST"))
	// Output:
	// [1 2 4]
	// [1 2]
	// [1 2]
	// keqing
	// 2024-02-29 UTC
}
```

## 功能导航

- [List、Set 与集合工具](#listset-与集合工具)
- [切片工具](#切片工具)
- [日期与时间](#日期与时间)
- [空值判断](#空值判断)
- [字符串与数值解析](#字符串与数值解析)
- [类型断言与属性复制](#类型断言与属性复制)
- [YAML 配置取值](#yaml-配置取值)
- [金额计算](#金额计算)
- [文件工具](#文件工具)
- [MD5 摘要](#md5-摘要)
- [RSA 加解密](#rsa-加解密)
- [随机数与标识符](#随机数与标识符)

本文中“原地”表示修改传入的数据或接收者，“浅拷贝”表示复制容器但共享元素中的指针、切片、map 等引用。`panic` 与返回 `error` 是不同的失败方式，请按各接口的说明处理。

## List、Set 与集合工具

`List[T comparable]` 保持插入顺序，允许重复值；`Set[T comparable]` 保持首次插入顺序，按元素相等规则去重。`T` 必须满足 Go 的 `comparable` 约束，例如整数、字符串、可比较结构体或指针；包含 slice/map 字段的结构体可使用指针作为元素类型。

使用构造函数初始化集合，`hashFunc` 传 `nil` 使用默认 hash。不要直接使用零值集合，也不要按值复制已使用的集合结构体；需要副本时调用 `Copy` 或 `DeepCopy`。

| 创建函数 | 返回值与行为 |
| --- | --- |
| `NewList[T comparable](hashFunc func(T) uint64) *List[T]` | 创建空列表 |
| `NewSet[T comparable](hashFunc func(T) uint64) *Set[T]` | 创建空集合 |
| `AsList[T comparable](arr []T, hashFunc func(T) uint64) *List[T]` | 按输入顺序创建列表，保留重复值 |
| `AsSet[T comparable](arr []T, hashFunc func(T) uint64) *Set[T]` | 按输入顺序创建集合，保留每个相等值的第一次出现 |

`AsList` / `AsSet` 不共享输入切片的底层数组，但元素中的引用仍然共享。nil 输入切片会创建空集合。

### 读取与修改

下列方法在 List 和 Set 上都可使用；修改方法直接修改接收者。

| 方法 | 行为 |
| --- | --- |
| `Contains(value T) bool` | 判断是否存在相等元素 |
| `Remove(value T) bool` | 删除第一个匹配项，未找到返回 `false` |
| `RemoveAll(items ...T)` | 删除所有匹配项，保持剩余元素的相对顺序；不传参数不删除元素 |
| `RetainAll(items ...T)` | 仅保留匹配项及其原有相对顺序；不传参数会清空 |
| `GetData() []T` / `ToArray() []T` | 返回切片浅拷贝；修改切片位置不会修改集合；空集合返回 `nil` |
| `Size() int` / `IsEmpty() bool` | 返回元素数量 / 是否为空 |
| `Clear()` | 清空接收者 |
| `Iterator() <-chan T` | 返回调用时的浅快照通道，按集合顺序读取；可提前 `break`，可在遍历中修改原集合 |
| `ToString() string` | 对元素切片进行 JSON 编码；空集合返回 `"null"`，编码失败返回 `""` |

List 的其他方法：

| 方法 | 行为 |
| --- | --- |
| `Add(value T)` | 追加元素，无返回值 |
| `Get(index int) T` / `Set(index int, value T)` | 从 0 开始读取 / 替换；越界 `panic` |
| `GetFirst() T` / `GetLast() T` | 读取首尾元素；空列表 `panic` |
| `IndexOf(value T) int` / `LastIndexOf(value T) int` | 返回第一 / 最后匹配位置；不存在为 `-1` |
| `RemoveAt(index int) bool` | 删除指定位置，越界返回 `false` |
| `SubList(indexs ...int) *List[T]` | `SubList(start)` 取到结尾；`SubList(start, end)` 取 `[start,end)`，返回独立浅拷贝 |
| `Sort(less func(a, b T) bool)` | 原地稳定排序，`less(a,b)` 表示 a 应排在 b 前面 |
| `SortASC()` / `SortDESC()` | 原地稳定升序 / 降序；支持具体类型为有符号整数、无符号整数（不含 `uintptr`）、浮点数、字符串及相应命名类型；其他类型（包括 `any`）`panic` |
| `ToSet() *Set[T]` | 创建新 Set，按首次出现顺序去重 |

`SubList` 必须传 1 或 2 个参数，并满足 `0 <= start < Size()`、`start <= end <= Size()`，否则 `panic`。因此空列表不可调用，`SubList(Size())` 也不合法。

Set 的其他方法：

| 方法 | 行为 |
| --- | --- |
| `Add(value T) bool` | 插入新值返回 `true`；相等值已存在时保留原值，返回 `false` |
| `GetAll() []T` | 旧的浅快照方法，已弃用；使用 `GetData` / `ToArray` |
| `ToList() *List[T]` | 创建新 List，保持元素顺序 |

### 集合运算、过滤和复制

以下方法返回新集合，不修改输入。`other` 必须为同类型的非 nil 集合指针。

| List 方法 | Set 方法 | 行为 |
| --- | --- | --- |
| `Union(other *List[T]) *List[T]` | `Union(other *Set[T]) *Set[T]` | 左侧元素在前，依次追加右侧未出现的值；List 保留左侧原有重复值 |
| `Intersect(other *List[T]) *List[T]` | `Intersect(other *Set[T]) *Set[T]` | 保留左侧中也存在于右侧的元素及顺序；List 保留左侧重复值 |
| `Difference(other *List[T]) *List[T]` | `Difference(other *Set[T]) *Set[T]` | 保留左侧中不存在于右侧的元素及顺序；List 保留左侧重复值 |
| `Filter(predicate func(T) bool) *List[T]` | `Filter(predicate func(T) bool) *Set[T]` | 基于调用时的浅快照，保留 predicate 返回 `true` 的元素 |
| `Copy() *List[T]` | `Copy() *Set[T]` | 创建浅拷贝 |
| `DeepCopy(cloneFunc func(T) T) *List[T]` | `DeepCopy(cloneFunc func(T) T) *Set[T]` | 逐项调用 cloneFunc 后构造新集合；新 Set 仍按相等规则去重 |

浅拷贝只隔离集合结构，元素里的指针、slice、map 等引用仍指向原对象。`DeepCopy` 的复制深度由 `cloneFunc` 决定，不会自动递归复制，也不会自动调用 `Clone()`。

### 相等规则、hash 和并发

hash 用于查找候选元素，相同 hash 不代表元素相等。元素实现 `Equal(T) bool` 时使用该方法，否则可比较值按 Go 的 `==` 比较，指针默认比较地址。`T=any` 的动态值无法直接比较时使用 `reflect.DeepEqual`；nil 值不会调用其 `Equal` 方法。自定义 `Equal` 应符合自反、对称、传递的相等关系。

自定义 hash 应保持确定性，并使相等元素获得相同 hash。实现 `Hash() uint64` 不会被自动调用，需要显式传入 hash 函数。复制、过滤及 List/Set 转换会保留原集合的 hash 策略，二元集合运算的结果沿用左侧策略。普通值使用 hash 桶查找，自定义 Equal 或动态接口值等复杂情况可能线性扫描，批量操作最坏可达到 O(n²)。

建议元素留在集合期间保持参与 Equal/hash 的字段稳定。Set 不会监测外部字段修改并自动重新去重；需要修改相等键时，应先删除该元素，修改后再添加，避免留下多个相等项或在后续删除时命中其他项。

集合方法对内部结构加锁，但不会保护元素引用对象的字段；多个方法组成的操作不保证整体原子性，两个集合的运算分别读取快照。hash、Equal 和排序比较回调应无副作用，不要在回调中再次操作当前集合，以免锁重入。`Filter`、`DeepCopy` 及下面的 map 转换回调基于快照执行；传入会被实际调用的 nil 回调会 `panic`。

### List/Set 转 map

| 函数 | 返回值与行为 |
| --- | --- |
| `List2Map[T comparable, K comparable](list *List[T], keyFunc func(T) K, onConflict int) map[K]T` | 按 List 顺序计算 key，创建新 map |
| `Set2Map[T comparable, K comparable](set *Set[T], keyFunc func(T) K, onConflict int) map[K]T` | 按 Set 顺序计算 key，创建新 map |

输入集合必须非 nil；转换不会修改集合，map 的 value 中的引用仍共享。重复 key 按以下策略处理：

| `onConflict` | 行为 |
| --- | --- |
| `OverwriteNone` | 遇到重复 key 时 `panic` |
| `OverwriteFirst` | 保留第一个值 |
| `OverwriteLast` | 保留最后一个值；其他未定义的整数也按此策略处理 |

map 键使用 Go 原生相等规则，不使用集合的 `Equal` 或 hash；map 遍历顺序没有保证。当 `K=any` 时，动态键也必须可比较，否则 `panic`。

下面的完整示例展示基本操作、集合运算、自定义相等规则、显式克隆和 map 冲突处理：

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

type User struct {
	ID   int
	Name string
	Tags []string
}

func (u *User) Equal(other *User) bool { return u.ID == other.ID }

func main() {
	list := keqing.AsList([]int{3, 1, 3}, nil)
	list.Add(2)
	fmt.Println(list.IndexOf(3), list.LastIndexOf(3))
	list.Remove(3) // 只删除第一个 3
	list.SortASC()
	snapshot := list.GetData()
	snapshot[0] = 99
	fmt.Println(list.GetData()) // 修改快照不影响列表

	set := keqing.AsSet([]int{3, 1, 3}, nil)
	fmt.Println(set.Add(1), set.Add(2))
	for value := range set.Iterator() {
		if value == 3 {
			set.Remove(value)
		}
	}
	fmt.Println(set.GetData())

	left := keqing.AsList([]int{1, 1, 2}, nil)
	right := keqing.AsList([]int{1, 3, 3}, nil)
	fmt.Println(left.Union(right).GetData(), left.Intersect(right).GetData())

	users := keqing.NewSet[*User](func(u *User) uint64 {
		if u == nil {
			return 0
		}
		return uint64(u.ID)
	})
	alice := &User{ID: 1, Name: "Alice", Tags: []string{"go"}}
	fmt.Println(users.Add(alice), users.Add(&User{ID: 1, Name: "Alicia"}))
	cloned := users.DeepCopy(func(u *User) *User {
		if u == nil {
			return nil
		}
		v := *u
		v.Tags = append([]string(nil), u.Tags...)
		return &v
	})
	cloned.GetData()[0].Tags[0] = "java"
	fmt.Println(alice.Tags, cloned.GetData()[0].Tags)

	records := keqing.AsList([]*User{alice, {ID: 1, Name: "Alicia"}}, nil)
	byID := keqing.List2Map(records, func(u *User) int { return u.ID }, keqing.OverwriteLast)
	fmt.Println(byID[1].Name)
	// Output:
	// 0 2
	// [1 2 3]
	// false true
	// [1 2]
	// [1 1 2 3] [1 1]
	// true false
	// [go] [java]
	// Alicia
}
```

## 切片工具

名称含 `Array` 的接口主要操作 Go 切片 `[]T`，不是固定长度数组 `[N]T`。

| 函数 | 返回值、修改行为与边界 |
| --- | --- |
| `ArrayContains(obj interface{}, value interface{}) bool` | 仅接受切片，用 `reflect.DeepEqual` 逐项比较；nil、非切片、数组、切片指针均返回 `false` |
| `ArrayMap[S any, T any](slice []S, fn func(S) T) []T` | 按顺序执行 fn，返回新切片；不修改输入，nil 输入返回 nil；实际调用 nil fn 时 `panic` |
| `GetOne[T comparable](arr []T) (T, bool)` | 取第一个元素；空/nil 切片返回 T 零值及 `false` |
| `Shuffle[T any](slice []T)` | 原地随机打乱，无返回值；影响共享底层数组的切片，空/nil 切片不变 |
| `Reverse[T any](s []T) []T` | 返回倒序排列的新切片，不修改 s；nil 输入返回非 nil 空切片 |
| `Difference[T comparable](a, b []T) []T` | a 中不存在于 b 的值；去重，结果顺序不保证 |
| `Union[T comparable](a, b []T) []T` | a、b 中所有值；去重，结果顺序不保证 |
| `Intersection[T comparable](a, b []T) []T` | 同时出现在 a、b 中的值；去重，按 b 中首次匹配的顺序返回 |
| `SortArray(obj interface{})` / `SortArrayDesc(obj interface{})` | 原地升序 / 降序，无返回值；接受切片或非 nil 切片指针，类型正确的 nil 切片也可使用 |

切片的交、并、差集不修改输入，无结果时返回 nil；它们使用 Go map 键规则，不调用 `Equal`。`T` 必须满足 `comparable`，当 `T=any` 时动态元素也必须可比较，否则 `panic`。`Reverse` 等新切片中的元素引用仍然共享；`ArrayMap` 如何转换、复制元素由 fn 决定。

`SortArray` / `SortArrayDesc` 仅支持元素底层类型为 `int`、`int32`、`int64`、`float32`、`float64`、`string`，包括对应的命名元素类型和命名切片类型。其他输入 `panic`，包括数组、nil 接口、nil 切片指针，以及 `int8`、`int16`、所有 `uint`、`uintptr`、`bool`、`any`、结构体元素类型。此范围与 List 的 `SortASC` / `SortDESC` 不同。

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	values := []int{3, 1, 2}
	reversed := keqing.Reverse(values)
	keqing.SortArray(values)
	fmt.Println(values, reversed)

	mapped := keqing.ArrayMap([]int{1, 2}, func(v int) string {
		return fmt.Sprintf("id-%d", v)
	})
	fmt.Println(mapped)
	fmt.Println(keqing.GetOne(values))

	a, b := []int{1, 1, 2, 3}, []int{3, 1, 3, 4}
	fmt.Println(keqing.Intersection(a, b))
	diff, union := keqing.Difference(a, b), keqing.Union(a, b)
	keqing.SortArray(diff) // 需要确定展示顺序时显式排序
	keqing.SortArray(union)
	fmt.Println(diff, union)
	// Output:
	// [1 2 3] [2 1 3]
	// [id-1 id-2]
	// 1 true
	// [3 1]
	// [2] [1 2 3 4]
}
```

## 日期与时间

格式参数使用 **Go 的参考时间布局**，例如 `"2006-01-02 15:04:05"`，不是 `yyyy-MM-dd`。解析失败会 panic；需要返回错误的解析流程可使用标准库 `time.Parse` / `time.ParseInLocation`。

| 函数 | 说明 |
| --- | --- |
| `NowDate() time.Time` / `NowUTCDate() time.Time` | 当前本地时间 / 当前 UTC 时间，均包含时分秒，不是当天零点 |
| `NowDateStr() string` / `NowUTCDateStr() string` | 使用默认本地 / UTC 布局格式化当前时间 |
| `DateInfo(date time.Time) TimeVO` | 在 date 自带时区中拆分年月日、时分秒和毫秒 |
| `FormatDate(date time.Time, format string) string` | 按布局格式化，保留 date 的时区，不自行转 UTC |
| `ParseDate(dateStr, format string, location *time.Location) time.Time` | 按给定布局和时区解析，`location == nil` 使用 UTC |
| `ParseLocal(local string) time.Time` / `ParseUTC(utc string) time.Time` | 使用默认本地 / UTC 布局和对应时区解析 |
| `UTC2Local(utc string) string` / `Local2UTC(local string) string` | 按默认布局解析、转换时区并输出字符串 |
| `UTC2LocalCustom(utc, utcFormat, localFormat string) string` | 自定义输入、输出布局，UTC 转本地 |
| `Local2UTCCustom(local, localFormat, utcFormat string) string` | 自定义输入、输出布局，本地转 UTC |
| `DateAddSecond(date time.Time, seconds int64) time.Time` / `DateSubSecond(date time.Time, seconds int64) time.Time` | 增加 / 减少固定秒数 |
| `DateAddMinute(date time.Time, minute int64) time.Time` / `DateSubMinute(date time.Time, minute int64) time.Time` | 增加 / 减少固定分钟数 |
| `DateAddHour(date time.Time, hour int64) time.Time` / `DateSubHour(date time.Time, hour int64) time.Time` | 增加 / 减少固定小时数 |
| `DateAddDay(date time.Time, day int64) time.Time` / `DateSubDay(date time.Time, day int64) time.Time` | 在 date 所在时区按日历天增减，夏令时切换时一天可能不是 24 小时 |
| `CurrentTimeMillis() int64` / `CurrentTimeSeconds() int64` | 当前 Unix 时间戳，分别以毫秒 / 秒为单位 |
| `GetDaysInMonth(data any, years ...int) int` | 返回公历指定月份天数，详见下文 |

`ParseDate` 的 location 用于解释未写时区的时间；输入与布局含显式时区偏移时，按输入偏移解析，不会强制转到 location。日期偏移允许负数，不可表示或溢出的偏移会 panic。`TimeVO` 包含 `Year/Month/Day/Hour/Minute/Second/Millisecond` 整数字段及对应的 `...Str` 字段；月份到秒的字符串补足两位，毫秒补足三位。

默认布局变量：

| 名称 | 初始值 |
| --- | --- |
| `DEFAULT_LOCAL_FORMAT` | `2006-01-02 15:04:05` |
| `LOCAL_FORMAT_DATE` / `LOCAL_FORMAT_TIME` | `2006-01-02` / `15:04:05` |
| `DEFAULT_UTC_FORMAT` | `2006-01-02T15:04:05.999999Z`，小数秒最多六位，省略尾随零 |
| `UTC_FORMAT_SHORT` | `2006-01-02T15:04:05Z` |

`UTC`、`LOCAL` 初始分别为 `time.UTC`、`time.Local`，供相应的解析函数使用；本地输出转换使用 `time.Time.Local()`，当天/当月范围使用系统本地时区。设置 `LOCAL` 不会统一改变所有“本地”接口的行为。可用 `ParseDate` 和 `time.Time.In` 显式指定业务时区，避免在并发运行时修改包级变量。默认 UTC 布局里的 `Z` 是字面量，不会自动将任意时间转换到 UTC，也不匹配 `+08:00` 这样的输入。

`GetDaysInMonth` 接受 `time.Time`、`time.Month`、内置 `int/int16/int32/int64`：传月份时必须同时传年份；传 `time.Time` 时默认取其年月，也可以用第一个额外年份覆盖。月份必须为 1～12；缺年份、类型不支持或月份越界会 panic。额外年份超过一个时只使用第一个。

### 今天、昨天与当月范围

| 本地边界函数 | 同一边界转 UTC 的函数 |
| --- | --- |
| `GetStartAndEnd4Today()` | `GetStartAndEndUTC4Today()` |
| `GetStartAndEnd4Yesterday()` | `GetStartAndEndUTC4Yesterday()` |
| `GetStartAndEnd4CurrentMonth()` | `GetStartAndEndUTC4CurrentMonth()` |

以上函数都无参数，返回 `(start time.Time, end time.Time)` **闭区间**，终点为下一日/下一月开始前的最后一纳秒。UTC 版本是将本地日历边界转换到 UTC，并非重新计算 UTC 的今天或当月。

时区跳变使午夜不存在时取首个有效时刻，午夜重复时取第一次出现。若昨天整个日期被时区调整跳过，返回空区间：`end.Before(start)` 且 `end == start.Add(-time.Nanosecond)`。查询前应检查区间是否为空；采用右开区间的存储查询可在非空时使用 `[start, end.Add(time.Nanosecond))`。

```go
package main

import (
	"fmt"
	"time"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	zone := time.FixedZone("UTC+8", 8*60*60)
	date := keqing.ParseDate("2024-02-29 12:30:00", keqing.DEFAULT_LOCAL_FORMAT, zone)
	fmt.Println(keqing.FormatDate(date.UTC(), time.RFC3339))
	fmt.Println(keqing.DateAddDay(date, 1).Format("2006-01-02"))
	fmt.Println(keqing.GetDaysInMonth(time.February, 2024))
	fmt.Println(keqing.DateInfo(date).MonthStr)
	// Output:
	// 2024-02-29T04:30:00Z
	// 2024-03-01
	// 29
	// 02
}
```

## 空值判断

这组函数按库内约定判断“空”，其中 bool、指针和 struct 的行为需要特别留意。

| 函数 | 作用 |
| --- | --- |
| `IsEmpty(obj interface{}) bool` | 按下表判断输入是否为空。 |
| `IsNotEmpty(obj interface{}) bool` | 对 `IsEmpty` 的结果取反；同样不支持 struct。 |

| 输入 | `IsEmpty` 的结果 |
| --- | --- |
| nil interface | `true`。 |
| 内置 string | 空串或忽略大小写后为 `"null"` 时返回 `true`；不去空白，`" "`、`" null "` 返回 `false`。 |
| 内置 `int`、`int8/16/32/64`、`uint`、`uint8/16/32/64`、`float32/64` | 等于零时返回 `true`。 |
| bool | **`true` 返回 `true`，`false` 返回 `false`**。 |
| slice / map | 长度为零时返回 `true`，包括 nil slice/map。 |
| struct | panic。 |
| 指针、数组和其他未列类型 | `false`，包括 nil 指针和零长度数组。 |

自定义具名 string、数字和 bool 不匹配上述内置类型分支；自定义 slice/map 仍按长度判断。
需要判断 bool 状态或指针是否为 nil 时，直接使用对应条件更清楚。

## 字符串与数值解析

字符串函数返回处理结果，不修改传入字符串。
查找位置使用 UTF-8 **字节下标**；按 Unicode 码点拆分时使用 `ToStringArray`。

### 查找、替换与拆分

| 函数 | 返回值与边界行为 |
| --- | --- |
| `StartWith(str, prefix string) bool` | 判断前缀；空前缀也匹配。 |
| `EndWith(str, suffix string) bool` | 判断后缀；空后缀也匹配。 |
| `Contains(mainStr, substr string) bool` | 判断是否包含子串；空子串也匹配。 |
| `IndexOf(mainStr, substr string) int` | 第一次匹配的字节下标，未找到返回 `-1`；空子串返回 `0`。 |
| `LastIndexOf(mainStr, substr string) int` | 最后一次匹配的字节下标，未找到返回 `-1`；空子串返回 `len(mainStr)`。 |
| `Replace(mainStr, old, new string) string` | 替换全部字面量匹配。`old == ""` 时在开头和每个 UTF-8 序列之后插入 `new`。 |
| `Replace2Regex(mainStr, pattern, new string) string` | 按 Go 正则替换全部匹配；非法正则会 panic。替换串支持 `$1`、`${name}` 捕获组展开，字面量 `$` 写成 `$$`。 |
| `GetBytes(str string) []byte` | 返回字节切片，中文通常占多个 UTF-8 字节。 |
| `ToStringArray(str string) []string` | 按 rune 拆分，不解析 JSON；空串返回 nil 切片。组合字符和 emoji 序列可能被拆成多个元素。 |
| `Split(str, sep string) []string` | 按分隔符拆分，无匹配时返回只含原串的切片。空串始终返回 `[]string{""}`；非空串搭配空分隔符时按 UTF-8 序列拆分。 |

### 大小写、空白与命名转换

| 函数 | 返回值与边界行为 |
| --- | --- |
| `ToUpperCase(s string) string` | 按 Unicode 规则将全部字符串转为大写。 |
| `ToLowerCase(s string) string` | 按 Unicode 规则将全部字符串转为小写。 |
| `ToFirstUpperCase(s string) string` | 只将第一个 rune 转为大写，且它必须是字母；空串原样返回。 |
| `ToFirstLowerCase(s string) string` | 只将第一个 rune 转为小写，且它必须是字母；空串原样返回。 |
| `Trim(str string) string` | 删除两端 Unicode 空白，保留中间空白。 |
| `TrimSuffix(s, suffix string) string` | 匹配时仅删除一次整个后缀，不会反复删除。 |
| `Snake2Camel(s string) string` | 按 `_` 分段，第一段保留原样，后续各段用 `strings.Title` 转换后连接。 |
| `Snake2BigCamel(s string) string` | 按 `_` 分段，各段用 `strings.Title` 转换后连接。 |
| `Camel2Snake(s string) string` | 在非首位大写字母前加 `_` 并转小写；`HTTPServer` 会得到 `h_t_t_p_server`。 |
| `Camel2Hyphen(s string) string` | 先执行 `Camel2Snake`，再将所有 `_` 改成 `-`。 |
| `Camel2WebURL(s string) string` | 先将首字母转小写，再执行 `Camel2Snake` 并将 `_` 改成 `-`；不进行 URL 编码或校验。 |

首字符大小写函数按 rune 处理，支持 `élan → Élan`，但不会跳过开头的空格或标点去寻找字母。
命名转换适合 `user_name`、`userName` 等常见标识符：不会统一既有大小写，也不会识别连续大写缩写，不保证对任意 Unicode 文本做一致的命名规范化。
`Split` 和 `Camel2Snake` 把 `"null"` 当作普通文本，不采用 `IsEmpty` 的空值约定。

### 数字与文本转换

| 函数 | 输入、返回与失败行为 |
| --- | --- |
| `Num2Str(obj interface{}) string` | 支持内置有符号/无符号整数和 `float32/float64`；整数输出十进制，浮点输出不使用指数形式的最短可往返小数表示。nil 或不支持的类型返回 `""`。 |
| `Float2Str(obj interface{}, precision int) string` | 只支持内置 `float32/float64`。非负 precision 表示固定小数位并舍入，负数等同 `Num2Str`；nil、整数等不支持类型返回 `""`。 |
| `ParseInt8(str string) int8` | 十进制整数，范围 `[-128, 127]`。 |
| `ParseInt16(str string) int16` | 十进制整数，范围 `[-32768, 32767]`。 |
| `ParseInt32(str string) int32` | 十进制整数，范围 `[-2147483648, 2147483647]`。 |
| `ParseInt64(str string) int64` | 十进制整数，范围 `[-9223372036854775808, 9223372036854775807]`。 |
| `ParseInt(str string) int` | 十进制整数，按运行平台的 32/64 位 int 范围检查。 |
| `ParseFloat32(str string) float32` | 按 `strconv.ParseFloat` 语法解析并舍入到 float32。 |
| `ParseFloat64(str string) float64` | 按 `strconv.ParseFloat` 语法解析并舍入到 float64。 |

整数解析允许正负号，不自动去空白，不接受小数、科学计数法或 `0x` 前缀；无效输入或越界会 panic。
浮点解析的无效输入和有限数值溢出也会 panic；极小数可能舍入为零，显式 `NaN` / `Inf` 受支持。
`Num2Str` 支持的整数为 `int/int8/int16/int32/int64` 和 `uint/uint8/uint16/uint32/uint64`，不包括 `uintptr` 或自定义具名数字类型。
浮点格式化基于实际二进制浮点值；需要金额计算时参照 MoneyUtils 的分精度约定。

### 调试字符串与 JSON

| 函数 | 输入、返回与失败行为 |
| --- | --- |
| `ToString(obj interface{}) string` | 自定义格式化 slice、map、struct；跳过 struct 的不可导出字段。nil 返回 `"null"`，nil slice/map 返回 `"[]"` / `"{}"`；顶层标量、数组和指针等不支持类型返回 `""`。 |
| `ToJsonString(data any) string` | 调用 `json.Marshal`；成功返回 JSON 字符串，Marshal 返回错误时返回 `""`。需要获取具体错误时直接使用标准库。 |

`ToString` 用于调试展示：map 顺序不保证稳定，键/字段名和字符串转义不遵循 JSON 规则，也不使用 JSON tag 或检测递归引用。
交换或保存 JSON 数据时使用 `ToJsonString` 或标准库 `json.Marshal`。

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	fmt.Println(keqing.IndexOf("刻晴", "晴")) // 字节下标
	fmt.Printf("%q\n", keqing.ToStringArray("刻晴"))
	fmt.Println(keqing.ToFirstUpperCase("élan"))
	fmt.Println(keqing.Snake2Camel("user_name"))
	fmt.Println(keqing.Camel2Snake("HTTPServer"))
	fmt.Println(keqing.Replace2Regex("user-42", `([a-z]+)-(\d+)`, "${1}:${2}"))
	fmt.Printf("%q\n", keqing.Split("", ","))
	fmt.Println(keqing.ParseInt8("127")) // "128" 会 panic
	fmt.Println(keqing.Float2Str(12.5, 2))
	fmt.Println(keqing.ToString([]any{"刻晴", nil, 42}))
	fmt.Printf("%q\n", keqing.ToString("刻晴")) // 顶层字符串不受支持
	fmt.Println(keqing.ToJsonString(map[string]any{"name": "刻晴"}))

	// Output:
	// 3
	// ["刻" "晴"]
	// Élan
	// userName
	// h_t_t_p_server
	// user:42
	// [""]
	// 127
	// 12.50
	// ["刻晴", null, 42]
	// ""
	// {"name":"刻晴"}
}
```

## 类型断言与属性复制

`ToArray` 和 `ToObject` 主要做类型断言，不是通用类型转换：不会解析 JSON、将 map 映射为结构体或将 `int` 自动转换成 `int64`。

| 函数 | 输入、返回值与边界 |
| --- | --- |
| `ToObject[T any](data interface{}) T` | 尝试 `data.(T)`，失败或输入 nil 返回 T 的零值；不修改输入。需要区分失败与合法零值时直接使用 `v, ok := data.(T)` |
| `ToArray[T any](data interface{}) []T` | 支持 `[]T`、`*[]T`、`[]*T`。前两者返回原切片，共享底层数组；`[]*T` 新建切片逐项解引用，nil 元素变为 T 零值。nil 接口、nil 切片指针、类型不匹配返回 nil |
| `CopyProperties(dst, src interface{}, ignoreFields ...string) error` | dst 必须是非 nil 结构体指针，src 为结构体或非 nil 结构体指针；将同名、同类型且可设置的字段浅赋值给 dst，忽略指定名称的字段，成功返回 nil |

`ToArray` 使用准确的 Go 类型断言：定义 `type Ints []int` 后，`ToArray[int](Ints{1,2})` 不会自动转成 `[]int`。nil `[]T` 返回 nil；nil `[]*T` 返回非 nil 空切片。解引用只是复制 T 的值，不会递归复制其内部引用。

`CopyProperties` 修改 dst，不做数值转换或递归深拷贝。未导出、目标缺失或类型不同的字段会跳过，只处理源结构体直接声明的字段。非法 dst/src，或目标字段位于无法到达的 nil 嵌入指针中时返回 error；返回 error 前可能已复制部分字段。`ignoreFields` 按字段名精确匹配。

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

type Source struct {
	Name string
	Age  int
	Tags []string
}
type Target struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	data := []int{1, 2}
	same := keqing.ToArray[int](data)
	same[0] = 9
	fmt.Println(data)                           // 共享底层数组
	fmt.Println(keqing.ToObject[int64](int(7))) // 不执行数值转换
	fmt.Println(keqing.ToArray[int]([]int64{1, 2}) == nil)

	number := 7
	fmt.Println(keqing.ToArray[int]([]*int{&number, nil}))

	src := Source{Name: "Alice", Age: 18, Tags: []string{"go"}}
	dst := Target{Age: 20}
	if err := keqing.CopyProperties(&dst, src, "Age"); err != nil {
		panic(err)
	}
	fmt.Println(dst.Name, dst.Age)
	dst.Tags[0] = "java"
	fmt.Println(src.Tags) // slice 字段仍共享底层数组
	// Output:
	// [9 2]
	// 0
	// true
	// [7 0]
	// Alice 20
	// [java]
}
```

## YAML 配置取值

这组函数读取已经解析好的 map，本身不读取文件或解析 YAML。
配置中的嵌套对象需要是 `map[string]interface{}`，不支持 `map[interface{}]interface{}`。

| 函数 | 作用 |
| --- | --- |
| `LoadYml(parent, child map[string]interface{}) error` | 设置基础配置和覆盖配置；两者都不能是 nil，但可以是空 map。nil 输入返回 error，保留之前成功加载的配置。 |

取值键必须写成 `${server.port}`，否则 panic；`.` 表示逐层查找，不支持列表下标或转义键名中的点。
先查 child 的完整路径，值非 nil 就使用它，缺失或 nil 时再查 parent；`false`、`0`、`""`、空列表/空 map 都是有效覆盖值。
路径中间遇到标量或列表视为路径不存在。这里的 nil 指 interface 本身为 nil，类型化 nil 值不会触发回退。

| 函数 | 支持的底层值与转换 | 缺失或最终为 nil |
| --- | --- | --- |
| `ValueObject(key string) interface{}` | 任意值，直接返回。 | nil |
| `ValueInt(key string) int` | 有符号整数/浮点数，按平台 int 位宽检查。 | `0` |
| `ValueInt32(key string) int32` | 有符号整数/浮点数，按 int32 范围检查。 | `0` |
| `ValueInt64(key string) int64` | 有符号整数/浮点数，按 int64 范围检查。 | `0` |
| `ValueFloat32(key string) float32` | 有符号整数/浮点数；有限值转为无穷时 panic，显式 NaN/Inf 保留。 | `0` |
| `ValueFloat64(key string) float64` | 有符号整数/浮点数；转换可能损失精度，NaN/Inf 保留。 | `0` |
| `ValueBool(key string) bool` | bool，不解析文本 `"true"`。 | `false` |
| `ValueString(key string) string` | string 或有符号整数/浮点数，数字通过 `Num2Str` 转换。 | `""` |
| `ValueStringArr(key string) []string` | **`[]interface{}`，且每项都是 string**；不接受 `[]string`。返回新切片，空源列表返回 nil 切片。 | 空的非 nil 切片 |

表中数值类型限于内置 `int/int8/int16/int32/int64` 和 `float32/float64`，不包含无符号整数、自定义具名类型或数字字符串。
整数取值先把小数向零截断，再检查范围；溢出、NaN、Inf 会 panic。浮点转换可能舍入或损失整数精度。
除 `ValueObject` 外，非 nil 值的类型不符合上表时会 panic；需要区分缺失和零值时先检查 `ValueObject`。

`LoadYml` 保存 map 引用，没有复制、深合并或并发同步；建议启动阶段加载后只读，并发重载或修改需要调用方同步。
`ValueObject` 也返回底层值，修改其中的 map/slice 会影响之后的读取。
没有环境覆盖时，传入空 map 作为 child 即可。下面演示覆盖、回退和字符串列表读取：

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	parent := map[string]any{
		"server": map[string]any{"port": 8080, "enabled": true, "host": "localhost"},
		"tags":   []any{"go", "keqing"},
	}
	child := map[string]any{
		"server": map[string]any{"port": 9090, "enabled": false, "host": nil},
	}
	if err := keqing.LoadYml(parent, child); err != nil {
		panic(err)
	}
	fmt.Println(keqing.ValueInt("${server.port}"))
	fmt.Println(keqing.ValueBool("${server.enabled}")) // false 也覆盖父配置
	fmt.Println(keqing.ValueString("${server.host}"))  // nil 回退至父配置
	fmt.Println(keqing.ValueStringArr("${tags}"))
	fmt.Println(keqing.ValueObject("${missing}") == nil)

	// Output:
	// 9090
	// false
	// localhost
	// [go keqing]
	// true
}
```

## 金额计算

金额输入以“元”为单位，使用十进制字符串；计算精度为“分”，不足一分的部分**向零截断，不四舍五入**。
内部按十进制有理数和整数分计算，避免先转 float64 再乘 100 产生的尾差。

| 函数 | 返回与精度约定 |
| --- | --- |
| `Yuan2Fen(yuan string) int64` | 元转整数分，允许正负号、小数和十进制科学计数法，如 `"1.23"`、`".5"`、`"1e2"`；不足一分向零截断。 |
| `Fen2Yuan(cents int64) string` | 整数分转元，完整支持 int64 范围，固定输出两位小数。 |
| `YuanAdd(arr ...string) string` | 每项先分别转为整数分再求和，结果固定两位小数；没有参数时返回 `"0"`。 |
| `YuanSub(a, b string) string` | a、b 分别转为整数分后计算 a − b，结果固定两位小数。 |
| `YuanMul(a string, b float64) string` | a 先转为整数分再乘 b，结果不足一分向零截断，固定输出两位小数。 |
| `YuanDiv(a string, b float64) string` | a 先转为整数分再除以 b，结果不足一分向零截断，固定输出两位小数；零除数会 panic。 |

输入不自动去空白，不接受空串、NaN、Inf、十六进制或 `"1/2"`。
无效输入、截断后的单个输入分值或最终整数分结果超出 int64 范围，均会 panic；可表示的整数分对应 `-92233720368547758.08` 至 `92233720368547758.07` 元。
极大指数还受 float64 有限范围检查限制，下溢到零的极小金额返回零。
乘除倍率必须是有限 float64，NaN/Inf 会 panic；内部使用其最短十进制表示，但不能消除调用方先前浮点计算产生的误差。

```go
package main

import (
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	fmt.Println(keqing.Yuan2Fen("19.999"))
	fmt.Println(keqing.Yuan2Fen("-0.019")) // 向零截断
	fmt.Println(keqing.Fen2Yuan(-1))
	fmt.Println(keqing.YuanAdd("0.009", "0.009")) // 每项先截断到分
	fmt.Println(keqing.YuanSub("10", "3.25"))
	fmt.Println(keqing.YuanMul("0.03", 0.5))
	fmt.Println(keqing.YuanDiv("1.00", 3))
	fmt.Println(keqing.YuanAdd())

	// Output:
	// 1999
	// -1
	// -0.01
	// 0.00
	// 6.75
	// 0.01
	// 0.33
	// 0
}
```

## 文件工具

| 函数 | 用法与返回值 |
| --- | --- |
| `MkdirAll(path string)` | 递归创建目录，已有目录不报错；失败 `panic`。 |
| `FileExists(path string) bool` | 文件或目录存在且查询成功时为 `true`；不存在或查询失败均为 `false`。 |
| `FileExistsWithError(path string) (bool, error)` | 不存在返回 `false, nil`；其他查询失败返回错误。与 `FileExists` 一样，会跟随符号链接查询目标。 |
| `ReadDirFiles(dir string) ([]string, error)` | 递归返回所有非目录路径，不跟随符号链接遍历目标。失败时可能同时返回部分结果和错误。路径相对或绝对与入参一致。 |
| `FileMD5(filePath string) (string, error)` | 流式计算文件内容的 MD5，返回 32 位小写十六进制摘要，不返回文件名。 |
| `SaveFile(path, data string) error` | 写入字符串；已有文件会被截断覆盖，不创建父目录。新文件权限参数为 `0644`。 |
| `CopyFile(filePath, savePath string) error` | 复制普通文件，自动创建目标父目录并覆盖目标内容；拒绝复制到同一文件（含硬链接）。不复制权限、时间等元数据，失败时目标可能只写入部分内容。 |
| `CopyFile2IO(dst io.Writer, src io.Reader) error` | 从 `src` 复制到 `dst`，返回复制错误，不关闭流。nil 接口返回错误；带类型的 nil 指针由 Reader/Writer 的实现处理。 |
| `RemoveFile(path string) error` | 删除文件或空目录，不递归；路径不存在也返回错误。 |
| `RemoveDir(path string) error` | 删除路径及全部内容，也能删除文件；路径不存在不报错。 |
| `ReNameFile(oldPath, newPath string) error` | 重命名或移动，目标父目录需存在；覆盖目标、跨文件系统移动等限制由操作系统决定。 |
| `HomePath() string` | 返回当前用户主目录；查询失败 `panic`。 |
| `FileSize(path string) int64` | 返回 `os.Stat` 的大小；不存在或查询失败返回 `0`，不特别排除目录。需区分空文件与失败时使用 `os.Stat`。 |

下面的示例只在新建的临时目录中写入文件，并在结束时清理：

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	dir, err := os.MkdirTemp("", "keqing-files-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	src := filepath.Join(dir, "hello.txt")
	if err := keqing.SaveFile(src, "hello"); err != nil {
		panic(err)
	}
	dst := filepath.Join(dir, "backup", "hello.txt")
	if err := keqing.CopyFile(src, dst); err != nil {
		panic(err)
	}
	exists, err := keqing.FileExistsWithError(dst)
	if err != nil {
		panic(err)
	}
	digest, err := keqing.FileMD5(dst)
	if err != nil {
		panic(err)
	}
	fmt.Println(exists, keqing.FileSize(dst))
	fmt.Println(digest == keqing.MD5("hello"))
	fmt.Println(keqing.MD5Salt("hello", "salt") == keqing.MD5("hellosalt"))
	// Output:
	// true 5
	// true
	// true
}
```

## MD5 摘要

| 函数 | 用法与返回值 |
| --- | --- |
| `MD5(str string) string` | 对字符串字节计算 MD5，返回 32 位小写十六进制字符串。例如 `MD5("hello")` 为 `5d41402abc4b2a76b9719d911017c592`。 |
| `MD5Salt(str, salt string) string` | 等同于 `MD5(str + salt)`，先把 salt 拼在原字符串后再计算摘要，参见上面的文件示例。 |

MD5 是摘要算法，不能解密；`MD5Salt` 不是密码存储专用的哈希方案。

## RSA 加解密

默认的 `RsaEncrypt` / `RsaDecrypt`、`Custom` 和 `4RsaKey` 系列均使用 **RSA-OAEP + SHA-256**，MGF1 使用 SHA-256，label 为空。只有 `RsaEncrypt4PKCS1` / `RsaDecrypt4PKCS1` 使用 **PKCS#1 v1.5 填充**。密钥编码格式与加密填充是不同概念；当前接口没有 OAEP-SHA1 模式。

加密返回标准 Base64 字符串，解密接收该 Base64 字符串。不会自动分块：2048 位密钥下，OAEP-SHA256 明文最多 190 字节，PKCS#1 v1.5 最多 245 字节；中文也按字节计算。

| 函数 | 用法与返回值 |
| --- | --- |
| `RsaGenerateKey(path string) string` | 生成并保存 2048 位密钥对，返回传入目录；失败 `panic`。 |
| `RsaGenerateKeyWithError(path string) (string, error)` | 相同生成操作；失败返回空路径和错误。 |
| `RsaLoadPublicKey(publicKeyPEM string) *rsa.PublicKey` | 解析 PEM 文本中的 PKIX 公钥（`PUBLIC KEY`）；格式错误或非 RSA 密钥 `panic`。 |
| `RsaLoadPrivateKey(privateKeyPEM string) *rsa.PrivateKey` | 解析 PEM 文本中的 PKCS#8 私钥；格式错误或非 RSA 密钥 `panic`。参数不是文件路径。 |
| `RsaLoadKey(publicKey, privateKey, keyType string) error` | 校验并加载匹配的默认密钥对。`RSA_KEY_FILE_TYPE` 表示两个参数为文件路径，`RSA_KEY_STRING_TYPE` 表示 PEM 文本。读取失败、类型无效、密钥不匹配返回错误；空参数和 PEM 解析错误仍会 `panic`。失败保留原默认密钥对。 |
| `RsaEncrypt(data string) string` / `RsaDecrypt(ciphertextString string) string` | 使用 `RsaLoadKey` 加载的进程级默认密钥对。 |
| `RsaEncryptCustom(publicKey, data string) string` / `RsaDecryptCustom(privateKey, ciphertextString string) string` | 每次传入 PEM 文本，不改变默认密钥对。 |
| `RsaEncrypt4RsaKey(publicKey *rsa.PublicKey, data string) string` / `RsaDecrypt4RsaKey(key *rsa.PrivateKey, ciphertextString string) string` | 使用已经解析的密钥对象，便于重复使用。 |
| `RsaEncrypt4PKCS1(publicKey *rsa.PublicKey, data string) string` / `RsaDecrypt4PKCS1(key *rsa.PrivateKey, ciphertextString string) string` | 使用密钥对象做 PKCS#1 v1.5 加解密，需与对端使用相同填充。 |

生成函数创建 `privateKey.pem`（PKCS#8）和 `publicKey.pem`（PKIX）；任一同名文件已存在都会失败，轮换应使用新目录。Unix 私钥权限为 `0600`，Windows 访问控制依赖目录 ACL。写入、同步和关闭错误都会被处理；失败仅清理本次创建且身份仍匹配的文件，不承诺两个文件跨崩溃原子发布。

加解密的 `string` 返回值不能承载错误：Base64 解码、长度超限、密钥/填充不匹配等失败会打印日志并返回空字符串；空参数、nil 密钥、未加载默认密钥或无效 PEM 会 `panic`。字符串判空沿用 `IsEmpty`，所以 `"null"`（忽略大小写）也会被拒绝。`RSA_PKCS1` 和 `RSA_PKCS8_SHA_256` 不是 `RsaLoadKey` 的 `keyType`，也不是加解密函数的模式选择参数。

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	ciphertext := keqing.RsaEncrypt4RsaKey(&privateKey.PublicKey, "hello")
	if ciphertext == "" {
		panic("RSA encryption failed")
	}
	plaintext := keqing.RsaDecrypt4RsaKey(privateKey, ciphertext)
	if plaintext == "" {
		panic("RSA decryption failed")
	}
	fmt.Println(plaintext)
	// Output:
	// hello
}
```

## 随机数与标识符

| 函数 | 用法与返回值 |
| --- | --- |
| `Random() float64` | 返回 `[0, 1)` 的伪随机数。 |
| `RandomNextInt(args ...interface{}) int` | 一个参数 `max` 得到 `[0, max)`，两个参数 `min, max` 得到 `[min, max)`；右端点不包含。 |
| `WinningRate(num float64) bool` | `num` 是百分比，`25` 表示 25%；`num <= 0` 为 false，`num >= 100` 为 true。范围内按 0.000001 个百分点向下量化，请传有限数值。 |
| `UUID() string` | 返回小写十六进制、带连字符的 36 字符 ID（`8-4-4-4-12`）。版本位为 7，包含毫秒时间、进程内递增序列及加密随机部分；随机源失败 `panic`。 |
| `ConfigureRandomID(workerID int64) error` | 首次生成雪花 ID 前设置 `0..1023` 的实例编号；越界或开始生成后再配置返回错误。 |
| `RandomId() int64` | 返回雪花 ID；默认初始化失败或时钟回拨时 `panic`。 |
| `RandomIdWithError() (int64, error)` | 使用同一个雪花生成器，通过 error 返回初始化或时钟回拨错误。 |

`Random`、`RandomNextInt` 和 `WinningRate` 共享内部加锁的 `math/rand` 生成器，适合抽样与模拟，不用于令牌或密钥生成。

`RandomNextInt` 只接受内置 `int`、`int8`、`int16`、`int32`、`int64`，不接受无符号整数、浮点数或自定义整数类型。零参数 `panic`；不支持的类型或超过两个参数返回 `0`；`max <= 0`、`min >= max` 或整数超出当前平台 `int` 范围会 `panic`。`0` 也可能是合法随机结果，不能用于判断参数是否有效。

雪花 ID 默认使用网卡 MAC 的哈希作为 workerID，仅适用于单进程。多进程/多主机部署必须为同时运行的实例分配不同的 workerID，并在任何 goroutine 生成 ID 前完成配置。库不负责分布式分配，也不持久化重启前的时间和序列。每毫秒序列耗尽会等待下一毫秒；等待期间回拨也返回错误。

```go
package main

import (
	"fmt"
	"regexp"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
	// 7 仅为示例；应用启动时从部署配置中读取，每个实例应不同。
	if err := keqing.ConfigureRandomID(7); err != nil {
		panic(err)
	}
	_, err := keqing.RandomIdWithError()
	if err != nil {
		panic(err)
	}
	fmt.Println(err == nil)
	x, y := keqing.Random(), keqing.RandomNextInt(-3, 4)
	fmt.Println(x >= 0 && x < 1, y >= -3 && y < 4)
	fmt.Println(keqing.WinningRate(100))
	id := keqing.UUID()
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	fmt.Println(len(id), regexp.MustCompile(pattern).MatchString(id))
	// Output:
	// true
	// true true
	// true
	// 36 true
}
```

## 示例与测试

更多示例见 [List](keqing/demo/List_test.go)、[Set](keqing/demo/Set_test.go)、[其他工具](keqing/demo/Demo_test.go) 和 [可执行 Example](keqing/Example_test.go)。在仓库根目录运行：

```sh
go test ./...
```

需要显式处理错误时，可选用 `RandomIdWithError`、`RsaGenerateKeyWithError`、`FileExistsWithError`；其他入口的失败方式见各模块说明。
