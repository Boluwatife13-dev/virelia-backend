package guardians

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrNoGuardian: 404,
	ErrExists:     409,
	ErrNotFound:   404,
	ErrNotChild:   400,
}

func Routes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	g := r.Group("/guardians", guards...)

	g.POST("/children", func(c *gin.Context) {
		var req FrontDeskChild

		if !httpx.Bind(c, &req) {
			return
		}

		me := c.GetString("user_id")
		facility := c.GetString("facility_id")

		id, err := s.RegisterByFrontDesk(c, req, me, facility)

		added(c, id, err)
	})

	g.GET("", func(c *gin.Context) {
		list, err := s.GuardiansOf(c, c.Query("child"))

		data := gin.H{"guardians": list}

		httpx.Respond(c, 200, data, err, statusFor)
	})

	g.POST("/:id/verify", func(c *gin.Context) {
		err := s.Verify(c, c.Param("id"))

		reply(c, "link verified", err)
	})

	g.POST("/:id/revoke", func(c *gin.Context) {
		err := s.Revoke(c, c.Param("id"))

		reply(c, "link revoked", err)
	})
}

func added(c *gin.Context, networkID string, err error) {
	data := gin.H{
		"message":    "child added",
		"network_id": networkID,
	}

	httpx.Respond(c, 201, data, err, statusFor)
}

func reply(c *gin.Context, msg string, err error) {
	data := gin.H{"message": msg}

	httpx.Respond(c, 200, data, err, statusFor)
}
