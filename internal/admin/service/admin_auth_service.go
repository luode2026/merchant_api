package service

import (
	"context"
	"errors"
	"fmt"
	"merchant_api/internal/dao"
	"merchant_api/internal/model"
	"merchant_api/internal/pkg/jwt"
	"merchant_api/internal/pkg/utils"
	"merchant_api/pkg/config"
	"merchant_api/pkg/database"
	"merchant_api/pkg/redis"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

type AdminAuthService struct {
	ctx context.Context
}

func NewAdminAuthService(ctx context.Context) *AdminAuthService {
	return &AdminAuthService{
		ctx: ctx,
	}
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string                  `json:"token"`
	AdminInfo *model.MerMerchantAdmin `json:"admin_info"`
	ExpiresIn int                     `json:"expires_in"`
}

// Login 管理员登录
func (s *AdminAuthService) Login(account, password, ip string) (*LoginResponse, error) {
	// 初始化 DAO
	dao.SetDefault(database.GetDB())
	adminDAO := dao.MerMerchantAdmin

	// 查询管理员（支持 account 或 phone 登录）
	admin, err := adminDAO.WithContext(s.ctx).
		Where(
			adminDAO.Account.Eq(account),
		).
		Or(adminDAO.Phone.Eq(account)).
		Where(adminDAO.IsDel.Eq(0)).
		First()

	if err != nil {
		return nil, errors.New("账号或密码错误")
	}

	// 检查账号状态
	if admin.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(password, admin.Pwd) {
		return nil, errors.New("账号或密码错误")
	}

	// 生成 JWT Token
	cfg := config.GlobalConfig
	token, err := jwt.GenerateToken(
		uint(admin.MerchantAdminID),
		uint(admin.MerID),
		admin.Account,
		"admin",
		cfg.JWT.Secret,
		cfg.JWT.Expire,
	)
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	// 存储 Token 到 Redis (key: admin:token:{token}, value: {admin_id}:{ip})
	redisKey := fmt.Sprintf("admin:token:%s", token)
	redisValue := fmt.Sprintf("%d:%s", admin.MerchantAdminID, ip)
	rdb := redis.GetRedis()
	err = rdb.Set(s.ctx, redisKey, redisValue, time.Duration(cfg.JWT.Expire)*time.Second).Err()
	if err != nil {
		return nil, fmt.Errorf("存储 Token 失败: %w", err)
	}

	// 更新登录信息
	now := time.Now()
	loginCount := admin.LoginCount + 1
	_, err = adminDAO.WithContext(s.ctx).
		Where(adminDAO.MerchantAdminID.Eq(admin.MerchantAdminID)).
		Updates(map[string]interface{}{
			"last_ip":     ip,
			"last_time":   now,
			"login_count": loginCount,
		})
	if err != nil {
		// 更新失败不影响登录，只记录日志
		fmt.Printf("更新登录信息失败: %v\n", err)
	}

	// 隐藏密码
	admin.Pwd = ""

	return &LoginResponse{
		Token:     token,
		AdminInfo: admin,
		ExpiresIn: cfg.JWT.Expire,
	}, nil
}

// VerifyToken 验证 Token
func (s *AdminAuthService) VerifyToken(token, currentIP string) (*jwt.Claims, error) {
	// 验证 JWT Token
	cfg := config.GlobalConfig
	claims, err := jwt.ParseToken(token, cfg.JWT.Secret)
	if err != nil {
		return nil, errors.New("Token 无效或已过期")
	}

	// 从 Redis 验证 Token
	redisKey := fmt.Sprintf("admin:token:%s", token)
	rdb := redis.GetRedis()
	redisValue, err := rdb.Get(s.ctx, redisKey).Result()
	if err != nil {
		if err == redisv8.Nil {
			return nil, errors.New("Token 已失效")
		}
		return nil, fmt.Errorf("验证 Token 失败: %w", err)
	}

	// 验证 IP 地址
	// redisValue 格式: {admin_id}:{ip}
	var storedAdminID int
	var storedIP string
	_, err = fmt.Sscanf(redisValue, "%d:%s", &storedAdminID, &storedIP)
	if err != nil {
		return nil, errors.New("Token 数据格式错误")
	}

	if storedIP != currentIP {
		return nil, errors.New("登录 IP 已变更，请重新登录")
	}

	if uint(storedAdminID) != claims.UserID {
		return nil, errors.New("Token 数据不匹配")
	}

	return claims, nil
}

// Logout 登出
func (s *AdminAuthService) Logout(token string) error {
	redisKey := fmt.Sprintf("admin:token:%s", token)
	rdb := redis.GetRedis()
	err := rdb.Del(s.ctx, redisKey).Err()
	if err != nil {
		return fmt.Errorf("登出失败: %w", err)
	}
	return nil
}
