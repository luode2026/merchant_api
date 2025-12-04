#!/bin/bash

# 管理员登录系统测试脚本

BASE_URL="http://localhost:8080"
ADMIN_ACCOUNT="admin"
ADMIN_PASSWORD="password123"

echo "========================================="
echo "管理员登录系统测试"
echo "========================================="
echo ""

# 测试 1: 登录
echo "测试 1: 登录接口"
echo "----------------------------------------"
LOGIN_RESPONSE=$(curl -s -X POST ${BASE_URL}/mer_admin/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"account\": \"${ADMIN_ACCOUNT}\",
    \"password\": \"${ADMIN_PASSWORD}\"
  }")

echo "响应: ${LOGIN_RESPONSE}"
echo ""

# 提取 token
TOKEN=$(echo ${LOGIN_RESPONSE} | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ 登录失败，无法获取 token"
  exit 1
else
  echo "✅ 登录成功，Token: ${TOKEN:0:50}..."
fi

echo ""
echo "========================================="
echo ""

# 测试 2: 访问受保护接口
echo "测试 2: 访问受保护接口"
echo "----------------------------------------"
PROFILE_RESPONSE=$(curl -s -X GET ${BASE_URL}/mer_admin/profile \
  -H "Authorization: Bearer ${TOKEN}")

echo "响应: ${PROFILE_RESPONSE}"

if echo ${PROFILE_RESPONSE} | grep -q "admin_id" && echo ${PROFILE_RESPONSE} | grep -q "mer_id"; then
  echo "✅ 认证成功，可以访问受保护接口，且包含 mer_id"
else
  echo "❌ 认证失败或缺少 mer_id"
fi

echo ""
echo "========================================="
echo ""

# 测试 3: 无效 Token
echo "测试 3: 使用无效 Token"
echo "----------------------------------------"
INVALID_RESPONSE=$(curl -s -X GET ${BASE_URL}/mer_admin/profile \
  -H "Authorization: Bearer invalid_token_123")

echo "响应: ${INVALID_RESPONSE}"

if echo ${INVALID_RESPONSE} | grep -q "401"; then
  echo "✅ 正确拒绝无效 Token"
else
  echo "❌ 应该拒绝无效 Token"
fi

echo ""
echo "========================================="
echo ""

# 测试 4: 登出
echo "测试 4: 登出接口"
echo "----------------------------------------"
LOGOUT_RESPONSE=$(curl -s -X POST ${BASE_URL}/mer_admin/auth/logout \
  -H "Authorization: Bearer ${TOKEN}")

echo "响应: ${LOGOUT_RESPONSE}"

if echo ${LOGOUT_RESPONSE} | grep -q "登出成功"; then
  echo "✅ 登出成功"
else
  echo "❌ 登出失败"
fi

echo ""
echo "========================================="
echo ""

# 测试 5: 登出后访问受保护接口
echo "测试 5: 登出后访问受保护接口"
echo "----------------------------------------"
AFTER_LOGOUT_RESPONSE=$(curl -s -X GET ${BASE_URL}/mer_admin/profile \
  -H "Authorization: Bearer ${TOKEN}")

echo "响应: ${AFTER_LOGOUT_RESPONSE}"

if echo ${AFTER_LOGOUT_RESPONSE} | grep -q "已失效"; then
  echo "✅ 正确拒绝已登出的 Token"
else
  echo "⚠️  Token 可能仍然有效（检查 Redis）"
fi

echo ""
echo "========================================="
echo "测试完成"
echo "========================================="
