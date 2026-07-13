package handler

import (
	"net/http"
	"strings"

	"authserver/internal/auth"
	"authserver/internal/config"

	"github.com/gin-gonic/gin"
)

const ClaimsKey = "claims"

func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")
		if !strings.HasPrefix(value, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := auth.Parse(cfg.JWTSecret, strings.TrimPrefix(value, "Bearer "))
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
