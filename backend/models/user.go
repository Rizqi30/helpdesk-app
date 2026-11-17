package models

import "time"

const (
	RoleRequester = "requester"
	RoleAgent     = "agent"
	RoleAdmin     = "admin"
)

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	Email        string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null"`
	Role         string `json:"role" gorm:"type:varchar(50);not null;default:'requester'"`
	CreatedAt    time.Time `json:"created_at"`
}