package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var statusFor = map[error]int{
	ErrExists:        http.StatusConflict,
	ErrBadDOB:        http.StatusBadRequest,
	ErrInvalidCode:   http.StatusBadRequest,
	ErrWrongPassword: http.StatusBadRequest,
	ErrWrongLogin:    http.StatusUnauthorized,
	ErrInvalidToken:  http.StatusUnauthorized,
	ErrNotVerified:   http.StatusForbidden,
	ErrInactive:      http.StatusForbidden,
	ErrUseKey:        http.StatusForbidden,
}

func on[T any](
	status int,
	run func(*gin.Context, T) (any, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		data, err := run(c, req)

		respond(c, status, data, err)
	}
}

func respond(
	c *gin.Context,
	status int,
	data any,
	err error,
) {
	if err == nil {
		c.JSON(status, data)
		return
	}

	for known, code := range statusFor {
		if errors.Is(err, known) {
			c.JSON(code, gin.H{"error": known.Error()})
			return
		}
	}

	c.JSON(500, gin.H{"error": "something went wrong"})
}

func msg(text string) gin.H {
	return gin.H{"message": text}
}
