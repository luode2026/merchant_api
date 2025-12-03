# Merchant API

基于 Go + Gin + GORM 的多应用系统架构项目，参考 GoAdmin 设计理念。
基于 Go 的多应用系统架构项目，采用 Gin + GORM + Viper + Zap 技术栈。

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM (MySQL)
- **配置管理**: Viper
- **日志**: Zap
- **认证**: JWT
- **缓存**: Redis
- **密码加密**: bcrypt

## 项目结构

```
merchant_api/
├── cmd/                    # 应用入口
│   ├── admin/             # Admin 端服务
│   └── app/               # App 端服务
├── internal/              # 内部代码
│   ├── admin/            # Admin 端特定代码
│   │   └── router/       # 路由配置
│   ├── app/              # App 端特定代码
│   │   └── router/       # 路由配置
│   ├── model/            # 数据模型
│   ├── dao/              # 数据访问层（自动生成）
│   ├── middleware/       # 中间件
│   └── pkg/              # 内部工具包
│       ├── response/     # 统一响应
│       ├── jwt/          # JWT 工具
│       └── utils/        # 工具函数
├── pkg/                   # 公共包
│   ├── config/           # 配置加载
│   ├── logger/           # 日志
│   ├── database/         # 数据库连接
│   └── redis/            # Redis 连接
├── configs/              # 配置文件
├── docs/                 # 文档
└── Makefile             # 构建脚本
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 2. 配置

复制配置文件并修改：

```bash
cp configs/config.yaml.example configs/config.yaml
```

编辑 `configs/config.yaml`，配置数据库和 Redis 连接信息。

### 3. 安装依赖

```bash
make tidy
```

### 4. 启动服务

**启动 Admin 端（端口 8080）：**
```bash
make admin
```

**启动 App 端（端口 8081）：**
```bash
make app
```

### 5. 健康检查

```bash
# Admin 端
curl http://localhost:8080/health

# App 端
curl http://localhost:8081/health
```

## 代码生成

使用 `gorm.io/gen` 从数据库表自动生成模型和 DAO 代码：

```bash
# 生成所有表
make gen

# 生成指定表
make gen table=eb_merchant

# 生成多个表
make gen table=eb_merchant,eb_merchant_admin
```

生成的文件：
- Model 结构体 → `internal/model/`
- DAO 查询代码 → `internal/dao/`

详细说明请参考：[模型生成器文档](docs/model_generator.md)

## 开发指南

### 添加新接口

1. 在 `internal/admin/router/router.go` 或 `internal/app/router/router.go` 中添加路由
2. 创建对应的 controller 处理函数
3. 如需数据库操作，使用生成的 DAO 代码或创建 service 层

### 中间件

- `Logger()`: 请求日志
- `Recovery()`: 错误恢复
- `CORS()`: 跨域处理
- `JWTAuth()`: JWT 认证

### 统一响应格式

使用 `internal/pkg/response` 包：

```go
import "merchant_api/internal/pkg/response"

// 成功响应
response.Success(c, data)

// 失败响应
response.Error(c, "错误信息")

// 分页响应
response.SuccessWithPagination(c, data, total, page, pageSize)
```

## 文档

- [GORM Gen 使用指南](docs/gorm_gen_usage.md)
- [模型生成器文档](docs/model_generator.md)

## License

MIT
