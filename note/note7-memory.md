# 七、内存管理与 GC 原理

## 1️⃣ 内存分配策略（高频）

### 面试题

**Go 的内存分配器是如何设计的？**

### 标准答案

Go 的内存分配器基于 **TCMalloc**（Thread-Caching Malloc）思想，核心特点：

1. **多级缓存**：减少锁竞争
   * **mcache**（线程缓存）：每个 P 独享，无锁分配
   * **mcentral**（中心缓存）：按 size class 划分，有锁
   * **mheap**（堆）：全局，管理大对象和向 OS 申请内存

2. **Size Class**（规格类）：
   * 67 个规格（8B ~ 32KB）
   * 减少内存碎片
   * 快速匹配合适的块

3. **对象分类分配**：
   * **微对象**（< 16B）：组合分配
   * **小对象**（16B ~ 32KB）：从 mcache 的 mspan 分配
   * **大对象**（> 32KB）：直接从 mheap 分配

---

### 内存分配流程

```
                  ┌──────────────┐
                  │  Object Size │
                  └───────┬──────┘
                          │
          ┌───────────────┼───────────────┐
          │               │               │
      < 16B          16B ~ 32KB        > 32KB
     微对象             小对象           大对象
          │               │               │
          ↓               ↓               ↓
   ┌──────────┐    ┌──────────┐    ┌──────────┐
   │  微对象   │    │  mcache  │    │  mheap   │
   │  组合分配 │    │ (P独享)  │    │ (全局)   │
   └──────────┘    └────┬─────┘    └──────────┘
                        │ miss
                        ↓
                  ┌──────────┐
                  │ mcentral │
                  │ (size类) │
                  └────┬─────┘
                       │ miss
                       ↓
                  ┌──────────┐
                  │  mheap   │
                  │ (向OS申请)│
                  └──────────┘
```

---

### 1.1 微对象分配（< 16B）

**特点**：多个微对象**组合到一个 16B 的内存块**，减少浪费。

```go
type microBlock struct {
    ptr    uintptr  // 当前 16B 块的位置
    offset uintptr  // 已用字节数
}

// 分配 8B 微对象
func allocMicro(size uintptr) unsafe.Pointer {
    if mcache.microBlock.offset + size <= 16 {
        ptr := mcache.microBlock.ptr + mcache.microBlock.offset
        mcache.microBlock.offset += size
        return unsafe.Pointer(ptr)
    }
    // 当前块满了，申请新的 16B 块
    mcache.microBlock = newMicroBlock()
    return allocMicro(size)
}
```

**示例**：

```go
// 3 个小对象组合到 1 个 16B 块
var a int8    // 1B
var b int32   // 4B
var c int64   // 8B
// 总共 13B，分配到同一个 16B 的 mspan
```

---

### 1.2 小对象分配（16B ~ 32KB）

**流程**：

1. **计算 size class**：根据对象大小映射到 67 个规格之一
2. **从 mcache 获取 mspan**：每个 size class 对应一个 mspan
3. **从 mspan 分配对象**：找到空闲位置
4. **mspan 满了**：从 mcentral 获取新的 mspan
5. **mcentral 也没有**：从 mheap 获取

**Size Class 表**（部分）：

| class | bytes/obj | objects | waste bytes |
|-------|-----------|---------|-------------|
| 1     | 8         | 512     | 0           |
| 2     | 16        | 256     | 0           |
| 3     | 24        | 170     | 8           |
| ...   | ...       | ...     | ...         |
| 67    | 32768     | 1       | 0           |

**代码示例**：

```go
// 分配 20B 对象
// 1. 映射到 size class 3（24B）
// 2. 从 mcache.alloc[3] 获取 mspan
// 3. 从 mspan 分配 24B（浪费 4B）

func mallocSmall(size uintptr) unsafe.Pointer {
    sizeClass := getSizeClass(size)
    span := mcache.alloc[sizeClass]
    
    if span == nil || span.isFull() {
        span = mcentral.get(sizeClass)
        mcache.alloc[sizeClass] = span
    }
    
    return span.allocate()
}
```

---

### 1.3 大对象分配（> 32KB）

**特点**：直接从 **mheap** 分配，不经过 mcache 和 mcentral。

```go
func mallocLarge(size uintptr) unsafe.Pointer {
    // 向 mheap 申请多个连续的页（页大小 8KB）
    npages := (size + pageSize - 1) / pageSize
    span := mheap.allocSpan(npages)
    return unsafe.Pointer(span.base)
}
```

