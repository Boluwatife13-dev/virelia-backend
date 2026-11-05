package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"virelia-backend/config"
	authctl "virelia-backend/controllers/auth"
	"virelia-backend/middleware"
	authrepo "virelia-backend/repositories/auth"
	authsvc "virelia-backend/services/auth"
	"virelia-backend/utils"
)

func SetupRoutes(
	router *gin.Engine,
	db *mongo.Database,
	cfg config.Config,
) {
	router.GET("/health", health)

	service := authsvc.NewService(
		authrepo.NewUserRepository(db),
		authrepo.NewRefreshTokenRepository(db),
		authrepo.NewOTPRepository(db),
		utils.Mailer{APIKey: cfg.BrevoAPIKey, From: cfg.EmailFrom},
		cfg.JWTSecret,
	)

	controller := authctl.NewController(service)
	requireLogin := middleware.Authenticate(cfg.JWTSecret)

	auth := router.Group("/api/v1/auth")

	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
	auth.POST("/refresh", controller.Refresh)
	auth.POST("/logout", controller.Logout)
	auth.GET("/me", requireLogin, controller.Me)

	auth.POST("/otp/send", requireLogin, controller.SendOTP)
	auth.POST("/otp/verify", requireLogin, controller.VerifyOTP)

	auth.POST("/password/forgot", controller.ForgotPassword)
	auth.POST("/password/reset", controller.ResetPassword)
	auth.POST(
		"/password/change",
		requireLogin,
		controller.ChangePassword,
	)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app":     "Virelia",
		"status":  "ok",
		"message": "Health check successful",
	})
}