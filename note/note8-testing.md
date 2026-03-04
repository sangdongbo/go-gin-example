# 八、测试与性能分析

## 1️⃣ 单元测试基础（必知）

### 面试题

**Go 的单元测试有什么特点？如何编写？**

### 1.1 测试文件命名规则

* 文件名：`xxx_test.go`
* 函数签名：`func TestXxx(t *testing.T)`
* 包名：与被测试代码相同，或 `xxx_test`（外部测试）

**示例**：

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(1, 2)
    if result != 3 {
        t.Errorf("Add(1, 2) = %d; want 3", result)
    }
}
```

---

### 1.2 testing.T 的常用方法

| 方法 | 作用 | 示例 |
|------|------|------|
| `t.Errorf()` | 记录错误，继续执行 | `t.Errorf("got %d, want %d", got, want)` |
| `t.Fatalf()` | 记录错误，终止当前测试 | `t.Fatalf("setup failed: %v", err)` |
| `t.Logf()` | 记录日志（-v 显示） | `t.Logf("testing with input %d", input)` |
| `t.Skip()` | 跳过测试 | `t.Skip("skipping in short mode")` |
| `t.Parallel()` | 并行执行 | `t.Parallel()` |
| `t.Helper()` | 标记辅助函数 | `t.Helper()` |

**区别**：
* `Error` vs `Fatal`：Fatal 会立即终止，Error 继续
* `Logf`：只在 `-v` 或测试失败时显示

```go
func TestExample(t *testing.T) {
    t.Logf("this is a log")  // 不会显示（除非 -v）
    
    got := compute()
    if got != want {
        t.Errorf("got %d, want %d", got, want)  // 记录错误，继续
        t.Logf("additional info")  // 会显示（因为测试失败）
    }
}
```

---

### 1.3 运行测试

```bash
# 运行所有测试
go test

# 运行指定包
go test ./pkg/math

# 运行指定测试函数
go test -run TestAdd

# 显示详细输出
go test -v

# 显示覆盖率
go test -cover

# 并行测试（指定并发数）
go test -parallel 4

# 短模式（跳过耗时测试）
go test -short
```

---

## 2️⃣ 表驱动测试（最佳实践）

### 面试题

**什么是表驱动测试？有什么优势？**

### 2.1 表驱动测试的原理

**核心思想**：用**数据表**定义测试用例，循环执行。

**优势**：
* 添加测试用例更简单（只需添加数据）
* 测试逻辑统一
* 易于维护

---

### 2.2 基本示例

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 2, 1},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

**运行结果**：

```bash
$ go test -v
=== RUN   TestAdd
=== RUN   TestAdd/positive
=== RUN   TestAdd/negative
=== RUN   TestAdd/zero
=== RUN   TestAdd/mixed
--- PASS: TestAdd (0.00s)
    --- PASS: TestAdd/positive (0.00s)
    --- PASS: TestAdd/negative (0.00s)
    --- PASS: TestAdd/zero (0.00s)
    --- PASS: TestAdd/mixed (0.00s)
```

---

### 2.3 t.Run 的作用

* **子测试**：每个用例独立
* **隔离失败**：一个用例失败不影响其他
* **并行执行**：`t.Parallel()`

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
    }
    
    for _, tt := range tests {
        tt := tt  // 捕获循环变量（避免并行测试问题）
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // 并行执行
            
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

**注意**：`tt := tt` 是为了**避免并行测试时的循环变量捕获问题**。

---

### 2.4 高级表驱动测试

**包含错误的测试**：

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name      string
        a, b      int
        want      int
        wantError bool
    }{
        {"normal", 6, 2, 3, false},
        {"divide by zero", 6, 0, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            if (err != nil) != tt.wantError {
                t.Errorf("Divide() error = %v, wantError %v", err, tt.wantError)
                return
            }
            if got != tt.want {
                t.Errorf("Divide() = %d, want %d", got, tt.want)
            }
        })
    }
}
```

---

## 3️⃣ Mock 和依赖注入（实战）

### 面试题

**如何测试依赖外部服务的代码？**