**为什么不缓存？**
* 大对象很少，缓存收益小
* 占用内存多，缓存浪费内存

---

## 2️⃣ 逃逸分析（必考）

### 面试题

**什么是逃逸分析？哪些情况会发生逃逸？如何查看？**

### 2.1 逃逸分析的定义

**逃逸分析**（Escape Analysis）：编译器分析变量的**作用域**，决定分配到**栈**还是**堆**。

* **栈分配**：函数返回后自动释放，快（无 GC）
* **堆分配**：需要 GC 回收，慢

**原则**：
* 能在栈上分配就不分配到堆
* 逃逸到堆的条件：变量**超出函数作用域**

---

### 2.2 常见逃逸场景

#### 场景 1：返回局部变量的指针

```go
// ❌ 逃逸到堆
func NewUser() *User {
    u := User{Name: "Alice"}
    return &u  // u 逃逸（生命周期超出函数）
}

// ✅ 不逃逸
func NewUser() User {
    u := User{Name: "Alice"}
    return u  // 值拷贝，u 在栈上
}
```

#### 场景 2：interface 参数

```go
// ❌ 逃逸到堆
func Print(v interface{}) {
    fmt.Println(v)  // v 逃逸
}

func main() {
    x := 42
    Print(x)  // x 逃逸到堆（interface 需要装箱）
}

// ✅ 不逃逸
func Print(v int) {
    fmt.Println(v)
}
```

**原因**：`interface{}` 需要存储**类型信息**，通常在堆上。

#### 场景 3：闭包引用外部变量

```go
// ❌ 逃逸到堆
func outer() func() int {
    x := 10
    return func() int {
        return x  // x 逃逸（被闭包捕获）
    }
}

// ✅ 不逃逸
func outer() func() int {
    return func() int {
        x := 10
        return x  // x 在栈上
    }
}
```

#### 场景 4：切片/map 动态扩容

```go
// ❌ 可能逃逸
func makeSlice() []int {
    s := make([]int, 0, 10)
    for i := 0; i < 100; i++ {
        s = append(s, i)  // 扩容时逃逸到堆
    }
    return s
}

// ✅ 不逃逸（大小确定）
func makeSlice() [10]int {
    var arr [10]int
    for i := 0; i < 10; i++ {
        arr[i] = i
    }
    return arr
}
```

#### 场景 5：变量太大（> 64KB）

```go
// ❌ 逃逸到堆
func largeArray() {
    arr := [1000000]int{}  // 4MB，逃逸到堆
    _ = arr
}
```

**原因**：栈空间有限（默认 8KB ~ 2MB），大对象在堆上。

#### 场景 6：发送指针或包含指针的值到 channel

```go
// ❌ 逃逸到堆
func sendToChan(ch chan *int) {
    x := 10
    ch <- &x  // x 逃逸
}

// ✅ 不逃逸
func sendToChan(ch chan int) {
    x := 10
    ch <- x
}
```

---

### 2.3 查看逃逸分析

**编译时查看**：

```bash
go build -gcflags="-m" main.go
```

**输出**：

```
./main.go:5:6: moved to heap: u
./main.go:6:9: &u escapes to heap
```

**详细模式**：

```bash
go build -gcflags="-m -m" main.go  # 更详细
go build -gcflags="-m -l" main.go  # 禁用内联
```

**示例**：

```go
package main

import "fmt"

func main() {
    x := 10
    fmt.Println(x)
}
```

**输出**：

```
./main.go:6:13: ... argument does not escape
./main.go:6:13: x escapes to heap
```

---

### 2.4 逃逸分析的性能影响

**benchmark**：

```go
// 栈分配
func BenchmarkStack(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = User{Name: "Alice"}
    }
}

// 堆分配
func BenchmarkHeap(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = &User{Name: "Alice"}
    }
}
```

**结果**：

```
BenchmarkStack-8   1000000000   0.3 ns/op   0 B/op   0 allocs/op
BenchmarkHeap-8     50000000   30 ns/op   48 B/op   1 allocs/op
```

**差距 100 倍**！

---

## 3️⃣ GC 原理（重点）

### 面试题

**Go 的 GC 算法是什么？三色标记法如何工作？**

### 3.1 三色标记法（Tri-color Marking）

**原理**：将对象标记为三种颜色：

