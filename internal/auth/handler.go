package auth

import "github.com/gin-gonic/gin"

type ctx = *gin.Context

const sent = "if the email exists, a code was sent"

func Routes(
	r *gin.RouterGroup,
	s *Service,
	requireLogin gin.HandlerFunc,
) {
	p := r.Group("/auth/patients")

	p.POST("/register", on(201,
		func(c ctx, req RegisterRequest) (any, error) {
			err := s.Register(c, req)

			return msg("account created, check your email"), err
		},
	))

	p.POST("/verify-email", on(200,
		func(c ctx, req VerifyRequest) (any, error) {
			err := s.VerifyEmail(c, req)

			return msg("email verified"), err
		},
	))

	p.POST("/resend-code", on(200,
		func(c ctx, req EmailRequest) (any, error) {
			err := s.ResendCode(c, req.Email)

			return msg(sent), err
		},
	))

	p.POST("/login", on(200,
		func(c ctx, req LoginRequest) (any, error) {
			agent := c.Request.UserAgent()

			return s.Login(c, req, agent, c.ClientIP())
		},
	))

	st := r.Group("/auth/staff")

	st.POST("/login", on(200,
		func(c ctx, req LoginRequest) (any, error) {
			return s.StaffLogin(c, req)
		},
	))

	st.POST("/verify-2fa", on(200,
		func(c ctx, req TwoFactorRequest) (any, error) {
			agent := c.Request.UserAgent()

			return s.VerifyTwoFactor(c, req, agent, c.ClientIP())
		},
	))

	a := r.Group("/auth")

	a.POST("/refresh", on(200,
		func(c ctx, req TokenRequest) (any, error) {
			agent := c.Request.UserAgent()

			return s.Refresh(c, req.RefreshToken, agent, c.ClientIP())
		},
	))

	a.POST("/logout", on(200,
		func(c ctx, req TokenRequest) (any, error) {
			err := s.Logout(c, req.RefreshToken)

			return msg("logged out"), err
		},
	))

	a.POST("/password/forgot", on(200,
		func(c ctx, req EmailRequest) (any, error) {
			err := s.ForgotPassword(c, req.Email)

			return msg(sent), err
		},
	))

	a.POST("/password/reset", on(200,
		func(c ctx, req ResetRequest) (any, error) {
			err := s.ResetPassword(c, req)

			return msg("password reset"), err
		},
	))

	a.POST("/password/change", requireLogin, on(200,
		func(c ctx, req ChangePasswordRequest) (any, error) {
			id := c.GetString("user_id")

			err := s.ChangePassword(c, id, req)

			return msg("password changed"), err
		},
	))

	a.GET("/me", requireLogin, func(c ctx) {
		id := c.GetString("user_id")

		user, patient, err := s.Me(c, id)

		data := gin.H{"user": user, "patient": patient}

		respond(c, 200, data, err)
	})
}
