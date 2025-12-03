# 模型代码生成工具

## 功能说明

使用 `gorm.io/gen` 根据数据库表自动生成 Go 模型文件和查询接口。

### GORM 模型生成器使用指南

本工具基于 `gorm.io/gen` 自动从数据库表生成 Go 模型代码和类型安全的 DAO 查询代码。

## 特性

- ✅ 自动从数据库表生成 Model 结构体
- ✅ 自动生成类型安全的 DAO 查询代码
- ✅ 支持指定表名或生成所有表
- ✅ 自动添加 JSON、GORM、Form 标签
- ✅ DAO 文件自动添加 `_dao` 后缀便于识别

## 目录结构

```
internal/
├── model/          # Model 结构体（数据库表映射）
│   ├── base.go     # 手写基础模型
│   ├── user.go     # 手写用户模型
│   └── eb_merchant.go  # 生成的商户模型
└── dao/            # DAO 查询代码（类型安全查询）
    ├── gen.go      # DAO 入口文件
    └── eb_merchant_dao.go  # 商户表 DAO
```

## 使用方法

### 1. 确保数据库连接配置正确

检查 `configs/config.yaml` 中的数据库配置。

### 2. 运行生成器

**生成所有表（默认）：**
```bash
make gen
```

**生成指定表：**
```bash
make gen table=eb_merchant
```

**生成多个表：**
```bash
make gen table=eb_merchant,eb_merchant_admin
```

### 3. 查看生成的文件

生成的文件位于 `internal/model/generated/` 目录：

- `query.go` - 查询接口
- `*.gen.go` - 各个表的模型文件

## 生成配置

### 字段标签

生成的模型会包含以下标签：

- `gorm` - GORM 标签（字段类型、索引等）
- `json` - JSON 序列化标签
- `form` - 表单绑定标签

### 字段类型映射

- `tinyint` → `int32`
- `smallint` → `int32`
- `mediumint` → `int32`
- `int` → `int32`
- `bigint` → `int64`

## 使用生成的模型

### 基础查询

```go
import (
    "merchant_api/internal/model/generated"
    "merchant_api/pkg/database"
)

// 初始化查询接口
q := generated.Use(database.GetDB())

// 查询所有用户
users, err := q.User.Find()

// 根据 ID 查询
user, err := q.User.Where(q.User.ID.Eq(1)).First()

// 条件查询
users, err := q.User.Where(q.User.Status.Eq(1)).Find()

// 分页查询
users, err := q.User.Limit(10).Offset(0).Find()
```

### 高级查询

```go
// 联表查询
var result []struct {
    generated.User
    MerchantName string
}

err := q.User.
    Select(q.User.ALL, q.Merchant.Name.As("merchant_name")).
    LeftJoin(q.Merchant, q.Merchant.UserID.EqCol(q.User.ID)).
    Scan(&result)
```

### 事务操作

```go
err := q.Transaction(func(tx *generated.Query) error {
    // 创建用户
    user := &generated.User{
        Username: "test",
        Password: "123456",
    }
    if err := tx.User.Create(user); err != nil {
        return err
    }
    
    // 创建商户
    merchant := &generated.Merchant{
        UserID: user.ID,
        Name:   "Test Merchant",
    }
    return tx.Merchant.Create(merchant)
})
```

## 自定义配置

如需自定义生成配置，可以修改 `cmd/generator/main.go`：

### 只生成指定表

```go
g.ApplyBasic(
    g.GenerateModel("users"),
    g.GenerateModel("merchants"),
)
```

### 添加自定义方法

```go
user := g.GenerateModel("users",
    gen.FieldJSONTag,
    gen.FieldGORMTag,
)

// 添加自定义方法
g.ApplyInterface(func(method gen.Method) {
    // 自定义查询方法
}, user)
```

### 修改输出目录

```go
g := gen.NewGenerator(gen.Config{
    OutPath: "./your/custom/path",
    // ...
})
```

## 注意事项

1. **备份现有模型**：生成前请备份 `internal/model/` 目录下的手写模型
2. **不要修改生成的文件**：生成的文件会被覆盖，自定义逻辑应写在其他文件中
3. **版本控制**：建议将生成的文件加入 `.gitignore`，或者只提交必要的文件

## 常见问题

### Q: 生成的模型字段类型不对？

A: 修改 `cmd/generator/main.go` 中的 `dataMap` 配置

### Q: 如何添加自定义字段？

A: 在 `internal/model/` 目录创建扩展文件，使用组合方式

### Q: 生成失败怎么办？

A: 检查数据库连接配置和表结构是否正确
