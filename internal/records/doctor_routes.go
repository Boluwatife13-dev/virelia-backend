package records

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/audit"
)

func DoctorRoutes(
	r *gin.RouterGroup,
	s *Service,
	a *audit.Service,
	guards ...gin.HandlerFunc,
) {
	h := doctorHandlers{s: s, a: a}
	d := r.Group("/doctor", guards...)

	d.POST("/records", h.add)
	d.POST("/records/:id/correct", h.correct)

	d.GET("/patients/:id/records", h.patient)
	d.GET("/records/:id/history", h.history)
	d.GET("/appointments/:id/records", h.attached)
}
