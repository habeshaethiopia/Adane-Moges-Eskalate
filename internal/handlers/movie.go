package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"eskalate-movie-api/internal/config"
	"eskalate-movie-api/internal/models"
	"eskalate-movie-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MovieRequest struct {
	Title       string   `form:"title" binding:"required,min=1,max=40"`
	Description string   `form:"description" binding:"required,min=10,max=1000"`
	Genres      []string `form:"genres" binding:"required"`
	Actors      []string `form:"actors" binding:"required"`
	Trailer     string   `form:"trailerUrl" binding:"required,youtubeurl"`
}

// RegisterMovieRoutes registers movie endpoints
func RegisterMovieRoutes(rg *gin.RouterGroup, movieService services.MovieService, cfg *config.Config) {
	rg.POST("/", CreateMovie(movieService, cfg))
	rg.PUT("/:id", UpdateMovie(movieService, cfg))
	rg.GET("/", GetMovies(movieService, cfg))
	rg.GET("/search", SearchMovies(movieService, cfg))
	rg.GET("/:id", MovieDetails(movieService, cfg))
	rg.DELETE("/:id", DeleteMovie(movieService, cfg))
}

// CreateMovie godoc
// @Summary      Create a new movie
// @Description  Create a new movie (auth required)
// @Tags         movies
// @Accept       multipart/form-data
// @Produce      json
// @Param        title formData string true "Title"
// @Param        description formData string true "Description"
// @Param        genres formData []string true "Genres"
// @Param        actors formData []string true "Actors"
// @Param        trailerUrl formData string true "Trailer URL"
// @Param        poster formData file true "Poster"
// @Success      201 {object} BaseResponse
// @Failure      400 {object} BaseResponse
// @Failure      401 {object} BaseResponse
// @Security     BearerAuth
// @Router       /api/movies [post]
func CreateMovie(movieService services.MovieService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MovieRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Invalid input", Errors: []string{err.Error()}})
			return
		}
		validate := validator.New()
		RegisterCustomValidators(validate)
		if err := validate.Struct(req); err != nil {
			errs := []string{}
			for _, e := range err.(validator.ValidationErrors) {
				errs = append(errs, e.Error())
			}
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Validation failed", Errors: errs})
			return
		}
		posterFile, posterHeader, err := c.Request.FormFile("poster")
		if err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Poster image is required", Errors: []string{err.Error()}})
			return
		}
		defer posterFile.Close()
		posterURL, err := UploadPosterToCloudinary(posterFile, posterHeader, cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, BaseResponse{Success: false, Message: "Failed to upload poster", Errors: []string{err.Error()}})
			return
		}
		userID, _ := c.Get("userID")
		uuidUser, _ := uuid.Parse(userID.(string))
		movie := &models.Movie{
			Title:       req.Title,
			Description: req.Description,
			Genres:      req.Genres,
			Actors:      req.Actors,
			Trailer:     req.Trailer,
			Poster:      posterURL,
			UserID:      uuidUser,
		}
		if err := movieService.Create(movie); err != nil {
			c.JSON(http.StatusInternalServerError, BaseResponse{Success: false, Message: "Failed to create movie", Errors: []string{err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, BaseResponse{Success: true, Message: "Movie created", Object: movie})
	}
}

// UpdateMovie godoc
// @Summary      Update a movie
// @Description  Update a movie (auth required, must own movie)
// @Tags         movies
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Movie ID"
// @Param        title formData string true "Title"
// @Param        description formData string true "Description"
// @Param        genres formData []string true "Genres"
// @Param        actors formData []string true "Actors"
// @Param        trailerUrl formData string true "Trailer URL"
// @Param        poster formData file false "Poster"
// @Success      200 {object} BaseResponse
// @Failure      400 {object} BaseResponse
// @Failure      401 {object} BaseResponse
// @Failure      403 {object} BaseResponse
// @Failure      404 {object} BaseResponse
// @Security     BearerAuth
// @Router       /api/movies/{id} [put]
func UpdateMovie(movieService services.MovieService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		movieID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Invalid movie ID", Errors: []string{err.Error()}})
			return
		}
		userID, _ := c.Get("userID")
		uuidUser, _ := uuid.Parse(userID.(string))
		var req MovieRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Invalid input", Errors: []string{err.Error()}})
			return
		}
		validate := validator.New()
		RegisterCustomValidators(validate)
		if err := validate.Struct(req); err != nil {
			errs := []string{}
			for _, e := range err.(validator.ValidationErrors) {
				errs = append(errs, e.Error())
			}
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Validation failed", Errors: errs})
			return
		}
		movie, err := movieService.GetByID(movieID)
		if err != nil {
			c.JSON(http.StatusNotFound, BaseResponse{Success: false, Message: "Movie not found", Errors: []string{"Movie not found"}})
			return
		}
		if movie.UserID != uuidUser {
			c.JSON(http.StatusForbidden, BaseResponse{Success: false, Message: "Forbidden", Errors: []string{"You do not own this movie"}})
			return
		}
		posterFile, posterHeader, err := c.Request.FormFile("poster")
		if err == nil {
			defer posterFile.Close()
			posterURL, err := UploadPosterToCloudinary(posterFile, posterHeader, cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, BaseResponse{Success: false, Message: "Failed to upload poster", Errors: []string{err.Error()}})
				return
			}
			movie.Poster = posterURL
		}
		movie.Title = req.Title
		movie.Description = req.Description
		movie.Genres = req.Genres
		movie.Actors = req.Actors
		movie.Trailer = req.Trailer
		if err := movieService.Update(movie, uuidUser); err != nil {
			c.JSON(http.StatusInternalServerError, BaseResponse{Success: false, Message: "Failed to update movie", Errors: []string{err.Error()}})
			return
		}
		c.JSON(http.StatusOK, BaseResponse{Success: true, Message: "Movie updated", Object: movie})
	}
}

