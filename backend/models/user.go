package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User matches existing PostgreSQL table `users`
// columns: id (uuid PK), full_name, email, password_hash, google_id, created_at, updated_at
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName     string    `json:"full_name" gorm:"column:full_name;not null"`
	Email        string    `json:"email" gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	GoogleID     *string   `json:"google_id,omitempty" gorm:"column:google_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string { return "users" }

// BeforeCreate ensures UUID is set when inserting via GORM
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
