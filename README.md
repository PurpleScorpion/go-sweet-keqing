# go-sweet-keqing

## 适用于go语言的工具类(刻晴版)
## 剑光如我，斩尽牛杂！
### 基础支持
 - go 1.20及以上环境
 - github地址: https://github.com/PurpleScorpion/go-sweet-keqing
### 使用方式

 - 1 引入包 
   ```text
    go get github.com/PurpleScorpion/go-sweet-keqing
    
   ```
   ```text
    使用以下语句来引入包
    import "github.com/PurpleScorpion/go-sweet-keqing/keqing"
   ```
### Set集合
   ```text
    Set集合
    说明: 与Java中的Set集合类似 , 集合中的元素不能重复 , 但是集合中的元素是存取有序的
    使用方式: 
    set := keqing.NewSet[string](nil) // 基本使用方式
    set := keqing.NewSet[int](nil) // 基本使用方式
    ...
    // 更多复杂使用方式,请查阅 demo/Set_test.go#TestSet3 和 demo/Set_test.go#TestSet4
   
    Add(T) bool // 添加元素 , 返回是否添加成功
    Contains(T) bool // 判断元素是否存在
    Remove(T) bool // 删除元素 , 返回是否删除成功
    Clear() // 清空集合
    Size() int // 集合大小
    已弃用 - GetAll() []T // 获取所有元素 , 返回一个切片(例如: []int) 
    GetData() []T // 获取所有元素 , 返回一个切片(例如: []int) 
    ToArray() []T // 获取所有元素 , 返回一个切片(例如: []int) 
    IsEmpty() bool // 判断集合是否为空
    Iterator() <-chan T // 迭代器 , 具体使用方式请查阅 demo/Set_test.go#TestSet1
    Union(other *Set[T]) *Set[T] // 集合并集
    Intersect(other *Set[T]) *Set[T] // 交集
    Difference(other *Set[T]) *Set[T] // 求两个集合的差集（属于 set 不属于 other 的元素）
    ToString() string // 符合json规则的格式化输出
    RemoveAll(items ...T) // 批量删除多个元素
    RetainAll(items ...T) // 保留指定元素
    Filter(f func(item T) bool) *Set[T] // 按条件筛选集合中的元素 , 不会修改原数组
    Copy() *Set[T] // 返回一个新的 Set，包含当前 Set 的所有元素（浅拷贝）
    DeepCopy(cloneFunc func(T) T) *Set[T] // 创建一个新的 Set，包含当前 Set 的所有元素（深拷贝）
    ToList() *List[T] // 将列表转换为 List
    
    完整demo示例请查阅 demo/Set_test.go
    注意: 包含切片类型的实体类请参考 demo/List_test.go#TestList_Sort_CustomStruct
   ```
### List 集合
  ```text
    List集合
    说明: 与Java中的List集合类似 , 集合中的元素可以重复 , 集合中的元素是存取有序的
    使用方式:
    list := keqing.NewList[string](nil)
    list := keqing.NewList[int](nil)
    ...
    // 更多复杂使用方式,请查阅 demo/List_test.go
    Add(value T) // 添加元素，无返回值
    Get(index) T // 获取指定索引位置的元素
    Set(index int, value T) // 替换指定索引位置的元素，无返回值
    Contains(value) bool // 判断元素是否存在
    IndexOf(value) int // 获取指定元素在集合中的索引位置
    LastIndexOf(value) int // 返回最后一个匹配项的索引位置，不存在则返回 -1
    GetData() []T // 获取所有元素 , 返回一个切片(例如: []int) 
    ToArray() []T // 获取所有元素 , 返回一个切片(例如: []int) 
    Size() int // 集合大小
    GetFirst()  T // 获取第一个元素
    GetLast()  T // 获取最后一个元素
    IsEmpty() bool // 判断集合是否为空
    Remove(T) bool // 删除元素 , 返回是否删除成功
    RemoveAt(index) bool // 删除指定索引位置的元素
    Iterator() <-chan T // 迭代器
    Union(other *List[T]) *List[T] // 集合并集
    Intersect(other *List[T]) *List[T] // 交集
    Difference(other *List[T]) *List[T] // 求两个集合的差集（属于 list 不属于 other 的元素）
    ToString() string // 符合json规则的格式化输出
    RemoveAll(items ...T) // 批量删除多个元素
    RetainAll(items ...T) // 保留指定元素
    Filter(f func(item T) bool) *List[T] // 按条件筛选集合中的元素 , 不会修改原数组
    Copy() *List[T] // 返回一个新的 List，包含当前 List 的所有元素（浅拷贝）
    DeepCopy(cloneFunc func(T) T) *List[T] // 创建一个新的 List，包含当前 List 的所有元素（深拷贝）
    SubList(start int) *List[T] // 返回一个新的 List，表示原列表的一个子集, 从 start 开始到末尾
    SubList(start, end int) *List[T] // 返回一个新的 List，表示原列表的一个子集, 从 start（包含）到 end（不包含）
    Sort(less func(a, b T) bool) // 对列表进行排序，接受一个比较函数 less(a, b T) bool，表示 a 是否应排在 b 前面
    SortDESC() // 对基本类型列表进行降序排序（仅支持 string/int/float）
    SortASC() // 对基本类型列表进行升序排序（仅支持 string/int/float）
    ToSet() *Set[T] // 将列表转换为 Set
    
    完整demo示例请查阅 demo/List_test.go
    注意: 包含切片类型的实体类请参考 demo/List_test.go#TestList_Sort_CustomStruct
  ```
