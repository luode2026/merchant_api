package response

import (
	"merchant_api/internal/pkg/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, msg)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, msg)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, msg)
}

// NotFound 404 错误
func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, msg)
}

// InternalServerError 500 错误
func InternalServerError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, msg)
}

// ============= i18n Support Functions =============

// SuccessWithKey 成功响应（使用 i18n 消息键）
func SuccessWithKey(c *gin.Context, messageKey string, data interface{}) {
	msg := i18n.TSimple(c, messageKey)
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  msg,
		Data: data,
	})
}

// ErrorWithKey 错误响应（使用 i18n 消息键）
func ErrorWithKey(c *gin.Context, code int, messageKey string, templateData map[string]interface{}) {
	msg := i18n.T(c, messageKey, templateData)
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequestWithKey 400 错误（使用 i18n 消息键）
func BadRequestWithKey(c *gin.Context, messageKey string, templateData map[string]interface{}) {
	ErrorWithKey(c, http.StatusBadRequest, messageKey, templateData)
}

// UnauthorizedWithKey 401 错误（使用 i18n 消息键）
func UnauthorizedWithKey(c *gin.Context, messageKey string) {
	ErrorWithKey(c, http.StatusUnauthorized, messageKey, nil)
}

// ForbiddenWithKey 403 错误（使用 i18n 消息键）
func ForbiddenWithKey(c *gin.Context, messageKey string) {
	ErrorWithKey(c, http.StatusForbidden, messageKey, nil)
}

// NotFoundWithKey 404 错误（使用 i18n 消息键）
func NotFoundWithKey(c *gin.Context, messageKey string) {
	ErrorWithKey(c, http.StatusNotFound, messageKey, nil)
}

// InternalServerErrorWithKey 500 错误（使用 i18n 消息键）
func InternalServerErrorWithKey(c *gin.Context, messageKey string) {
	ErrorWithKey(c, http.StatusInternalServerError, messageKey, nil)
}
