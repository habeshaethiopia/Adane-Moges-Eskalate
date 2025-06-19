package handlers

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("youtubeurl", validateYouTubeURL)
	v.RegisterValidation("customemail", validateEmail)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return false
	}

	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return false
	}

	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()]`).MatchString(password)
	if !hasSpecial {
		return false
	}

	return true
}

func validateYouTubeURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	pattern := `^https?://(www\.)?(youtube\.com/watch\?v=|youtu\.be/)[\w-]{11}$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Convert to lowercase for consistent validation
	email = strings.ToLower(email)

	// Check minimum and maximum length
	if len(email) < 5 || len(email) > 254 { // minimum: a@b.c
		return false
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Validate local part
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// Check for invalid characters in local part
	localPartRegex := `^[a-z0-9.!#$%&'*+/=?^_\x60{|}~-]+$`
	if matched, _ := regexp.MatchString(localPartRegex, localPart); !matched {
		return false
	}

	// Check for consecutive special characters
	if strings.Contains(localPart, "..") ||
		strings.HasPrefix(localPart, ".") ||
		strings.HasSuffix(localPart, ".") {
		return false
	}

	// Validate domain part
	if len(domainPart) < 3 || len(domainPart) > 255 {
		return false
	}

	// Check domain format (including subdomains)
	domainRegex := `^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,}$`
	if matched, _ := regexp.MatchString(domainRegex, domainPart); !matched {
		return false
	}

	// Additional domain validations
	if strings.Contains(domainPart, "..") ||
		strings.HasPrefix(domainPart, "-") ||
		strings.HasSuffix(domainPart, "-") ||
		strings.Contains(domainPart, ".-") ||
		strings.Contains(domainPart, "-.") {
		return false
	}

	return true
}
