# go-channel-concurrency-notes Specification

## Purpose
TBD - created by archiving change add-go-interview-notes. Update Purpose after archive.
## Requirements
### Requirement: 覆盖 channel 基础机制

文档 SHALL 详细解释 Go channel 的核心机制，包括无缓冲和有缓冲 channel 的区别、发送和接收操作的阻塞行为、channel 的关闭规则。

#### Scenario: 解释无缓冲 channel 阻塞行为
- **WHEN** 读者查阅无缓冲 channel 相关内容
- **THEN** 文档清晰说明发送和接收都会阻塞直到对方准备好

#### Scenario: 解释有缓冲 channel 特性
- **WHEN** 读者查阅有缓冲 channel 相关内容
- **THEN** 文档说明缓冲区未满时发送不阻塞，未空时接收不阻塞

#### Scenario: 说明 channel 关闭规则
- **WHEN** 读者查阅 channel 关闭相关内容
- **THEN** 文档说明关闭后接收者可继续读取剩余值，重复关闭会 panic，向已关闭 channel 发送会 panic

### Requirement: 讲解 select 多路复用

文档 SHALL 详细讲解 select 语句的使用，包括多个 case 的选择规则、default 分支的作用、select 的常见模式。

#### Scenario: 解释 select 随机选择机制
- **WHEN** 读者查阅 select 多个 case 同时就绪的行为
- **THEN** 文档说明 select 会随机选择一个可执行的 case

#### Scenario: 说明 default 的作用
- **WHEN** 读者查阅 select default 分支
- **THEN** 文档说明 default 使 select 非阻塞，所有 case 都未就绪时执行 default

#### Scenario: 展示超时模式
- **WHEN** 读者需要实现超时控制
- **THEN** 文档提供 select 配合 time.After 的代码示例

### Requirement: 说明常见并发模式

文档 SHALL 介绍 Go 中常见的并发模式，包括 worker pool、fan-in、fan-out、pipeline 等。

#### Scenario: 展示 worker pool 模式
- **WHEN** 读者需要控制并发数量
- **THEN** 文档提供 worker pool 的实现示例和使用场景

#### Scenario: 展示 fan-in 模式
- **WHEN** 读者需要合并多个 channel
- **THEN** 文档提供 fan-in 的实现示例

#### Scenario: 展示 pipeline 模式
- **WHEN** 读者需要构建数据处理流水线
- **THEN** 文档提供 pipeline 的实现示例

### Requirement: 分析死锁场景

文档 SHALL 分析常见的死锁场景，包括无缓冲 channel 自发自收、循环依赖、忘记关闭 channel 等问题。

#### Scenario: 说明自发自收死锁
- **WHEN** 读者查阅死锁问题
- **THEN** 文档展示 goroutine 向自己的无缓冲 channel 发送导致死锁的示例及解决方法

#### Scenario: 说明循环依赖死锁
- **WHEN** 读者查阅死锁问题
- **THEN** 文档展示多个 goroutine 相互等待对方的 channel 导致死锁的示例

#### Scenario: 说明忘记关闭 channel
- **WHEN** 读者查阅 range channel 相关问题
- **THEN** 文档说明 range 遍历 channel 时如果不关闭会导致死锁

### Requirement: 提供面试高频题

文档 SHALL 提供 channel 和并发相关的面试高频题，每题包含标准答案和追问点。

#### Scenario: 无缓冲 vs 有缓冲 channel 面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供该问题的标准答案和至少 2 个追问点

#### Scenario: select 原理面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 select 底层实现和使用场景的标准答案

#### Scenario: channel 关闭规则面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 channel 关闭的标准答案和异常场景分析

