# 密码问题排查和解决方案

## 问题描述
登录时提示"账号或密码错误"

## 原因分析

bcrypt 哈希值每次生成都不同（因为包含随机盐值），所以需要使用实际生成的哈希值。

## 解决方案

### 1. 生成正确的密码哈希

运行密码哈希生成工具：
```bash
cd /Users/luode/Documents/www/cb_finance/merchant_api
go run scripts/hash_password.go
```

输出示例：
```
原始密码: admin888
加密后的哈希: $2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa
哈希长度: 60
```

### 2. 更新数据库

执行更新后的 SQL 脚本：
```bash
mysql -h 127.0.0.1 -P 8306 -u root -proot finance < scripts/test_data.sql
```

或者手动执行：
```sql
-- 更新 lucky 账号的密码为 admin888
UPDATE mer_merchant_admin 
SET pwd = '$2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa'
WHERE account = 'lucky';
```

### 3. 测试登录

```bash
curl -X POST http://127.0.0.1:8080/mer_admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "account": "lucky",
    "password": "admin888"
  }'
```

## 测试账号信息

| 账号 | 密码 | 说明 |
|------|------|------|
| lucky | admin888 | 主测试账号 |
| admin | password123 | 备用测试账号 |

## 常见问题

### Q1: 为什么每次生成的哈希都不一样？
A: bcrypt 会在哈希中包含随机盐值，这是正常的安全特性。验证时 bcrypt 会自动提取盐值进行比对。

### Q2: 数据库中的密码字段长度够吗？
A: bcrypt 哈希固定为 60 个字符，数据库字段定义为 `char(64)` 足够。

### Q3: 如何验证密码是否正确？
A: 可以使用 `hash_password.go` 工具生成哈希，然后用 bcrypt 在线工具验证，或者直接测试登录。

## 调试技巧

如果仍然无法登录，可以在 Repository 层添加调试日志：

```go
// 在 admin_auth_repository.go 的 Login 方法中添加
fmt.Printf("数据库密码哈希: %s\n", admin.Pwd)
fmt.Printf("输入的密码: %s\n", req.Password)
fmt.Printf("验证结果: %v\n", utils.CheckPassword(req.Password, admin.Pwd))
```
