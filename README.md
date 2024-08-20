# go-sweet-keqing

## 适用于go语言的工具类(刻晴版)
## 剑光如我，斩尽牛杂！
### 基础支持
 - go 1.20及以上环境
### 使用方式

 - 1 引入包 
   ```text
    go get github.com/PurpleScorpion/go-sweet-keqing
    
   ```
   ```text
    使用以下语句来引入包
    import "github.com/PurpleScorpion/go-sweet-keqing/keqing"
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
    Shuffle(arr) []T // 随机打乱数组
    SortArray(arr) []T // 升序排序
    SortArrayDesc(arr) []T // 降序排序
   ```

### 日期工具类
  ```text
    DateInfo(time.Time) TimeVO // 获取日期详细信息
    NowDate() time.Time // 获取当前日期
    NowUTCDate() time.Time // 获取当前UTC日期
    NowDateStr() string // 获取当前日期字符串 使用默认格式 2006-01-02 15:04:05
    NowUTCDateStr() string // 获取当前UTC日期字符串 使用默认格式 2006-01-02T15:04:05.999999Z
    FormatDate(time.Time,string) string // 格式化日期 , 根据自定义格式化字符串format
    ParseDate(dateStr,string) time.Time // 解析日期 , 根据自定义格式化字符串format
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
    // TODO 待开发
  ```

### 软妹币工具类
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

### 字符串工具类
   ```text
    Replace2Regex(str,oldStr,newStr) string // 将字符串中的oldStr替换成newStr , 使用正则表达式
    Replace(str,oldStr,newStr) string // 将字符串中的oldStr替换成newStr
    IndexOf(str,subStr) int // 获取字符串中subStr的位置 , -1表示不存在
    LastIndexOf(str,subStr) int // 获取字符串中最后一个subStr的位置 , -1表示不存在
    GetBytes(str) []byte // 将字符串解析成byte数组
    Contains(str,subStr) bool // 判断字符串是否包含subStr
    ToString(value) string // 将 数组 / 结构体 转换为标准字符串打印
    ToStringArray(arr) []string // 将数组转换为字符串数组
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
    Float2Str(num) string // 浮点数转字符串 (可选择保留的小数位数) 若不需要精度，则直接调用Num2Str即可
    ParseInt64(str) int64 // 将字符串解析成int64
    ParseInt32(str) int32 // 将字符串解析成int32
    ParseInt16(str) int16 // 将字符串解析成int16
    ParseInt8(str) int8 // 将字符串解析成int8
    ParseInt(str) int // 将字符串解析成int
    ParseFloat32(str) float32 // 将字符串解析成float32
    ParseFloat64(str) float64 // 将字符串解析成float64
    Split(str,sep) []string // 将字符串按照sep分割成字符串数组
   ```