### Collections工具类
      
   ```text
        集合工具类 , 对List和Set集合进行增强操作
        List2Map[T comparable, K comparable]( list *List[T], keyFunc func(T) K, onConflict int) map[K]T // List集合转为Map集合 , 支持 key 冲突策略
        AsList[T comparable](arr []T, hashFunc func(T) uint64) *List[T] // 将数组转换为 List
        Set2Map[T comparable, K comparable]( set *Set[T], keyFunc func(T) K, onConflict int) map[K]T // Set集合转为Map集合 , 支持 key 冲突策略
        AsSet[T comparable](arr []T, hashFunc func(T) uint64) *Set[T] // 将数组转换为 Set
        
        完整demo示例请查阅 demo/List_test.go 和 demo/Set_test.go
   ```

### 数组工具类
   ```text
    ArrayContains(arr,value) bool //判断数组是否包含某个值
    GetOne(arr) T,bool // 获取数组中第0个位置的值
        value: 和数组的值类型一致
        bool: 是否存在
    Difference(arr1,arr2) []T // 求两个数组的差集
    Intersection(arr1,arr2) []T // 求两个数组的交集
    Union(arr1,arr2) []T // 求两个数组的并集
    Reverse(arr) []T // 反转数组
    Shuffle[T any](arr []T) // 原地随机打乱切片，无返回值
    SortArray(arr interface{}) // 原地升序排序，无返回值
    SortArrayDesc(arr interface{}) // 原地降序排序，无返回值
   ```

### 日期工具类
  ```text
    DateInfo(time.Time) TimeVO // 获取日期详细信息
    NowDate() time.Time // 获取当前日期
    NowUTCDate() time.Time // 获取当前UTC日期
    NowDateStr() string // 获取当前日期字符串 使用默认格式 2006-01-02 15:04:05
    NowUTCDateStr() string // 获取当前UTC日期字符串 使用默认格式 2006-01-02T15:04:05.999999Z
    FormatDate(time.Time,string) string // 格式化日期 , 根据自定义格式化字符串format
    ParseDate(dateStr string, format string, location *time.Location) time.Time // 解析日期；nil 表示 UTC
    DateAddSecond(time.Time,int64) time.Time // 添加秒数
    DateAddMinute(time.Time,int64) time.Time // 添加分钟数
    DateAddHour(time.Time,int64) time.Time // 添加小时数
    DateAddDay(time.Time,int64) time.Time // 添加天数
    DateSubSecond(time.Time,int64) time.Time // 减少秒数
    DateSubMinute(time.Time,int64) time.Time // 减少分钟数
    DateSubHour(time.Time,int64) time.Time // 减少小时数
    DateSubDay(time.Time,int64) time.Time // 减少天数
    ParseUTC(dateStr) time.Time // 解析UTC日期 使用默认格式 2006-01-02T15:04:05.999999Z
    ParseLocal(dateStr) time.Time // 解析本地日期 使用默认格式 2006-01-02 15:04:05
    
    UTC2Local(string) string // UTC日期转换为本地日期 , 使用默认格式
    UTC2LocalCustom(utc string, utcFormat string, localFormat string) string // UTC日期转换为本地日期,根据自定义格式化字符串format
    Local2UTC(string) string // 本地日期转换为UTC日期 , 使用默认格式
    Local2UTCCustom(local string, localFormat string, utcFormat string) string // 本地日期转换为UTC日期,根据自定义格式化字符串format
    
    CurrentTimeMillis() int64 // 获取当前时间戳(毫秒)
    CurrentTimeSeconds() int64 // 获取当前时间戳(秒)
    
    // 获取指定月份有多少天 , 其中month可以是 int(包含int/int16/int32/int64,以下简称为int)或者time.Month或time.Time类型
    // 当month是time.Month或int类型时, year参数必填
    // 当month是time.Time类型时, year参数可以不填, 如果填了year参数, 则会根据填写的year参数和month计算出当月的天数
    GetDaysInMonth(month any,year... int) int 
  ```

