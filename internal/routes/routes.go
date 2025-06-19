package routes

import (
	"eskalate-movie-api/internal/config"
	"eskalate-movie-api/internal/handlers"
	"eskalate-movie-api/internal/repository"
	"eskalate-movie-api/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, db)
	handlers.RegisterAuthRoutes(r.Group("/api/auth"), authService, cfg)

	movieRepo := repository.NewMovieRepository(db)
	movieService := services.NewMovieService(movieRepo)
	handlers.RegisterMovieRoutes(r.Group("/api/movies"), movieService, cfg)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
