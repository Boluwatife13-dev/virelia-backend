package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"
)

func CORS(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		h := c.Writer.Header()

		h.Add("Vary", "Origin")

		if slices.Contains(allowed, origin) {
			h.Set("Access-Control-Allow-Origin", origin)

			h.Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)

			h.Set(
				"Access-Control-Allow-Headers",
				"Authorization, Content-Type",
			)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
