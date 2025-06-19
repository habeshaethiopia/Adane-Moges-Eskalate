package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "eskalate-movie-api/docs"
	"eskalate-movie-api/internal/config"
	"eskalate-movie-api/internal/handlers"
	"eskalate-movie-api/internal/routes"
)

// @title Eskalate Movie API
// @description REST API for a personal movie collection
// @version 1.0
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
func main() {
	cfg := config.LoadConfig()

	if cfg.Port == "" {
		cfg.Port = ":8080"
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	// Register custom validators globally for Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		handlers.RegisterCustomValidators(v)
	}

	// Enable uuid-ossp extension
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Auto-migrate models
	db.AutoMigrate()

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	routes.RegisterRoutes(r, db, cfg)

	log.Printf("Server running on %s", cfg.Port)
	r.Run(cfg.Port)
}
