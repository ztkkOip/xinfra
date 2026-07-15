package router

import (
	"github.com/1024XEngineer/xinfra/server/internal/handler"
	"github.com/1024XEngineer/xinfra/server/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 初始化路由
// @title           xinfra API
// @version         1.0
// @description     xinfra 平台后端 API 文档
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 Bearer Token（例如：Bearer xxx）
func Setup() *gin.Engine {
	r := gin.Default()

	// 加载中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 健康检查
	r.GET("/api/v1/ping", handler.Ping)

	// 测试路由组 — 用于验证统一响应格式
	setupTestRoutes(r)

	return r
}

// setupTestRoutes 注册测试路由
func setupTestRoutes(r *gin.Engine) {
	testGroup := r.Group("/api/v1/test")
	{
		// 成功响应
		testGroup.GET("/success", handler.TestSuccess)

		// 业务错误 — 传入错误码
		// GET /api/v1/test/error/10001 → LDAP 认证失败
		// GET /api/v1/test/error/10002 → 用户名或密码错误
		// GET /api/v1/test/error/10003 → Token 已过期
		// GET /api/v1/test/error/20001 → 任务创建失败
		// GET /api/v1/test/error/20002 → 任务执行超时
		testGroup.GET("/error/:code", handler.TestError)

		// HTTP 500 错误
		testGroup.GET("/500", handler.Test500)

		// HTTP 401 未授权
		testGroup.GET("/401", handler.Test401)

		// HTTP 403 禁止访问
		testGroup.GET("/403", handler.Test403)

		// 模拟超时（sleep 超过前端 10s timeout）
		testGroup.GET("/timeout", handler.TestTimeout)

		// 返回分页数据
		testGroup.GET("/paginated", handler.TestPaginated)
	}
}
