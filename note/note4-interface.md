# 四、类型系统内功：Interface 与反射

## 1️⃣ Interface 底层结构（高频）

### 面试题

**Go 的 interface 底层是如何实现的？为什么 nil interface 判断会出错？**

### 标准答案

Go 的 interface 有两种底层结构：

#### 1. `eface`（空接口 interface{}）

```go
type eface struct {
    _type *_type        // 动态类型信息
    data  unsafe.Pointer // 动态值指针
}
```

* 用于 `interface{}`
* 只有类型和数据，无方法表

#### 2. `iface`（非空接口）

```go
type iface struct {
    tab  *itab          // 类型信息 + 方法表
    data unsafe.Pointer // 动态值指针
}

type itab struct {
    inter *interfacetype // 接口类型
    _type *_type         // 实际类型
    fun   [1]uintptr     // 方法地址数组（动态大小）
}
```

* 用于有方法的接口
* `itab` 包含方法表，用于动态派发

### 追问 1：什么时候 interface 为 nil？

**只有动态类型和动态值都为 nil，interface 才为 nil**

```go
var i interface{}
fmt.Println(i == nil)  // true（类型和值都是 nil）

var p *int = nil
i = p
fmt.Println(i == nil)  // false！（类型不为 nil）

// 底层结构：
// i := eface{_type: *int 的类型信息, data: nil}
// 类型信息不是 nil，所以 interface 不是 nil
```

### 追问 2：如何正确判断 interface 是否为 nil？

**方法 1**：使用反射

```go
func isNil(i interface{}) bool {
    if i == nil {
        return true
    }
    v := reflect.ValueOf(i)
    return v.Kind() == reflect.Ptr && v.IsNil()
}
```

**方法 2**：类型断言

```go
var i interface{} = (*int)(nil)

// 错误判断
fmt.Println(i == nil)  // false

// 正确判断
if i != nil {
    if p, ok := i.(*int); ok && p == nil {
        fmt.Println("pointer is nil")
    }
}
```

### 追问 3：interface 赋值时发生了什么？

```go
type Animal interface {
    Speak() string
}

type Dog struct {}
func (d Dog) Speak() string { return "Woof" }

var a Animal = Dog{}
```

**步骤**：
1. 查找 `Dog` 是否实现了 `Animal` 的方法（编译期）
2. 创建 `itab`（首次运行时，之后缓存）
3. 生成 `iface{tab: itab指针, data: Dog实例指针}`

**性能**：
* 首次赋值：需要创建 itab，有开销
* 后续赋值：itab 被缓存，开销小

---

## 2️⃣ 类型断言与判断（必考）

### 面试题

**类型断言的两种形式有什么区别？type switch 怎么用？**

### 标准答案

#### 类型断言的两种形式

**单返回值**：

```go
var i interface{} = "hello"
s := i.(string)  // 成功：s = "hello"
n := i.(int)     // panic: interface conversion: interface {} is string, not int
```

* 断言失败会 **panic**
* 适用于确定类型的场景

**双返回值**：

```go
var i interface{} = "hello"
s, ok := i.(string)  // 成功：s = "hello", ok = true
n, ok := i.(int)     // 失败：n = 0, ok = false（不 panic）
```

* 断言失败返回 **零值 + false**
* 适用于不确定类型的场景
* **推荐使用**

### type switch 用法

```go
func checkType(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %s\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}
```

**注意**：
* `i.(type)` 只能在 `switch` 中使用
* `v` 在每个 `case` 中是具体类型，不是 `interface{}`

### 追问 1：类型断言的性能如何？

* 比反射快（10倍+）
* 但比直接调用慢（需要查表）
* 频繁断言会有开销，考虑缓存结果

**benchmark 示例**：

```go
// 直接调用：1ns
// 类型断言：5ns
// 反射调用：100ns
```

### 追问 2：可以断言具体类型，能断言接口类型吗？

**可以！**

```go
type Reader interface {
    Read() string
}

type Writer interface {
    Write(string)
}

type ReadWriter interface {
    Reader
    Writer
}

var rw ReadWriter = &MyStruct{}

// 断言为子接口
r, ok := rw.(Reader)  // ok = true
w, ok := rw.(Writer)  // ok = true
```

---

## 3️⃣ 反射基本概念（核心）

### 面试题

**解释反射三定律，以及如何通过反射修改变量值？**

### 反射三定律

#### 定律 1：从 interface{} 到 reflection object

```go
var x float64 = 3.4
v := reflect.ValueOf(x)   // Value
t := reflect.TypeOf(x)    // Type

fmt.Println("type:", t)   // float64
fmt.Println("value:", v)  // 3.4
```

* `TypeOf` 获取类型信息
* `ValueOf` 获取值信息

#### 定律 2：从 reflection object 到 interface{}

```go
v := reflect.ValueOf(3.4)
x := v.Interface().(float64)  // 还原为 float64
fmt.Println(x)  // 3.4
```

* `Interface()` 方法返回 `interface{}`
* 需要类型断言还原具体类型

#### 定律 3：要修改 reflection object，值必须可寻址

