package model

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RolePatient Role = "patient"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Doctor struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Specialization string    `json:"specialization"`
	CreatedAt      time.Time `json:"created_at"`
}

type AppointmentStatus string

const (
	StatusScheduled AppointmentStatus = "scheduled"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	ID        int               `json:"id"`
	PatientID int               `json:"patient_id"`
	DoctorID  int               `json:"doctor_id"`
	Time      time.Time         `json:"time"`
	Status    AppointmentStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`

	DoctorName     string `json:"doctor_name,omitempty"`
	Specialization string `json:"specialization,omitempty"`
}
