package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"virelia/internal/platform/config"
	"virelia/internal/platform/database"
	mw "virelia/internal/platform/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	router := gin.Default()

	router.Use(mw.CORS(cfg.FrontendURLs))

	setupRoutes(router, db, cfg)

	log.Println("Virelia API running on port " + cfg.Port)

	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Server failed: ", err)
	}
}
