package models

import "time"

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

type Doctor struct {
	ID             int
	Name           string
	Specialization string
}
type Appointment struct {
	ID        int
	UserID    int
	DoctorID  int
	DateTime  string
	CreatedAt time.Time
}
