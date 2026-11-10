package main

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"virelia/internal/auth"
	"virelia/internal/facilities"
	"virelia/internal/guardians"
	"virelia/internal/keys"
	"virelia/internal/notify"
	"virelia/internal/patients"
	"virelia/internal/platform/config"
	"virelia/internal/platform/crypt"
	"virelia/internal/staff"
)

// services holds one of each module, ready to use.
type services struct {
	auth       *auth.Service
	staff      *staff.Service
	facilities *facilities.Service
	patients   *patients.Service
	guardians  *guardians.Service
	keys       *keys.Service
}

func buildServices(
	db *pgxpool.Pool,
	cfg config.Config,
) services {
	mailer := notify.Brevo{
		APIKey: cfg.BrevoAPIKey,
		From:   cfg.EmailFrom,
	}

	box, err := crypt.New(cfg.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewService(
		auth.NewRepository(db),
		mailer,
		cfg.JWTSecret,
		box,
	)

	staffService := staff.NewService(
		staff.NewRepository(db),
		mailer,
		authService,
		cfg.FrontendURLs[0],
	)

	facilityService := facilities.NewService(
		facilities.NewRepository(db),
		mailer,
		staffService,
	)

	patientService := patients.NewService(
		patients.NewRepository(db),
		authService,
	)

	guardianService := guardians.NewService(
		guardians.NewRepository(db),
		authService,
	)

	// The browser page's address must be listed,
	// or the key login is refused.
	origins := append(
		cfg.FrontendURLs,
		"http://localhost:"+cfg.Port,
	)

	keyService, err := keys.NewService(
		keys.NewRepository(db),
		authService,
		cfg.WebAuthnRPID,
		origins,
	)
	if err != nil {
		log.Fatal(err)
	}

	return services{
		auth:       authService,
		staff:      staffService,
		facilities: facilityService,
		patients:   patientService,
		guardians:  guardianService,
		keys:       keyService,
	}
}
