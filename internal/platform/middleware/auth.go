package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"virelia/internal/auth"
)

func RequireLogin(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, _ := strings.CutPrefix(header, "Bearer ")

		claims, err := auth.ParseAccessToken(token, secret)
		if err != nil {
			stop(c, 401, "please log in")
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("role", claims.Role)
		c.Set("facility_id", claims.FacilityID)

		c.Next()
	}
}

func stop(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