### 空值判断工具类
  ```text
    IsEmpty(value) bool // 判断是否为空
    IsNotEmpty(value) bool // 判断是否不为空
    
    可以判断的值有
    1. string
    2. int int8 int32 int64
    3. uint uint8 uint32 uint64
    4. float32 float64
    5. bool
    6. 数组切片
    7. map
    注意: 不支持struct类型的数据判断 , 会直接抛出异常(问: 你是故意的还是不小心的? 答:是故意不小心的)
    int,float类型的数据 0 为空
    bool其实没啥意义 , 原样返回给你
  ```

### MD5工具类
  ```text
    MD5(str) string // md5加密
    MD5Salt(str,salt) string // md5加盐加密
  ```

### RSA工具类
  ```text
    RsaGenerateKey(path string) string // 生成 RSA 密钥对，返回目录；失败 panic
    RsaGenerateKeyWithError(path string) (string, error) // 同上，失败返回 error
    RsaLoadKey(publicKey string, privateKey string, keyType string) error // 加载RSA密钥对
    RsaLoadPrivateKey(privateKey string) *rsa.PrivateKey // 加载RSA私钥
    RsaLoadPublicKey(publicKey string) *rsa.PublicKey // 加载RSA公钥
    RsaEncrypt(data string) string // RSA加密 - 前置条件需使用RsaLoadKey加载密钥对
    RsaEncryptCustom(publicKey string, data string) string // RSA加密 - 自定义密钥字符串加密
    RsaEncrypt4RsaKey(publicKey *rsa.PublicKey, data string) string // RSA加密 - 自定义公钥对象加密
    RsaDecrypt(data string) string // RSA解密 - 前置条件需使用RsaLoadKey加载密钥对
    RsaDecryptCustom(privateKey string, data string) string // RSA解密 - 自定义密钥字符串解密
    RsaDecrypt4RsaKey(privateKey *rsa.PrivateKey, data string) string // RSA解密 - 自定义私钥对象解密
    
    // PKCS1 的使用方式 , 这属于Java常用的加密方式 , 在与Java系统对接时可以使用
    * <li>{@code RSA/ECB/PKCS1Padding} (1024, 2048)</li>
    * <li>{@code RSA/ECB/OAEPWithSHA-1AndMGF1Padding} (1024, 2048)</li>
    * <li>{@code RSA/ECB/OAEPWithSHA-256AndMGF1Padding} (1024, 2048)</li>
    
    RsaEncrypt4PKCS1(publicKey *rsa.PublicKey, data string) string // RSA加密 - PKCS1加密
    RsaDecrypt4PKCS1(key *rsa.PrivateKey, ciphertextString string) string // RSA解密 - PKCS1解密
    
  ```

### File工具类
  ```text
   MkdirAll(path) // 创建文件夹
   FileExists(path string) bool // 仅查询成功且存在时为 true
   FileExistsWithError(path string) (bool, error) // 不存在为 false,nil；其他查询失败返回 error
   ReadDirFiles(path) ([]string, error) // 获取文件夹下的所有文件
   FileMD5(path) (string, error) // 获取文件的MD5值
   SaveFile(path,content) error // 将字符串保存到文件
   CopyFile(src,dest string) error // 复制文件 - 使用字符串路径
   CopyFile2IO(dst io.Writer, src io.Reader) error // 复制文件 - 使用io.Reader和io.Writer
   RemoveFile(path) error // 删除文件
   RemoveDir(path) error // 删除文件夹
   ReNameFile(oldPath,newPath) error // 重命名文件
   FileSize(path string) int64 // 获取文件大小；不存在或查询失败返回 0
   HomePath() string // 获取当前用户的主目录
  ```

