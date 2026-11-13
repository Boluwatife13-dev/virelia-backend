package audit

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

func PatientRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	p := r.Group("/me", guards...)

	p.GET("/access-log", func(c *gin.Context) {
		list, err := s.ForPatientUser(c, c.GetString("user_id"))

		send(c, list, err)
	})
}

func HospitalRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	h := r.Group("/audit", guards...)

	h.GET("", func(c *gin.Context) {
		f := filterFrom(c)
		f.FacilityID = c.GetString("facility_id")

		list, err := s.Find(c, f)

		send(c, list, err)
	})
}

func PlatformRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	p := r.Group("/platform/audit", guards...)

	p.GET("", func(c *gin.Context) {
		f := filterFrom(c)
		f.FacilityID = c.Query("facility_id")

		list, err := s.Find(c, f)

		send(c, list, err)
	})
}

func filterFrom(c *gin.Context) Filter {
	return Filter{
		PatientID: c.Query("patient_id"),
		ActorID:   c.Query("actor_id"),
	}
}

func send(c *gin.Context, list []Record, err error) {
	data := gin.H{"entries": list}

	httpx.Respond(c, 200, data, err, nil)
}
