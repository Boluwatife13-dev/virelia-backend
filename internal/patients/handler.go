package patients

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrExists:       409,
	ErrBadBirthDate: 400,
	ErrNoSearch:     400,
}

func Routes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	p := r.Group("/patients", guards...)

	p.POST("", func(c *gin.Context) {
		var req RegisterRequest

		if !httpx.Bind(c, &req) {
			return
		}

		me := c.GetString("user_id")
		facility := c.GetString("facility_id")

		id, err := s.Register(c, req, me, facility)

		data := gin.H{
			"message":    "patient registered",
			"network_id": id,
		}

		httpx.Respond(c, 201, data, err, statusFor)
	})

	p.GET("/search", func(c *gin.Context) {
		q := Search{
			Phone:     c.Query("phone"),
			NIN:       c.Query("nin"),
			NetworkID: c.Query("network_id"),
		}

		list, err := s.Find(c, q)

		data := gin.H{"patients": list}

		httpx.Respond(c, 200, data, err, statusFor)
	})
}
