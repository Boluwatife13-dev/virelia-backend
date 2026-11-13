package directory

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrBadHours: 400,
}

func Routes(
	r *gin.RouterGroup,
	s *Service,
	doctorGuards ...gin.HandlerFunc,
) {
	r.GET("/doctors", func(c *gin.Context) {
		q := Search{
			Specialty:  c.Query("specialty"),
			State:      c.Query("state"),
			FacilityID: c.Query("facility_id"),
			Name:       c.Query("name"),
		}

		list, err := s.Search(c, q)

		httpx.Respond(c, 200, gin.H{"doctors": list}, err, statusFor)
	})

	r.GET("/doctors/:id/hours", func(c *gin.Context) {
		hours, err := s.Hours(c, c.Param("id"))

		httpx.Respond(c, 200, gin.H{"hours": hours}, err, statusFor)
	})

	d := r.Group("/doctor", doctorGuards...)

	d.POST("/profile", func(c *gin.Context) {
		var req ProfileRequest

		if !httpx.Bind(c, &req) {
			return
		}

		err := s.SaveProfile(c, c.GetString("user_id"), req)

		reply(c, "profile saved", err)
	})

	d.POST("/hours", func(c *gin.Context) {
		var req HoursRequest

		if !httpx.Bind(c, &req) {
			return
		}

		err := s.SetHours(c, c.GetString("user_id"), req.Hours)

		reply(c, "working hours saved", err)
	})
}

func reply(c *gin.Context, msg string, err error) {
	httpx.Respond(c, 200, gin.H{"message": msg}, err, statusFor)
}
