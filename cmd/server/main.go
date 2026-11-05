package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"virelia-backend/config"
	"virelia-backend/database"
	"virelia-backend/middleware"
	"virelia-backend/routes"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.ConnectDatabase(
		cfg.MongoURI,
		cfg.DatabaseName,
	)
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	if err := database.EnsureIndexes(db); err != nil {
		log.Fatal("Database index setup failed: ", err)
	}

	router := gin.Default()

	router.Use(middleware.CORS(cfg.FrontendURLs))
	routes.SetupRoutes(router, db, cfg)

	log.Println("Virelia server running on port " + cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Server failed: ", err)
	}
}