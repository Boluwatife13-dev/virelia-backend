package audit

import "github.com/gin-gonic/gin"

func Middleware(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !worthLogging(c) {
			return
		}

		e := From(c, c.Request.Method+" "+c.FullPath())

		params := map[string]any{}

		for _, p := range c.Params {
			params[p.Key] = p.Value

			if e.TargetID == "" {
				e.TargetID = p.Value
			}
		}

		e.Details = map[string]any{
			"params": params,
			"status": c.Writer.Status(),
		}

		s.Log(c, e)
	}
}

func worthLogging(c *gin.Context) bool {
	path := c.FullPath()

	if path == "" || c.Writer.Status() >= 400 {
		return false
	}

	if path == "/health" || path == "/api/v1/auth/refresh" {
		return false
	}

	if c.Request.Method != "GET" {
		return true
	}

	role := c.GetString("role")

	return role != "" && role != "patient"
}
