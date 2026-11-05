package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"virelia-backend/utils"
)

func Authenticate(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		token, found := strings.CutPrefix(header, "Bearer ")
		if !found {
			unauthorized(c)
			return
		}

		claims, err := utils.ParseAccessToken(token, secret)
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{"error": "login required"},
	)
}