### 3.1 为什么需要 Mock？

**场景**：测试依赖数据库、HTTP API 等外部服务。

**问题**：
* 外部服务慢
* 外部服务不稳定
* 外部服务难以模拟异常情况

**解决**：用 **Mock** 替代真实依赖。

---

### 3.2 通过 Interface 抽象

```go
// 定义接口
type UserRepository interface {
    GetUser(id int) (*User, error)
}

// 真实实现
type DBUserRepository struct {
    db *sql.DB
}

func (r *DBUserRepository) GetUser(id int) (*User, error) {
    // 查询数据库
}

// 业务逻辑依赖接口
type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUserName(id int) (string, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}
```

---

### 3.3 手动 Mock

```go
// Mock 实现
type MockUserRepository struct {
    GetUserFunc func(id int) (*User, error)
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    if m.GetUserFunc != nil {
        return m.GetUserFunc(id)
    }
    return nil, nil
}

// 测试
func TestUserService_GetUserName(t *testing.T) {
    mockRepo := &MockUserRepository{
        GetUserFunc: func(id int) (*User, error) {
            if id == 1 {
                return &User{ID: 1, Name: "Alice"}, nil
            }
            return nil, errors.New("not found")
        },
    }
    
    service := &UserService{repo: mockRepo}
    
    // 测试正常情况
    name, err := service.GetUserName(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if name != "Alice" {
        t.Errorf("got %s, want Alice", name)
    }
    
    // 测试异常情况
    _, err = service.GetUserName(999)
    if err == nil {
        t.Error("expected error, got nil")
    }
}
```

---

### 3.4 使用 gomock（推荐）

**安装**：

```bash
go install github.com/golang/mock/mockgen@latest
```

**生成 Mock**：

```bash
mockgen -source=user.go -destination=mock_user.go -package=main
```

**使用**：

```go
func TestUserService_GetUserName(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := NewMockUserRepository(ctrl)
    
    // 设置期望
    mockRepo.EXPECT().GetUser(1).Return(&User{ID: 1, Name: "Alice"}, nil)
    
    service := &UserService{repo: mockRepo}
    name, err := service.GetUserName(1)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if name != "Alice" {
        t.Errorf("got %s, want Alice", name)
    }
}
```

**优势**：
* 自动生成代码
* 验证调用次数和参数
* 支持复杂 mock 逻辑

---

## 4️⃣ Benchmark 性能测试（重点）

### 面试题

**如何进行性能测试？b.N 是什么？**

### 4.1 Benchmark 基础

**函数签名**：

```go
func BenchmarkXxx(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 被测试的代码
    }
}
```

**b.N 的作用**：
* 自动调整循环次数
* 运行足够长时间以获得准确结果
* 通常 1-10 亿次

---

### 4.2 基本示例

```go
func Fibonacci(n int) int {
    if n < 2 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}

func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(10)
    }
}
```

**运行**：

```bash
$ go test -bench=.
BenchmarkFibonacci-8   3000000   450 ns/op
```

**含义**：
* `-8`：8 个 CPU 核心
* `3000000`：b.N 的值（运行 300 万次）
* `450 ns/op`：每次操作耗时 450 纳秒

---

### 4.3 b.ResetTimer 的作用

**问题**：setup 代码会影响测试结果。

```go
func BenchmarkProcess(b *testing.B) {
    // setup：耗时 1 秒
    data := generateLargeData()  // 慢
    
    b.ResetTimer()  // 重置计时器，排除 setup 时间
    
    for i := 0; i < b.N; i++ {
        process(data)
    }
}
```

---

### 4.4 测试内存分配

```bash
$ go test -bench=. -benchmem
BenchmarkFibonacci-8   3000000   450 ns/op   0 B/op   0 allocs/op
```

**含义**：
* `0 B/op`：每次操作分配 0 字节
* `0 allocs/op`：每次操作分配 0 次

**示例**：

```go
func BenchmarkStringConcat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = "hello" + "world"
    }
}

func BenchmarkStringBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        sb.WriteString("hello")
        sb.WriteString("world")
        _ = sb.String()
    }
}
```

