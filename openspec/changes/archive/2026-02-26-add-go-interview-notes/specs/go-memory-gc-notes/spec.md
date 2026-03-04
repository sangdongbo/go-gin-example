## ADDED Requirements

### Requirement: 解释 Go 内存分配策略

文档 SHALL 解释 Go 的内存分配策略，包括 TCMalloc 思想、size class、span、mcache、mcentral、mheap 的层次结构。

#### Scenario: 说明 TCMalloc 思想
- **WHEN** 读者查阅 Go 内存分配器设计
- **THEN** 文档说明 Go 借鉴 TCMalloc 的多级缓存思想，减少锁竞争

#### Scenario: 说明 size class
- **WHEN** 读者查阅小对象分配
- **THEN** 文档说明小对象按 size class 分配，减少内存碎片

#### Scenario: 说明三级缓存结构
- **WHEN** 读者查阅内存分配流程
- **THEN** 文档说明 mcache（P 级）→ mcentral（全局）→ mheap（全局）→ 操作系统的分配层次

#### Scenario: 说明小对象、大对象分配
- **WHEN** 读者查阅不同大小对象的分配策略
- **THEN** 文档说明小对象（< 32KB）从 mcache 分配，大对象直接从 mheap 分配

### Requirement: 讲解逃逸分析

文档 SHALL 讲解 Go 的逃逸分析，包括逃逸的定义、常见逃逸场景、如何查看逃逸分析结果。

#### Scenario: 说明逃逸分析的定义
- **WHEN** 读者查阅逃逸分析概念
- **THEN** 文档说明逃逸分析决定变量分配在栈还是堆，逃逸指变量从栈转移到堆

#### Scenario: 列举常见逃逸场景
- **WHEN** 读者查阅哪些情况会逃逸
- **THEN** 文档列举返回局部变量指针、闭包捕获、interface 赋值、slice 扩容、大对象等逃逸场景

#### Scenario: 展示查看逃逸分析
- **WHEN** 读者需要查看逃逸分析结果
- **THEN** 文档提供 `go build -gcflags="-m"` 命令示例和输出解读

#### Scenario: 说明逃逸的性能影响
- **WHEN** 读者关注逃逸对性能的影响
- **THEN** 文档说明逃逸增加 GC 压力，栈分配比堆分配快，提供减少逃逸的优化建议

### Requirement: 详解 GC 原理

文档 SHALL 详解 Go 的垃圾回收原理，包括三色标记法、写屏障、STW 时机、GC 触发条件。

#### Scenario: 说明三色标记法
- **WHEN** 读者查阅 GC 算法
- **THEN** 文档说明白色（待回收）、灰色（待扫描）、黑色（已扫描）三色标记的工作过程

#### Scenario: 说明写屏障
- **WHEN** 读者查阅并发标记问题
- **THEN** 文档说明写屏障用于解决并发标记时对象引用变化导致的漏标问题

#### Scenario: 说明 STW 时机
- **WHEN** 读者查阅 STW（Stop The World）发生时机
- **THEN** 文档说明 Go 1.5+ STW 仅在标记开始和结束时短暂发生，标记过程并发进行

#### Scenario: 说明 GC 触发条件
- **WHEN** 读者查阅 GC 何时触发
- **THEN** 文档说明 GOGC 环境变量（默认 100）控制堆增长到上次 GC 后的 2 倍时触发，以及手动触发 runtime.GC()

### Requirement: 分析 GC 性能调优

文档 SHALL 分析 GC 性能调优方法，包括减少堆分配、对象池复用、调整 GOGC、监控 GC 指标。

#### Scenario: 说明减少堆分配
- **WHEN** 读者需要减少 GC 压力
- **THEN** 文档提供预分配 slice、避免逃逸、减少指针使用等优化建议

#### Scenario: 说明对象池复用
- **WHEN** 读者需要复用对象
- **THEN** 文档提供 sync.Pool 的使用示例和注意事项

#### Scenario: 说明 GOGC 调优
- **WHEN** 读者需要调整 GC 频率
- **THEN** 文档说明增大 GOGC 减少 GC 频率但增加内存使用，减小 GOGC 相反

#### Scenario: 说明 GC 监控
- **WHEN** 读者需要监控 GC 性能
- **THEN** 文档提供 GODEBUG=gctrace=1、runtime.ReadMemStats、pprof heap 等监控方法

### Requirement: 说明内存泄漏排查

文档 SHALL 说明 Go 中内存泄漏的常见原因和排查方法，包括 goroutine 泄漏、全局变量持有引用、time.Ticker 未关闭。

#### Scenario: 列举常见泄漏场景
- **WHEN** 读者排查内存泄漏
- **THEN** 文档列举 goroutine 泄漏、slice 持有大数组引用、map 持有大对象、time.Ticker 未 Stop 等场景

#### Scenario: 展示 pprof heap 分析
- **WHEN** 读者需要分析内存使用
- **THEN** 文档提供 pprof heap 的使用示例，包括 inuse_space、alloc_objects 等指标

#### Scenario: 展示 goroutine 泄漏排查
- **WHEN** 读者怀疑 goroutine 泄漏
- **THEN** 文档提供 pprof goroutine 和 runtime.NumGoroutine() 的使用方法

### Requirement: 提供面试高频题

文档 SHALL 提供内存管理和 GC 相关的面试高频题，每题包含标准答案和追问点。

#### Scenario: 逃逸分析面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供逃逸分析的定义、常见场景和性能影响的标准答案

#### Scenario: 三色标记法面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供三色标记法的工作流程和写屏障作用的标准答案

#### Scenario: GC 调优面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 GC 调优方法和监控指标的标准答案
