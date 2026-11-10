package keys

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrNotAdmin: 403,
	ErrNoKey:    400,
	ErrHasKey:   400,
	ErrBadKey:   401,
	ErrExpired:  400,
}

// Routes for platform admins. The step-1 login
// token goes in a header, because the body is
// what the key signed.
func Routes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	k := r.Group("/auth/staff/key")

	k.POST("/setup/begin", func(c *gin.Context) {
		options, err := s.BeginSetup(c, challenge(c))

		httpx.Respond(c, 200, options, err, statusFor)
	})

	k.POST("/setup/finish", func(c *gin.Context) {
		agent := c.Request.UserAgent()

		tokens, err := s.FinishSetup(c, challenge(c), c.Request, agent, c.ClientIP())

		httpx.Respond(c, 201, tokens, err, statusFor)
	})

	k.POST("/begin", func(c *gin.Context) {
		options, err := s.BeginLogin(c, challenge(c))

		httpx.Respond(c, 200, options, err, statusFor)
	})

	k.POST("/finish", func(c *gin.Context) {
		agent := c.Request.UserAgent()

		tokens, err := s.FinishLogin(c, challenge(c), c.Request, agent, c.ClientIP())

		httpx.Respond(c, 200, tokens, err, statusFor)
	})

	extraKeyRoutes(r, s, guards...)
}

func challenge(c *gin.Context) string {
	return c.GetHeader("X-Challenge-Token")
}
