package controller

import (
	"merchant_api/internal/admin/service"
	"merchant_api/internal/pkg/response"
	"merchant_api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AdminAuthController struct{}

func NewAdminAuthController() *AdminAuthController {
	return &AdminAuthController{}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`  // 账号或手机号
	Password string `json:"password" binding:"required"` // 密码
}

// Login 管理员登录
func (ctrl *AdminAuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithKey(c, "error.invalid_params", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	// 获取客户端 IP
	ip := utils.GetClientIP(c)

	// 调用服务层
	authService := service.NewAdminAuthService(c.Request.Context())
	resp, err := authService.Login(req.Account, req.Password, ip)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, resp)
}

// Logout 管理员登出
func (ctrl *AdminAuthController) Logout(c *gin.Context) {
	// 从 header 获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		response.UnauthorizedWithKey(c, "error.auth.no_credentials")
		return
	}

	token := authHeader[7:] // 去掉 "Bearer "

	// 调用服务层
	authService := service.NewAdminAuthService(c.Request.Context())
	err := authService.Logout(token)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithKey(c, "success.logout", nil)
}
