package httpx

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Bind reads the JSON body. On failure it replies 400.
func Bind(c *gin.Context, req any) bool {
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	return err == nil
}

func Respond(
	c *gin.Context, status int, data any,
	err error, known map[error]int,
) {
	if err == nil {
		c.JSON(status, data)
		return
	}

	for e, code := range known {
		if errors.Is(err, e) {
			c.JSON(code, gin.H{"error": e.Error()})
			return
		}
	}

	c.JSON(500, gin.H{"error": "something went wrong"})
}
