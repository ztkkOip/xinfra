package handler

import (
	"net/http"
	"strings"

	"github.com/1024XEngineer/xinfra/server/internal/auth"
	"github.com/1024XEngineer/xinfra/server/internal/config"

	"github.com/gin-gonic/gin"
)

const ClaimsKey = "claims"

func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")
		tokenValue := ""
		if strings.HasPrefix(value, "Bearer ") {
			tokenValue = strings.TrimPrefix(value, "Bearer ")
		} else {
			tokenValue = strings.TrimSpace(c.Query("access_token"))
		}
		if tokenValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := auth.Parse(cfg.JWTSecret, tokenValue)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func CurrentClaims(c *gin.Context) (*auth.Claims, bool) {
	value, ok := c.Get(ClaimsKey)
	if !ok {
		return nil, false
	}
	claims, ok := value.(*auth.Claims)
	return claims, ok
}
