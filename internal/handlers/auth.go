package handlers

import (
	"net/http"

	"eskalate-movie-api/internal/config"
	"eskalate-movie-api/internal/models"
	"eskalate-movie-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,password"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RegisterAuthRoutes registers auth endpoints
func RegisterAuthRoutes(rg *gin.RouterGroup, authService services.AuthService, cfg *config.Config) {
	rg.POST("/signup", Signup(authService, cfg))
	rg.POST("/login", Login(authService, cfg))
	rg.POST("/refresh", RefreshToken(authService, cfg))
	rg.POST("/logout", Logout(authService))

	// Debug endpoint - REMOVE IN PRODUCTION
	rg.GET("/debug/:email", func(c *gin.Context) {
		email := c.Param("email")
		user, err := authService.GetUserByEmail(email)
		if err != nil {
			c.JSON(http.StatusNotFound, BaseResponse{
				Success: false,
				Message: "User not found",
				Errors:  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, BaseResponse{
			Success: true,
			Message: "User found",
			Object: gin.H{
				"id":            user.ID,
				"email":         user.Email,
				"username":      user.Username,
				"password_hash": user.Password,
				"password_len":  len(user.Password),
				"created_at":    user.CreatedAt,
				"updated_at":    user.UpdatedAt,
			},
		})
	})
}

// Signup godoc
// @Summary      Register a new user
// @Description  Register a new user with email, username, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        signupRequest body SignupRequest true "Signup request"
// @Success      201 {object} BaseResponse
// @Failure      400 {object} BaseResponse
// @Router       /api/auth/signup [post]
func Signup(authService services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignupRequest
		validate := validator.New()
		RegisterCustomValidators(validate)

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Invalid input format",
				Errors:  []string{err.Error()},
			})
			return
		}

		if err := validate.Struct(req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errs := make([]string, len(validationErrors))
				for i, err := range validationErrors {
					errs[i] = getValidationErrorMsg(err)
				}
				c.JSON(http.StatusBadRequest, BaseResponse{
					Success: false,
					Message: "Validation failed",
					Errors:  errs,
				})
				return
			}
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Validation failed",
				Errors:  []string{err.Error()},
			})
			return
		}

		user := &models.User{
			Email:    req.Email,
			Username: req.Username,
			Password: req.Password,
		}
		if err := authService.Signup(user); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{Success: false, Message: err.Error(), Errors: []string{err.Error()}})
			return
		}
		user.Password = ""
		c.JSON(http.StatusCreated, BaseResponse{Success: true, Message: "Signup successful", Object: user})
	}
}

// Login godoc
// @Summary      Login a user
// @Description  Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest body LoginRequest true "Login request"
// @Success      200 {object} BaseResponse
// @Failure      401 {object} BaseResponse
// @Router       /api/auth/login [post]
func Login(authService services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		validate := validator.New()
		RegisterCustomValidators(validate)

		if err := c.ShouldBindJSON(&req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errs := make([]string, len(validationErrors))
				for i, err := range validationErrors {
					errs[i] = getValidationErrorMsg(err)
				}
				c.JSON(http.StatusBadRequest, BaseResponse{
					Success: false,
					Message: "Validation failed",
					Errors:  errs,
				})
				return
			}
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Invalid input format",
				Errors:  []string{err.Error()},
			})
			return
		}

		if err := validate.Struct(req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errs := make([]string, len(validationErrors))
				for i, err := range validationErrors {
					errs[i] = getValidationErrorMsg(err)
				}
				c.JSON(http.StatusBadRequest, BaseResponse{
					Success: false,
					Message: "Validation failed",
					Errors:  errs,
				})
				return
			}
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Validation failed",
				Errors:  []string{err.Error()},
			})
			return
		}

		accessToken, refreshToken, err := authService.LoginWithRefresh(req.Email, req.Password, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, BaseResponse{
				Success: false,
				Message: "Login failed",
				Errors:  []string{"Invalid email or password"},
			})
			return
		}

		c.JSON(http.StatusOK, BaseResponse{
			Success: true,
			Message: "Login successful",
			Object:  TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken},
		})
	}
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        refreshRequest body RefreshRequest true "Refresh token request"
// @Success      200 {object} BaseResponse
// @Failure      401 {object} BaseResponse
// @Router       /api/auth/refresh [post]
func RefreshToken(authService services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Invalid input format",
				Errors:  []string{err.Error()},
			})
			return
		}
		accessToken, err := authService.RefreshAccessToken(req.RefreshToken, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, BaseResponse{
				Success: false,
				Message: "Invalid or expired refresh token",
				Errors:  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, BaseResponse{
			Success: true,
			Message: "Token refreshed successfully",
			Object:  gin.H{"accessToken": accessToken},
		})
	}
}

// Logout godoc
// @Summary      Logout (revoke refresh token)
// @Description  Revoke a refresh token (logout)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        logoutRequest body LogoutRequest true "Logout request"
// @Success      200 {object} BaseResponse
// @Failure      400 {object} BaseResponse
// @Router       /api/auth/logout [post]
func Logout(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LogoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Invalid input format",
				Errors:  []string{err.Error()},
			})
			return
		}
		if err := authService.RevokeRefreshToken(req.RefreshToken); err != nil {
			c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Message: "Failed to revoke refresh token",
				Errors:  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, BaseResponse{
			Success: true,
			Message: "Logged out successfully",
		})
	}
}

// getValidationErrorMsg returns a user-friendly error message for validation errors
func getValidationErrorMsg(fieldError validator.FieldError) string {
	switch fieldError.Field() {
	case "Email":
		switch fieldError.Tag() {
		case "required":
			return "Email address is required"
		case "customemail":
			return "Please enter a valid email address"
		default:
			return "Invalid email format"
		}
	case "Username":
		switch fieldError.Tag() {
		case "required":
			return "Username is required"
		case "alphanum":
			return "Username must contain only letters and numbers"
		case "min":
			return "Username must be at least 3 characters long"
		case "max":
			return "Username cannot be longer than 20 characters"
		default:
			return "Invalid username format"
		}
	case "Password":
		switch fieldError.Tag() {
		case "required":
			return "Password is required"
		case "min":
			return "Password must be at least 8 characters long"
		case "containsany":
			if fieldError.Param() == "!@#$%^&*()" {
				return "Password must contain at least one special character (!@#$%^&*())"
			} else if fieldError.Param() == "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
				return "Password must contain at least one uppercase letter"
			} else if fieldError.Param() == "abcdefghijklmnopqrstuvwxyz" {
				return "Password must contain at least one lowercase letter"
			}
		}
		return "Invalid password format"
	}
	return fieldError.Error()
}
