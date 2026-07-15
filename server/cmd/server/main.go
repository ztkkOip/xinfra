package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/handler"
	"github.com/1024XEngineer/xinfra/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 加载中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/api/v1/ping", func(c *gin.Context) {
		handler.Success(c, gin.H{"time": time.Now().Format(time.RFC3339)})
	})

	// 测试路由组 — 用于验证统一响应格式
	testGroup := r.Group("/api/v1/test")
	{
		// 成功响应
		testGroup.GET("/success", func(c *gin.Context) {
			handler.Success(c, gin.H{
				"username": "testuser",
				"role":     "admin",
			})
		})

		// 业务错误 — 传入错误码
		// GET /api/v1/test/error/10001 → LDAP 认证失败
		// GET /api/v1/test/error/10002 → 用户名或密码错误
		// GET /api/v1/test/error/10003 → Token 已过期
		// GET /api/v1/test/error/20001 → 任务创建失败
		// GET /api/v1/test/error/20002 → 任务执行超时
		testGroup.GET("/error/:code", func(c *gin.Context) {
			codeStr := c.Param("code")
			code, err := strconv.Atoi(codeStr)
			if err != nil {
				handler.BadRequest(c, "错误码必须是数字")
				return
			}
			details := fmt.Sprintf("这是错误码 %d 的测试响应", code)
			handler.Error(c, http.StatusOK, code, details)
		})

		// HTTP 500 错误
		testGroup.GET("/500", func(c *gin.Context) {
			handler.InternalError(c, "模拟服务器内部错误")
		})

		// HTTP 401 未授权
		testGroup.GET("/401", func(c *gin.Context) {
			handler.Unauthorized(c, "模拟 Token 过期")
		})

		// HTTP 403 禁止访问
		testGroup.GET("/403", func(c *gin.Context) {
			handler.Forbidden(c, "模拟无权限")
		})

		// 模拟超时（sleep 超过前端 10s timeout）
		testGroup.GET("/timeout", func(c *gin.Context) {
			time.Sleep(15 * time.Second)
			handler.Success(c, "这条消息不应该出现")
		})

		// 返回分页数据
		testGroup.GET("/paginated", func(c *gin.Context) {
			items := []gin.H{
				{"id": 1, "name": "item-1"},
				{"id": 2, "name": "item-2"},
				{"id": 3, "name": "item-3"},
			}
			handler.SuccessWithPaginated(c, 100, items)
		})
	}

	fmt.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
