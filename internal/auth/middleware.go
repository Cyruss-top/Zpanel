package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/model"
)

// Middleware JWT 认证中间件
func (s *Service) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Fail("unauthorized", "UNAUTHORIZED"))
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := s.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Fail("invalid token", "UNAUTHORIZED"))
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
