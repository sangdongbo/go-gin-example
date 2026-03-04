# 五、错误处理内功：defer / panic / recover

## 1️⃣ defer 执行顺序（高频）

### 面试题

**defer 的执行顺序是什么？参数什么时候求值？它对返回值有什么影响？**

### 标准答案

#### defer 执行顺序：LIFO（后进先出）

```go
func main() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
}
// 输出：3 2 1
```

* 多个 defer 按**栈**的方式执行
* 最后声明的最先执行

#### 参数求值时机：defer 语句执行时

```go
func main() {
    x := 1
    defer fmt.Println(x)  // 此时 x = 1，参数确定为 1
    
    x = 2
}
// 输出：1（不是 2）
```

**原因**：
* defer 语句执行时就确定参数值
* 不是函数执行时才求值

---

### defer 与 return 的关系（重要）

**执行顺序**：

1. 赋值返回值
2. 执行 defer
3. 返回

#### 示例 1：无命名返回值

```go
func f() int {
    x := 5
    defer func() {
        x++  // 修改局部变量
    }()
    return x  // 返回 5
}

fmt.Println(f())  // 5
```

**原因**：
* `return x` → 将 x 的值（5）赋给匿名返回值
* defer 修改的是局部变量 x，不是返回值
* 返回 5

#### 示例 2：命名返回值

```go
func f() (result int) {
    defer func() {
        result++  // 修改返回值
    }()
    return 5  // result = 5
}

fmt.Println(f())  // 6
```

**原因**：
* `return 5` → `result = 5`
* defer 修改 `result`（返回值） → `result = 6`
* 返回 6

#### 示例 3：指针返回值

```go
func f() *int {
    x := 5
    defer func() {
        x++  // 修改局部变量
    }()
    return &x  // 返回 x 的地址
}

fmt.Println(*f())  // 6
```

**原因**：
* `return &x` → 将 x 的地址赋给返回值
* defer 修改 x → x = 6
* 通过指针读取，得到 6

### 追问 1：defer 为什么能修改命名返回值？

* 命名返回值是**函数作用域变量**
* defer 在函数返回前执行，可访问所有函数变量
* 匿名返回值是编译器生成的临时变量，defer 无法访问

### 追问 2：多个 return 和 defer 的顺序？

```go
func f() (result int) {
    defer func() {
        fmt.Println("defer 1")
        result++
    }()
    
    defer func() {
        fmt.Println("defer 2")
        result += 10
    }()
    
    return 5
}

// 执行顺序：
// 1. result = 5
// 2. defer 2: result = 15
// 3. defer 1: result = 16
// 4. 返回 16

// 输出：
// defer 2
// defer 1
// 返回值：16
```

---

### defer 与 panic 的关系

```go
func main() {
    defer fmt.Println("defer 1")
    defer fmt.Println("defer 2")
    
    panic("error")
    
    defer fmt.Println("defer 3")  // 不执行
}

// 输出：
// defer 2
// defer 1
// panic: error
```

**规则**：
* panic 触发后，按 LIFO 顺序执行**已注册**的 defer
* panic 之后的 defer 不会注册，不会执行

---

## 2️⃣ panic 和 recover 机制（必考）

### 面试题

**panic 的传播过程是怎样的？recover 必须怎么用才有效？**

### panic 传播过程

```go
func level3() {
    defer fmt.Println("defer in level3")
    panic("error")
}

func level2() {
    defer fmt.Println("defer in level2")
    level3()
    fmt.Println("after level3")  // 不执行
}

func level1() {
    defer fmt.Println("defer in level1")
    level2()
    fmt.Println("after level2")  // 不执行
}

func main() {
    defer fmt.Println("defer in main")
    level1()
    fmt.Println("after level1")  // 不执行
}

// 输出：
// defer in level3
// defer in level2
// defer in level1
// defer in main
// panic: error
```

**传播规则**：

1. panic 发生后，当前函数停止执行
2. 按 LIFO 执行当前函数的 defer
3. 返回调用者，继续执行调用者的 defer
4. 逐层向上传播，直到：
   - 被 recover 捕获
   - 或程序崩溃

---

### recover 使用规则（重要）

#### 规则 1：必须在 defer 函数中**直接**调用

```go
// ✅ 正确
func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered:", r)
        }
    }()
    
    panic("error")
}

// ❌ 错误：间接调用
func f() {
    defer recoverFunc()  // recover 在这里无效
    panic("error")
}

func recoverFunc() {
    if r := recover(); r != nil {
        fmt.Println("recovered:", r)
    }
}
```

**原因**：
* recover 检查**调用它的函数**是否是 defer 函数
* 间接调用时，recover 在 recoverFunc 中，不是 defer 函数

#### 规则 2：只捕获当前 goroutine 的 panic

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered in main:", r)  // 不会执行
        }
    }()
    
    go func() {
        panic("error in goroutine")
    }()
    
    time.Sleep(time.Second)
}
// panic: error in goroutine（main 的 recover 无法捕获）
```

**原因**：
* 每个 goroutine 有独立的调用栈
* recover 只对当前 goroutine 有效

**正确做法**：

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered in goroutine:", r)
        }
    }()
    
    panic("error")
}()
```

