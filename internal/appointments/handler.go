package appointments

import (
	"github.com/gin-gonic/gin"

	"virelia/internal/platform/httpx"
)

var statusFor = map[error]int{
	ErrNoDoctor:      404,
	ErrNoPatient:     404,
	ErrNotFound:      404,
	ErrTaken:         409,
	ErrBadDate:       400,
	ErrBadTime:       400,
	ErrOtherHospital: 403,
}

func PublicRoutes(r *gin.RouterGroup, s *Service) {
	r.GET("/doctors/:id/slots", func(c *gin.Context) {
		slots, err := s.Slots(c, c.Param("id"), c.Query("date"))

		httpx.Respond(c, 200, gin.H{"slots": slots}, err, statusFor)
	})
}

func PatientRoutes(
	r *gin.RouterGroup,
	s *Service,
	guards ...gin.HandlerFunc,
) {
	m := r.Group("/me/appointments", guards...)

	m.GET("", func(c *gin.Context) {
		list, err := s.MyAppointments(c, c.GetString("user_id"))

		send(c, list, err)
	})

	m.POST("", func(c *gin.Context) {
		var req BookRequest

		if !httpx.Bind(c, &req) {
			return
		}

		id, err := s.BookForMe(c, c.GetString("user_id"), req)

		booked(c, id, err)
	})

	m.POST("/:id/cancel", func(c *gin.Context) {
		var req CancelRequest

		if !httpx.Bind(c, &req) {
			return
		}

		me := c.GetString("user_id")
		err := s.CancelMine(c, me, c.Param("id"), req.Reason)

		reply(c, "appointment cancelled", err)
	})
}

func send(c *gin.Context, list []Appointment, err error) {
	httpx.Respond(c, 200, gin.H{"appointments": list}, err, statusFor)
}

func booked(c *gin.Context, id string, err error) {
	data := gin.H{"message": "appointment booked", "id": id}

	httpx.Respond(c, 201, data, err, statusFor)
}

func reply(c *gin.Context, msg string, err error) {
	httpx.Respond(c, 200, gin.H{"message": msg}, err, statusFor)
}
