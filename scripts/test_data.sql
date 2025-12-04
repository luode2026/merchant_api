-- 管理员登录系统 - 正确的测试数据
-- 使用实际的 bcrypt 哈希值

-- 方案1: 如果表中已有数据，更新密码
UPDATE mer_merchant_admin 
SET pwd = '$2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa'  -- admin888
WHERE account = 'lucky';

-- 方案2: 如果表中没有数据，插入新记录
-- 账号: lucky, 密码: admin888
INSERT INTO mer_merchant_admin (
  mer_id, 
  account, 
  pwd, 
  real_name, 
  phone, 
  login_count, 
  level, 
  is_del, 
  status
) VALUES (
  1, 
  'lucky', 
  '$2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa',  -- admin888
  '幸运管理员',
  '13800138001',
  0, 
  1, 
  0, 
  1
)
ON DUPLICATE KEY UPDATE 
  pwd = '$2a$10$mSQTjoKQEy2s9Sqw0oZaF.mk5JOFQND9/tE5xMG08Odk0y/hYCXpa';

-- 额外创建一个测试账号
-- 账号: admin, 密码: password123
INSERT INTO mer_merchant_admin (
  mer_id, 
  account, 
  pwd, 
  real_name, 
  phone, 
  login_count, 
  level, 
  is_del, 
  status
) VALUES (
  1, 
  'admin', 
  '$2a$10$n8d1uPbjOA4KIftZ5nZaxO05ls1vvSCl0L4ae2UcS4oBhZLeG8voq',  -- password123
  '测试管理员',
  '13800138000',
  0, 
  1, 
  0, 
  1
)
ON DUPLICATE KEY UPDATE 
  pwd = '$2a$10$n8d1uPbjOA4KIftZ5nZaxO05ls1vvSCl0L4ae2UcS4oBhZLeG8voq';

-- 查询验证
SELECT 
  merchant_admin_id,
  account,
  real_name,
  phone,
  status,
  login_count,
  last_ip,
  last_time,
  LENGTH(pwd) as pwd_length,
  LEFT(pwd, 20) as pwd_preview
FROM mer_merchant_admin 
WHERE account IN ('lucky', 'admin')
ORDER BY account;
