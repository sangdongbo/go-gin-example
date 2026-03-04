# Go 性能分析实战指南

本项目已集成 pprof 性能分析工具和常见性能问题演示接口，用于学习和实践 Go 性能优化。

## 快速开始

### 1. 启动服务

```bash
go run main.go
```

### 2. 访问 pprof Web UI

浏览器打开：http://localhost:8000/debug/pprof/

你会看到以下分析选项：
- **allocs**: 内存分配采样
- **block**: 阻塞分析
- **cmdline**: 命令行参数
- **goroutine**: goroutine 堆栈追踪
- **heap**: 堆内存分析
- **mutex**: 互斥锁竞争分析
- **profile**: CPU 性能分析（30秒采样）
- **threadcreate**: 线程创建分析
- **trace**: 执行追踪

## 四大性能问题演示

### 一、CPU 性能问题：slice 未预分配

#### 问题场景
访问 http://localhost:8000/api/v1/perf/cpu-slow

**问题代码：**
```go
data := make([]int, 0)  // 未预分配容量
for i := 0; i < 100000; i++ {
    data = append(data, i)  // 多次扩容
}
```

**现象：**
- CPU 占用高
- 热点在 `runtime.growslice` 和 `runtime.memmove`

#### 优化方案
访问 http://localhost:8000/api/v1/perf/cpu-fast

**优化代码：**
```go
data := make([]int, 0, 100000)  // 预分配容量
for i := 0; i < 100000; i++ {
    data = append(data, i)  // 无需扩容
}
```

#### pprof 分析

1. **压测慢接口：**
```bash
# 使用 ab 或 wrk 压测
ab -n 1000 -c 10 http://localhost:8000/api/v1/perf/cpu-slow
```

2. **采集 CPU profile：**
```bash
go tool pprof http://localhost:8000/debug/pprof/profile?seconds=30
```

3. **查看热点：**
```bash
(pprof) top
(pprof) list TestCPUSlow
(pprof) web  # 生成火焰图（需要 graphviz）
```

**预期结果：** 看到 `runtime.growslice` 占用大量 CPU

---

### 二、GC 压力问题：interface{} 滥用

#### 问题场景
访问 http://localhost:8000/api/v1/perf/gc-pressure

**问题代码：**
```go
result := make([]interface{}, 0)
for i := 0; i < 10000; i++ {
    item := map[string]interface{}{  // 装箱，逃逸到堆
        "id": i,
        "name": fmt.Sprintf("item_%d", i),
    }
    result = append(result, item)
}
```

**现象：**
- QPS 不高但延迟抖动
- GC 频繁触发

#### 优化方案
访问 http://localhost:8000/api/v1/perf/gc-optimized

**优化代码：**
```go
type Item struct {
    ID   int
    Name string
}
result := make([]Item, 0, 10000)  // 具体类型
```

#### pprof 分析

1. **查看堆内存分配：**
```bash
go tool pprof http://localhost:8000/debug/pprof/heap
```

2. **查看 GC 情况：**
```bash
# 启动时开启 GC trace
GODEBUG=gctrace=1 go run main.go
```

**预期结果：**
- 慢接口：大量 `interface{}` 对象在堆上
- 快接口：堆分配大幅减少

---

### 三、锁竞争问题：mutex vs atomic

#### 问题场景
访问 http://localhost:8000/api/v1/perf/mutex-slow

**问题代码：**
```go
var mu sync.Mutex
var total int

// 1000 个 goroutine 竞争
mu.Lock()
total++
mu.Unlock()
```

**现象：**
- goroutine 很多但 CPU 跑不满
- 大量时间花在锁等待

#### 优化方案
访问 http://localhost:8000/api/v1/perf/atomic-fast

**优化代码：**
```go
var total int64
atomic.AddInt64(&total, 1)  // 无锁操作
```

#### pprof 分析

1. **开启 mutex profiling：**
```go
import _ "net/http/pprof"
import "runtime"

func main() {
    runtime.SetMutexProfileFraction(1)  // 开启 mutex profiling
    // ...
}
```

2. **分析 mutex 竞争：**
```bash
go tool pprof http://localhost:8000/debug/pprof/mutex
```

**预期结果：**
- 慢接口：`sync.(*Mutex).Lock` 占比很高
- 快接口：无锁等待

---

### 四、goroutine 泄漏问题

#### 问题场景
访问 http://localhost:8000/api/v1/perf/goroutine-leak

**问题代码：**
```go
ch := make(chan int)
go func() {
    <-ch  // 永远阻塞，channel 永远不关闭
}()
```

**现象：**
- 服务越跑越慢
- 内存持续增长
- 最终 OOM

#### 优化方案
访问 http://localhost:8000/api/v1/perf/goroutine-clean

**优化代码：**
```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    select {
    case <-ctx.Done():
        return  // 正常退出
    }
}()
// ...
cancel()  // 通知退出
```

