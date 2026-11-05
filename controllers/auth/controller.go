package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	repo "virelia-backend/repositories/auth"
	service "virelia-backend/services/auth"
	"virelia-backend/validators"
)

var knownErrors = map[error]int{
	repo.ErrUserExists:         http.StatusConflict,
	service.ErrInvalidLogin:    http.StatusUnauthorized,
	service.ErrInvalidToken:    http.StatusUnauthorized,
	service.ErrAccountInactive: http.StatusForbidden,
	service.ErrInvalidOTP:      http.StatusBadRequest,
	service.ErrWrongPassword:   http.StatusBadRequest,
	service.ErrSMSUnavailable:  http.StatusServiceUnavailable,
}

type Controller struct {
	svc *service.Service
}

func NewController(svc *service.Service) *Controller {
	return &Controller{svc}
}

func (ctl *Controller) Register(c *gin.Context) {
	run := func(req validators.Register) (any, error) {
		user, err := ctl.svc.RegisterPatient(
			c.Request.Context(), service.RegisterInput(req),
		)
		return gin.H{"user": user}, err
	}

	handle(c, http.StatusCreated, run)
}

func (ctl *Controller) Login(c *gin.Context) {
	handle(c, http.StatusOK, func(req validators.Login) (any, error) {
		return ctl.svc.Login(
			c.Request.Context(), req.Email, req.Password,
			c.Request.UserAgent(), c.ClientIP(),
		)
	})
}

func (ctl *Controller) Refresh(c *gin.Context) {
	handle(c, http.StatusOK, func(req validators.Token) (any, error) {
		return ctl.svc.Refresh(
			c.Request.Context(), req.RefreshToken,
			c.Request.UserAgent(), c.ClientIP(),
		)
	})
}

func (ctl *Controller) Logout(c *gin.Context) {
	handle(c, http.StatusOK, func(req validators.Token) (any, error) {
		err := ctl.svc.Logout(c.Request.Context(), req.RefreshToken)
		return gin.H{"message": "logged out"}, err
	})
}

func (ctl *Controller) Me(c *gin.Context) {
	user, err := ctl.svc.GetUser(
		c.Request.Context(), c.GetString("user_id"),
	)

	respond(c, http.StatusOK, gin.H{"user": user}, err)
}

func handle[T any](
	c *gin.Context, status int, run func(T) (any, error),
) {
	var req T

	if !bind(c, &req) {
		return
	}

	data, err := run(req)

	respond(c, status, data, err)
}

func bind(c *gin.Context, req any) bool {
	err := c.ShouldBindJSON(req)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	return err == nil
}

func respond(c *gin.Context, status int, data any, err error) {
	if err == nil {
		c.JSON(status, data)
		return
	}

	for known, code := range knownErrors {
		if errors.Is(err, known) {
			c.JSON(code, gin.H{"error": known.Error()})
			return
		}
	}

	c.JSON(500, gin.H{"error": "something went wrong"})
}