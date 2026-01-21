# Go-Gin-Example 进阶功能训练指南

本文档介绍了 Stock 模块中实现的四个进阶功能示例，帮助你深入学习 Go 和 GORM 的高级特性。

## 1. 事务 (Transaction)

### 1.1 批量创建产品及明细

**接口:** `POST /api/v1/stock/transaction/batch-create`

**功能:** 演示如何使用事务保证批量操作的原子性（要么全部成功，要么全部失败）

**请求示例:**
```json
{
  "products": [
    {
      "name": "产品A",
      "description": "描述A",
      "details": [
        {
          "quantity": 100,
          "price": 10.5,
          "location": "仓库A",
          "supplier": "供应商A",
          "batch_number": "BATCH001"
        }
      ]
    }
  ]
}
```

**学习要点:**
- `tx := models.Db.Begin()` - 开始事务
- `tx.Commit()` - 提交事务
- `tx.Rollback()` - 回滚事务
- `defer` 配合 `recover()` 处理异常回滚

**代码位置:** `service/stock_service/stock_advanced.go` - `BatchCreateProductWithDetails`

---

### 1.2 库存转移

**接口:** `POST /api/v1/stock/transaction/transfer`

**功能:** 演示事务的 ACID 特性，确保库存转移过程的一致性

**请求示例:**
```json
{
  "from_detail_id": 1,
  "to_detail_id": 2,
  "quantity": 50
}
```

**学习要点:**
- 事务隔离级别的应用
- 并发情况下的数据一致性保证
- 业务逻辑验证在事务中的处理

**代码位置:** `service/stock_service/stock_advanced.go` - `TransferStock`

---

## 2. Hooks (钩子函数)

### 2.1 创建产品（演示 Hooks）

**接口:** `POST /api/v1/stock/hooks/product`

**功能:** 演示 GORM 的生命周期钩子，在记录创建前后自动执行特定逻辑

**请求示例:**
```json
{
  "name": "测试产品",
  "description": "测试描述"
}
```

**学习要点:**
GORM 提供以下 Hooks：

1. **BeforeSave** - 创建或更新之前
2. **BeforeCreate** - 创建之前
   - 自动生成字段值（如 SKU）
   - 数据验证
   - 记录操作日志
3. **AfterCreate** - 创建之后
   - 发送通知
   - 记录审计日志
   - 触发其他业务逻辑
4. **BeforeUpdate** - 更新之前
   - 记录变更历史
   - 数据验证
5. **AfterUpdate** - 更新之后
   - 清除缓存
   - 同步到其他系统
6. **BeforeDelete** - 删除之前
   - 权限检查
   - 级联删除处理
7. **AfterDelete** - 删除之后
   - 清理关联数据
8. **AfterFind** - 查询之后
   - 数据解密或格式化

**代码位置:** `models/stock/stock_product.go` - Hooks 部分

---

## 3. 关联查询优化

### 3.1 预加载关联数据 (Preload)

**接口:** `GET /api/v1/stock/optimize/products-with-details`

**功能:** 使用 Preload 避免 N+1 查询问题

**查询参数:**
- `page`: 页码
- `page_size`: 每页数量

**N+1 问题说明:**
```go
// ❌ 不使用 Preload - 会产生 N+1 查询
products := GetProducts() // 1次查询
for _, product := range products {
    details := GetDetails(product.ID) // N次查询
}

// ✅ 使用 Preload - 只需要 2次查询
db.Preload("StockProductDetails").Find(&products)
// 查询1：SELECT * FROM stock_products
// 查询2：SELECT * FROM stock_product_details WHERE stock_product_id IN (...)
```

**学习要点:**
- `Preload()` 的使用方法
- 嵌套关联的预加载：`Preload("Details.Supplier")`
- 条件预加载：`Preload("Details", "quantity > ?", 10)`

**代码位置:** `service/stock_service/stock_advanced.go` - `GetProductsWithDetailsOptimized`

---

### 3.2 使用 Joins 查询

**接口:** `GET /api/v1/stock/optimize/products-join`

**功能:** 使用 Joins 优化需要关联条件过滤的查询

**查询参数:**
- `min_quantity`: 最小库存数量（只返回库存大于此值的产品）

**Preload vs Joins:**

| 特性 | Preload | Joins |
|------|---------|-------|
| SQL | 多次查询 | 单次 JOIN 查询 |
| 适用场景 | 需要完整关联数据 | 需要根据关联条件筛选 |
| 性能 | 数据量大时更好 | 筛选条件多时更好 |
| 返回数据 | 主表记录和关联记录 | 可能产生重复行 |

**代码位置:** `service/stock_service/stock_advanced.go` - `GetProductsWithJoin`

---

## 4. 事务嵌套 (Nested Transaction)

### 4.1 创建订单（嵌套事务）

**接口:** `POST /api/v1/stock/nested-transaction/order`

**功能:** 演示如何在事务中处理多个相关操作，使用 SAVEPOINT 实现嵌套事务

