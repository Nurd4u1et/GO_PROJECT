package model

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RolePatient Role = "patient"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never return password hash in JSON
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
