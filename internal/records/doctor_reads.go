package records

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

func (h doctorHandlers) patient(c *gin.Context) {
	d := doctorFrom(c)
	list, patient, err := h.s.AtMyHospital(c, d, c.Param("id"))
	h.done(c, "record.view", patient, err)

	data := gin.H{"records": list}
	httpx.Respond(c, 200, data, err, statusFor)
}

func (h doctorHandlers) history(c *gin.Context) {
	d := doctorFrom(c)
	list, patient, err := h.s.History(c, d, c.Param("id"))
	h.done(c, "record.history", patient, err)

	data := gin.H{"versions": list}
	httpx.Respond(c, 200, data, err, statusFor)
}

func (h doctorHandlers) attached(c *gin.Context) {
	d := doctorFrom(c)
	list, err := h.s.Attached(c, d, c.Param("id"))

	patient := ""
	if len(list) > 0 {
		patient = list[0].PatientID
	}
	h.done(c, "record.view_attached", patient, err)

	data := gin.H{"records": list}
	httpx.Respond(c, 200, data, err, statusFor)
}
