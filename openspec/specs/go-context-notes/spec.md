# go-context-notes Specification

## Purpose
TBD - created by archiving change add-go-interview-notes. Update Purpose after archive.
## Requirements
### Requirement: 解释 context 的设计目的

文档 SHALL 解释 context 的设计目的，包括为什么需要 context、context 解决的问题、context 的核心功能。

#### Scenario: 说明 context 的作用
- **WHEN** 读者查阅 context 的用途
- **THEN** 文档说明 context 用于跨 API 边界和 goroutine 传递取消信号、超时时间、截止时间和请求范围的值

#### Scenario: 说明 context 的设计原则
- **WHEN** 读者理解 context 设计理念
- **THEN** 文档说明 context 不可变、可派生、取消信号单向传播的设计原则

#### Scenario: 说明 context 的典型场景
- **WHEN** 读者需要了解使用场景
- **THEN** 文档说明 HTTP 请求处理、数据库查询、RPC 调用等需要超时控制的场景

### Requirement: 讲解 context 的四种创建方式

文档 SHALL 讲解 context 的四种创建方式，包括 Background、TODO、WithCancel、WithTimeout、WithDeadline、WithValue。

#### Scenario: 说明 Background 和 TODO
- **WHEN** 读者查阅根 context 创建
- **THEN** 文档说明 Background 用于 main、init、测试，TODO 用于不确定使用哪个 context 的过渡阶段

#### Scenario: 说明 WithCancel
- **WHEN** 读者需要手动取消 context
- **THEN** 文档提供 WithCancel 的使用示例，说明如何通过 cancel 函数主动取消

#### Scenario: 说明 WithTimeout 和 WithDeadline
- **WHEN** 读者需要超时控制
- **THEN** 文档说明 WithTimeout 指定相对时间，WithDeadline 指定绝对时间，超时后自动取消

#### Scenario: 说明 WithValue
- **WHEN** 读者需要传递请求范围的值
- **THEN** 文档说明 WithValue 用于传递 traceID、userID 等请求信息，key 应使用自定义类型避免冲突

### Requirement: 说明 context 的取消传播机制

文档 SHALL 说明 context 的取消传播机制，包括父 context 取消后子 context 的行为、Done channel 的使用、Err 方法的返回值。

#### Scenario: 说明取消的单向传播
- **WHEN** 读者查阅取消传播规则
- **THEN** 文档说明父 context 取消会级联取消所有子 context，但子 context 取消不影响父和兄弟 context

#### Scenario: 展示 Done channel 的使用
- **WHEN** 读者需要监听 context 取消
- **THEN** 文档提供通过 select 监听 ctx.Done() 的代码示例

#### Scenario: 说明 Err 方法
- **WHEN** 读者需要判断取消原因
- **THEN** 文档说明 Err 返回 Canceled（手动取消）或 DeadlineExceeded（超时），未取消时返回 nil

### Requirement: 说明 context 的使用规范

文档 SHALL 说明 context 的使用规范，包括作为函数首参数、不存储在结构体、不传 nil、命名约定。

#### Scenario: 说明 context 应作为首参数
- **WHEN** 读者设计函数签名
- **THEN** 文档说明 context 应作为函数首个参数，命名为 ctx

#### Scenario: 说明不应存储 context
- **WHEN** 读者考虑在结构体中存储 context
- **THEN** 文档说明不应在结构体中存储 context，每个方法应接收 context 参数

#### Scenario: 说明不应传递 nil context
- **WHEN** 读者不确定传什么 context
- **THEN** 文档说明应传递 context.Background() 或 context.TODO()，不要传 nil

### Requirement: 分析 context.Value 的使用和限制

文档 SHALL 分析 context.Value 的使用场景和限制，说明不应滥用 Value 传递参数。

#### Scenario: 说明 Value 的适用数据
- **WHEN** 读者评估是否使用 context.Value
- **THEN** 文档说明 Value 仅用于请求范围的元数据（如 traceID、认证信息），不应传递业务参数

#### Scenario: 说明 Value 的性能考虑
- **WHEN** 读者关注性能
- **THEN** 文档说明 Value 查找需要遍历 context 链，频繁访问应在入口处提取并直接传递

#### Scenario: 说明 key 的类型选择
- **WHEN** 读者选择 context.Value 的 key
- **THEN** 文档说明应使用自定义类型作为 key 避免冲突，提供代码示例

### Requirement: 提供面试高频题

文档 SHALL 提供 context 相关的面试高频题，每题包含标准答案和追问点。

#### Scenario: context 设计目的面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 context 解决的问题和设计原则的标准答案

#### Scenario: context 取消传播面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供取消单向传播的标准答案和代码示例

#### Scenario: context.Value 使用规范面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 Value 适用场景和限制的标准答案

