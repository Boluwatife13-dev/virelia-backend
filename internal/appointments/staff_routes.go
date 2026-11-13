package appointments

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

func DeskRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	d := r.Group("/appointments", guards...)

	d.GET("", func(c *gin.Context) {
		facility := c.GetString("facility_id")

		list, err := s.Day(c, "a.facility_id", facility, c.Query("date"))

		send(c, list, err)
	})

	d.POST("", func(c *gin.Context) {
		var req DeskBookRequest

		if !httpx.Bind(c, &req) {
			return
		}

		me := c.GetString("user_id")
		facility := c.GetString("facility_id")

		id, err := s.BookAtDesk(c, me, facility, req)

		booked(c, id, err)
	})

	d.POST("/:id/check-in", deskStatus(s, "checked_in", "checked in"))
	d.POST("/:id/no-show", deskStatus(s, "no_show", "marked as no-show"))
	d.POST("/:id/cancel", deskStatus(s, "cancelled", "appointment cancelled"))
}

func deskStatus(s *Service, to, msg string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CancelRequest

		c.ShouldBindJSON(&req)

		facility := c.GetString("facility_id")
		me := c.GetString("user_id")

		err := s.AtDesk(c, facility, me, c.Param("id"), to, req.Reason)

		reply(c, msg, err)
	}
}

func DoctorRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	d := r.Group("/doctor/appointments", guards...)

	d.GET("", func(c *gin.Context) {
		me := c.GetString("user_id")

		list, err := s.Day(c, "a.doctor_id", me, c.Query("date"))

		send(c, list, err)
	})

	d.POST("/:id/complete", func(c *gin.Context) {
		err := s.Complete(c, c.GetString("user_id"), c.Param("id"))

		reply(c, "visit completed", err)
	})
}