```go
var x float64 = 3.4
v := reflect.ValueOf(x)
v.SetFloat(7.1)  // panic: reflect.Value.SetFloat using unaddressable value
```

**为什么？**

* `ValueOf(x)` 传的是 `x` 的**拷贝**
* 修改拷贝不影响原变量
* 必须传**指针**

**正确写法**：

```go
var x float64 = 3.4
p := reflect.ValueOf(&x)  // 传指针
v := p.Elem()             // 获取指针指向的值
v.SetFloat(7.1)           // 修改成功
fmt.Println(x)            // 7.1
```

### 追问 1：如何判断 Value 是否可修改？

```go
v := reflect.ValueOf(x)
fmt.Println(v.CanSet())  // false

v := reflect.ValueOf(&x).Elem()
fmt.Println(v.CanSet())  // true
```

* `CanSet()` 返回是否可修改
* 只有可寻址的 Value 才能 Set

### 追问 2：反射能修改未导出字段吗？

**不能！**

```go
type Person struct {
    name string  // 未导出
    Age  int     // 导出
}

p := Person{name: "Alice", Age: 30}
v := reflect.ValueOf(&p).Elem()

// 修改导出字段：成功
v.FieldByName("Age").SetInt(31)

// 修改未导出字段：panic
v.FieldByName("name").SetString("Bob")  // panic: reflect.Value.SetString using value obtained using unexported field
```

**绕过方法**（不推荐）：

```go
nameField := v.FieldByName("name")
reflect.NewAt(nameField.Type(), unsafe.Pointer(nameField.UnsafeAddr())).
    Elem().SetString("Bob")
```

---

## 4️⃣ 反射常见用法（实战）

### 4.1 读取 struct tag

```go
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"email"`
}

func printTags(u User) {
    t := reflect.TypeOf(u)
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        jsonTag := field.Tag.Get("json")
        validateTag := field.Tag.Get("validate")
        fmt.Printf("%s: json=%s, validate=%s\n", 
            field.Name, jsonTag, validateTag)
    }
}

// 输出：
// Name: json=name, validate=required
// Email: json=email, validate=email
```

**应用**：JSON 序列化、ORM 映射、参数验证

### 4.2 动态调用方法

```go
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
    return a + b
}

func callMethod() {
    c := Calculator{}
    v := reflect.ValueOf(c)
    
    // 获取方法
    method := v.MethodByName("Add")
    
    // 准备参数
    args := []reflect.Value{
        reflect.ValueOf(10),
        reflect.ValueOf(20),
    }
    
    // 调用
    results := method.Call(args)
    fmt.Println(results[0].Int())  // 30
}
```

**注意**：
* 方法名必须导出（首字母大写）
* 参数类型必须匹配
* 返回值是 `[]reflect.Value`

### 4.3 遍历结构体字段

```go
type Person struct {
    Name string
    Age  int
}

func printFields(p Person) {
    v := reflect.ValueOf(p)
    t := reflect.TypeOf(p)
    
    for i := 0; i < v.NumField(); i++ {
        fieldValue := v.Field(i)
        fieldType := t.Field(i)
        fmt.Printf("%s (%s) = %v\n", 
            fieldType.Name, 
            fieldType.Type, 
            fieldValue.Interface())
    }
}

// 输出：
// Name (string) = Alice
// Age (int) = 30
```

### 4.4 动态创建实例

```go
func newInstance(t reflect.Type) interface{} {
    // 创建指针类型的实例
    ptrValue := reflect.New(t)
    
    // 返回具体类型
    return ptrValue.Interface()
}

// 使用
t := reflect.TypeOf(Person{})
p := newInstance(t).(*Person)
p.Name = "Bob"
```

---

## 5️⃣ 反射性能影响（重要）

### 面试题

**反射为什么慢？什么场景适合用反射？**

### 性能开销来源

1. **动态类型检查**：运行时才知道类型
2. **方法查找**：需要遍历方法表
3. **参数装箱**：基本类型 → interface{}
4. **禁用内联优化**：编译器无法优化

**benchmark 对比**：

```go
// 直接调用
func BenchmarkDirect(b *testing.B) {
    c := Calculator{}
    for i := 0; i < b.N; i++ {
        _ = c.Add(1, 2)
    }
}
// 结果：1 ns/op

// 反射调用
func BenchmarkReflect(b *testing.B) {
    c := Calculator{}
    v := reflect.ValueOf(c)
    method := v.MethodByName("Add")
    args := []reflect.Value{
        reflect.ValueOf(1),
        reflect.ValueOf(2),
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        method.Call(args)
    }
}
// 结果：100 ns/op（慢 100 倍）
```

### 适用场景

**✅ 适合用反射**：

* **框架/库开发**：如 JSON 序列化、ORM、依赖注入
* **需要通用性**：处理任意类型
* **代码量远大于性能损失**：写 1000 个类型转换 vs 用反射

**❌ 不适合用反射**：

* **热点代码路径**：频繁调用的函数
* **性能敏感**：实时系统、高频交易
* **简单类型转换**：直接用类型断言

### 优化建议

#### 1. 缓存反射结果

```go
var typeCache = make(map[reflect.Type]*CachedInfo)

