package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"virelia/internal/appointments"
	"virelia/internal/audit"
	"virelia/internal/auth"
	"virelia/internal/directory"
	"virelia/internal/facilities"
	"virelia/internal/guardians"
	"virelia/internal/keys"
	"virelia/internal/patients"
	"virelia/internal/platform/config"
	mw "virelia/internal/platform/middleware"
	"virelia/internal/records"
	"virelia/internal/staff"
	"virelia/web"
)

func setupRoutes(
	router *gin.Engine,
	db *pgxpool.Pool,
	cfg config.Config,
) {
	s := buildServices(db, cfg)

	router.GET("/health", health(db))

	api := router.Group("/api/v1")

	api.Use(mw.RateLimit(60, time.Minute))

	api.Use(audit.Middleware(s.audit))

	login := mw.RequireLogin(cfg.JWTSecret)

	hospitalAdmin := mw.RequireRole(auth.RoleHospitalAdmin)
	platformAdmin := mw.RequireRole(auth.RolePlatformAdmin)
	frontDesk := mw.RequireRole(auth.RoleFrontDesk)
	patient := mw.RequireRole(auth.RolePatient)
	doctor := mw.RequireRole(auth.RoleDoctor)

	auth.Routes(api, s.auth, login)

	guardians.PatientRoutes(api, s.guardians, login, patient)

	staff.Routes(api, s.staff, login, hospitalAdmin)

	patients.Routes(api, s.patients, login, frontDesk)
	guardians.Routes(api, s.guardians, login, frontDesk)

	facilities.Routes(api, s.facilities, login, platformAdmin)
	staff.PlatformRoutes(api, s.staff, login, platformAdmin)

	keys.Routes(api, s.keys, login, platformAdmin)

	directory.Routes(api, s.directory, login, doctor)

	appointments.PublicRoutes(api, s.appointments)
	appointments.PatientRoutes(api, s.appointments, login, patient)
	appointments.DeskRoutes(api, s.appointments, login, frontDesk)
	appointments.DoctorRoutes(api, s.appointments, login, doctor)
	records.PatientRoutes(api, s.records, login, patient)
	records.DoctorRoutes(api, s.records, s.audit, login, doctor)
	audit.PatientRoutes(api, s.audit, login, patient)
	audit.HospitalRoutes(api, s.audit, login, hospitalAdmin)
	audit.PlatformRoutes(api, s.audit, login, platformAdmin)

	router.GET("/key-test", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", web.KeyTest)
	})
}

func health(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"

		if db.Ping(c) != nil {
			status = "down"
		}

		c.JSON(200, gin.H{"app": "Virelia", "database": status})
	}
}
