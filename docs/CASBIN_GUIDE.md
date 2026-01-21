# Casbin 权限管理使用指南

本文档介绍如何在 go-gin-example 项目中使用 Casbin 进行权限管理。

## 目录结构

```
├── conf/
│   ├── rbac_model.conf      # Casbin RBAC 模型配置
│   └── rbac_policy.csv      # 权限策略文件
├── middleware/
│   └── casbin/
│       └── casbin.go        # Casbin 中间件
├── pkg/
│   └── casbin/
│       └── casbin.go        # Casbin 核心功能封装
└── routers/
    └── api/
        └── v1/
            └── casbin.go    # 权限管理 API
```

## 权限校验位置

Casbin 权限校验主要在以下位置使用：

### 1. 中间件层（推荐）

**位置:** `middleware/casbin/casbin.go`

**用途:** 在路由中间件中进行权限校验，这是**最常用**的方式。

```go
// 在路由中使用
stock := apiv1.Group("/stock")
stock.Use(casbinMiddleware.Casbin()) // 添加权限中间件
{
    stock.GET("/products", v1.GetStockProducts)
    stock.POST("/product", v1.AddStockProduct)
    // ...
}
```

### 2. 核心功能层

**位置:** `pkg/casbin/casbin.go`

**用途:** 封装 Casbin 的核心功能，供中间件和业务逻辑调用。

### 3. API 管理层

**位置:** `routers/api/v1/casbin.go`

**用途:** 提供权限管理的 API 接口，用于动态管理权限。

---

## 快速开始

### 1. 存储说明

本项目使用文件适配器存储权限策略，权限数据保存在 `conf/rbac_policy.csv` 文件中。

**为什么使用文件适配器？**
- 项目使用的是旧版本 GORM (github.com/jinzhu/gorm)
- Casbin GORM adapter 需要新版本 GORM (gorm.io/gorm)
- 文件适配器简单、轻量，适合小型项目

### 2. 配置权限模型

模型文件已创建在 `conf/rbac_model.conf`：

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

**说明:**
- `sub`: 主体（用户/角色）
- `obj`: 对象（资源路径）
- `act`: 动作（HTTP 方法）

### 3. 添加角色和权限

#### 方式一：直接编辑 CSV 文件

编辑 `conf/rbac_policy.csv` 文件：

```csv
p, admin, /api/v1/*, *
p, editor, /api/v1/articles, POST
p, editor, /api/v1/articles/*, GET
g, user:1, admin
g, user:2, editor
```

**格式说明:**
- `p` 开头的行：权限策略（角色，资源路径，HTTP方法）
- `g` 开头的行：用户角色关系（用户ID，角色名）

#### 方式二：使用 API 接口

```bash
# 1. 为用户添加角色
curl -X POST http://localhost:8000/api/v1/casbin/add-role \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "role": "admin"}'

# 2. 为角色添加权限
curl -X POST http://localhost:8000/api/v1/casbin/add-policy \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "admin", "path": "/api/v1/stock/*", "method": "GET"}'
```

#### 方式二：直接在代码中添加

```go
import casbinPkg "github.com/EDDYCJY/go-gin-example/pkg/casbin"

// 添加角色
casbinPkg.AddRoleForUser("user:1", "admin")

// 添加权限
casbinPkg.AddPolicy("admin", "/api/v1/articles", "POST")
```

---

## 使用场景

### 场景1: 为整个路由组添加权限验证

```go
import casbinMiddleware "github.com/EDDYCJY/go-gin-example/middleware/casbin"

// 方式1: 对整个组应用权限中间件
stock := apiv1.Group("/stock")
stock.Use(casbinMiddleware.Casbin())
{
    stock.GET("/products", v1.GetStockProducts)
    stock.POST("/product", v1.AddStockProduct)
}
```

### 场景2: 为单个接口添加权限验证

```go
// 方式2: 对单个接口应用
stock.POST("/product", casbinMiddleware.Casbin(), v1.AddStockProduct)
```

### 场景3: 验证特定角色

```go
// 只允许 admin 或 stock_manager 角色访问
stock := apiv1.Group("/stock")
stock.Use(casbinMiddleware.CasbinWithRoles("admin", "stock_manager"))
{
    stock.DELETE("/product/:id", v1.DeleteStockProduct)
}
```

### 场景4: 在业务逻辑中检查权限

```go
import casbinPkg "github.com/EDDYCJY/go-gin-example/pkg/casbin"

func SomeBusinessLogic(userID int) error {
    subject := fmt.Sprintf("user:%d", userID)
    ok, err := casbinPkg.Enforce(subject, "/api/v1/articles", "POST")
    
    if err != nil {
        return err
    }
    
    if !ok {
        return errors.New("无权限")
    }
    
    // 执行业务逻辑
    return nil
}
```