func getInfo(t reflect.Type) *CachedInfo {
    if info, ok := typeCache[t]; ok {
        return info
    }
    
    // 首次创建
    info := &CachedInfo{
        Fields: parseFields(t),
        Methods: parseMethods(t),
    }
    typeCache[t] = info
    return info
}
```

#### 2. 避免重复 TypeOf/ValueOf

```go
// 慢
for i := 0; i < 1000; i++ {
    t := reflect.TypeOf(obj)  // 重复创建
    // ...
}

// 快
t := reflect.TypeOf(obj)
for i := 0; i < 1000; i++ {
    // 复用 t
}
```

#### 3. 使用代码生成替代反射

* 编译期生成类型特定的代码
* 工具：`go generate`、`protobuf`、`stringer`

```go
//go:generate stringer -type=Status
type Status int
```

---

## 6️⃣ 面试高频题汇总

### Q1：interface 为 nil 的判断陷阱？

**答**：

```go
var p *int = nil
var i interface{} = p
fmt.Println(i == nil)  // false！
```

* interface 包含**类型和值**
* 只有两者都为 nil，interface 才为 nil
* `p` 的类型是 `*int`（不是 nil），所以 `i != nil`

**正确判断**：

```go
func isNil(i interface{}) bool {
    return i == nil || reflect.ValueOf(i).IsNil()
}
```

### Q2：反射三定律是什么？

**答**：

1. **interface → reflection object**：`TypeOf/ValueOf`
2. **reflection object → interface**：`Interface()`
3. **修改值需要可寻址**：必须传指针，用 `Elem()` 获取

### Q3：反射为什么慢？

**答**：

* 动态类型检查（运行时）
* 方法查找开销
* 参数装箱/拆箱
* 禁用编译器优化

**慢多少**：比直接调用慢 10-100 倍

### Q4：什么时候用反射？

**答**：

* ✅ 框架开发（JSON/ORM/DI）
* ✅ 需要通用性
* ❌ 热点代码路径
* ❌ 性能敏感场景

### Q5：如何读取 struct tag？

**答**：

```go
type User struct {
    Name string `json:"name"`
}

t := reflect.TypeOf(User{})
field, _ := t.FieldByName("Name")
tag := field.Tag.Get("json")  // "name"
```

### Q6：反射能修改私有字段吗？

**答**：

* **不能**（正常方式）
* 可用 `unsafe.Pointer` 强制修改（不推荐，破坏封装）

### Q7：interface{} 的性能损失？

**答**：

* 装箱/拆箱开销
* 无法内联
* GC 压力增加（堆分配）

**优化**：
* 热点路径避免 interface{}
* 用泛型（Go 1.18+）替代

### Q8：反射获取方法名列表？

**答**：

```go
t := reflect.TypeOf(obj)
for i := 0; i < t.NumMethod(); i++ {
    method := t.Method(i)
    fmt.Println(method.Name)
}
```

### Q9：如何判断类型实现了某接口？

**答**：

```go
type Reader interface {
    Read() string
}

var r Reader
readerType := reflect.TypeOf((*Reader)(nil)).Elem()

t := reflect.TypeOf(obj)
if t.Implements(readerType) {
    fmt.Println("implements Reader")
}
```

---

## 🔥 追问升级：高级应用

### Q：如何用反射实现深拷贝？

**答**：

```go
func deepCopy(src interface{}) interface{} {
    srcVal := reflect.ValueOf(src)
    dstVal := reflect.New(srcVal.Type()).Elem()
    
    copyRecursive(dstVal, srcVal)
    return dstVal.Interface()
}

func copyRecursive(dst, src reflect.Value) {
    switch src.Kind() {
    case reflect.Ptr:
        if !src.IsNil() {
            dst.Set(reflect.New(src.Elem().Type()))
            copyRecursive(dst.Elem(), src.Elem())
        }
    case reflect.Struct:
        for i := 0; i < src.NumField(); i++ {
            copyRecursive(dst.Field(i), src.Field(i))
        }
    case reflect.Slice:
        dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
        for i := 0; i < src.Len(); i++ {
            copyRecursive(dst.Index(i), src.Index(i))
        }
    default:
        dst.Set(src)
    }
}
```

### Q：反射和泛型（Go 1.18+）如何选择？

**答**：

| 场景 | 反射 | 泛型 |
|-----|-----|-----|
| 编译期类型检查 | ❌ | ✅ |
| 性能 | 慢 | 快（接近直接调用） |
| 支持任意类型 | ✅ | ❌（需类型约束） |
| 代码生成 | 运行时 | 编译期 |

**建议**：
* 能用泛型的优先用泛型
* 需要处理完全未知类型才用反射

---

**本章总结**：

✅ interface 有两种底层结构：eface（空接口）和 iface（非空接口）  
✅ nil interface 陷阱：只有类型和值都为 nil 才是 nil  
✅ 类型断言用双返回值形式更安全  
✅ 反射三定律：interface ↔ reflection object，修改需可寻址  
✅ 反射慢 10-100 倍，适合框架开发而非热点代码  
✅ 缓存反射结果、避免重复 TypeOf、用代码生成替代反射
