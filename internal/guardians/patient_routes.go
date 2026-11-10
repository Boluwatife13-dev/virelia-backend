package guardians

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

func PatientRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	m := r.Group("/me/children", guards...)

	m.GET("", func(c *gin.Context) {
		list, err := s.MyChildren(c, c.GetString("user_id"))

		data := gin.H{"children": list}

		httpx.Respond(c, 200, data, err, statusFor)
	})

	m.POST("", func(c *gin.Context) {
		var req Child

		if !httpx.Bind(c, &req) {
			return
		}

		id, err := s.AddMyChild(c, req, c.GetString("user_id"))

		added(c, id, err)
	})

	m.POST("/:id/remove", func(c *gin.Context) {
		me := c.GetString("user_id")

		err := s.Leave(c, c.Param("id"), me)

		reply(c, "child removed from your account", err)
	})
}
