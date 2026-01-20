# (go get github.com/chenfeifan111/generics_crud@main
# 泛型数据库查询框架

基于 Gin + Viper + GORM 的泛型数据库查询框架，支持动态过滤、排序和分页。

## 功能特性

- ✅ **泛型查询**: 使用 Go 泛型实现类型安全的查询
- ✅ **动态过滤**: 通过 struct tag 定义过滤规则
- ✅ **灵活排序**: 支持多字段排序，可配置升序/降序
- ✅ **分页支持**: 可选的分页功能
- ✅ **Tag 驱动**: 基于 struct tag 的声明式查询

## 快速开始

### 1. 定义查询请求结构体
http://localhost:8080/users/query
```go
type UserQueryRequest struct {
    Page    *Page `json:"page"`
    Filters struct {
        Age  *int    `json:"age" filter:"gte"`      // 年龄 >= 值
        Name *string `json:"name" filter:"like"`    // 名称模糊查询
    } `json:"filters"`
    Orders struct {
        Age       *bool `json:"age" order:"desc"`         // 按年龄降序
        CreatedAt *bool `json:"created_at" order:"desc"`  // 按创建时间降序
    } `json:"orders"`
}
```

### 2. 实现 QueryRequest 接口

```go
func (r *UserQueryRequest) GetPage() interface{} {
    return r.Page
}

func (r *UserQueryRequest) GetFilters() interface{} {
    return r.Filters
}

func (r *UserQueryRequest) GetOrders() interface{} {
    return r.Orders
}
```

### 3. 在 Controller 中使用

```go
func QueryUsers(c *gin.Context) {
    var req request.UserQueryRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error("Invalid request: "+err.Error()))
        return
    }
    
    users, total, err := logic.ExecuteQuery(config.DB, entity.User{}, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
        return
    }
    
    if req.Page != nil && req.Page.IsValid() {
        c.JSON(http.StatusOK, response.SuccessWithPage(users, total, req.Page.PageNum, req.Page.PageSize))
    } else {
        c.JSON(http.StatusOK, response.Success(users))
    }
}
```

## 支持的过滤操作符

| Tag 值 | 说明 | SQL 示例 |
|--------|------|----------|
| `eq` 或 空 | 等于 | `column = ?` |
| `ne`, `neq` | 不等于 | `column != ?` |
| `gt` | 大于 | `column > ?` |
| `gte` | 大于等于 | `column >= ?` |
| `lt` | 小于 | `column < ?` |
| `lte` | 小于等于 | `column <= ?` |
| `like` | 模糊查询 | `column LIKE %?%` |
| `in` | 在列表中 | `column IN (?)` |
| `not_in` | 不在列表中 | `column NOT IN (?)` |
| `is_null` | 为空 | `column IS NULL` |
| `is_not_null` | 不为空 | `column IS NOT NULL` |

## 排序规则

- `order:"asc"` - 升序排序
- `order:"desc"` - 降序排序（默认）
- 字段值为 `true` 时启用该排序

## API 使用示例

### 1. 查询年龄大于等于 18 的用户（带分页）

```bash
curl -X POST http://localhost:8080/users/query \
  -H "Content-Type: application/json" \
  -d '{
    "page": {
      "page_num": 1,
      "page_size": 10
    },
    "filters": {
       "age": 18,
        "name":"zs"
    },
    "orders": { 
        "age":"asc", 
        "id": "desc"
    }
  }'
```

### 2. OR 条件查询

```bash
curl -X POST http://localhost:8080/users/query-or \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {
      "or": [
        {"age": 30},
        {"name": "admin"}
      ]
    },
    "orders": {
      "age": "desc"
    }
  }'
```

详细示例请查看 [API_EXAMPLES.md](API_EXAMPLES.md)

### 模糊查询名称包含 "张" 的用户

```bash
curl -X POST http://localhost:8080/users/query \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {
      "name": "张"
    }
  }'
```

### 查询所有用户（不分页）

```bash
curl -X POST http://localhost:8080/users/query \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {},
    "orders": {}
  }'
```

## 响应格式

### 带分页的响应

```json
{
  "code": 200,
  "msg": "success",
  "data": [...],
  "total": 100,
  "page": {
    "page_num": 1,
    "page_size": 10
  }
}
```

### 不带分页的响应

```json
{
  "code": 200,
  "msg": "success",
  "data": [...]
}
```

## 项目结构

```
.
├── config/           # 配置相关
│   ├── database.go   # 数据库连接
│   └── viper.go      # 配置加载
├── controller/       # 控制器
│   └── user_controller.go
├── entity/           # 实体模型
│   └── user.go
├── logic/            # 业务逻辑
│   ├── query_builder.go    # 查询构建器
│   └── generic_query.go    # 泛型查询执行器
├── request/          # 请求结构
│   ├── base.go       # 基础分页结构
│   └── user_request.go
├── response/         # 响应结构
│   └── base.go
├── config.yaml       # 配置文件
└── main.go           # 入口文件
```

## 扩展示例

### 创建新的查询请求

```go
type ProductQueryRequest struct {
    Page    *Page `json:"page"`
    Filters struct {
        Price    *float64 `json:"price" filter:"lte"`      // 价格 <= 值
        Category *string  `json:"category" filter:"eq"`    // 分类精确匹配
        Status   []*int     `json:"status" filter:"in"`      // 状态在列表中
    } `json:"filters"`
    Orders struct {
        Price     *bool `json:"price" order:"asc"`         // 价格升序
        CreatedAt *bool `json:"created_at" order:"desc"`   // 创建时间降序
    } `json:"orders"`
}

// 实现接口
func (r *ProductQueryRequest) GetPage() interface{}    { return r.Page }
func (r *ProductQueryRequest) GetFilters() interface{} { return r.Filters }
func (r *ProductQueryRequest) GetOrders() interface{}  { return r.Orders }
```

### 在 Controller 中使用

```go
func QueryProducts(c *gin.Context) {
    var req request.ProductQueryRequest
    c.ShouldBindJSON(&req)
    
    products, total, err := logic.ExecuteQuery(config.DB, entity.Product{}, &req)
    // ... 处理响应
}
```
