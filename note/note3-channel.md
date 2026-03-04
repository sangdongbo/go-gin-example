# 三、并发内功：Channel 与并发模式

## 1️⃣ Channel 基础机制（高频）

### 面试题

**无缓冲 channel 和有缓冲 channel 的区别？什么时候会阻塞？**

### 标准答案

#### 无缓冲 channel

* 容量为 0
* 发送操作：**立即阻塞**，直到有接收者
* 接收操作：**立即阻塞**，直到有发送者
* 特点：**同步通信**，必须配对进行

```go
ch := make(chan int)  // 无缓冲

// 如果没有接收者，下面会 deadlock
ch <- 1  // 阻塞
```

#### 有缓冲 channel

* 容量 > 0
* 发送操作：缓冲区**未满时不阻塞**，满了才阻塞
* 接收操作：缓冲区**非空时不阻塞**，空了才阻塞
* 特点：**异步通信**，允许暂存数据

```go
ch := make(chan int, 3)  // 缓冲容量 3

ch <- 1  // 不阻塞
ch <- 2  // 不阻塞
ch <- 3  // 不阻塞
ch <- 4  // 阻塞，缓冲区满
```

### 追问 1：关闭 channel 后会发生什么？

* **接收者**：可继续读取剩余值，读完后返回零值 + false
* **发送者**：panic（向已关闭 channel 发送）
* **重复关闭**：panic

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

v1, ok := <-ch  // 1, true
v2, ok := <-ch  // 2, true
v3, ok := <-ch  // 0, false（零值）

ch <- 3  // panic: send on closed channel
close(ch)  // panic: close of closed channel
```

### 追问 2：为什么建议由发送方关闭 channel？

* 接收方不知道发送方是否还会发送数据
* 如果接收方关闭，发送方继续发送会 panic
* **最佳实践**：谁创建谁关闭，由发送方负责关闭

### 追问 3：如何判断 channel 已关闭？

* 使用双返回值接收：`v, ok := <-ch`
* `ok == false` 表示 channel 已关闭且无数据
* 或者用 `range`：channel 关闭后自动退出

```go
for v := range ch {
    // channel 关闭后自动退出
}
```

---

## 2️⃣ select 多路复用（必考）

### 面试题

**select 语句的工作原理？有多个 case 就绪时怎么选择？**

### 标准答案

* select 用于监听多个 channel 操作
* **随机选择**：多个 case 同时就绪时，随机选择一个执行
* **阻塞行为**：所有 case 都阻塞时，select 阻塞
* **default 分支**：如果有 default，select 永不阻塞

```go
select {
case v := <-ch1:
    // ch1 有数据
case v := <-ch2:
    // ch2 有数据
case ch3 <- v:
    // 可以向 ch3 发送
default:
    // 所有 case 都未就绪
}
```

### 追问 1：为什么要随机选择？

* 避免**饥饿**：如果按顺序选择，后面的 case 可能永远得不到执行
* 公平性：每个 case 都有平等的机会
* runtime 实现：通过洗牌算法随机化 case 顺序

### 追问 2：select 的超时模式怎么写？

```go
select {
case v := <-ch:
    fmt.Println("received:", v)
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}
```

⚠️ 注意：`time.After` 会创建新 timer，频繁调用有 GC 压力

**优化版**：

```go
timer := time.NewTimer(1 * time.Second)
defer timer.Stop()

select {
case v := <-ch:
    fmt.Println("received:", v)
case <-timer.C:
    fmt.Println("timeout")
}
```

### 追问 3：空 select 会怎样？

```go
select {}  // 永久阻塞
```

* 常用于 main goroutine，防止程序退出
* 但更好的方式是 `<-ctx.Done()` 或 `<-make(chan struct{})`

---

## 3️⃣ 常见并发模式（实战）

### 面试题

**说说 Go 中常见的并发模式，如 worker pool、fan-in、pipeline？**

### 3.1 Worker Pool（限流）

**场景**：控制并发度，避免 goroutine 爆炸

```go
func workerPool(tasks []int, workerCount int) {
    taskChan := make(chan int, len(tasks))
    results := make(chan int, len(tasks))
    
    // 启动 worker
    for i := 0; i < workerCount; i++ {
        go func() {
            for task := range taskChan {
                // 处理任务
                results <- process(task)
            }
        }()
    }
    
    // 分发任务
    for _, task := range tasks {
        taskChan <- task
    }
    close(taskChan)
    
    // 收集结果
    for i := 0; i < len(tasks); i++ {
        <-results
    }
}
```

**追问**：如何优雅关闭 worker pool？

* 使用 `sync.WaitGroup` 等待所有 worker 完成
* 关闭 taskChan 让 worker 退出

```go
var wg sync.WaitGroup
for i := 0; i < workerCount; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for task := range taskChan {
            results <- process(task)
        }
    }()
}