#### 规则 3：recover 返回 panic 的参数值

```go
defer func() {
    if r := recover(); r != nil {
        switch v := r.(type) {
        case string:
            fmt.Println("string panic:", v)
        case error:
            fmt.Println("error panic:", v)
        default:
            fmt.Println("unknown panic:", v)
        }
    }
}()

panic("custom error")  // 输出：string panic: custom error
```

* `recover()` 返回 `panic(v)` 中的 `v`
* 如果没有 panic，返回 `nil`

---

## 3️⃣ 错误处理最佳实践（重要）

### 面试题

**什么时候用 error，什么时候用 panic？错误包装怎么做？**

### error vs panic 的选择

#### 使用 error 的场景（常规）

* **可预期的错误**：文件不存在、网络超时、参数错误
* **调用者可以处理**：重试、降级、返回默认值
* **业务逻辑错误**：库函数应该返回 error

```go
func ReadFile(path string) ([]byte, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("file not found: %s", path)
    }
    return os.ReadFile(path)
}
```

#### 使用 panic 的场景（异常）

* **不可恢复的错误**：内存耗尽、不变式被破坏
* **编程错误**：空指针、数组越界、类型断言失败
* **初始化失败**：配置加载失败、必需资源不可用

```go
func init() {
    cfg, err := loadConfig()
    if err != nil {
        panic("failed to load config: " + err.Error())
    }
}
```

**经验规则**：

* **库函数**：返回 error，不 panic
* **应用程序**：可在 main/init 中 panic
* **中间层**：recover panic，转换为 error

---

### 错误包装（Go 1.13+）

#### fmt.Errorf + %w

```go
func loadUser(id int) (*User, error) {
    data, err := db.Query(id)
    if err != nil {
        return nil, fmt.Errorf("load user %d: %w", id, err)
    }
    return parseUser(data)
}
```

* `%w` 包装原始错误
* 保留错误链

#### errors.Is 检查错误类型

```go
if errors.Is(err, os.ErrNotExist) {
    // 处理文件不存在
}
```

* 遍历错误链，检查是否包含目标错误

#### errors.As 提取错误类型

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("path:", pathErr.Path)
    fmt.Println("op:", pathErr.Op)
}
```

* 遍历错误链，提取第一个匹配类型

#### 自定义错误类型

```go
type ValidationError struct {
    Field string
    Value interface{}
    Err   error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %v", e.Field, e.Err)
}

func (e *ValidationError) Unwrap() error {
    return e.Err  // 支持 errors.Is/As
}
```

---

## 4️⃣ defer 性能分析（实战）

### 面试题

**defer 有性能开销吗？什么时候应该避免使用？**

### defer 的性能开销

**Go 1.13 之前**：
* defer 需要在堆上分配结构体
* 每次 defer 有 ~50ns 开销

**Go 1.13+ 优化**：
* defer 在栈上分配（open-coded defer）
* 开销降至 ~10ns（接近直接调用）

**benchmark 对比**：

```go
// 无 defer
func BenchmarkNormal(b *testing.B) {
    for i := 0; i < b.N; i++ {
        mu.Lock()
        // do something
        mu.Unlock()
    }
}
// 结果：10 ns/op

// 有 defer（Go 1.13+）
func BenchmarkDefer(b *testing.B) {
    for i := 0; i < b.N; i++ {
        mu.Lock()
        defer mu.Unlock()
        // do something
    }
}
// 结果：12 ns/op（几乎无差别）
```

### 使用建议

#### 推荐使用 defer 的场景

* **资源释放**：锁、文件、连接
* **代码清晰度**：确保一定执行

```go
// 推荐
func process() error {
    mu.Lock()
    defer mu.Unlock()
    
    // 复杂逻辑
    if err := step1(); err != nil {
        return err  // 自动解锁
    }
    if err := step2(); err != nil {
        return err  // 自动解锁
    }
    return nil
}

// 不推荐：容易忘记 unlock
func process() error {
    mu.Lock()
    
    if err := step1(); err != nil {
        mu.Unlock()
        return err
    }
    if err := step2(); err != nil {
        mu.Unlock()  // 重复代码
        return err
    }
    mu.Unlock()
    return nil
}
```

#### 避免 defer 的场景

**热点路径 + 极简操作**：

```go
// 如果循环体非常简单，且性能极度敏感
for i := 0; i < 1000000; i++ {
    mu.Lock()
    counter++
    mu.Unlock()
    
    // 不用 defer，避免开销
}
```

但 99% 的场景，defer 的清晰度 > 微小性能损失

---

### defer 与锁的配合

```go
type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
}

// 读锁
func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

// 写锁
func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

