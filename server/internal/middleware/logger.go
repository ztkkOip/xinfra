package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回请求日志中间件
// 记录请求方法、路径、状态码、耗时、客户端 IP
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 记录日志
		log.Printf("[REQUEST] %s %s | %d | %v | %s",
			method,
			path,
			statusCode,
			duration,
			clientIP,
		)
	}
}