// 发送任务
go func() {
    for _, task := range tasks {
        taskChan <- task
    }
    close(taskChan)
}()

wg.Wait()
close(results)
```

### 3.2 Fan-in（合并多个 channel）

**场景**：多个数据源，合并到一个 channel

```go
func fanIn(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    
    go func() {
        for {
            select {
            case v := <-ch1:
                out <- v
            case v := <-ch2:
                out <- v
            }
        }
    }()
    
    return out
}
```

**优化版**（支持关闭）：

```go
func fanIn(chs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range chs {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### 3.3 Fan-out（分发任务）

**场景**：一个输入源，多个 worker 处理

```go
func fanOut(in <-chan int, workers int) []<-chan int {
    outs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        out := make(chan int)
        outs[i] = out
        
        go func(ch chan int) {
            for v := range in {
                ch <- process(v)
            }
            close(ch)
        }(out)
    }
    
    return outs
}
```

### 3.4 Pipeline（流水线）

**场景**：多阶段数据处理

```go
// 阶段1：生成数据
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// 阶段2：平方
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// 阶段3：加倍
func double(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        close(out)
    }()
    return out
}

// 使用
func main() {
    in := generator(1, 2, 3, 4)
    out := double(square(in))
    
    for v := range out {
        fmt.Println(v)  // 2, 8, 18, 32
    }
}
```

**追问**：如何优雅取消 pipeline？

* 使用 `context.Context` 传递取消信号
* 每个阶段监听 `ctx.Done()`

```go
func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case n, ok := <-in:
                if !ok {
                    return
                }
                select {
                case out <- n * n:
                case <-ctx.Done():
                    return
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

---

## 4️⃣ 死锁场景分析（常见陷阱）

### 面试题

**在使用 channel 时，哪些情况会导致死锁？**

### 场景 1：无缓冲 channel 自发自收

```go
func main() {
    ch := make(chan int)
    ch <- 1  // 阻塞，等待接收者
    v := <-ch
    fmt.Println(v)
}
// fatal error: all goroutines are asleep - deadlock!
```

**原因**：main goroutine 在发送时阻塞，没有其他 goroutine 接收

**解决**：

1. 使用有缓冲 channel
2. 或在 goroutine 中发送

```go
ch := make(chan int, 1)
ch <- 1
v := <-ch

// 或
ch := make(chan int)
go func() { ch <- 1 }()
v := <-ch
```

### 场景 2：忘记关闭 channel

```go
func main() {
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    
    for v := range ch {
        fmt.Println(v)
    }
}
// deadlock: range 等待 channel 关闭
```

**原因**：`range` 会一直等待 channel，不关闭就无法退出

**解决**：发送完后关闭 channel

```go
ch <- 1
ch <- 2
ch <- 3
close(ch)  // 必须关闭

for v := range ch {
    fmt.Println(v)
}
```

### 场景 3：循环依赖

```go
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    go func() {
        v := <-ch1  // 等待 ch1
        ch2 <- v
    }()
    
    go func() {
        v := <-ch2  // 等待 ch2
        ch1 <- v
    }()
    
    time.Sleep(time.Second)
}
// deadlock: 两个 goroutine 相互等待
```

**原因**：goroutine1 等 ch1，goroutine2 等 ch2，形成环

**解决**：打破循环，先发送再接收

```go
go func() {
    ch1 <- 1  // 先发送
}()

go func() {
    v := <-ch1
    ch2 <- v
}()
```

### 场景 4：select 没有 default，所有 case 阻塞

```go
func main() {
    ch := make(chan int)
    
    select {
    case v := <-ch:  // ch 无数据，阻塞
        fmt.Println(v)
    }
}
// deadlock
```

**解决**：使用 default 或 timeout

```go
select {
case v := <-ch:
    fmt.Println(v)
default:
    fmt.Println("no data")
}
```

---

## 5️⃣ 面试高频题汇总

### Q1：无缓冲 vs 有缓冲 channel

**答**：
* 无缓冲：同步，发送和接收必须配对，阻塞直到对方准备好
* 有缓冲：异步，缓冲区未满/非空时不阻塞
* 使用场景：无缓冲用于强同步，有缓冲用于削峰、解耦

### Q2：channel 发送和接收的阻塞条件？

**答**：

| 操作 | 无缓冲 | 有缓冲 |
|-----|-------|-------|
| 发送 | 立即阻塞（无接收者） | 缓冲区满时阻塞 |
| 接收 | 立即阻塞（无发送者） | 缓冲区空时阻塞 |
| 关闭后发送 | panic | panic |
| 关闭后接收 | 返回零值 + false | 返回剩余值，读完返回零值 + false |

### Q3：select 的工作原理？

**答**：
1. 评估所有 case，检查是否就绪
2. 如果有多个就绪，**随机选择**一个执行
3. 如果所有都未就绪：
   - 有 default：执行 default
   - 无 default：阻塞，等待任一 case 就绪
4. 只执行一次，不循环

### Q4：如何实现超时控制？

**答**：

```go
select {
case v := <-ch:
    // 正常接收
case <-time.After(time.Second):
    // 超时处理
}
```

**注意**：`time.After` 每次调用创建新 timer，高频调用应使用 `time.NewTimer` 并复用

### Q5：channel 是否线程安全？

**答**：
* **是**，channel 内部有锁保护
* 多个 goroutine 同时读写 channel 不会有数据竞争
* 但关闭 channel 不是原子的：
  - 不能多次关闭
  - 不应从接收方关闭

### Q6：channel 的底层数据结构？

**答**：

```go
type hchan struct {
    qcount   uint           // 当前元素个数
    dataqsiz uint           // 缓冲区大小
    buf      unsafe.Pointer // 环形缓冲区
    elemsize uint16         // 元素大小
    closed   uint32         // 关闭标志
    sendx    uint           // 发送索引
    recvx    uint           // 接收索引
    recvq    waitq          // 接收等待队列
    sendq    waitq          // 发送等待队列
    lock     mutex          // 互斥锁
}
```

* 环形缓冲区：存储数据
* 两个等待队列：阻塞的 goroutine
* 锁保护：确保并发安全

### Q7：什么时候用有缓冲/无缓冲 channel？

**答**：

**无缓冲**：
* 需要**强同步**：确保发送者等到接收者
* 示例：任务完成确认、握手同步

**有缓冲**：
* **削峰填谷**：缓冲突发流量
* **解耦**：发送者和接收者独立工作
* 示例：任务队列、日志聚合

**经验**：
* 不确定用多大缓冲？从 0 开始
* 观察阻塞情况再调整
* 缓冲越大，内存占用越多

### Q8：如何安全关闭 channel？

**答**：

**原则**：
* 只由**发送方**关闭
* **不要**从接收方关闭
* **不要**关闭只读 channel

**多发送者场景**：
* 使用 `sync.Once` 确保只关闭一次
* 或使用额外的 `done` channel 广播关闭信号

```go
type SafeChannel struct {
    ch   chan int
    once sync.Once
}

func (s *SafeChannel) Close() {
    s.once.Do(func() {
        close(s.ch)
    })
}
```

---

## 🔥 追问升级：channel 性能优化

### Q：频繁 channel 操作有性能问题吗？

**答**：有，主要开销：

1. **锁开销**：每次发送/接收都要加锁
2. **goroutine 调度**：阻塞会触发调度
3. **内存分配**：创建 channel 需要堆分配

**优化方法**：

1. **批量操作**：减少 channel 操作次数
2. **使用缓冲**：减少阻塞和调度
3. **避免过度使用**：简单场景用 `sync.Mutex` 更快

### Q：channel 和锁哪个更快？

**答**：

* **锁更快**：无 goroutine 切换开销
* **channel 更清晰**：传递数据 + 同步信号

**经验**：
* 热点路径、高频操作：用锁
* 需要传递数据、跨 goroutine：用 channel
* "不要通过共享内存通信，而要通过通信共享内存"（Go 哲学）

---

**本章总结**：

✅ channel 是 Go 并发的核心，面试必考  
✅ 掌握无缓冲/有缓冲的区别和阻塞条件  
✅ 熟练使用 select 实现超时、多路复用  
✅ 理解 worker pool、fan-in/out、pipeline 等模式  
✅ 警惕死锁：自发自收、忘记关闭、循环依赖  
✅ 安全关闭原则：由发送方关闭，只关闭一次