**结果**：

```
BenchmarkStringConcat-8     100000000   10 ns/op   0 B/op   0 allocs/op
BenchmarkStringBuilder-8     50000000   30 ns/op  16 B/op   1 allocs/op
```

**分析**：
* 简单拼接编译器优化为常量，无分配
* StringBuilder 有分配，但适合循环拼接

---

### 4.5 比较 Benchmark

```go
func BenchmarkSliceAppend(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var s []int
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

func BenchmarkSlicePrealloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := make([]int, 0, 1000)
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}
```

**运行**：

```bash
$ go test -bench=. -benchmem
BenchmarkSliceAppend-8      50000   30000 ns/op   16384 B/op   10 allocs/op
BenchmarkSlicePrealloc-8   200000    8000 ns/op    8192 B/op    1 allocs/op
```

**结论**：预分配快 3.75 倍，内存分配少 10 倍。

---

## 5️⃣ 测试覆盖率（质量保证）

### 面试题

**如何查看测试覆盖率？覆盖率多少合适？**

### 5.1 生成覆盖率报告

**终端显示**：

```bash
$ go test -cover
PASS
coverage: 80.0% of statements
```

**生成详细报告**：

```bash
# 生成 coverage.out
$ go test -coverprofile=coverage.out

# 查看报告
$ go tool cover -func=coverage.out
example.go:10:  Add         100.0%
example.go:15:  Subtract    50.0%
example.go:20:  Multiply    0.0%
total:          (statements) 66.7%

# HTML 可视化
$ go tool cover -html=coverage.out
```

---

### 5.2 覆盖率目标

| 覆盖率 | 建议 |
|--------|------|
| < 50% | ❌ 不够 |
| 50-70% | ⚠️ 基本 |
| 70-85% | ✅ 良好 |
| > 85% | 🌟 优秀 |

**注意**：
* 覆盖率不是越高越好
* 重要逻辑必须覆盖
* 简单的 getter/setter 可以不测

---

## 6️⃣ pprof 性能分析（高级）

### 面试题

**如何排查性能问题？pprof 怎么用？**

### 6.1 CPU Profile

**生成 CPU profile**：

```bash
$ go test -bench=. -cpuprofile=cpu.prof
```

**分析**：

```bash
$ go tool pprof cpu.prof
(pprof) top
(pprof) list funcName
(pprof) web  # 生成调用图
```

**示例**：

```
(pprof) top
Showing nodes accounting for 2.5s, 90% of 2.8s total
      flat  flat%   sum%        cum   cum%
     1.2s 42.86% 42.86%      1.5s 53.57%  main.Fibonacci
     0.8s 28.57% 71.43%      0.8s 28.57%  runtime.cgocall
     0.5s 17.86% 89.29%      0.5s 17.86%  runtime.memmove
```

---

### 6.2 内存 Profile

**生成 mem profile**：

```bash
$ go test -bench=. -memprofile=mem.prof
```

**分析**：

```bash
$ go tool pprof mem.prof
(pprof) top
(pprof) list funcName
```

---

### 6.3 与 Benchmark 结合

```go
func BenchmarkProcess(b *testing.B) {
    b.ReportAllocs()  // 报告内存分配
    
    for i := 0; i < b.N; i++ {
        process()
    }
}
```

**运行**：

```bash
$ go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
```

**分析**：

```bash
# CPU 瓶颈
$ go tool pprof cpu.prof
(pprof) top

# 内存瓶颈
$ go tool pprof mem.prof
(pprof) top
```

---

## 7️⃣ 面试高频题汇总

### Q1：测试文件的命名规则？

**答**：
* 文件名：`xxx_test.go`
* 函数签名：`func TestXxx(t *testing.T)`

### Q2：testing.T 的常用方法？

| 方法 | 作用 |
|------|------|
| `t.Errorf()` | 记录错误，继续 |
| `t.Fatalf()` | 记录错误，终止 |
| `t.Logf()` | 记录日志 |
| `t.Skip()` | 跳过测试 |
| `t.Parallel()` | 并行执行 |

