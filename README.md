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
    RsaGenerateKey(path) // 生成RSA密钥对 , path为密钥对保存路径
    RsaLoadKey(publicKey string, privateKey string, keyType string) // 加载RSA密钥对
    RsaLoadPrivateKey(privateKey string) *rsa.PrivateKey // 加载RSA私钥
    RsaLoadPublicKey(publicKey string) *rsa.PublicKey // 加载RSA公钥
    RsaEncrypt(data string) string // RSA加密 - 前置条件需使用RsaLoadKey加载密钥对
    RsaEncryptCustom(publicKey string, data string) string // RSA加密 - 自定义密钥字符串加密
    RsaEncrypt4RsaKey(publicKey *rsa.PublicKey, data string) string // RSA加密 - 自定义公钥对象加密
    RsaDecrypt(data string) string // RSA解密 - 前置条件需使用RsaLoadKey加载密钥对
    RsaDecryptCustom(privateKey string, data string) string // RSA解密 - 自定义密钥字符串解密
    RsaDecrypt4RsaKey(privateKey *rsa.PrivateKey, data string) string // RSA解密 - 自定义私钥对象解密
  ```

### File工具类
  ```text
   MkdirAll(path) // 创建文件夹
   FileExists(path) bool // 判断文件是否存在
   ReadDirFiles(path) ([]string, error) // 获取文件夹下的所有文件
   FileMD5(path) (string, error) // 获取文件的MD5值
   SaveFile(path,content) error // 将字符串保存到文件
   CopyFile(src,dest string) error // 复制文件 - 使用字符串路径
   CopyFile2IO(dst io.Writer, src io.Reader) error // 复制文件 - 使用io.Reader和io.Writer
   RemoveFile(path) error // 删除文件
   RemoveDir(path) error // 删除文件夹
   ReNameFile(oldPath,newPath) error // 重命名文件
   FileSize(path) (int64, error) // 获取文件大小
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
