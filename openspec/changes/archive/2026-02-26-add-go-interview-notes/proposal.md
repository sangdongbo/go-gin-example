## Why

当前 note 目录下已有基础容器（slice/map/string）、GMP 调度模型和性能优化案例的总结（note1.md 和 note2.md），但缺少其他 Go 面试高频考点的系统性梳理。为了帮助面试者更全面地掌握 Go 语言核心知识点，需要补充更多详细的、面向面试场景的知识总结文档。

## What Changes

- 在 note 目录下新增多个 Go 知识总结文档
- 每个文档聚焦一个主题领域，包含面试题、标准答案、追问和实战案例
- 内容风格与现有 note1.md、note2.md 保持一致
- 覆盖 Go 面试中的高频知识点，形成完整的面试知识体系

## Capabilities

### New Capabilities
- `go-channel-concurrency-notes`: channel 机制、并发模式、select 使用、死锁问题等
- `go-interface-reflection-notes`: interface 底层实现、类型断言、反射机制、空接口使用
- `go-defer-panic-recover-notes`: defer 执行顺序、panic/recover 机制、错误处理最佳实践
- `go-context-notes`: context 使用场景、超时控制、值传递、取消信号传播
- `go-memory-gc-notes`: 内存分配策略、逃逸分析、GC 原理（三色标记）、性能调优
- `go-testing-benchmark-notes`: 单元测试、表驱动测试、benchmark 写法、pprof 分析

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增 6 个 markdown 文件到 note 目录
- 与现有 note1.md、note2.md 形成完整的 Go 面试知识体系
- 不影响项目代码和配置
- 文档总量预计增加 2000+ 行（每个文档 300-400 行左右）
