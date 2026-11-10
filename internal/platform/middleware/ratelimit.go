package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(max int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	hits := map[string][]time.Time{}

	return func(c *gin.Context) {
		mu.Lock()

		ip := c.ClientIP()
		cutoff := time.Now().Add(-window)

		var recent []time.Time
		for _, t := range hits[ip] {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}

		allowed := len(recent) < max
		if allowed {
			recent = append(recent, time.Now())
		}

		hits[ip] = recent
		mu.Unlock()

		if !allowed {
			stop(c, 429, "too many attempts, try again later")
			return
		}

		c.Next()
	}
}
