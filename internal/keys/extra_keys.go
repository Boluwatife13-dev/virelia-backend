package keys

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

// extraKeyRoutes let a logged-in platform admin add
// a backup key, in case the first one is lost.
func extraKeyRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	reg := r.Group("/auth/keys", guards...)

	reg.POST("/register/begin", func(c *gin.Context) {
		options, err := s.BeginRegister(c, c.GetString("user_id"))

		httpx.Respond(c, 200, options, err, statusFor)
	})

	reg.POST("/register/finish", func(c *gin.Context) {
		err := s.FinishRegister(c, c.GetString("user_id"), c.Request)

		data := gin.H{"message": "backup key added"}

		httpx.Respond(c, 201, data, err, statusFor)
	})
}
