package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"virelia-backend/validators"
)

func (ctl *Controller) ForgotPassword(c *gin.Context) {
	var req validators.ForgotPassword

	if !bind(c, &req) {
		return
	}

	err := ctl.svc.ForgotPassword(
		c.Request.Context(),
		req.Email,
	)

	message := gin.H{
		"message": "if the account exists, a code was sent",
	}

	respond(c, http.StatusOK, message, err)
}

func (ctl *Controller) ResetPassword(c *gin.Context) {
	var req validators.ResetPassword

	if !bind(c, &req) {
		return
	}

	err := ctl.svc.ResetPassword(
		c.Request.Context(),
		req.Email,
		req.Code,
		req.Password,
	)

	message := gin.H{"message": "password reset"}

	respond(c, http.StatusOK, message, err)
}

func (ctl *Controller) ChangePassword(c *gin.Context) {
	var req validators.ChangePassword

	if !bind(c, &req) {
		return
	}

	err := ctl.svc.ChangePassword(
		c.Request.Context(),
		c.GetString("user_id"),
		req.CurrentPassword,
		req.Password,
	)

	message := gin.H{"message": "password changed"}

	respond(c, http.StatusOK, message, err)
}