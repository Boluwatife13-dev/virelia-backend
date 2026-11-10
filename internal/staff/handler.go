package staff

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrBadInvite: 400,
	ErrExists:    409,
	ErrNotFound:  404,
	ErrSelf:      400,
}

func Routes(r *gin.RouterGroup, s *Service, guards ...gin.HandlerFunc) {
	r.POST("/staff/invites/accept", func(c *gin.Context) {
		var req AcceptRequest

		if httpx.Bind(c, &req) {
			reply(c, 201, "account created", s.Accept(c, req))
		}
	})

	a := r.Group("/staff", guards...)

	a.POST("/invites", func(c *gin.Context) {
		var in Invite

		if !httpx.Bind(c, &in) {
			return
		}

		in.FacilityID = c.GetString("facility_id")
		in.InvitedBy = c.GetString("user_id")

		reply(c, 201, "invite sent", s.SendInvite(c, in))
	})

	a.GET("", func(c *gin.Context) {
		staff, err := s.List(c, c.GetString("facility_id"))

		httpx.Respond(c, 200, gin.H{"staff": staff}, err, statusFor)
	})

	a.POST("/:userID/remove", func(c *gin.Context) {
		facility := c.GetString("facility_id")
		me := c.GetString("user_id")

		err := s.Remove(c, facility, c.Param("userID"), me)

		reply(c, 200, "staff member removed", err)
	})

	a.POST("/:userID/reset-2fa", func(c *gin.Context) {
		facility := c.GetString("facility_id")
		me := c.GetString("user_id")

		err := s.Reset2FA(c, facility, c.Param("userID"), me)

		reply(c, 200, "2FA reset", err)
	})
}

func PlatformRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	p := r.Group("/platform/staff", guards...)

	p.POST("/:userID/reset-2fa", func(c *gin.Context) {
		me := c.GetString("user_id")

		err := s.Reset2FA(c, "", c.Param("userID"), me)

		reply(c, 200, "2FA reset", err)
	})
}

func reply(c *gin.Context, status int, msg string, err error) {
	httpx.Respond(c, status, gin.H{"message": msg}, err, statusFor)
}
