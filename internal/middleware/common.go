package middleware

import (
	pkgi18n "merchant_api/internal/pkg/i18n"
	"merchant_api/internal/pkg/response"
	"merchant_api/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				response.InternalServerErrorWithKey(c, "error.internal")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LocaleMiddleware 国际化中间件
// 从 Accept-Language header 中提取语言偏好，创建对应的 localizer
func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Accept-Language header
		acceptLanguage := c.GetHeader("Accept-Language")

		// Parse language preference
		// Accept-Language format: "zh-CN,zh;q=0.9,en;q=0.8"
		var preferredLang string
		if acceptLanguage != "" {
			// Split by comma to get language preferences
			langs := strings.Split(acceptLanguage, ",")
			if len(langs) > 0 {
				// Get the first preference and remove quality value if present
				firstLang := strings.TrimSpace(langs[0])
				// Remove quality value (e.g., ";q=0.9")
				if idx := strings.Index(firstLang, ";"); idx != -1 {
					firstLang = firstLang[:idx]
				}
				preferredLang = firstLang
			}
		}

		// Default to English if no language specified
		if preferredLang == "" {
			preferredLang = "en"
		}

		// Normalize language code (e.g., "zh-CN" -> "zh", "en-US" -> "en")
		// This allows us to support both "zh" and "zh-CN" with the same language file
		langTag, err := language.Parse(preferredLang)
		if err != nil {
			// If parsing fails, default to English
			preferredLang = "en"
		} else {
			// Get base language (e.g., "zh-CN" -> "zh")
			base, _ := langTag.Base()
			preferredLang = base.String()
		}

		// Create localizer for the preferred language with English as fallback
		bundle := pkgi18n.GetBundle()
		localizer := i18n.NewLocalizer(bundle, preferredLang, "en")

		// Store localizer in context
		c.Set("localizer", localizer)

		c.Next()
	}
}
