package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !slices.Contains(roles, c.GetString("role")) {
			stop(c, 403, "you do not have access to this")
			return
		}

		c.Next()
	}
}