* **白色**：未访问（潜在垃圾）
* **灰色**：已访问，子对象未访问
* **黑色**：已访问，子对象已访问

**标记过程**：

```
1. 初始：所有对象都是白色
2. 根对象（栈、全局变量）标记为灰色
3. 循环：
   - 选一个灰色对象
   - 标记为黑色
   - 它引用的白色对象标记为灰色
4. 没有灰色对象时结束
5. 清理：回收所有白色对象
```

**示例**：

```
初始状态：
根 → A → B → C
全为白色

第 1 步：标记根
灰：根
白：A, B, C

第 2 步：扫描根
黑：根
灰：A
白：B, C

第 3 步：扫描 A
黑：根, A
灰：B
白：C

第 4 步：扫描 B
黑：根, A, B
灰：C
白：-

第 5 步：扫描 C
黑：根, A, B, C
灰：-
白：-

结束：无白色对象，无需清理
```

---

### 3.2 写屏障（Write Barrier）

**问题**：并发标记时，黑色对象引用白色对象会导致**漏标记**。

```
标记中：
黑：A
灰：B
白：C

程序执行：
A.ref = C  // 黑色对象 A 引用白色对象 C
B.ref = nil

结果：C 会被误回收！
```

**解决**：**写屏障**（在赋值时插入检查）

**Dijkstra 插入写屏障**：

```go
// 伪代码
func WriteBarrier(obj, ref) {
    if ref.color == White {
        ref.color = Grey  // 标记为灰色
    }
    obj.ref = ref
}
```

**Go 1.8+ 混合写屏障**：

* 结合 **Dijkstra 插入屏障** 和 **Yuasa 删除屏障**
* 所有新创建的对象标记为黑色
* 减少 STW 时间

---

### 3.3 STW（Stop The World）时机

**STW 阶段**：

1. **标记准备**（~10-30µs）：
   * 开启写屏障
   * 扫描栈（标记根对象）

2. **标记终止**（~10-30µs）：
   * 关闭写屏障
   * 清理白色对象

**大部分时间并发**：
* **并发标记**：程序继续运行
* **并发清理**：程序继续运行

**Go 1.5+**：STW < 1ms（目标）

---

### 3.4 GC 触发条件

**触发时机**：

1. **内存阈值**：`堆内存增长 > GOGC%`
   * 默认 GOGC=100（堆翻倍触发）
   * 上次 GC 后堆大小 50MB，GOGC=100，则堆达到 100MB 时触发

2. **时间间隔**：2 分钟未 GC，强制触发

3. **手动触发**：`runtime.GC()`

**公式**：

```
下次 GC 堆大小 = 上次 GC 后存活对象大小 * (1 + GOGC/100)
```

**示例**：

```
上次 GC 后：100MB 存活
GOGC=100：下次 GC 触发点 = 100 * (1 + 100/100) = 200MB
GOGC=50：下次 GC 触发点 = 100 * (1 + 50/100) = 150MB
```

---

## 4️⃣ GC 性能调优（实战）

### 面试题

**如何减少 GC 压力？有哪些调优方法？**

### 4.1 减少堆分配

#### 方法 1：使用栈分配

```go
// ❌ 堆分配
func process() {
    u := &User{Name: "Alice"}  // 逃逸到堆
    handle(u)
}

// ✅ 栈分配
func process() {
    u := User{Name: "Alice"}  // 在栈上
    handle(&u)  // 不逃逸
}
```

#### 方法 2：避免 interface{}

```go
// ❌ 频繁装箱
func log(v interface{}) {
    fmt.Println(v)  // v 逃逸
}

// ✅ 泛型（Go 1.18+）
func log[T any](v T) {
    fmt.Println(v)  // 可能不逃逸
}
```

#### 方法 3：预分配切片容量

```go
// ❌ 频繁扩容
func process() []int {
    var s []int
    for i := 0; i < 1000; i++ {
        s = append(s, i)  // 扩容多次
    }
    return s
}

// ✅ 预分配
func process() []int {
    s := make([]int, 0, 1000)  // 一次分配
    for i := 0; i < 1000; i++ {
        s = append(s, i)
    }
    return s
}
```

---

### 4.2 对象池复用（sync.Pool）

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    buf.WriteString("hello")
    // ...
}
```

**注意**：
* `sync.Pool` 不保证对象不被回收
* 适合临时对象（如 buffer、连接）
* 每次 GC 会清空 Pool

---

### 4.3 调整 GOGC

```bash
# 环境变量
export GOGC=200  # GC 触发更晚（减少 GC 次数，内存占用更多）

