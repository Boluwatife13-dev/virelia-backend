package facilities

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrBadExpiry:   400,
	ErrNotPending:  404,
	ErrNeedsSecond: 409,
}

func Routes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	r.POST("/facilities/apply", func(c *gin.Context) {
		var req ApplyRequest

		if !httpx.Bind(c, &req) {
			return
		}

		_, err := s.Apply(c, req)

		reply(c, 201, "application received", err)
	})

	p := r.Group("/platform/facilities", guards...)

	p.GET("", func(c *gin.Context) {
		list, err := s.List(c, c.Query("status"))

		data := gin.H{"facilities": list}

		httpx.Respond(c, 200, data, err, statusFor)
	})

	p.POST("/:id/check", func(c *gin.Context) {
		var req CheckRequest

		if !httpx.Bind(c, &req) {
			return
		}

		id := c.Param("id")
		me := c.GetString("user_id")

		err := s.Check(c, id, me, req.Note)

		reply(c, 200, "checked, waiting for a second admin", err)
	})

	p.POST("/:id/approve", func(c *gin.Context) {
		id := c.Param("id")
		me := c.GetString("user_id")

		err := s.Approve(c, id, me)

		reply(c, 200, "approved, admin invited", err)
	})

	p.POST("/:id/reject", func(c *gin.Context) {
		var req RejectRequest

		if !httpx.Bind(c, &req) {
			return
		}

		id := c.Param("id")
		me := c.GetString("user_id")

		err := s.Reject(c, id, me, req.Reason)

		reply(c, 200, "rejected", err)
	})
}

func reply(c *gin.Context, status int, msg string, err error) {
	data := gin.H{"message": msg}

	httpx.Respond(c, status, data, err, statusFor)
}
