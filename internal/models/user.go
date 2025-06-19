package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username string    `gorm:"unique;not null" json:"username" validate:"required,alphanum,min=3,max=20"`
	Email    string    `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,password"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