### Yml工具类
  ```text
   yml数据获取工具, 旨在模拟@Value注解的功能
   重要前提: 
   请先用 gopkg.in/yaml.v3 v3.0.1 版本的包对yml文件进行解析成 map[string]interface{} 
   一共需要解析application.yml文件 与 application-激活的环境.yml文件
   问: 为什么不直接将这个包导入进来?
   答: 该工具类目的是纯净化封装, 不想因为导入了其他第三方的包从而导致在其他项目中产生不兼容的情况
   
   gopkg.in/yaml.v3的使用方法:
   data, err := os.ReadFile(confPath)
   if err != nil {
      panic("Error reading configuration file: " + err.Error())
   }
   var conf1 = make(map[string]interface{})
   yaml.Unmarshal(data, &conf1)
   
   data2, err2 := os.ReadFile(confPath2)
   if err2 != nil {
      panic("Error reading configuration file: " + err2.Error())
   }
   var conf2 = make(map[string]interface{})
   yaml.Unmarshal(data2, &conf2)
   
   然后使用本工具类中的-LoadYml(parent map[string]interface{}, child map[string]interface{}) error 函数进行加载
   
   方法介绍: 
   ValueInt(str) int // 获取yml中int类型的值
   ValueInt64(str) int64 // 获取yml中int64类型的值
   ValueInt32(str) int32 // 获取yml中int32类型的值
   ValueFloat64(str) float64 // 获取yml中float64类型的值
   ValueFloat32(str) float32 // 获取yml中float32类型的值
   ValueBool(str) bool // 获取yml中bool类型的值
   ValueString(str) string // 获取yml中string类型的值
   ValueStringArr(str) []string // 获取yml中字符串数组类型的值
   ValueObject(obj) interface{} // 获取yml中object类型的值 , 但是不建议使用 , 除非你想打印数据看看
   
   使用示例:
   注意: 参数的格式应当与@Value注解一致 , 当找不到指定key时, 会返回go中对应类型的默认值
   1. 获取int类型的值
   val1 := ValueInt("${com.example.config.value}")
   2. 获取string类型的值
   val2 := ValueString("${com.example.config.value}")
   3. 获取字符串数组类型的值
   val3 := ValueStringArr("${com.example.config.value}")
   4. 获取bool类型的值
   val4 := ValueBool("${com.example.config.value}")
   
   注意: gopkg.in/yaml.v3 中的一些坑
   当时使用map接受yml数据的时候
   偶尔我们的密码是123456这种纯数字类型的 , 他这个转换会给转成数字 , 导致我们无法获取到正确的值
   这是需要手动转成string类型
   例如:
    pwd := keqing.ValueObject("${sweet.mysql.password}")
	password := ""
	switch pwd.(type) {
	case int:
		password = fmt.Sprintf("%d", pwd.(int))
	case string:
		password = pwd.(string)
	}
  ```

### 软妹币工具类(仅支持元和分)
  金额以整数分为精度边界，使用精确十进制运算；需要小于分的精度时请使用适合业务的十进制类型。
  ```text
    Yuan2Fen(string) int64 // 元转分
    Fen2Yuan(int64) string // 分转元
    YuanAdd(arr... string) string // 多个元相加 , 可变数组
    YuanSub(str1,str2) string // 元相减
    YuanMul(str1,float64) string // 用元乘以一个浮点数 , 通常用于打折或者乘计算
    YuanDiv(str1,float64) string // 用元除以一个浮点数 , 通常用于打折或者除计算
  ```

### UUID工具类
  ```text
   UUID() 获取uuid字符串 , 该函数改自google的uuid库
   仅供学习参考，更多详细使用请查看google的uuid实现: https://github.com/google/uuid
  ```