# 代码设置
debug.SetGCPercent(200)
```

**权衡**：
* **GOGC 大**：GC 少，延迟低，内存高
* **GOGC 小**：GC 多，延迟高，内存低

**场景**：
* 内存充足：GOGC=200
* 内存紧张：GOGC=50

---

### 4.4 监控 GC 指标

#### 方法 1：runtime.ReadMemStats

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)

fmt.Printf("Alloc = %v MB\n", m.Alloc / 1024 / 1024)
fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc / 1024 / 1024)
fmt.Printf("Sys = %v MB\n", m.Sys / 1024 / 1024)
fmt.Printf("NumGC = %v\n", m.NumGC)
fmt.Printf("PauseNs = %v\n", m.PauseNs[(m.NumGC+255)%256])
```

#### 方法 2：GODEBUG=gctrace=1

```bash
$ GODEBUG=gctrace=1 ./app

gc 1 @0.002s 3%: 0.018+1.2+0.003 ms clock, 0.14+0.35/1.0/0+0.027 ms cpu, 4->4->1 MB, 5 MB goal, 8 P
```

**含义**：
* `gc 1`：第 1 次 GC
* `3%`：GC 占用 CPU 3%
* `4->4->1 MB`：GC 前堆 4MB → 标记后 4MB → 存活 1MB
* `5 MB goal`：下次 GC 触发点

---

## 5️⃣ 内存泄漏排查（重要）

### 面试题

**Go 有哪些常见的内存泄漏场景？如何排查？**

### 5.1 常见泄漏场景

#### 场景 1：goroutine 泄漏

```go
// ❌ goroutine 永远不退出
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // 永远阻塞
        fmt.Println(val)
    }()
}

// ✅ 使用 context 控制
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return  // 退出
        }
    }()
}
```

#### 场景 2：time.Ticker 未停止

```go
// ❌ Ticker 泄漏
func leak() {
    ticker := time.NewTicker(time.Second)
    // forget to ticker.Stop()
}

// ✅ 停止 Ticker
func noLeak() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // ...
        }
    }
}
```

#### 场景 3：切片截取导致底层数组泄漏

```go
// ❌ 切片持有大数组引用
func leak() []byte {
    data := make([]byte, 1024*1024)  // 1MB
    return data[:10]  // 返回 10B，但底层数组 1MB 无法释放
}

// ✅ 拷贝需要的数据
func noLeak() []byte {
    data := make([]byte, 1024*1024)
    result := make([]byte, 10)
    copy(result, data[:10])
    return result  // 底层数组可被 GC
}
```

#### 场景 4：map 只增不减

```go
// ❌ map 无限增长
var cache = make(map[string][]byte)

func addCache(key string, val []byte) {
    cache[key] = val  // 从不删除
}

// ✅ 使用 LRU 或定期清理
type LRUCache struct {
    // ...
}

func (c *LRUCache) Add(key string, val []byte) {
    c.cache[key] = val
    if len(c.cache) > c.maxSize {
        delete(c.cache, c.oldest())
    }
}
```

---

### 5.2 排查工具：pprof heap

**生成 heap profile**：

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // your app
}
```

**抓取 heap**：

```bash
# 实时查看
go tool pprof http://localhost:6060/debug/pprof/heap

# 保存到文件
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

**常用命令**：

```bash
(pprof) top  # 内存占用 Top 10
(pprof) list funcName  # 查看函数内存分配
(pprof) web  # 生成调用图（需要 graphviz）
```

---

### 5.3 排查 goroutine 泄漏

```bash
# 查看 goroutine 数量
curl http://localhost:6060/debug/pprof/goroutine?debug=1

# 保存到文件
curl http://localhost:6060/debug/pprof/goroutine > goroutine.prof
go tool pprof goroutine.prof
```

**分析**：

```bash
(pprof) top
# 如果某个函数创建了大量 goroutine，说明可能泄漏

(pprof) list funcName
# 查看具体代码
```

---

## 6️⃣ 面试高频题汇总

### Q1：Go 的内存分配器基于什么思想？

**答**：**TCMalloc**（Thread-Caching Malloc）
* 多级缓存（mcache/mcentral/mheap）
* Size Class 规格
* 微对象/小对象/大对象分级

