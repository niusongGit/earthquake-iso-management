package middleware

import (
	"earthquake-iso-management/internal/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recover panic恢复中间件，捕获接口运行时的panic并记录日志
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 打印错误堆栈信息
				debug.PrintStack()

				// 记录panic日志
				logger.Log.Error("接口panic",
					zap.Any("err", r),
					zap.String("log_from", "recover_middleware"),
					zap.String("method", c.Request.Method),
					zap.String("uri", c.Request.RequestURI),
					zap.String("ip", c.ClientIP()),
					zap.ByteString("stack", debug.Stack()),
				)

				// 返回统一错误响应
				c.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": errorToString(r),
					"data":    nil,
				})
				// 终止后续接口调用
				c.Abort()
			}
		}()
		c.Next()
	}
}

// errorToString 将recover捕获的错误转为字符串
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "未知错误"
	}
}
