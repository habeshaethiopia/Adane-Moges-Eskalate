package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title" validate:"required,min=1,max=40"`
	Description string    `gorm:"not null" json:"description" validate:"required,min=10,max=1000"`
	Poster      string    `gorm:"not null" json:"poster"`
	Trailer     string    `gorm:"not null" json:"trailer" validate:"required,youtubeurl"`
	Actors      []string  `gorm:"type:text[]" json:"actors" validate:"required,min=1,dive,required"`
	Genres      []string  `gorm:"type:text[]" json:"genres" validate:"required,min=1,dive,required"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
