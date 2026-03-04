# 六、并发控制内功：Context 使用

## 1️⃣ Context 设计目的（高频）

### 面试题

**为什么需要 context？它解决了什么问题？**

### 标准答案

Context 是 Go 1.7 引入的标准库，用于**跨 API 边界和 goroutin传递**：

1. **取消信号**：告诉 goroutine 停止工作
2. **超时/截止时间**：自动取消超时的操作
3. **请求范围的值**：传递 traceID、userID 等元数据

### 核心问题

**问题 1**：如何优雅停止 goroutine？

```go
// ❌ 没有 context：无法通知 worker 停止
func worker(ch <-chan int) {
    for v := range ch {
        process(v)
    }
}

// ✅ 有 context：可以取消
func worker(ctx context.Context, ch <-chan int) {
    for {
        select {
        case v := <-ch:
            process(v)
        case <-ctx.Done():
            return  // 收到取消信号
        }
    }
}
```

**问题 2**：如何控制超时？

```go
// ❌ 没有 context：手动管理 timer
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()
select {
case result := <-ch:
    // ...
case <-timer.C:
    // timeout
}

// ✅ 有 context：统一超时管理
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case result := <-ch:
    // ...
case <-ctx.Done():
    // timeout or canceled
}
```

### Context 的设计原则

1. **不可变**：context 一旦创建，不能修改
2. **可派生**：可以从 context 派生新的 context
3. **取消单向传播**：父 context 取消，子 context 全部取消
4. **线程安全**：多个 goroutine 可以安全访问同一个 context

---

### 典型使用场景

1. **HTTP 请求处理**：请求超时控制
2. **数据库查询**：查询超时
3. **RPC 调用**：调用链超时传播
4. **后台任务**：优雅关闭

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // r.Context() 包含请求超时、取消信号
    ctx := r.Context()
    
    // 传递给下游
    data, err := fetchData(ctx)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    json.NewEncoder(w).Encode(data)
}

func fetchData(ctx context.Context) ([]byte, error) {
    // 带超时的数据库查询
    return db.QueryContext(ctx, "SELECT ...")
}
```

---

## 2️⃣ Context 的创建方式（必考）

### 面试题

**context 有哪几种创建方式？它们分别用于什么场景？**

### 2.1 根 Context

#### context.Background()

* **顶层 context**，无取消、无超时、无值
* 用于 **main 函数**、**init 函数**、**测试**
* 作为其他 context 的根

```go
func main() {
    ctx := context.Background()
    
    // 从根 context 派生
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    // ...
}
```

#### context.TODO()

* **占位符 context**，语义上表示"不确定用哪个 context"
* 用于**过渡期**，代码重构时临时使用
* 功能与 Background() 相同

```go
func oldFunc() {
    ctx := context.TODO()  // 未来会替换为真实 context
    newFunc(ctx)
}
```

---

### 2.2 WithCancel：手动取消

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()  // 确保释放资源
```

* 返回派生 context 和 cancel 函数
* 调用 `cancel()` 触发 `ctx.Done()` 关闭
* 适用于**手动控制取消**

**示例**：

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("worker stopped")
            return
        default:
            // do work
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx)
    
    time.Sleep(500 * time.Millisecond)
    cancel()  // 通知 worker 停止
    time.Sleep(100 * time.Millisecond)
}
```

---

### 2.3 WithTimeout：相对超时

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()
```

* **5 秒后**自动取消
* 返回 cancel 函数，可提前取消
* 适用于**有时间限制的操作**

**示例**：

```go
func query(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    ch := make(chan result)
    go func() {
        ch <- slowQuery()  // 可能很慢的查询
    }()
    
    select {
    case res := <-ch:
        return processResult(res)
    case <-ctx.Done():
        return ctx.Err()  // context.DeadlineExceeded
    }
}
```

---

### 2.4 WithDeadline：绝对截止时间

```go
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()
```

* 在**指定时间点**自动取消
* WithTimeout 本质是 WithDeadline 的封装

**区别**：

```go
// WithTimeout：相对时间
ctx, _ := context.WithTimeout(ctx, 5*time.Second)

// WithDeadline：绝对时间
ctx, _ := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
// 等价
```

---

### 2.5 WithValue：传递请求范围的值

```go
ctx := context.WithValue(parent, key, value)
```

* 传递**请求范围的元数据**
* key 应使用**自定义类型**避免冲突
* **不应传递业务参数**

**示例**：

```go
type contextKey string

const (
    traceIDKey contextKey = "traceID"
    userIDKey  contextKey = "userID"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    traceID := generateTraceID()
    ctx := context.WithValue(r.Context(), traceIDKey, traceID)
    
    // 传递给下游
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    traceID := ctx.Value(traceIDKey).(string)
    log.Printf("[%s] processing...", traceID)
}
```

**为什么 key 要用自定义类型？**