#### pprof 分析

1. **查看 goroutine 数量：**
```bash
go tool pprof http://localhost:8000/debug/pprof/goroutine
```

2. **查看堆栈：**
```bash
(pprof) top
(pprof) list TestGoroutineLeak
```

**预期结果：**
- 泄漏接口：大量 goroutine waiting on chan receive
- 正常接口：goroutine 正常退出

---

## 综合性能对比

访问 http://localhost:8000/api/v1/perf/benchmark

一次性对比所有优化效果，包括：
- Slice 预分配 vs 不预分配
- Mutex vs Atomic
- 性能提升倍数

---

## 实战练习流程

### 步骤 1：制造问题

1. 连续访问 `cpu-slow` 接口
2. 使用 `ab` 或 `wrk` 进行压测：
```bash
ab -n 10000 -c 100 http://localhost:8000/api/v1/perf/cpu-slow
```

### 步骤 2：采集 profile

```bash
# CPU profile
curl http://localhost:8000/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof

# Heap profile
curl http://localhost:8000/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Goroutine profile
curl http://localhost:8000/debug/pprof/goroutine > goroutine.prof
go tool pprof goroutine.prof
```

### 步骤 3：分析热点

```bash
(pprof) top        # 查看 top 函数
(pprof) list 函数名  # 查看具体代码
(pprof) web        # 火焰图（需要 graphviz）
(pprof) png        # 生成 PNG 图片
```

### 步骤 4：优化验证

访问对应的 `-fast` 或 `-optimized` 接口，对比性能差异

---

## pprof 常用命令速查

### Web UI
```bash
# 浏览器访问
http://localhost:8000/debug/pprof/
```

### 命令行分析

```bash
# CPU 分析（默认 30 秒采样）
go tool pprof http://localhost:8000/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:8000/debug/pprof/heap

# Goroutine 分析
go tool pprof http://localhost:8000/debug/pprof/goroutine

# Mutex 分析
go tool pprof http://localhost:8000/debug/pprof/mutex

# 交互式分析
(pprof) top          # 显示占用最多的函数
(pprof) top -cum     # 按累计时间排序
(pprof) list func    # 查看函数源码
(pprof) web          # 生成调用图
(pprof) png          # 导出 PNG
(pprof) pdf          # 导出 PDF
```

### 火焰图生成

```bash
# 1. 安装 graphviz
# Windows: choco install graphviz
# Mac: brew install graphviz
# Linux: apt-get install graphviz

# 2. 生成火焰图
go tool pprof -http=:8080 http://localhost:8000/debug/pprof/profile
```

---

## 性能优化检查清单

### ✅ CPU 优化
- [ ] Slice/Map 预分配容量
- [ ] 避免不必要的内存拷贝
- [ ] 减少反射使用
- [ ] 热路径避免字符串拼接（用 strings.Builder）

### ✅ 内存优化
- [ ] 使用具体类型代替 interface{}
- [ ] 复用对象（sync.Pool）
- [ ] 避免大对象频繁分配
- [ ] 及时释放不用的资源

### ✅ 并发优化
- [ ] 低竞争计数用 atomic
- [ ] 读多写少用 sync.RWMutex
- [ ] 分片锁减少竞争
- [ ] 合理控制 goroutine 数量

### ✅ goroutine 管理
- [ ] 使用 context 控制生命周期
- [ ] 及时关闭 channel
- [ ] 使用 WaitGroup 等待完成
- [ ] Worker Pool 限制并发数

---

## 面试常见问题

### Q1: 为什么预分配 slice 能提升性能？
**A:** 避免了多次扩容导致的 `runtime.growslice` 和 `memmove`，减少内存分配和拷贝开销。

### Q2: interface{} 为什么会导致 GC 压力？
**A:** interface{} 会引入装箱和逃逸，使原本栈上的对象进入堆，增加 GC 标记和扫描成本。

### Q3: mutex 和 atomic 怎么选？
**A:** 
- 简单计数/标志位 → atomic
- 复杂临界区 → mutex  
- 高并发热点 → 考虑分片或无锁结构

### Q4: goroutine 泄漏为什么会拖慢 GC？
**A:** 每个 goroutine 有栈和元数据，GC 需要扫描这些对象，数量增长会线性增加 GC 成本。

---

## 推荐学习路径

1. **理论学习**：阅读 `note/note2.md`
2. **实战演练**：访问本文档中的测试接口
3. **pprof 分析**：按照示例采集和分析 profile
4. **对比优化**：对比优化前后的性能差异
5. **总结归纳**：整理自己的性能优化经验

---

## 参考资料

- [Go pprof 官方文档](https://pkg.go.dev/net/http/pprof)
- [Go 性能优化实战](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
- [Understanding Allocations](https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/)
- [Go GC 原理](https://go.dev/blog/ismmkeynote)

---

祝学习愉快！🚀
