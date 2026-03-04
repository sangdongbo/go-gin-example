# go-defer-panic-recover-notes Specification

## Purpose
TBD - created by archiving change add-go-interview-notes. Update Purpose after archive.
## Requirements
### Requirement: 解释 defer 执行顺序

文档 SHALL 详细解释 defer 的执行顺序，包括后进先出（LIFO）规则、defer 参数求值时机、defer 在 return 和 panic 时的行为。

#### Scenario: 说明 defer 的 LIFO 顺序
- **WHEN** 读者查阅 defer 执行顺序
- **THEN** 文档说明多个 defer 按照后进先出的顺序执行，并提供代码示例

#### Scenario: 说明 defer 参数求值时机
- **WHEN** 读者查阅 defer 参数何时求值
- **THEN** 文档说明 defer 语句执行时就确定参数值，而非函数执行时

#### Scenario: 说明 defer 与 return 的关系
- **WHEN** 读者查阅 defer 对返回值的影响
- **THEN** 文档说明 defer 可以修改命名返回值，执行顺序是：赋值返回值 → 执行 defer → 返回

#### Scenario: 说明 defer 与 panic 的关系
- **WHEN** 读者查阅 panic 时 defer 的行为
- **THEN** 文档说明 panic 会触发当前函数的所有 defer 执行后再向上传播

### Requirement: 讲解 panic 和 recover 机制

文档 SHALL 讲解 panic 和 recover 的工作机制，包括 panic 的传播过程、recover 的使用规则、recover 能捕获的场景。

#### Scenario: 说明 panic 的传播过程
- **WHEN** 读者查阅 panic 如何传播
- **THEN** 文档说明 panic 会向上层调用栈传播，每层执行 defer 后继续向上，直到被 recover 或程序终止

#### Scenario: 说明 recover 的使用规则
- **WHEN** 读者查阅 recover 使用方法
- **THEN** 文档说明 recover 必须在 defer 函数中直接调用才有效，间接调用或 defer 外调用无效

#### Scenario: 说明 recover 能捕获的场景
- **WHEN** 读者查阅 recover 的限制
- **THEN** 文档说明 recover 只能捕获当前 goroutine 的 panic，无法跨 goroutine 捕获

#### Scenario: 展示 recover 的返回值
- **WHEN** 读者需要获取 panic 的值
- **THEN** 文档说明 recover 返回 panic 的参数值，无 panic 时返回 nil

### Requirement: 对比错误处理最佳实践

文档 SHALL 对比 Go 中错误处理的不同方式，包括 error 接口、panic/recover、错误包装（errors.Is/As）。

#### Scenario: 说明 error 和 panic 的使用场景
- **WHEN** 读者需要选择错误处理方式
- **THEN** 文档说明 error 用于可预期的错误，panic 用于不可恢复的异常情况

#### Scenario: 展示错误包装
- **WHEN** 读者需要包装和检查错误
- **THEN** 文档提供 fmt.Errorf("%w", err)、errors.Is 和 errors.As 的使用示例

#### Scenario: 说明 panic 的合理使用场景
- **WHEN** 读者评估是否使用 panic
- **THEN** 文档说明 init 函数失败、配置加载失败、不变式被破坏等场景可以使用 panic

### Requirement: 分析 defer 的性能影响

文档 SHALL 分析 defer 对性能的影响，说明 defer 的开销来源和优化方法。

#### Scenario: 说明 defer 的性能开销
- **WHEN** 读者关注 defer 性能
- **THEN** 文档说明 Go 1.13 前 defer 有较大开销，1.13+ 通过栈上分配优化后开销大幅降低

#### Scenario: 提供 defer 使用建议
- **WHEN** 读者评估是否使用 defer
- **THEN** 文档说明热点路径上可考虑手动调用替代 defer，但普通场景 defer 开销可接受且代码更清晰

#### Scenario: 展示 defer 与锁的配合
- **WHEN** 读者需要安全释放锁
- **THEN** 文档提供 defer mu.Unlock() 的最佳实践示例

### Requirement: 提供面试高频题

文档 SHALL 提供 defer、panic、recover 相关的面试高频题，每题包含标准答案和追问点。

#### Scenario: defer 执行顺序面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 defer 执行顺序的标准答案和代码示例，包含参数求值时机的陷阱

#### Scenario: defer 修改返回值面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 defer 如何修改命名返回值的标准答案和代码示例

#### Scenario: recover 使用规则面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 recover 必须在 defer 中直接调用的标准答案和常见错误示例