**请求示例:**
```json
{
  "order_sn": "ORDER20260121001",
  "items": [
    {
      "product_id": 1,
      "quantity": 10
    },
    {
      "product_id": 2,
      "quantity": 5
    }
  ]
}
```

**嵌套事务流程:**
```
外层事务开始 (BEGIN)
  ├─ 创建订单
  ├─ 遍历订单项
  │   ├─ 设置保存点 (SAVEPOINT sp1)
  │   ├─ 检查库存
  │   ├─ 扣减库存
  │   └─ 如果失败，回滚到 sp1 (ROLLBACK TO sp1)
  └─ 提交事务 (COMMIT)
```

**学习要点:**
- MySQL SAVEPOINT 的使用
- `tx.Exec("SAVEPOINT sp_name")` - 设置保存点
- `tx.Exec("ROLLBACK TO SAVEPOINT sp_name")` - 回滚到保存点
- `tx.Exec("RELEASE SAVEPOINT sp_name")` - 释放保存点
- 嵌套事务的应用场景

**代码位置:** `service/stock_service/stock_advanced.go` - `CreateOrderWithNestedTransaction`

---

## 测试建议

### 1. 事务测试
```bash
# 测试批量创建
curl -X POST http://localhost:8000/api/v1/stock/transaction/batch-create \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"products":[...]}'

# 测试库存转移
curl -X POST http://localhost:8000/api/v1/stock/transaction/transfer \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_detail_id":1,"to_detail_id":2,"quantity":50}'
```

### 2. Hooks 测试
- 创建产品后检查日志表
- 更新产品后检查变更历史表
- 删除产品前检查关联数据处理

### 3. 查询优化测试
```bash
# 测试 Preload（观察 SQL 日志）
curl -X GET "http://localhost:8000/api/v1/stock/optimize/products-with-details?page=1" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 测试 Joins
curl -X GET "http://localhost:8000/api/v1/stock/optimize/products-join?min_quantity=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. 嵌套事务测试
```bash
# 测试订单创建
curl -X POST http://localhost:8000/api/v1/stock/nested-transaction/order \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_sn":"ORDER001","items":[{"product_id":1,"quantity":10}]}'
```

---

## 性能优化建议

### 1. 事务优化
- 保持事务尽可能短小
- 避免在事务中执行耗时操作（如 HTTP 请求）
- 合理设置事务隔离级别

### 2. Hooks 优化
- Hooks 中避免复杂的业务逻辑
- 耗时操作应该异步处理（使用消息队列）
- 谨慎使用 AfterFind，可能导致性能问题

### 3. 查询优化
- 根据实际场景选择 Preload 或 Joins
- 使用 Select 只查询需要的字段
- 添加适当的数据库索引
- 使用分页避免一次加载过多数据

### 4. 嵌套事务优化
- 只在必要时使用 SAVEPOINT
- 避免过深的嵌套层级
- 考虑使用补偿事务代替嵌套事务

---

## 扩展练习

1. **事务练习**
   - 实现分布式事务（使用 2PC 或 TCC）
   - 实现事务超时控制
   - 处理死锁问题

2. **Hooks 练习**
   - 实现软删除的自动处理
   - 实现字段加密/解密
   - 实现审计日志系统

3. **查询优化练习**
   - 分析慢查询并优化
   - 实现查询结果缓存
   - 使用 Eager Loading 优化多层关联

4. **嵌套事务练习**
   - 实现复杂的业务流程（如秒杀系统）
   - 处理并发冲突
   - 实现事务重试机制

---

## 5. 协程并发处理

### 5.1 并发查询多个产品

**接口:** `POST /api/v1/stock/goroutine/batch-query`

**功能:** 演示基本协程使用，并发查询多个产品信息

**请求示例:**
```json
{
  "product_ids": [1, 2, 3, 4, 5]
}
```

**学习要点:**
- `go` 关键字启动协程
- `sync.WaitGroup` 等待所有协程完成
- `sync.Mutex` 保护共享资源
- 协程间的数据同步

**代码位置:** `service/stock_service/stock_goroutine.go` - `BatchQueryProductsWithGoroutine`

---

### 5.2 并发更新库存

**接口:** `POST /api/v1/stock/goroutine/batch-update`

**功能:** 演示带超时控制的协程，并发更新多个产品库存

**请求示例:**
```json
{
  "updates": [
    {"product_id": 1, "quantity": 100.5},
    {"product_id": 2, "quantity": 200.0}
  ]
}
```

**学习要点:**
- `context.WithTimeout` 超时控制
- Channel 进行协程间通信
- `select` 监听 context 取消
- 错误处理和结果收集

**代码位置:** `service/stock_service/stock_goroutine.go` - `BatchUpdateStockWithGoroutine`

---

### 5.3 Worker Pool 模式

**接口:** `POST /api/v1/stock/goroutine/worker-pool`

**功能:** 演示固定数量的工作协程池模式，避免创建过多协程

**请求示例:**
```json
{
  "product_ids": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
  "worker_count": 3
}
```

**Worker Pool 工作原理:**
```
任务队列 (jobChan)
    ↓
