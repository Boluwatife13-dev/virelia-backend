package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(frontendURLs string) gin.HandlerFunc {
	allowed := strings.Split(frontendURLs, ",")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		header := c.Writer.Header()

		header.Add("Vary", "Origin")

		if slices.Contains(allowed, origin) {
			header.Set("Access-Control-Allow-Origin", origin)
			header.Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)
			header.Set(
				"Access-Control-Allow-Headers",
				"Authorization, Content-Type",
			)
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}