### 随机数工具类
  ```text
   Random() float64 // 获取0-1之间的随机数, [0,1)
   RandomNextInt(max int) int // 获取0-max之间的随机数, [0,max)
   RandomNextInt(min,max int) int // 获取min-max之间的随机数, [min,max)
   RandomId() int64 // 雪花 ID；初始化失败或时钟回拨会 panic
   RandomIdWithError() (int64, error) // 使用同一生成器，失败返回 error
   ConfigureRandomID(workerID int64) error // 在首次生成 ID 前配置 0..1023 的 workerID
   WinningRate(num float64) bool // 根据输入的概率,  返回是否中奖 最高支持0.000001%精度
  ```
### 类型工具类
  ```text
    ToArray[T any](obj interface{}) []T // 将一个interface{}对象转换为指定类型的列表
    ToObject[T any](obj interface{}) T // 将一个interface{}对象转换为指定类型的对象
    ArrayMap(slice []S, fn func(S) T) []T // 将一个切片类型转换为另一个切片类型 (具体使用查看TestDemo35)
    CopyProperties(dst, src interface{}, ignoreFields ...string) error // 复制对象属性,其中ignoreFields是可选项, 表示忽略的字段 (具体使用查看TestDemo37)
  ```

### 字符串工具类
   ```text
    Replace2Regex(str,oldStr,newStr) string // 将字符串中的oldStr替换成newStr , 使用正则表达式
    Replace(str,oldStr,newStr) string // 将字符串中的oldStr替换成newStr
    IndexOf(str,subStr) int // 获取字符串中subStr的位置 , -1表示不存在
    LastIndexOf(str,subStr) int // 获取字符串中最后一个subStr的位置 , -1表示不存在
    GetBytes(str) []byte // 将字符串解析成byte数组
    Contains(str,subStr) bool // 判断字符串是否包含subStr
    ToString(value) string // 将 数组 / 结构体 转换为标准字符串打印
    ToJsonString(value) string // 将 数组 / 结构体 转换为标准Json字符串打印
    ToStringArray(str string) []string // 按 Unicode 字符将字符串拆成 []string
    ToFirstUpperCase(str) string // 将字符串的首字母转换为大写
    ToFirstLowerCase(str) string // 将字符串的首字母转换为小写
    ToUpperCase(str) string // 将字符串转换为大写
    ToLowerCase(str) string // 将字符串转换为小写
    Trim(str) string // 去除字符串两端的空白
    TrimSuffix(str,suffix) string // 去除字符串最后一个指定字符(如果有的话)
    Snake2Camel(str) string // 下划线转小驼峰
    Snake2BigCamel(str) string // 下划线转大驼峰
    Camel2Snake(str) string // 驼峰转下划线 (大小驼峰通吃)
    Num2Str(num) string // 数字转字符串 int家族/float家族 通用
    Float2Str(num interface{}, precision int) string // 浮点数转字符串 (可选择保留的小数位数) 若不需要精度，则直接调用Num2Str即可
    ParseInt64(str) int64 // 将字符串解析成int64
    ParseInt32(str) int32 // 将字符串解析成int32
    ParseInt16(str) int16 // 将字符串解析成int16
    ParseInt8(str) int8 // 将字符串解析成int8
    ParseInt(str) int // 将字符串解析成int
    ParseFloat32(str) float32 // 将字符串解析成float32
    ParseFloat64(str) float64 // 将字符串解析成float64
    Split(str,sep) []string // 将字符串按照sep分割成字符串数组
    StartWith(str,prefix) bool // 判断字符串是否以prefix开头
    EndWith(str,suffix) bool // 判断字符串是否以suffix结尾
    Camel2Hyphen(str) string // 驼峰转连字符(或称呼为中划线, 例如: keqing-yun)
    Camel2WebURL(str) string // 驼峰转web url标准格式的字符串(中划线拼接的方式)
   ```


### 可编译的最小示例

以下程序与 `keqing/Example_test.go` 中的 Example 一起维护；运行 `go test ./...` 会编译并验证示例。

