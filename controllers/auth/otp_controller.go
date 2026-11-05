package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"virelia-backend/models"
	"virelia-backend/validators"
)

func (ctl *Controller) SendOTP(c *gin.Context) {
	var req validators.SendOTP

	if !bind(c, &req) {
		return
	}

	err := ctl.svc.SendVerification(
		c.Request.Context(),
		c.GetString("user_id"),
		otpPurpose(req.Channel),
	)

	message := gin.H{"message": "code sent"}

	respond(c, http.StatusOK, message, err)
}

func (ctl *Controller) VerifyOTP(c *gin.Context) {
	var req validators.VerifyOTP

	if !bind(c, &req) {
		return
	}

	err := ctl.svc.VerifyContact(
		c.Request.Context(),
		c.GetString("user_id"),
		otpPurpose(req.Channel),
		req.Code,
	)

	message := gin.H{"message": "verified"}

	respond(c, http.StatusOK, message, err)
}

func otpPurpose(channel string) models.OTPPurpose {
	if channel == "email" {
		return models.OTPVerifyEmail
	}

	return models.OTPVerifyPhone
}