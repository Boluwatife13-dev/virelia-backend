package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"virelia/internal/auth"
	"virelia/internal/facilities"
	"virelia/internal/guardians"
	"virelia/internal/keys"
	"virelia/internal/patients"
	"virelia/internal/platform/config"
	mw "virelia/internal/platform/middleware"
	"virelia/internal/staff"
	"virelia/web"
)

// setupRoutes connects every URL, grouped by
// who is allowed to use it.
func setupRoutes(
	router *gin.Engine,
	db *pgxpool.Pool,
	cfg config.Config,
) {
	s := buildServices(db, cfg)

	router.GET("/health", health(db))

	api := router.Group("/api/v1")

	api.Use(mw.RateLimit(60, time.Minute))

	login := mw.RequireLogin(cfg.JWTSecret)

	hospitalAdmin := mw.RequireRole(auth.RoleHospitalAdmin)
	platformAdmin := mw.RequireRole(auth.RolePlatformAdmin)
	frontDesk := mw.RequireRole(auth.RoleFrontDesk)
	patient := mw.RequireRole(auth.RolePatient)

	// Sign up, log in, passwords.
	auth.Routes(api, s.auth, login)

	// Patients manage their children.
	guardians.PatientRoutes(api, s.guardians, login, patient)

	// Hospital admins manage their own staff.
	// Accepting an invite is open to everyone.
	staff.Routes(api, s.staff, login, hospitalAdmin)

	// Front desk registers and finds walk-in patients.
	patients.Routes(api, s.patients, login, frontDesk)
	guardians.Routes(api, s.guardians, login, frontDesk)

	// Platform admins review hospitals and reset
	// hospital admins' 2FA. Applying is open to everyone.
	facilities.Routes(api, s.facilities, login, platformAdmin)
	staff.PlatformRoutes(api, s.staff, login, platformAdmin)

	// Platform admins add and use security keys.
	keys.Routes(api, s.keys, login, platformAdmin)

	// A page to try security keys in the browser.
	router.GET("/key-test", func(c *gin.Context) { c.Data(200, "text/html; charset=utf-8", web.KeyTest) })
}

// health shows whether the app and database are up.
func health(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"

		if db.Ping(c) != nil {
			status = "down"
		}

		c.JSON(200, gin.H{"app": "Virelia", "database": status})
	}
}