### Q3：如何运行指定测试？

**答**：`go test -run TestName`

### Q4：什么是表驱动测试？

**答**：用数据表定义测试用例，循环执行。

### Q5：t.Run 的作用？

**答**：
* 创建子测试
* 隔离失败
* 支持并行

### Q6：为什么需要 Mock？

**答**：
* 外部服务慢/不稳定
* 难以模拟异常
* 隔离依赖

### Q7：如何实现 Mock？

**答**：
* 通过 interface 抽象
* 手动 Mock 或使用 gomock

### Q8：Benchmark 的函数签名？

**答**：`func BenchmarkXxx(b *testing.B)`

### Q9：b.N 是什么？

**答**：
* 自动调整的循环次数
* 保证测试运行足够长时间
* 通常 1-10 亿次

### Q10：b.ResetTimer 的作用？

**答**：排除 setup 代码的时间影响。

### Q11：如何查看内存分配？

**答**：`go test -bench=. -benchmem`

### Q12：如何生成覆盖率报告？

**答**：

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Q13：覆盖率多少合适？

**答**：
* 70-85%：良好
* > 85%：优秀

### Q14：如何使用 pprof？

**答**：

```bash
# 生成 profile
go test -bench=. -cpuprofile=cpu.prof

# 分析
go tool pprof cpu.prof
```

### Q15：pprof 有哪些 profile 类型?

| 类型 | 作用 |
|------|------|
| cpu | CPU 瓶颈 |
| mem | 内存分配 |
| heap | 堆内存 |
| goroutine | goroutine 数量 |
| block | 阻塞分析 |
| mutex | 锁竞争 |

---

## 🔥 追问升级：实战场景

### Q：如何测试 HTTP Handler？

**答**：

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/user/123", nil)
    w := httptest.NewRecorder()
    
    handler(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
    }
    
    body := w.Body.String()
    if !strings.Contains(body, "Alice") {
        t.Errorf("body missing Alice: %s", body)
    }
}
```

### Q：如何测试并发代码？

**答**：

```go
func TestConcurrent(t *testing.T) {
    var wg sync.WaitGroup
    counter := NewSafeCounter()
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc()
        }()
    }
    
    wg.Wait()
    
    if counter.Value() != 1000 {
        t.Errorf("got %d, want 1000", counter.Value())
    }
}
```

### Q：如何测试 time.After？

**答**：用接口抽象时间

```go
type Clock interface {
    After(d time.Duration) <-chan time.Time
}

type RealClock struct{}

func (c *RealClock) After(d time.Duration) <-chan time.Time {
    return time.After(d)
}

type MockClock struct {
    ch chan time.Time
}

func (c *MockClock) After(d time.Duration) <-chan time.Time {
    return c.ch
}

// 测试
func TestTimeout(t *testing.T) {
    mockClock := &MockClock{ch: make(chan time.Time, 1)}
    mockClock.ch <- time.Now()  // 立即触发
    
    // 使用 mockClock 测试超时逻辑
}
```

### Q：如何提高 Benchmark 准确性？

**答**：

```bash
# 多次运行
$ go test -bench=. -benchtime=10s -count=5

# 禁用内联
$ go test -bench=. -gcflags="-l"

# 固定 CPU
$ GOMAXPROCS=1 go test -bench=.
```

---

**本章总结**：

✅ 单元测试：文件命名 `xxx_test.go`/函数签名 `TestXxx(t *testing.T)`/testing.T 方法  
✅ 表驱动测试：数据表定义用例/t.Run 创建子测试/并行执行 t.Parallel()  
✅ Mock：interface 抽象依赖/手动 Mock/gomock 自动生成  
✅ Benchmark：函数签名 `BenchmarkXxx(b *testing.B)`/b.N 自动调整/b.ResetTimer 排除 setup  
✅ 覆盖率：`go test -cover`/HTML 可视化/70-85% 良好  
✅ pprof：CPU profile 找瓶颈/内存 profile 找泄漏/与 Benchmark 结合  
✅ 实战：HTTP 测试 httptest/并发测试 sync/时间测试 Mock Clock