// GetMovies godoc
// @Summary      Get all movies
// @Description  Get a paginated list of movies
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        pageNumber query int false "Page number"
// @Param        pageSize query int false "Page size"
// @Success      200 {object} PaginatedResponse
// @Router       /api/movies [get]
func GetMovies(movieService services.MovieService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNumber := 1
		pageSize := 10
		if pn := c.Query("pageNumber"); pn != "" {
			fmt.Sscanf(pn, "%d", &pageNumber)
		}
		if ps := c.Query("pageSize"); ps != "" {
			fmt.Sscanf(ps, "%d", &pageSize)
		}
		movies, total, err := movieService.GetAll(pageNumber, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, PaginatedResponse{Success: false, Message: "Failed to fetch movies", Errors: []string{err.Error()}})
			return
		}
		c.JSON(http.StatusOK, PaginatedResponse{
			Success:    true,
			Message:    "Movies fetched",
			Object:     movies,
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalSize:  total,
		})
	}
}

// SearchMovies godoc
// @Summary      Search movies by title
// @Description  Search movies by title (case-insensitive)
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        title query string false "Title substring"
// @Param        pageNumber query int false "Page number"
// @Param        pageSize query int false "Page size"
// @Success      200 {object} PaginatedResponse
// @Router       /api/movies/search [get]
func SearchMovies(movieService services.MovieService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title")
		pageNumber := 1
		pageSize := 10
		if pn := c.Query("pageNumber"); pn != "" {
			fmt.Sscanf(pn, "%d", &pageNumber)
		}
		if ps := c.Query("pageSize"); ps != "" {
			fmt.Sscanf(ps, "%d", &pageSize)
		}
		movies, total, err := movieService.Search(strings.ToLower(title), pageNumber, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, PaginatedResponse{Success: false, Message: "Failed to search movies", Errors: []string{err.Error()}})
			return
		}
		c.JSON(http.StatusOK, PaginatedResponse{
			Success:    true,
			Message:    "Movies fetched",
			Object:     movies,
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalSize:  total,
		})
	}
}

// MovieDetails godoc
// @Summary      Get movie details
// @Description  Get details for a single movie by ID
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id path string true "Movie ID"
// @Success      200 {object} BaseResponse
// @Failure      404 {object} BaseResponse
// @Router       /api/movies/{id} [get]
func MovieDetails(movieService services.MovieService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		movieID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Invalid movie ID", Errors: []string{err.Error()}})
			return
		}
		movie, err := movieService.GetByID(movieID)
		if err != nil {
			c.JSON(http.StatusNotFound, BaseResponse{Success: false, Message: "Movie not found", Errors: []string{"Movie not found"}})
			return
		}
		c.JSON(http.StatusOK, BaseResponse{Success: true, Message: "Movie found", Object: movie})
	}
}

// DeleteMovie godoc
// @Summary      Delete a movie
// @Description  Delete a movie (auth required, must own movie)
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id path string true "Movie ID"
// @Success      200 {object} BaseResponse
// @Failure      401 {object} BaseResponse
// @Failure      403 {object} BaseResponse
// @Failure      404 {object} BaseResponse
// @Security     BearerAuth
// @Router       /api/movies/{id} [delete]
func DeleteMovie(movieService services.MovieService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		movieID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: "Invalid movie ID", Errors: []string{err.Error()}})
			return
		}
		userID, _ := c.Get("userID")
		uuidUser, _ := uuid.Parse(userID.(string))
		if err := movieService.Delete(movieID, uuidUser); err != nil {
			if err.Error() == "forbidden" {
				c.JSON(http.StatusForbidden, BaseResponse{Success: false, Message: "Forbidden", Errors: []string{"You do not own this movie"}})
				return
			}
			c.JSON(http.StatusInternalServerError, BaseResponse{Success: false, Message: "Failed to delete movie", Errors: []string{err.Error()}})
			return
		}
		c.JSON(http.StatusOK, BaseResponse{Success: true, Message: "Movie deleted"})
	}
}
