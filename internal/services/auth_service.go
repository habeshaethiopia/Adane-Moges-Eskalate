package services

import (
	"errors"
	"log"
	"time"

	"eskalate-movie-api/internal/models"
	"eskalate-movie-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Signup(user *models.User) error
	Login(email, password, jwtSecret string) (string, error)
	LoginWithRefresh(email, password, jwtSecret string) (string, string, error)
	RefreshAccessToken(refreshToken, jwtSecret string) (string, error)
	RevokeRefreshToken(refreshToken string) error
}

type authService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewAuthService(userRepo repository.UserRepository, db *gorm.DB) AuthService {
	return &authService{userRepo, db}
}

func (s *authService) Signup(user *models.User) error {
	if _, err := s.userRepo.FindByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}
	if _, err := s.userRepo.FindByUsername(user.Username); err == nil {
		return errors.New("username already exists")
	}

	// Hash the password with a cost of 10
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	log.Printf("Original password: %s", user.Password)
	log.Printf("Generated hash: %s", string(hash))

	user.Password = string(hash)
	return s.userRepo.Create(user)
}

func (s *authService) Login(email, password, jwtSecret string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (s *authService) LoginWithRefresh(email, password, jwtSecret string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("User not found with email: %s", email)
		return "", "", errors.New("invalid credentials")
	}

	log.Printf("Login attempt - Email: %s", email)
	log.Printf("Stored hash: %s", user.Password)
	log.Printf("Provided password: %s", password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return "", "", errors.New("invalid credentials")
	}

	log.Printf("Password comparison successful")

	// Access token (short-lived)
	accessClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}
	// Refresh token (longer-lived)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":    "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}
	// Store refresh token in DB
	refreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
	}
	err = s.db.Create(refreshTokenModel).Error
	if err != nil {
		return "", "", err
	}
	return accessTokenStr, refreshTokenStr, nil
}

func (s *authService) RefreshAccessToken(refreshToken, jwtSecret string) (string, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil || claims["type"] != "refresh" {
		return "", errors.New("invalid refresh token claims")
	}
	userID := claims["user_id"].(string)
	// Check if refresh token is in DB and not revoked/expired
	var dbToken models.RefreshToken
	err = s.db.Where("token = ? AND revoked = false AND expires_at > ?", refreshToken, time.Now()).First(&dbToken).Error
	if err != nil {
		return "", errors.New("refresh token not recognized or expired")
	}
	if dbToken.UserID.String() != userID {
		return "", errors.New("refresh token does not belong to user")
	}
	// Issue new access token
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return accessTokenStr, nil
}

func (s *authService) RevokeRefreshToken(refreshToken string) error {
	return s.db.Model(&models.RefreshToken{}).Where("token = ?", refreshToken).Update("revoked", true).Error
}
