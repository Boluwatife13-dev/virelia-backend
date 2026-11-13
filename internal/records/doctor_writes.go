package records

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/audit"
	"virelia/internal/platform/httpx"
)

type doctorHandlers struct {
	s *Service
	a *audit.Service
}

func (h doctorHandlers) done(
	c *gin.Context,
	action, patient string,
	err error,
) {
	if err == nil && patient != "" {
		logPatient(c, h.a, action, patient)
	}
}

func (h doctorHandlers) add(c *gin.Context) {
	var req AddRequest
	if !httpx.Bind(c, &req) {
		return
	}

	id, patient, err := h.s.Add(c, doctorFrom(c), req)
	h.done(c, "record.add", patient, err)

	data := gin.H{"message": "record saved", "id": id}
	httpx.Respond(c, 201, data, err, statusFor)
}

func (h doctorHandlers) correct(c *gin.Context) {
	var req CorrectRequest
	if !httpx.Bind(c, &req) {
		return
	}

	d := doctorFrom(c)
	id, patient, err := h.s.Correct(c, d, c.Param("id"), req)
	h.done(c, "record.correct", patient, err)

	data := gin.H{"message": "correction saved", "id": id}
	httpx.Respond(c, 201, data, err, statusFor)
}