---

## 权限管理 API

### 1. 为用户添加角色

**接口:** `POST /api/v1/casbin/add-role`

```json
{
  "user_id": 1,
  "role": "admin"
}
```

### 2. 删除用户角色

**接口:** `DELETE /api/v1/casbin/delete-role`

```json
{
  "user_id": 1,
  "role": "editor"
}
```

### 3. 获取用户角色

**接口:** `GET /api/v1/casbin/user-roles?user_id=1`

**响应:**
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "user_id": "1",
    "roles": ["admin", "editor"]
  }
}
```

### 4. 添加权限策略

**接口:** `POST /api/v1/casbin/add-policy`

```json
{
  "role": "editor",
  "path": "/api/v1/articles",
  "method": "POST"
}
```

### 5. 删除权限策略

**接口:** `DELETE /api/v1/casbin/delete-policy`

```json
{
  "role": "editor",
  "path": "/api/v1/articles",
  "method": "DELETE"
}
```

### 6. 检查权限

**接口:** `GET /api/v1/casbin/check-permission?user_id=1&path=/api/v1/articles&method=POST`

**响应:**
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "user_id": "1",
    "path": "/api/v1/articles",
    "method": "POST",
    "has_access": true
  }
}
```

---

## 预定义角色

项目中预定义了以下角色（在 `conf/rbac_policy.csv` 中）：

### 1. admin（管理员）
- 权限：所有 API 的所有操作

### 2. editor（编辑者）
- 权限：
  - 文章：查看、创建、编辑
  - 标签：查看

### 3. viewer（查看者）
- 权限：
  - 文章：查看
  - 标签：查看
  - 库存：查看

### 4. stock_manager（库存管理员）
- 权限：
  - 库存：所有操作

---

## 工作流程

```
1. 用户请求 API
   ↓
2. JWT 中间件验证 token（获取用户ID）
   ↓
3. Casbin 中间件检查权限
   ├─ 从 token 获取用户标识（user:1）
   ├─ 获取请求路径和方法
   └─ 调用 Casbin.Enforce() 验证权限
   ↓
4. 权限验证通过 → 继续处理请求
   权限验证失败 → 返回 403 Forbidden
```

---

## 高级用法

### 1. 通配符匹配

Casbin 支持路径通配符：

```go
// 允许访问所有 stock 相关的 GET 请求
casbinPkg.AddPolicy("viewer", "/api/v1/stock/*", "GET")

// 允许访问所有 API
casbinPkg.AddPolicy("admin", "/api/v1/*", "*")
```

### 2. 多租户支持

```go
// 为不同租户设置不同权限
casbinPkg.AddPolicy("tenant:1:admin", "/api/v1/tenant/1/*", "*")
casbinPkg.AddPolicy("tenant:2:admin", "/api/v1/tenant/2/*", "*")
```

### 3. 动态权限更新

所有权限修改都会实时生效，不需要重启服务：

```go
// 添加新权限后立即生效
casbinPkg.AddPolicy("new_role", "/api/v1/new_resource", "POST")
```

---

## 注意事项

### 1. JWT Token 中必须包含用户ID

确保你的 JWT token claims 中包含 `user_id` 字段：

```go
claims := util.Claims{
    UserID: 1,  // 必须包含
    Name:   "admin",
    StandardClaims: jwt.StandardClaims{
        // ...
    },
}
```

### 2. 路径匹配

Casbin 默认使用精确匹配，如需通配符匹配，需要在 model 中配置：

```conf
[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
```

### 3. 性能优化

- Casbin 会将权限规则缓存在内存中
- 使用数据库适配器时，只在修改时才更新数据库
- 对于高并发场景，Casbin 的性能足够好

---

## 故障排查

### 问题1: 权限验证总是失败

**检查清单:**
1. 确认 JWT token 有效
2. 确认用户已分配角色
3. 确认角色有对应的权限策略
4. 检查路径和方法是否完全匹配

```bash
# 查看用户角色
curl "http://localhost:8000/api/v1/casbin/user-roles?user_id=1" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 检查权限
curl "http://localhost:8000/api/v1/casbin/check-permission?user_id=1&path=/api/v1/articles&method=GET" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 问题2: Casbin 初始化失败

检查 `main.go` 中的初始化代码和数据库连接。

### 问题3: 无法从 token 获取用户ID

修改 `middleware/casbin/casbin.go` 中的 `getUserIDFromToken` 函数，确保正确解析你的 JWT token 结构。

---

## 参考资料

- [Casbin 官方文档](https://casbin.org/zh/)
- [RBAC 模型](https://casbin.org/docs/rbac)
- [GORM Adapter](https://casbin.org/docs/adapters)
