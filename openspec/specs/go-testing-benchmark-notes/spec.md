# go-testing-benchmark-notes Specification

## Purpose
TBD - created by archiving change add-go-interview-notes. Update Purpose after archive.
## Requirements
### Requirement: 讲解单元测试基础

文档 SHALL 讲解 Go 单元测试的基础知识，包括测试文件命名、测试函数签名、testing.T 的常用方法、运行测试的命令。

#### Scenario: 说明测试文件命名规范
- **WHEN** 读者创建测试文件
- **THEN** 文档说明测试文件以 _test.go 结尾，与被测代码在同一包

#### Scenario: 说明测试函数签名
- **WHEN** 读者编写测试函数
- **THEN** 文档说明测试函数以 Test 开头，参数为 *testing.T，提供代码示例

#### Scenario: 说明 testing.T 常用方法
- **WHEN** 读者使用 testing.T
- **THEN** 文档说明 Error/Errorf、Fatal/Fatalf、Log/Logf、Skip/Skipf 等方法的区别

#### Scenario: 说明运行测试命令
- **WHEN** 读者需要运行测试
- **THEN** 文档提供 go test、go test -v、go test -run 等命令示例

### Requirement: 讲解表驱动测试

文档 SHALL 讲解表驱动测试的写法和优势，包括测试用例结构、子测试 t.Run 的使用。

#### Scenario: 展示表驱动测试结构
- **WHEN** 读者编写多组测试用例
- **THEN** 文档提供表驱动测试的代码模板，包含 name、input、want 等字段

#### Scenario: 展示子测试用法
- **WHEN** 读者需要运行单个测试用例
- **THEN** 文档说明 t.Run 的使用，提供 go test -run TestFunc/case_name 的示例

#### Scenario: 说明表驱动测试优势
- **WHEN** 读者评估测试写法
- **THEN** 文档说明表驱动测试易于添加用例、代码简洁、便于维护

### Requirement: 讲解 mock 和依赖注入

文档 SHALL 讲解测试中的 mock 技术，包括 interface 抽象、gomock 使用、依赖注入测试。

#### Scenario: 说明 interface 抽象利于测试
- **WHEN** 读者设计可测试代码
- **THEN** 文档说明通过 interface 抽象依赖，便于 mock 替换

#### Scenario: 展示手动 mock
- **WHEN** 读者需要 mock 简单依赖
- **THEN** 文档提供手动实现 interface 的 mock 结构体示例

#### Scenario: 展示 gomock 使用
- **WHEN** 读者需要复杂 mock
- **THEN** 文档提供 gomock 生成 mock 代码和使用示例

### Requirement: 讲解 benchmark 性能测试

文档 SHALL 讲解 Go 的 benchmark 性能测试，包括 benchmark 函数签名、testing.B 的使用、运行和分析 benchmark 结果。

#### Scenario: 说明 benchmark 函数签名
- **WHEN** 读者编写 benchmark
- **THEN** 文档说明 benchmark 函数以 Benchmark 开头，参数为 *testing.B

#### Scenario: 说明 b.N 的使用
- **WHEN** 读者编写 benchmark 循环
- **THEN** 文档说明循环 b.N 次，Go 会自动调整 N 使测试运行足够长

#### Scenario: 说明 ResetTimer 和 StopTimer
- **WHEN** 读者需要排除初始化时间
- **THEN** 文档提供 b.ResetTimer() 和 b.StopTimer() 的使用示例

#### Scenario: 说明运行和分析 benchmark
- **WHEN** 读者需要运行 benchmark
- **THEN** 文档提供 go test -bench、-benchmem、-cpuprofile 等命令示例

### Requirement: 讲解测试覆盖率

文档 SHALL 讲解 Go 的测试覆盖率统计，包括生成覆盖率报告、查看覆盖率详情、提高覆盖率的方法。

#### Scenario: 说明生成覆盖率报告
- **WHEN** 读者需要统计覆盖率
- **THEN** 文档提供 go test -cover、-coverprofile 命令示例

#### Scenario: 说明查看覆盖率详情
- **WHEN** 读者需要查看未覆盖代码
- **THEN** 文档提供 go tool cover -html 生成可视化报告的方法

#### Scenario: 说明覆盖率目标
- **WHEN** 读者评估测试质量
- **THEN** 文档说明 80% 是常见目标，但覆盖率高不等于测试质量高

### Requirement: 讲解 pprof 性能分析

文档 SHALL 讲解使用 pprof 进行性能分析，包括 CPU profile、内存 profile、goroutine profile 的生成和分析。

#### Scenario: 说明生成 CPU profile
- **WHEN** 读者需要分析 CPU 热点
- **THEN** 文档提供 go test -cpuprofile、pprof 命令和 top、list 等分析方法

#### Scenario: 说明生成内存 profile
- **WHEN** 读者需要分析内存分配
- **THEN** 文档提供 go test -memprofile 和 inuse_space、alloc_objects 等指标的分析方法

#### Scenario: 说明 benchmark 与 pprof 结合
- **WHEN** 读者需要对 benchmark 做深度分析
- **THEN** 文档提供 benchmark 生成 profile 并用 pprof 分析的完整流程

### Requirement: 提供面试高频题

文档 SHALL 提供测试和 benchmark 相关的面试高频题，每题包含标准答案和追问点。

#### Scenario: 单元测试最佳实践面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供测试文件组织、表驱动测试、mock 使用等最佳实践的标准答案

#### Scenario: benchmark 使用面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 benchmark 编写、b.N 使用、结果分析的标准答案

#### Scenario: pprof 分析面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 pprof 类型、分析方法、常见优化手段的标准答案

