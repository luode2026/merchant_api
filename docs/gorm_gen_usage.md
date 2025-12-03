# GORM Gen 使用指南

本文档详细介绍了 `gorm.io/gen` 生成的代码与普通 GORM 模型的区别，以及如何使用生成的类型安全 API 进行开发。

## 1. 核心区别：GORM vs GORM Gen

### 普通 GORM 模型
- **特点**：仅定义结构体，查询依赖字符串硬编码。
- **缺点**：容易写错字段名（运行时报错），类型不安全，重构困难。

```go
// 普通 GORM 写法
var user User
// ❌ 风险：字段名 "username" 写错不会编译报错
err := db.Where("username = ?", "admin").First(&user).Error
```

### GORM Gen 生成代码
- **特点**：生成了结构体 + **类型安全的查询对象 (DAO)**。
- **优点**：编译期检查字段名，IDE 智能补全，类型安全，支持复杂查询构建。

```go
// GORM Gen 写法
u := q.User
// ✅ 安全：u.Username 是对象属性，写错直接编译报错
user, err := u.Where(u.Username.Eq("admin")).First()
```

---

## 2. 初始化

在使用生成的代码前，需要先初始化全局查询对象。

```go
import (
    "merchant_api/internal/model/generated"
    "merchant_api/pkg/database"
)

// 建议在服务启动时或 Service 层初始化
q := generated.Use(database.GetDB())
```

---

## 3. 基础 CRUD 操作

假设我们有一个 `User` 表，生成的查询对象为 `q.User`。

### 3.1 创建 (Create)

```go
u := q.User

// 单条插入
user := &model.User{
    Username: "zhangsan",
    Age:      18,
}
err := u.Create(user)

// 批量插入
users := []*model.User{
    {Username: "u1"},
    {Username: "u2"},
}
err := u.CreateInBatches(users, 100) // 每次插入 100 条
```

### 3.2 查询 (Read)

```go
u := q.User

// 1. 根据主键查询
user, err := u.Where(u.ID.Eq(1)).First()

// 2. 多条件查询 (AND)
// SELECT * FROM users WHERE username = 'admin' AND age > 18
user, err := u.Where(u.Username.Eq("admin"), u.Age.Gt(18)).First()

// 3. OR 查询
// SELECT * FROM users WHERE username = 'admin' OR age > 18
user, err := u.Where(u.Username.Eq("admin")).Or(u.Age.Gt(18)).First()

// 4. 模糊查询 (LIKE)
users, err := u.Where(u.Username.Like("%admin%")).Find()

// 5. IN 查询
users, err := u.Where(u.ID.In(1, 2, 3)).Find()

// 6. 排序与分页
users, err := u.Order(u.CreateTime.Desc()).Limit(10).Offset(0).Find()

// 7. 只查询指定字段
var result []model.User
err := u.Select(u.Username, u.Age).Scan(&result)
```

### 3.3 更新 (Update)

```go
u := q.User

// 1. 更新单个字段
// UPDATE users SET age = 20 WHERE id = 1
_, err := u.Where(u.ID.Eq(1)).Update(u.Age, 20)

// 2. 更新多个字段 (使用 Map)
_, err := u.Where(u.ID.Eq(1)).Updates(map[string]interface{}{
    "age": 20,
    "email": "new@example.com",
})

// 3. 更新多个字段 (使用结构体，注意：零值字段不会被更新)
_, err := u.Where(u.ID.Eq(1)).Updates(&model.User{
    Age: 20,
    Email: "new@example.com",
})
```

### 3.4 删除 (Delete)

```go
u := q.User

// 1. 根据条件删除
// DELETE FROM users WHERE id = 1
_, err := u.Where(u.ID.Eq(1)).Delete()

// 2. 物理删除 (如果有软删除)
_, err := u.Unscoped().Where(u.ID.Eq(1)).Delete()
```

---

## 4. 高级查询功能

### 4.1 关联查询 (Joins)

假设 `User` 表和 `Order` 表有关联。

```go
u := q.User
o := q.Order

var result []struct {
    model.User
    OrderNo string
}

// SELECT users.*, orders.order_no 
// FROM users 
// LEFT JOIN orders ON orders.user_id = users.id
err := u.Select(u.ALL, o.OrderNo).
    LeftJoin(o, o.UserID.EqCol(u.ID)).
    Scan(&result)
```

### 4.2 子查询 (SubQuery)

```go
u := q.User
o := q.Order

// 查询有订单的用户
// SELECT * FROM users WHERE id IN (SELECT user_id FROM orders)
subQuery := o.Select(o.UserID)
users, err := u.Where(u.ID.In(subQuery)).Find()
```

### 4.3 事务 (Transaction)

```go
err := q.Transaction(func(tx *generated.Query) error {
    // 在事务中使用 tx 而不是 q
    if err := tx.User.Create(&user); err != nil {
        return err
    }
    
    if err := tx.Order.Create(&order); err != nil {
        return err
    }
    
    return nil
})
```

---

## 5. 常用条件方法速查表

| 方法 | SQL 示例 | 说明 |
| :--- | :--- | :--- |
| `Eq` | `=` | 等于 |
| `Neq` | `<>` | 不等于 |
| `Gt` | `>` | 大于 |
| `Gte` | `>=` | 大于等于 |
| `Lt` | `<` | 小于 |
| `Lte` | `<=` | 小于等于 |
| `In` | `IN` | 包含 |
| `NotIn` | `NOT IN` | 不包含 |
| `Like` | `LIKE` | 模糊匹配 |
| `Between` | `BETWEEN` | 区间 |
| `IsNull` | `IS NULL` | 为空 |
| `IsNotNull` | `IS NOT NULL` | 不为空 |

## 6. 最佳实践

1.  **始终使用生成的对象**：尽量避免使用 `db.Where("string")`，除非生成的 API 无法满足需求。
2.  **Service 层初始化**：建议在 Service 层持有 `q.User` 等对象的引用，方便调用。
3.  **不要修改生成文件**：`internal/model/generated` 目录下的文件是自动生成的，不要手动修改，否则下次生成会丢失。
