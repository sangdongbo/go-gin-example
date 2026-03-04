# 案例总览（你会学到什么）

我们会拆 4 类**最常见 Go 慢因**：

1. **CPU 慢**：忙在“你以为不慢”的地方
2. **内存慢 / GC 抖动**：QPS 不高却卡顿
3. **锁慢 / 并发退化**：goroutine 很多但跑不满
4. **goroutine 泄漏**：越跑越慢，最终 OOM

每一类我都给你：

* ❌ 问题代码
* 🔍 pprof / runtime 如何看
* 🧠 背后的 runtime 原因
* ✅ 正确修法
* 🎯 面试官追问标准答案

---

# 一、CPU 很高，但你“什么也没做”

## ❌ 问题场景（真实）

> 一个 API QPS 只有几百，但 CPU 占满

### 示例代码（简化版）

```go
func handle(w http.ResponseWriter, r *http.Request) {
    data := make([]int, 0)
    for i := 0; i < 100000; i++ {
        data = append(data, i)
    }
    fmt.Fprintln(w, "ok")
}
```

你可能会说：

> “不就 append 吗？能有多慢？”

---

## 🔍 pprof 怎么看

```go
import _ "net/http/pprof"
```

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

你看到的热点很可能是：

```
runtime.growslice
runtime.memmove
```

---

## 🧠 真正慢在哪（runtime 视角）

* slice **从 0 cap 开始**
* 多次扩容
* 每次扩容：

  * 新数组
  * memmove 拷贝旧数据
* CPU 全花在 **内存搬运**

---

## ✅ 正确修法

```go
data := make([]int, 0, 100000)
```

### 修复后 pprof

* growslice 消失
* CPU 使用率大幅下降

---

## 🎯 面试官追问

**Q：为什么预分配能显著提升性能？**

> 因为避免了 slice 多次扩容导致的内存拷贝和 cache miss，本质是减少 runtime.growslice 与 memmove 的开销。

---

# 二、QPS 不高，但服务“抖”——GC 在搞你

## ❌ 问题场景

* QPS 300
* 延迟周期性飙高
* CPU 有空闲

---

## 🔍 heap pprof

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

你看到：

```
大量 *interface{} / map[string]interface{}
```

---

## 🧠 runtime 真相

### 问题代码模式

```go
func process(v interface{}) {
    // ...
}
```

* interface 装箱
* 逃逸到堆
* 大量短生命周期对象
* GC 被频繁触发

---

## 🔍 进一步验证（GC trace）

```bash
GODEBUG=gctrace=1 ./app
```

你会看到：

```
gc 7 @0.45s 12%: 0.8+1.2+0.4 ms
```

👉 GC 在频繁跑，但你业务并不重

---

## ✅ 正确修法

* 避免无意义 interface{}
* 用具体类型
* 热路径避免 map[string]interface{}
* 可复用对象 → `sync.Pool`

---

## 🎯 面试官追问

**Q：为什么 interface 会导致 GC 压力？**

> interface 会引入装箱和逃逸，导致原本可在栈上分配的对象进入堆，增加 GC 标记与扫描成本。

---

# 三、并发很多，但跑不满 —— 锁在拖后腿

## ❌ 问题代码

```go
var mu sync.Mutex
var total int

func work() {
    mu.Lock()
    total++
    mu.Unlock()
}
```

1000 个 goroutine 跑这个。

---

## 🔍 mutex pprof

```bash
go tool pprof http://localhost:6060/debug/pprof/mutex
```

你看到：

```
sync.(*Mutex).Lock
```

占比极高。

---

## 🧠 runtime 原因

* mutex 竞争
* goroutine 挂起
* P 空转
* 上下文切换成本

---

## ✅ 正确修法

### 方案一：atomic

```go
atomic.AddInt64(&total, 1)
// 相比于使用 sync.Mutex（互斥锁），原子操作不需要上下文切换或内核态转换，性能更高，属于无锁（Lock-free）编程。
```

### 方案二：局部聚合 + 合并

```go
local := 0
// goroutine 内累加
```

---

## 🎯 面试官追问

**Q：mutex 和 atomic 怎么选？**

> 低冲突计数 → atomic
> 复杂临界区 → mutex
> 高并发热点 → 分片或局部化

## 📋 Mutex vs Atomic 快速对比

| 维度 | **atomic (原子操作)** | **sync.Mutex (互斥锁)** |
| --- | --- | --- |
| **底层原理** | CPU 指令级支持（CAS/Lock XADD） | 操作系统信号量 + 调度器挂起 |
| **保护范围** | 单个基础变量（int, pointer等） | 一段代码块或多个相关联的变量 |
| **开销** | **极低**（无系统调用，无调度） | **较高**（涉及协程切换、上下文损耗） |
| **编程模型** | 乐观锁策略（非阻塞） | 悲观锁策略（阻塞等待） |
| **易用性** | 难，容易写出逻辑错误 | 简单，符合直觉 |

* 能用 Atomic 解决的简单计数问题，不用 Mutex。
* 涉及多个变量关联操作，或者包含 IO、函数调用，直接用 Mutex。
* 如果追求极致性能且逻辑非常简单（如单写多读的标志位），用 Atomic。
* 如果不确定，请先用 Mutex，因为代码的可读性和正确性比那几纳秒的性能更重要。

---

# 四、越跑越慢 —— goroutine 泄漏

## ❌ 问题代码

```go
go func() {
    for {
        select {
        case <-ch:
            return
        }
    }
}()
```

但 `ch` 永远不 close。

---

## 🔍 goroutine pprof

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

你看到：

```
100000 goroutines waiting on select
```

---

## 🧠 runtime 视角

* goroutine 不退出
* 栈 + 调度成本持续存在
* GC 扫描成本上升

---

## ✅ 正确修法

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            return
        }
    }
}()
```

---

## 🎯 面试官追问

**Q：goroutine 泄漏为什么会拖慢 GC？**

> 每个 goroutine 都有栈和元数据，GC 需要扫描这些对象，数量增长会线性增加 GC 成本。

