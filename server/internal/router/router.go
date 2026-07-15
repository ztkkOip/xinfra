package router

import (
	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/handler"
	"github.com/1024XEngineer/xinfra/server/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Dependencies struct {
	Config config.Config
	DB     *gorm.DB
}

// New 初始化路由。
// @title           xinfra API
// @version         1.0
// @description     xinfra 平台后端 API 文档
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 Bearer Token（例如：Bearer xxx）
func New(deps Dependencies) *gin.Engine {
	if deps.Config.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	registerSwaggerRoutes(r)
	registerTestRoutes(r)
	registerAuthServerRoutes(r, deps)

	return r
}

func registerSwaggerRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
}

func registerTestRoutes(r *gin.Engine) {
	r.GET("/api/v1/ping", handler.Ping)

	testGroup := r.Group("/api/v1/test")
	{
		testGroup.GET("/success", handler.TestSuccess)
		testGroup.GET("/error/:code", handler.TestError)
		testGroup.GET("/500", handler.Test500)
		testGroup.GET("/401", handler.Test401)
		testGroup.GET("/403", handler.Test403)
		testGroup.GET("/timeout", handler.TestTimeout)
		testGroup.GET("/paginated", handler.TestPaginated)
	}
}

func registerAuthServerRoutes(r *gin.Engine, deps Dependencies) {
	auditService := service.NewAuditService(deps.DB)
	authService := service.NewAuthService(deps.Config, deps.DB, auditService)
	wayenService := service.NewWayenService(deps.Config, deps.DB)

	healthHandler := handler.NewHealthHandler(deps.DB)
	userHandler := handler.NewUserHandler()
	wayenHandler := handler.NewWayenHandler(deps.DB, wayenService, auditService)
	clouddmHandler := handler.NewCloudDMHandler(deps.Config, auditService)
	samlHandler := handler.NewSAMLHandler(deps.Config, authService)
	oauthHandler := handler.NewOAuthHandler(deps.Config, deps.DB, auditService)

	r.GET("/healthz", healthHandler.Healthz)
	r.GET("/readyz", healthHandler.Readyz)
	r.GET("/auth/.well-known/openid-configuration", oauthHandler.Discovery)
	r.GET("/auth/oauth/authorize", oauthHandler.Authorize)
	r.POST("/auth/oauth/token", oauthHandler.Token)
	r.GET("/auth/oauth/jwks", oauthHandler.JWKS)
	r.GET("/auth/oauth/userinfo", oauthHandler.UserInfo)

	v1 := r.Group("/auth/api/v1")
	{
		v1.GET("/login/internal-sso", samlHandler.Login)
		v1.GET("/saml/metadata", samlHandler.Metadata)
		v1.POST("/saml/acs", samlHandler.ACS)

		protected := v1.Group("")
		protected.Use(handler.AuthMiddleware(deps.Config))
		protected.GET("/users/me", userHandler.Me)
		protected.GET("/wayen/login", wayenHandler.Login)
		protected.GET("/wayen/credential", wayenHandler.GetCredential)
		protected.PUT("/wayen/credential", wayenHandler.SaveCredential)
		protected.GET("/clouddm/login", clouddmHandler.Login)
	}
}