```go
// ❌ 错误：用 string 作 key，容易冲突
ctx = context.WithValue(ctx, "userID", 123)
ctx = context.WithValue(ctx, "userID", 456)  // 覆盖

// ✅ 正确：用自定义类型
type myKey string
const userIDKey myKey = "userID"
ctx = context.WithValue(ctx, userIDKey, 123)
```

---

## 3️⃣ Context 取消传播（重要）

### 面试题

**父 context 取消后，子 context 会怎样？反过来呢？**

### 取消的单向传播

**规则**：
* **父取消 → 子取消**（级联）
* **子取消 ✗ 父不受影响**（隔离）

```go
// 父 context
parentCtx, parentCancel := context.WithCancel(context.Background())

// 子 context 1
childCtx1, childCancel1 := context.WithCancel(parentCtx)

// 子 context 2
childCtx2, childCancel2 := context.WithCancel(parentCtx)

// 孙 context
grandCtx, grandCancel := context.WithCancel(childCtx1)

// 取消父 context
parentCancel()

// 结果：childCtx1, childCtx2, grandCtx 全部被取消
// 但兄弟之间不受影响
```

**示意图**：

```
Background
    ↓
parentCtx ----【cancel】
    ├─→ childCtx1 ----【被取消】
    │       ↓
    │   grandCtx ----【被取消】
    │
    └─→ childCtx2 ----【被取消】
```

---

### Done() channel 的使用

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    <-ctx.Done()  // 阻塞，直到 context 取消
    fmt.Println("context canceled")
    fmt.Println("reason:", ctx.Err())
}()

time.Sleep(time.Second)
cancel()  // 触发 Done() 关闭
```

**Done() 特性**：
* 返回只读 channel
* context 取消时关闭
* 重复读取返回零值（不阻塞）

---

### Err() 方法

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

time.Sleep(2 * time.Second)

if err := ctx.Err(); err != nil {
    if err == context.DeadlineExceeded {
        fmt.Println("timeout")
    } else if err == context.Canceled {
        fmt.Println("canceled")
    }
}
```

**返回值**：
* `nil`：未取消
* `context.Canceled`：调用了 cancel()
* `context.DeadlineExceeded`：超时

---

## 4️⃣ Context 使用规范（最佳实践）

### 面试题

**使用 context 有哪些注意事项和规范？**

### 规范 1：作为函数首参数

```go
// ✅ 正确
func DoSomething(ctx context.Context, arg string) error {
    // ...
}

// ❌ 错误：不是首参数
func DoSomething(arg string, ctx context.Context) error {
    // ...
}
```

**原因**：
* Go 社区约定
* 便于识别和传递
* 参数命名统一为 `ctx`

---

### 规范 2：不要存储在结构体中

```go
// ❌ 错误：存储 context
type Server struct {
    ctx context.Context  // 不要这样做
}

// ✅ 正确：每个方法接收 context
type Server struct {
    // 其他字段
}

func (s *Server) HandleRequest(ctx context.Context) {
    // ...
}
```

**原因**：
* context 生命周期与请求绑定，不应全局共享
* 结构体可能被多个请求使用

**例外**：
* 测试代码可以存储 background context
* 极少数场景（如长期运行的后台worker）

---

### 规范 3：不要传递 nil context

```go
// ❌ 错误
func DoSomething(ctx context.Context) {
    // ctx 可能是 nil，导致 panic
    select {
    case <-ctx.Done():  // panic if ctx == nil
        return
    }
}

// ✅ 正确：使用 Background 或 TODO
func DoSomething(ctx context.Context) {
    if ctx == nil {
        ctx = context.Background()
    }
    // ...
}

// 更好：调用者传正确的 context
DoSomething(context.Background())
```

---

### 规范 4：WithValue 仅用于请求范围的元数据

```go
// ✅ 适合的值
ctx = context.WithValue(ctx, traceIDKey, "abc123")
ctx = context.WithValue(ctx, userIDKey, 456)
ctx = context.WithValue(ctx, requestIDKey, "req-789")

// ❌ 不适合的值
ctx = context.WithValue(ctx, "db", dbConn)  // 不要传递依赖
ctx = context.WithValue(ctx, "config", cfg)  // 不要传递配置
ctx = context.WithValue(ctx, "result", data)  // 不要传递业务数据
```

**原因**：
* 业务参数应该显式传递（类型安全）
* Value 查找有性能开销
* 依赖应通过依赖注入或参数传递

---

## 5️⃣ Context.Value 的性能考虑（实战）

### 面试题

**context.Value 有性能问题吗？如何优化？**

### Value 的性能开销

**原理**：
* context 是链表结构
* 查找 Value 需要**遍历链**

```go
parentCtx := context.Background()
ctx1 := context.WithValue(parentCtx, key1, "value1")
ctx2 := context.WithValue(ctx1, key2, "value2")
ctx3 := context.WithValue(ctx2, key3, "value3")

// 查找 key1 需要遍历 ctx3 → ctx2 → ctx1
v := ctx3.Value(key1)
```

**benchmark**：

