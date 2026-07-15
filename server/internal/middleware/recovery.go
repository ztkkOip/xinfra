package middleware

import (
	"log"
	"runtime/debug"

	"github.com/1024XEngineer/xinfra/server/internal/handler"
	"github.com/gin-gonic/gin"
)

// Recovery 返回 panic 恢复中间件
// 捕获请求处理过程中的 panic，返回 500 错误而不是崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 堆栈信息
				log.Printf("[PANIC] %v\n%s", err, debug.Stack())

				// 返回 500 错误
				handler.InternalError(c, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}
