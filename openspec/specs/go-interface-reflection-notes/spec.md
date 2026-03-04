# go-interface-reflection-notes Specification

## Purpose
TBD - created by archiving change add-go-interview-notes. Update Purpose after archive.
## Requirements
### Requirement: 解释 interface 底层结构

文档 SHALL 详细解释 Go interface 的底层实现，包括 eface 和 iface 的区别、动态类型和动态值的概念。

#### Scenario: 说明 eface 结构
- **WHEN** 读者查阅空接口（interface{}）底层实现
- **THEN** 文档说明 eface 包含 type 和 data 两个字段

#### Scenario: 说明 iface 结构
- **WHEN** 读者查阅非空接口底层实现
- **THEN** 文档说明 iface 包含 itab（类型信息和方法表）和 data 两个字段

#### Scenario: 解释动态类型和动态值
- **WHEN** 读者查阅 interface nil 判断问题
- **THEN** 文档说明只有动态类型和动态值都为 nil 时 interface 才为 nil

### Requirement: 讲解类型断言和类型判断

文档 SHALL 讲解类型断言的用法，包括单返回值和双返回值形式、type switch 的使用。

#### Scenario: 说明类型断言的两种形式
- **WHEN** 读者查阅类型断言
- **THEN** 文档说明单返回值形式断言失败会 panic，双返回值形式返回 false

#### Scenario: 展示 type switch 用法
- **WHEN** 读者需要判断 interface 的具体类型
- **THEN** 文档提供 type switch 的代码示例和使用场景

#### Scenario: 说明类型断言的性能
- **WHEN** 读者关注性能问题
- **THEN** 文档说明类型断言比反射快，但频繁使用仍有开销

### Requirement: 说明反射的基本概念

文档 SHALL 说明 Go 反射的基本概念，包括 reflect.Type 和 reflect.Value、反射三定律。

#### Scenario: 解释反射三定律
- **WHEN** 读者学习反射基础
- **THEN** 文档说明反射三定律：interface 到 reflection object、reflection object 到 interface、要修改 reflection object 需要可寻址

#### Scenario: 展示 Type 和 Value 的获取
- **WHEN** 读者需要使用反射
- **THEN** 文档提供 reflect.TypeOf 和 reflect.ValueOf 的使用示例

#### Scenario: 说明可寻址性
- **WHEN** 读者需要通过反射修改值
- **THEN** 文档说明只有通过指针获取的 Value 才可以修改，并提供示例

### Requirement: 讲解反射的常见用法

文档 SHALL 讲解反射的常见用法，包括结构体 tag 读取、动态调用方法、类型转换。

#### Scenario: 展示 struct tag 读取
- **WHEN** 读者需要读取结构体 tag
- **THEN** 文档提供 StructField.Tag.Get 的使用示例

#### Scenario: 展示动态调用方法
- **WHEN** 读者需要动态调用对象方法
- **THEN** 文档提供 Value.MethodByName 和 Call 的使用示例

#### Scenario: 展示类型转换
- **WHEN** 读者需要将 reflect.Value 转换回具体类型
- **THEN** 文档提供 Interface() 方法和类型断言的组合使用

### Requirement: 分析反射的性能影响

文档 SHALL 分析反射对性能的影响，说明反射的开销来源和使用场景。

#### Scenario: 说明反射的性能开销
- **WHEN** 读者关注反射性能
- **THEN** 文档说明反射比直接调用慢 10-100 倍，开销来自动态类型检查和方法查找

#### Scenario: 说明反射的适用场景
- **WHEN** 读者评估是否使用反射
- **THEN** 文档说明反射适用于框架、序列化库等需要通用性的场景，不适合热点代码路径

#### Scenario: 提供优化建议
- **WHEN** 读者需要优化反射性能
- **THEN** 文档提供缓存 reflect.Type、避免重复反射、使用代码生成替代反射等优化建议

### Requirement: 提供面试高频题

文档 SHALL 提供 interface 和反射相关的面试高频题，每题包含标准答案和追问点。

#### Scenario: interface nil 判断面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供 interface 为 nil 的判断标准和常见陷阱

#### Scenario: 反射三定律面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供反射三定律的详细解释和代码示例

#### Scenario: 反射性能面试题
- **WHEN** 读者复习面试题
- **THEN** 文档提供反射性能分析和优化方法