// 锁升级（读锁 → 写锁）
func (c *Cache) GetOrCreate(key string, fn func() string) string {
    // 先尝试读锁
    c.mu.RLock()
    if v, ok := c.data[key]; ok {
        c.mu.RUnlock()
        return v
    }
    c.mu.RUnlock()
    
    // 升级为写锁
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 双重检查
    if v, ok := c.data[key]; ok {
        return v
    }
    
    v := fn()
    c.data[key] = v
    return v
}
```

---

## 5️⃣ 面试高频题汇总

### Q1：defer 的执行顺序？

**答**：LIFO（后进先出），像栈一样

### Q2：defer 参数何时求值？

**答**：defer **语句执行时**，不是函数执行时

```go
x := 1
defer fmt.Println(x)  // 此时确定为 1
x = 2
// 输出：1
```

### Q3：defer 如何修改返回值？

**答**：只能修改**命名返回值**

```go
func f() (result int) {
    defer func() { result++ }()
    return 5  // result = 5 → defer执行 → result = 6
}
```

### Q4：panic 的传播过程？

**答**：

1. 当前函数停止执行
2. 执行当前函数的 defer（LIFO）
3. 返回调用者，继续传播
4. 直到被 recover 或程序崩溃

### Q5：recover 的使用规则？

**答**：

* ✅ 必须在 defer 函数中**直接**调用
* ✅ 只捕获当前 goroutine 的 panic
* ❌ 不能在 defer 外调用
* ❌ 不能间接调用

```go
// ✅ 正确
defer func() {
    if r := recover(); r != nil {
        // handle
    }
}()

// ❌ 错误
defer recoverFunc()  // 间接调用无效
```

### Q6：什么时候用 error，什么时候用 panic？

**答**：

* **error**：可预期、可处理的错误（常规）
* **panic**：不可恢复、编程错误、初始化失败（异常）

**经验**：
* 库函数返回 error
* 应用程序可在 main/init 中 panic
* 中间层 recover panic → error

### Q7：如何包装错误？

**答**：

```go
return fmt.Errorf("context: %w", err)  // 包装

errors.Is(err, target)  // 检查
errors.As(err, &target)  // 提取
```

### Q8：defer 的性能开销？

**答**：

* Go 1.13 前：~50ns
* Go 1.13+：~10ns（栈上分配）
* 几乎可忽略，优先考虑代码清晰度

### Q9：多次 defer 同一个 unlock 会怎样？

**答**：panic（double unlock）

```go
mu.Lock()
defer mu.Unlock()
defer mu.Unlock()  // panic: sync: unlock of unlocked mutex
```

### Q10：defer 中 panic，原 panic 会怎样？

**答**：**新 panic 覆盖旧 panic**

```go
defer func() {
    panic("new panic")
}()
panic("old panic")
// 最终传播的是 "new panic"
```

---

## 🔥 追问升级：高级场景

### Q：如何捕获 goroutine 的 panic 并上报？

**答**：

```go
func SafeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                // 获取堆栈
                buf := make([]byte, 4096)
                n := runtime.Stack(buf, false)
                stack := string(buf[:n])
                
                // 上报到监控系统
                logError(fmt.Sprintf("panic: %v\nstack:\n%s", r, stack))
            }
        }()
        
        fn()
    }()
}

// 使用
SafeGo(func() {
    // 可能 panic 的代码
    panic("error")
})
```

### Q：如何实现 try-catch-finally？

**答**：

```go
func Try(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    fn()
    return nil
}

func Finally(fn func()) {
    defer fn()
}

// 使用
err := Try(func() {
    // try block
    panic("error")
})
if err != nil {
    // catch block
    fmt.Println("caught:", err)
}
Finally(func() {
    // finally block
    fmt.Println("cleanup")
})
```

### Q：defer 中使用闭包的陷阱？

**答**：

```go
// ❌ 错误：闭包捕获的是变量引用
for i := 0; i < 3; i++ {
    defer fmt.Println(i)
}
// 输出：2 2 2（i 最后是 3，但循环结束前是 2）

// ✅ 正确：传参
for i := 0; i < 3; i++ {
    defer func(n int) {
        fmt.Println(n)
    }(i)
}
// 输出：2 1 0

// ✅ 正确：创建新变量
for i := 0; i < 3; i++ {
    i := i  // 创建新变量
    defer fmt.Println(i)
}
// 输出：2 1 0
```

---

**本章总结**：

✅ defer 按 LIFO 执行，参数在 defer 语句执行时求值  
✅ defer 可修改命名返回值，执行顺序是：赋值返回值 → defer → 返回  
✅ panic 逐层向上传播，执行每层的 defer，直到被 recover 或崩溃  
✅ recover 必须在 defer 函数中直接调用，只捕获当前 goroutine  
✅ error 用于可预期错误，panic 用于不可恢复的异常  
✅ 使用 fmt.Errorf("%w", err) 包装错误，errors.Is/As 检查错误链  
✅ Go 1.13+ defer 开销大幅降低，几乎可忽略，优先考虑代码清晰度