```go
// 查找第一个 Value：~10ns
// 查找第 10 个 Value：~50ns
// 查找第 100 个 Value：~500ns
```

---

### 优化方法

#### 1. 减少 Value 数量

```go
// ❌ 每个字段单独存
ctx = context.WithValue(ctx, "traceID", id)
ctx = context.WithValue(ctx, "userID", uid)
ctx = context.WithValue(ctx, "requestID", rid)

// ✅ 用结构体聚合
type RequestMeta struct {
    TraceID   string
    UserID    int
    RequestID string
}

ctx = context.WithValue(ctx, metaKey, RequestMeta{...})
```

#### 2. 在入口处提取，避免重复查找

```go
// ❌ 每次调用都查找
func process(ctx context.Context) {
    traceID := ctx.Value(traceIDKey).(string)  // 查找
    step1(ctx)
    step2(ctx)
}

func step1(ctx context.Context) {
    traceID := ctx.Value(traceIDKey).(string)  // 重复查找
    // ...
}

// ✅ 提取后直接传递
func process(ctx context.Context) {
    traceID := ctx.Value(traceIDKey).(string)  // 查找一次
    step1(traceID)
    step2(traceID)
}

func step1(traceID string) {
    // 直接使用，无需查找
}
```

#### 3. 使用自定义类型作为 key

```go
// ❌ string key：容易冲突，难以优化
ctx = context.WithValue(ctx, "key", value)

// ✅ 自定义类型：类型安全，可优化比较
type myKey struct{}
var key myKey
ctx = context.WithValue(ctx, key, value)
```

---

## 6️⃣ 面试高频题汇总

### Q1：context 解决了什么问题？

**答**：
* 跨 goroutine 传递取消信号
* 统一的超时/截止时间管理
* 传递请求范围的元数据

### Q2：context 的创建方式有哪些？

**答**：

| 函数 | 用途 |
|------|------|
| Background() | 根 context（main/init/test） |
| TODO() | 占位符 context |
| WithCancel() | 手动取消 |
| WithTimeout() | 相对超时 |
| WithDeadline() | 绝对截止时间 |
| WithValue() | 传递值 |

### Q3：取消如何传播？

**答**：
* **单向传播**：父取消 → 子全部取消
* 子取消不影响父和兄弟

### Q4：Done() channel 的特性？

**答**：
* 只读 channel
* context 取消时关闭
* 重复读取不阻塞（返回零值）

### Q5：Err() 的返回值？

**答**：
* `nil`：未取消
* `context.Canceled`：手动取消
* `context.DeadlineExceeded`：超时

### Q6：WithValue 应该存什么？

**答**：
* ✅ 请求元数据：traceID、userID、requestID
* ❌ 业务参数、依赖、配置

### Q7：为什么不要在结构体中存储 context？

**答**：
* context 生命周期与请求绑定
* 结构体可能被多个请求共享
* 违反 context 的设计原则

### Q8：如何避免传递 nil context？

**答**：
* 使用 `context.Background()`
* 或 `context.TODO()`

### Q9：context.Value 的性能如何？

**答**：
* 需要遍历 context 链
* 查找开销随链长度增加
* 优化：入口提取，减少 Value 数量，直接传递

### Q10：为什么 WithValue 的 key 要用自定义类型？

**答**：避免冲突

```go
type myKey string
const userIDKey myKey = "userID"
```

---

## 🔥 追问升级：实战场景

### Q：如何实现请求链路追踪？

**答**：

```go
type trace struct {
    TraceID string
    SpanID  string
}

type traceKey struct{}

func StartTrace(ctx context.Context) context.Context {
    return context.WithValue(ctx, traceKey{}, &trace{
        TraceID: generateTraceID(),
        SpanID:  generateSpanID(),
    })
}

func GetTrace(ctx context.Context) *trace {
    if t, ok := ctx.Value(traceKey{}).(*trace); ok {
        return t
    }
    return nil
}

// 使用
ctx := StartTrace(context.Background())
callService1(ctx)
callService2(ctx)
```

### Q：如何优雅关闭服务？

**答**：

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // 启动多个 worker
    go worker1(ctx)
    go worker2(ctx)
    go worker3(ctx)
    
    // 监听信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    
    <-sigCh
    fmt.Println("shutting down...")
    
    cancel()  // 通知所有 worker 停止
    time.Sleep(2 * time.Second)  // 等待清理
}

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            cleanup()
            return
        default:
            doWork()
        }
    }
}
```

---

**本章总结**：

✅ context 用于跨 API 边界传递取消、超时、元数据  
✅ Background 用于 main/init，TODO 用于占位  
✅ WithCancel 手动取消，WithTimeout 相对超时，WithDeadline 绝对时间  
✅ 取消单向传播：父取消→子全部取消，子取消不影响父  
✅ Done() 关闭表示取消，Err() 返回取消原因  
✅ 作为首参数，不存结构体，不传 nil，Value 仅用于元数据  
✅ WithValue 有遍历开销，优化：入口提取、减少数量、直接传递
