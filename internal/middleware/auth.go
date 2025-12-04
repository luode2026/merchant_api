package middleware

import (
	"merchant_api/internal/admin/service"
	"merchant_api/internal/pkg/jwt"
	"merchant_api/internal/pkg/response"
	"merchant_api/internal/pkg/utils"
	"merchant_api/pkg/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.UnauthorizedWithKey(c, "error.auth.no_credentials")
			c.Abort()
			return
		}

		// 解析 Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.UnauthorizedWithKey(c, "error.auth.invalid_format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString, config.GlobalConfig.JWT.Secret)
		if err != nil {
			response.UnauthorizedWithKey(c, "error.auth.invalid_token")
			c.Abort()
			return
		}

		// 检查角色权限
		if requiredRole != "" && claims.Role != requiredRole {
			response.ForbiddenWithKey(c, "error.auth.insufficient_permissions")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("merchant_admin_id", claims.UserID)
		c.Set("mer_id", claims.MerID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件（包含 Redis 和 IP 验证）
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.UnauthorizedWithKey(c, "error.auth.no_credentials")
			c.Abort()
			return
		}

		// 解析 Token (格式: Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.UnauthorizedWithKey(c, "error.auth.invalid_format")
			c.Abort()
			return
		}

		token := parts[1]

		// 获取客户端 IP
		currentIP := utils.GetClientIP(c)

		// 验证 Token（包含 Redis 和 IP 验证）
		authService := service.NewAdminAuthService(c.Request.Context())
		claims, err := authService.VerifyToken(token, currentIP)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("admin_id", claims.UserID)
		c.Set("mer_id", claims.MerID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
