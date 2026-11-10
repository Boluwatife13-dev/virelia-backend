package main

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"virelia/internal/platform/config"
	"virelia/internal/platform/database"
)

type Admin struct {
	Email    string
	Phone    string
	First    string
	Last     string
	Password string
}

const insertAdmin = `
	WITH u AS (
		INSERT INTO users (
			email,
			phone,
			password_hash,
			email_verified_at
		)
		VALUES ($1, $2, $3, now())
		RETURNING id
	)
	INSERT INTO staff (user_id, role, first_name, last_name)
	SELECT id, 'platform_admin', $4, $5
	FROM u`

func save(a Admin) error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(a.Password),
		12,
	)
	if err != nil {
		return err
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		context.Background(),
		insertAdmin,
		a.Email,
		a.Phone,
		string(hash),
		a.First,
		a.Last,
	)

	return err
}
