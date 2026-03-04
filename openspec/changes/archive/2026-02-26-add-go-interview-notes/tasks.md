## 1. 创建 channel 和并发文档 (go-channel-concurrency-notes)

- [x] 1.1 创建 note/note3-channel.md 文件
- [x] 1.2 编写 channel 基础机制章节（无缓冲/有缓冲/关闭规则）
- [x] 1.3 编写 select 多路复用章节（随机选择/default/超时模式）
- [x] 1.4 编写常见并发模式章节（worker pool/fan-in/fan-out/pipeline）
- [x] 1.5 编写死锁场景分析章节（自发自收/循环依赖/忘记关闭）
- [x] 1.6 编写面试高频题章节（无缓冲 vs 有缓冲/select 原理/关闭规则）
- [x] 1.7 添加代码示例和追问点
- [x] 1.8 验证文档完整性和格式统一

## 2. 创建 interface 和反射文档 (go-interface-reflection-notes)

- [x] 2.1 创建 note/note4-interface.md 文件
- [x] 2.2 编写 interface 底层结构章节（eface/iface/动态类型和值）
- [x] 2.3 编写类型断言和判断章节（单返回值/双返回值/type switch）
- [x] 2.4 编写反射基本概念章节（Type/Value/反射三定律）
- [x] 2.5 编写反射常见用法章节（struct tag/动态调用方法/类型转换）
- [x] 2.6 编写反射性能分析章节（性能开销/适用场景/优化建议）
- [x] 2.7 编写面试高频题章节（interface nil 判断/反射三定律/反射性能）
- [x] 2.8 添加代码示例和追问点
- [x] 2.9 验证文档完整性和格式统一

## 3. 创建 defer/panic/recover 文档 (go-defer-panic-recover-notes)

- [x] 3.1 创建 note/note5-defer.md 文件
- [x] 3.2 编写 defer 执行顺序章节（LIFO/参数求值时机/与 return 的关系/与 panic 的关系）
- [x] 3.3 编写 panic 和 recover 机制章节（panic 传播/recover 使用规则/跨 goroutine 限制）
- [x] 3.4 编写错误处理最佳实践章节（error vs panic/错误包装/panic 合理场景）
- [x] 3.5 编写 defer 性能分析章节（性能开销/使用建议/与锁的配合）
- [x] 3.6 编写面试高频题章节（defer 执行顺序/修改返回值/recover 使用规则）
- [x] 3.7 添加代码示例和追问点
- [x] 3.8 验证文档完整性和格式统一

## 4. 创建 context 使用文档 (go-context-notes)

- [x] 4.1 创建 note/note6-context.md 文件
- [x] 4.2 编写 context 设计目的章节（作用/设计原则/典型场景）
- [x] 4.3 编写四种创建方式章节（Background/TODO/WithCancel/WithTimeout/WithDeadline/WithValue）
- [x] 4.4 编写取消传播机制章节（单向传播/Done channel/Err 方法）
- [x] 4.5 编写使用规范章节（首参数/不存储/不传 nil/命名约定）
- [x] 4.6 编写 context.Value 分析章节（适用数据/性能考虑/key 类型选择）
- [x] 4.7 编写面试高频题章节（设计目的/取消传播/Value 使用规范）
- [x] 4.8 添加代码示例和追问点
- [x] 4.9 验证文档完整性和格式统一

## 5. 创建内存管理和 GC 文档 (go-memory-gc-notes)

- [x] 5.1 创建 note/note7-memory.md 文件
- [x] 5.2 编写内存分配策略章节（TCMalloc 思想/size class/三级缓存/小对象大对象）
- [x] 5.3 编写逃逸分析章节（逃逸定义/常见场景/查看方法/性能影响）
- [x] 5.4 编写 GC 原理章节（三色标记法/写屏障/STW 时机/触发条件）
- [x] 5.5 编写 GC 性能调优章节（减少堆分配/对象池复用/GOGC 调整/监控指标）
- [x] 5.6 编写内存泄漏排查章节（常见场景/pprof heap/goroutine 泄漏）
- [x] 5.7 编写面试高频题章节（逃逸分析/三色标记法/GC 调优）
- [x] 5.8 添加代码示例和追问点
- [x] 5.9 验证文档完整性和格式统一

## 6. 创建测试和 benchmark 文档 (go-testing-benchmark-notes)

- [x] 6.1 创建 note/note8-testing.md 文件
- [x] 6.2 编写单元测试基础章节（文件命名/函数签名/testing.T 方法/运行命令）
- [x] 6.3 编写表驱动测试章节（测试结构/子测试 t.Run/优势）
- [x] 6.4 编写 mock 和依赖注入章节（interface 抽象/手动 mock/gomock 使用）
- [x] 6.5 编写 benchmark 性能测试章节（函数签名/b.N 使用/ResetTimer/运行分析）
- [x] 6.6 编写测试覆盖率章节（生成报告/查看详情/覆盖率目标）
- [x] 6.7 编写 pprof 性能分析章节（CPU profile/内存 profile/与 benchmark 结合）
- [x] 6.8 编写面试高频题章节（单元测试最佳实践/benchmark 使用/pprof 分析）
- [x] 6.9 添加代码示例和追问点
- [x] 6.10 验证文档完整性和格式统一

## 7. 最终审查

- [x] 7.1 检查所有 6 个文档是否与现有 note1.md、note2.md 风格一致
- [x] 7.2 确保每个文档包含"面试题 → 标准答案 → 追问"结构
- [x] 7.3 验证中文表达自然，关键术语标注英文
- [x] 7.4 检查代码示例的正确性和完整性
- [x] 7.5 确认每个文档长度在 300-400 行范围内
- [x] 7.6 验证所有 6 个 spec 的 requirements 都被覆盖
