package records

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/audit"
	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrNoPatient:  404,
	ErrNotFound:   404,
	ErrNotAllowed: 403,
	ErrOldVersion: 409,
	ErrTooBig:     400,
}

func doctorFrom(c *gin.Context) Doctor {
	return Doctor{
		ID:         c.GetString("user_id"),
		FacilityID: c.GetString("facility_id"),
	}
}

func logPatient(c *gin.Context,
	a *audit.Service,
	action,
	patientID string) {
	e := audit.From(c, action)
	e.PatientID = patientID
	e.TargetID = c.Param("id")

	a.Log(c, e)
}

func PatientRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	m := r.Group("/me/records", guards...)

	m.GET("", func(c *gin.Context) {
		list, err := s.MyRecord(c, c.GetString("user_id"))

		httpx.Respond(c, 200, gin.H{"records": list}, err, statusFor)
	})
}