Worker 1 ← 
Worker 2 ← 处理任务
Worker 3 ← 
    ↓
结果队列 (resultChan)
```

**学习要点:**
- 控制并发数量，避免资源耗尽
- 任务分发和结果收集
- Channel 的缓冲区使用
- 优雅关闭 Channel

**代码位置:** `service/stock_service/stock_goroutine.go` - `ProcessWithWorkerPool`

---

### 5.4 Pipeline 模式

**接口:** `GET /api/v1/stock/goroutine/pipeline`

**功能:** 演示数据流处理的 Pipeline 模式

**查询参数:**
- `limit`: 处理的产品数量限制

**Pipeline 流程:**
```
Stage 1: 生成产品ID
    ↓ (channel)
Stage 2: 获取产品信息
    ↓ (channel)
Stage 3: 计算统计信息
    ↓ (channel)
最终结果
```

**学习要点:**
- 数据流式处理
- 多阶段处理管道
- Channel 作为管道连接
- 每个阶段独立的协程

**代码位置:** `service/stock_service/stock_goroutine.go` - `ProcessWithPipeline`

---

### 5.5 Fan-out/Fan-in 模式

**接口:** `POST /api/v1/stock/goroutine/fan-out-in`

**功能:** 演示扇出-扇入模式，多个协程并行处理，单个协程汇总结果

**请求示例:**
```json
{
  "product_ids": [1, 2, 3, 4, 5]
}
```

**Fan-out/Fan-in 流程:**
```
输入数据
    ↓
Fan-out (分发)
    ├─→ Worker 1
    ├─→ Worker 2
    └─→ Worker 3
         ↓
Fan-in (汇总)
    ↓
合并结果
```

**学习要点:**
- 任务并行处理
- 多通道合并
- 动态协程数量
- 结果聚合统计

**代码位置:** `service/stock_service/stock_goroutine.go` - `ProcessWithFanOutFanIn`

---

## 协程最佳实践

### 1. 避免协程泄漏
```go
// ❌ 错误：channel 未关闭，协程永远阻塞
func bad() {
    ch := make(chan int)
    go func() {
        for v := range ch { // 永远等待
            process(v)
        }
    }()
}

// ✅ 正确：关闭 channel
func good() {
    ch := make(chan int)
    go func() {
        for v := range ch {
            process(v)
        }
    }()
    // ... 处理完后
    close(ch)
}
```

### 2. 使用 Context 控制生命周期
```go
func ProcessWithContext(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return // 及时退出
        default:
            // 正常处理
        }
    }
}
```

### 3. 合理设置 Worker 数量
```go
// 根据 CPU 核心数设置
workerCount := runtime.NumCPU()

// 或根据任务类型设置
// IO 密集型可以更多
// CPU 密集型建议等于核心数
```

### 4. 使用 sync.Once 保证只执行一次
```go
var once sync.Once
var instance *Resource

func GetInstance() *Resource {
    once.Do(func() {
        instance = &Resource{}
    })
    return instance
}
```

### 5. 使用 Buffer Channel 减少阻塞
```go
// 无缓冲：发送方会阻塞到接收方准备好
ch := make(chan int)

// 有缓冲：可以发送多个而不阻塞
ch := make(chan int, 100)
```

---

## 性能测试建议

### 协程测试
```bash
# 并发查询测试
curl -X POST http://localhost:8000/api/v1/stock/goroutine/batch-query \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_ids":[1,2,3,4,5]}'

# Worker Pool 测试（比较不同 worker 数量的性能）
curl -X POST http://localhost:8000/api/v1/stock/goroutine/worker-pool \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_ids":[1,2,3,4,5,6,7,8,9,10],"worker_count":3}'

# Pipeline 测试
curl -X GET "http://localhost:8000/api/v1/stock/goroutine/pipeline?limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 协程进阶练习

1. **并发控制**
   - 实现令牌桶限流
   - 实现信号量控制并发数
   - 实现超时重试机制

2. **数据竞争**
   - 使用 `go run -race` 检测数据竞争
   - 使用原子操作 `sync/atomic`
   - 正确使用读写锁 `sync.RWMutex`

3. **高级模式**
   - 实现生产者-消费者模式
   - 实现发布-订阅模式
   - 实现协程池复用

4. **错误处理**
   - 协程中的 panic 捕获
   - 错误聚合和上报
   - 部分失败的处理策略

---

## 参考资料

- [GORM 官方文档](https://gorm.io/zh_CN/docs/)
- [MySQL 事务文档](https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html)
- [Go 并发编程](https://go.dev/doc/effective_go#concurrency)
- [Go Channel 详解](https://go.dev/ref/spec#Channel_types)
- [数据库事务隔离级别](https://en.wikipedia.org/wiki/Isolation_(database_systems))
