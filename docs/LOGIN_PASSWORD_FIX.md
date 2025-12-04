# 登录密码错误问题 - 完整解决方案

## 问题现象

```
POST /mer_admin/auth/login
{
  "account": "lucky",
  "password": "admin888"
}

响应: {"code": 401, "msg": "账号或密码错误"}
```

## 根本原因

数据库中存储的密码哈希值与实际密码不匹配。bcrypt 每次生成的哈希都包含随机盐值，所以需要使用正确生成的哈希值。

## 快速解决方案

### 步骤 1: 生成正确的密码哈希

```bash
cd /Users/luode/Documents/www/cb_finance/merchant_api

# 为 admin888 生成哈希
go run scripts/gen_hash.go admin888
```

输出示例：
```
✅ 密码加密成功
原始密码: admin888
Bcrypt 哈希: $2a$10$BO./4wpEJMfrDZqdACIoyOjjLjOXqs4p3SwUvlfdolYVeSeIBWgoW

SQL 更新语句:
UPDATE mer_merchant_admin SET pwd = '$2a$10$BO./4wpEJMfrDZqdACIoyOjjLjOXqs4p3SwUvlfdolYVeSeIBWgoW' WHERE account = 'YOUR_ACCOUNT';
```

### 步骤 2: 更新数据库

**方案 A: 使用提供的 SQL 脚本**

```bash
mysql -h 127.0.0.1 -P 8306 -u root -proot finance < scripts/test_data.sql
```

**方案 B: 手动执行 SQL**

```sql
-- 连接数据库
mysql -h 127.0.0.1 -P 8306 -u root -proot finance

-- 更新 lucky 账号密码为 admin888
UPDATE mer_merchant_admin 
SET pwd = '$2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa'
WHERE account = 'lucky';

-- 验证更新
SELECT account, real_name, LENGTH(pwd) as pwd_length, status 
FROM mer_merchant_admin 
WHERE account = 'lucky';
```

### 步骤 3: 测试登录

```bash
curl -X POST http://127.0.0.1:8080/mer_admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "account": "lucky",
    "password": "admin888"
  }'
```

预期成功响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin_info": {
      "merchant_admin_id": 1,
      "account": "lucky",
      "real_name": "幸运管理员",
      ...
    },
    "expires_in": 7200
  }
}
```

## 测试账号列表

| 账号 | 密码 | Bcrypt 哈希（示例） |
|------|------|-------------------|
| lucky | admin888 | `$2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa` |
| admin | password123 | `$2a$10$n8d1uPbjOA4KIftZ5nZaxO05ls1vvSCl0L4ae2UcS4oBhZLeG8voq` |

> ⚠️ **注意**: bcrypt 哈希每次生成都不同（因为包含随机盐值），这是正常的。请使用 `gen_hash.go` 工具生成最新的哈希值。

## 工具使用

### 生成任意密码的哈希

```bash
# 生成自定义密码的哈希
go run scripts/gen_hash.go YOUR_PASSWORD

# 示例
go run scripts/gen_hash.go admin888
go run scripts/gen_hash.go mypassword123
```

### 批量生成常用密码哈希

```bash
# 不带参数运行，显示常用密码的哈希
go run scripts/gen_hash.go
```

## 调试技巧

如果更新后仍然无法登录，可以添加调试日志：

### 1. 在 Repository 层添加日志

编辑 `internal/admin/repository/admin_auth_repository.go`：

```go
// 在 Login 方法的密码验证前添加
fmt.Printf("=== 密码验证调试 ===\n")
fmt.Printf("账号: %s\n", req.Account)
fmt.Printf("输入密码: %s\n", req.Password)
fmt.Printf("数据库哈希: %s\n", admin.Pwd)
fmt.Printf("哈希长度: %d\n", len(admin.Pwd))

// 验证密码
isValid := utils.CheckPassword(req.Password, admin.Pwd)
fmt.Printf("验证结果: %v\n", isValid)
fmt.Printf("==================\n")

if !isValid {
    return nil, errors.New("账号或密码错误")
}
```

### 2. 检查数据库字段

```sql
-- 检查密码字段定义
SHOW CREATE TABLE mer_merchant_admin;

-- 应该看到类似：
-- `pwd` char(64) NOT NULL COMMENT '商户管理员密码'

-- 检查实际数据
SELECT 
  account,
  LENGTH(pwd) as pwd_length,
  LEFT(pwd, 30) as pwd_preview,
  status,
  is_del
FROM mer_merchant_admin 
WHERE account = 'lucky';
```

### 3. 验证 bcrypt 功能

```bash
# 运行密码测试工具
go run scripts/hash_password.go
```

应该看到类似输出：
```
原始密码: admin888
加密后的哈希: $2a$10$...
哈希长度: 60
密码验证结果: true
```

## 常见问题 FAQ

### Q1: 为什么每次生成的哈希都不一样？
**A**: bcrypt 算法会在哈希中包含随机盐值，这是安全特性。虽然哈希值不同，但验证时会自动提取盐值进行正确比对。

### Q2: 数据库密码字段长度够吗？
**A**: bcrypt 哈希固定 60 字符，数据库定义为 `char(64)` 足够。如果是 `varchar(64)` 也可以。

### Q3: 可以用其他工具生成 bcrypt 哈希吗？
**A**: 可以，但建议使用项目提供的 `gen_hash.go` 工具，确保使用相同的 bcrypt 配置（cost=10）。

### Q4: 如何批量更新多个账号的密码？
**A**: 
```bash
# 生成多个密码的哈希
for pwd in admin888 password123 test123; do
  echo "密码: $pwd"
  go run scripts/gen_hash.go $pwd
  echo ""
done
```

### Q5: 登录成功后如何验证 Token？
**A**:
```bash
# 保存登录返回的 token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 访问受保护接口
curl -X GET http://127.0.0.1:8080/mer_admin/profile \
  -H "Authorization: Bearer $TOKEN"
```

## 完整测试流程

```bash
# 1. 生成密码哈希
go run scripts/gen_hash.go admin888

# 2. 复制输出的 SQL 语句，更新数据库
mysql -h 127.0.0.1 -P 8306 -u root -proot finance
# 粘贴并执行 UPDATE 语句

# 3. 测试登录
curl -X POST http://127.0.0.1:8080/mer_admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"lucky","password":"admin888"}'

# 4. 如果成功，会返回 token，保存它
# 5. 使用 token 访问受保护接口
curl -X GET http://127.0.0.1:8080/mer_admin/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 总结

密码错误的主要原因是数据库中的 bcrypt 哈希值不正确。使用提供的工具生成正确的哈希值并更新数据库即可解决问题。

**关键文件：**
- `scripts/gen_hash.go` - 密码哈希生成工具（推荐）
- `scripts/hash_password.go` - 密码测试工具
- `scripts/test_data.sql` - 测试数据 SQL
- `scripts/PASSWORD_FIX.md` - 问题排查文档
