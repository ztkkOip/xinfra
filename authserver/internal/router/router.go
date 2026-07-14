package router

import (
	"authserver/internal/config"
	"authserver/internal/handler"
	"authserver/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Dependencies struct {
	Config config.Config
	DB     *gorm.DB
}

func New(deps Dependencies) *gin.Engine {
	if deps.Config.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

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

	return r
}