```go
package main

import (
    "fmt"

    "github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

func main() {
    list := keqing.AsList([]int{3, 1}, nil) // nil 使用默认 hash
    list.Add(2)                           // 无返回值
    list.Set(0, 4)                        // 无返回值
    copied := list.DeepCopy(func(v int) int { return v })
    values := copied.ToArray()
    keqing.Shuffle(values)                // 以下三个函数原地修改切片
    keqing.SortArrayDesc(values)
    keqing.SortArray(values)
    fmt.Println(values)                  // [1 2 4]

    set := keqing.AsSet([]int{1, 1, 2}, nil)
    fmt.Println(set.DeepCopy(func(v int) int { return v }).ToArray()) // [1 2]
    fmt.Println(keqing.ToArray[int]([]int{1, 2}))                     // [1 2]
    fmt.Println(keqing.ToObject[string]("keqing"))                    // keqing
    date := keqing.ParseDate("2024-02-29", "2006-01-02", nil)
    fmt.Println(date.Format("2006-01-02 MST")) // 2024-02-29 UTC
}
```

`DeepCopy` 必须传入有效的 `cloneFunc`，由它负责复制元素内部的指针、切片、map 等引用；需要浅拷贝时使用 `Copy()`。`CopyProperties` 应检查返回的 error；`FileSize` 只有一个返回值，若需区分空文件与查询失败，请使用 `os.Stat`。

### 修复后的行为与兼容说明

现有导出函数的名称、参数和返回值保持不变。上面的参数与返回值修订纠正的是旧 README，已有能够编译的调用无需因此改写。新增的 `RandomIdWithError`、`RsaGenerateKeyWithError`、`FileExistsWithError` 可按需替换旧入口来处理错误。

- **金额和数值**：元字符串按十进制精确转换（支持科学计数法），每个金额参数先转换为整数分，超出分的部分向零截断；乘除结果也向零截断。例如 `Yuan2Fen("0.29") == 29`、`Yuan2Fen("-1.239") == -123`。金额 NaN/Inf、格式错误、除零和超出 int64 分范围会 panic。整数解析和 YAML 数值缩窄转换越界会 panic；通用 ParseFloat 对解析溢出报错，但仍允许显式的 NaN/Inf 字面量。
- **集合相等与迭代**：hash 碰撞不会合并不同元素。实现了 `Equal(T) bool` 的类型使用该结果（包括 false）；其他可比较值使用 Go `==`，指针默认按地址比较，`any` 内不可比较的值使用 `reflect.DeepEqual`。业务上按某几个字段去重时应实现 Equal，不能只靠 hash。Iterator 返回浅快照，可以提前停止或在迭代中修改原集合；两个集合的运算分别读取快照，不保证两份快照来自同一时刻。快照不会复制指针所指对象，外部修改仍需自行同步。
- **随机 ID**：默认 MAC 哈希 workerID 仅适用于单进程。多实例部署必须由部署方分配互不重复的 workerID，并在每个进程第一次生成 ID 前调用 `ConfigureRandomID`，检查返回错误；首次生成后不可重新配置。此接口不提供分布式分配或跨进程重启的状态持久化。旧入口在初始化失败、时钟回拨时 panic；需要恢复错误时使用 `RandomIdWithError`。
- **RSA 密钥**：生成到目录内的 `privateKey.pem`、`publicKey.pem`；任一同名文件已存在都会失败，不覆盖旧密钥。轮换应使用新目录。Unix 私钥权限为 0600，Windows 的访问控制依赖目录 ACL。写入、同步、关闭失败会返回错误或由旧入口 panic；失败清理仅处理本次创建且身份仍匹配的文件，不承诺两个文件跨崩溃原子发布。`RsaLoadKey` 在解析及公私钥匹配成功后一次替换全局密钥对，失败保留旧状态；错误 PEM 仍沿用原有 panic 行为。
- **文件与日期**：CopyFile 拒绝复制到同一文件（包括硬链接），避免截断源文件。FileExists 查询失败为 false，可用 WithError 版本区分原因。DateAddDay/DateSubDay 按日历天计算，夏令时切换日可能是 23 或 25 小时；时分秒仍是固定时长。今天、昨天、当月的范围终点均包含最后一纳秒；对应 UTC 方法将本地日历边界转换为 UTC。
- **反射与配置**：ToArray 对 nil 切片指针返回 nil，对 `[]*T` 内 nil 元素返回 T 的零值；CopyProperties 遇无法访问的 nil 嵌入指针返回 error。YAML 子配置的完整非 nil 路径覆盖父配置，缺失路径或 null 回退父配置；false、0、空字符串仍覆盖。Split/Camel2Snake 将字面量 `"null"` 当普通文本。