### Q2：什么是 Size Class？

**答**：67 个内存规格（8B ~ 32KB），对象分配时向上取整到最近的规格。

### Q3：mcache、mcentral、mheap 的区别？

| 层级 | 范围 | 锁 | 作用 |
|------|------|-----|------|
| mcache | 每个 P | 无锁 | 快速分配小对象 |
| mcentral | 全局（按 size class） | 有锁 | 补充 mcache |
| mheap | 全局 | 有锁 | 管理大对象和向 OS 申请 |

### Q4：什么是逃逸分析？

**答**：编译器分析变量作用域，决定分配到栈还是堆。

### Q5：哪些情况会逃逸到堆？

**答**：
1. 返回局部变量指针
2. interface{} 参数
3. 闭包捕获外部变量
4. 切片/map 动态扩容
5. 变量太大（> 64KB）
6. 发送指针到 channel

### Q6：如何查看逃逸分析？

**答**：`go build -gcflags="-m" main.go`

### Q7：Go 的 GC 算法是什么？

**答**：**三色标记法** + **并发标记** + **混合写屏障**

### Q8：三色标记法的三种颜色？

**答**：
* 白色：未访问（潜在垃圾）
* 灰色：已访问，子对象未访问
* 黑色：已访问，子对象已访问

### Q9：什么时候触发 GC？

**答**：
1. 堆内存增长超过 GOGC%（默认 100%）
2. 2 分钟未 GC
3. 手动调用 `runtime.GC()`

### Q10：如何减少 GC 压力？

**答**：
1. 减少堆分配（栈分配、避免 interface、预分配容量）
2. 对象池复用（sync.Pool）
3. 调整 GOGC
4. 减少指针数量

### Q11：GOGC 的作用？

**答**：控制 GC 触发频率
* GOGC=100：堆翻倍触发
* GOGC=200：GC 少，内存高
* GOGC=50：GC 多，内存低

### Q12：如何监控 GC 性能？

**答**：
* `runtime.ReadMemStats`
* `GODEBUG=gctrace=1`
* pprof

### Q13：常见的内存泄漏场景？

**答**：
1. goroutine 泄漏（未退出）
2. time.Ticker 未停止
3. 切片截取持有大数组
4. map 只增不减

### Q14：如何排查内存泄漏？

**答**：
* `pprof heap`：查看内存分配
* `pprof goroutine`：查看 goroutine 数量
* `GODEBUG=gctrace=1`：监控 GC

---

## 🔥 追问升级：实战场景

### Q：sync.Pool 的原理？

**答**：

```go
type Pool struct {
    local unsafe.Pointer  // 每个 P 的本地缓存
    New func() interface{}  // 创建新对象
}

// Get 流程
func (p *Pool) Get() interface{} {
    // 1. 先从当前 P 的 local 获取
    // 2. 没有则从其他 P steal
    // 3. 还没有调用 New 创建
}
```

**注意**：每次 GC 会清空 Pool

### Q：哪些标准库使用了 sync.Pool？

**答**：
* `fmt.Printf`：复用 buffer
* `encoding/json`：复用 encoder/decoder
* `net/http`：复用 buffer

### Q：如何优化高频小对象分配？

**答**：

```go
// 方法 1：对象池
var pool = sync.Pool{
    New: func() interface{} {
        return &Object{}
    },
}

// 方法 2：批量分配
type ObjectPool struct {
    objects []*Object
    index   int
}

func (p *ObjectPool) Get() *Object {
    if p.index >= len(p.objects) {
        p.objects = make([]*Object, 1000)  // 批量分配
        p.index = 0
    }
    obj := p.objects[p.index]
    p.index++
    return obj
}
```

---

**本章总结**：

✅ 内存分配器：TCMalloc 思想/三级缓存/Size Class/微小大对象  
✅ 逃逸分析：栈 vs 堆/6 种逃逸场景/gcflags -m 查看/性能差 100 倍  
✅ GC 原理：三色标记法/写屏障/并发标记/STW < 1ms  
✅ GC 触发：GOGC%/2 分钟/手动触发  
✅ 性能调优：减少堆分配/sync.Pool 复用/调整 GOGC/监控指标  
✅ 内存泄漏：goroutine 泄漏/Ticker 未停/切片底层数组/map 无限增长  
✅ 排查工具：pprof heap/pprof goroutine/gctrace
