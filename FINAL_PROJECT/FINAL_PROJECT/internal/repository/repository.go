package repository

import (
	"clinic-cli/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string, role model.Role) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
}

type DoctorRepository interface {
	Create(ctx context.Context, name, specialization string) (*model.Doctor, error)
	GetAll(ctx context.Context) ([]model.Doctor, error)
	GetByID(ctx context.Context, id int) (*model.Doctor, error)
}

type AppointmentRepository interface {
	Create(ctx context.Context, patientID, doctorID int, timeStr string) (*model.Appointment, error)
	GetByPatientID(ctx context.Context, patientID int) ([]model.Appointment, error)
	GetByDoctorID(ctx context.Context, doctorID int) ([]model.Appointment, error)
	Cancel(ctx context.Context, id int) error
}

type Registry struct {
	User        UserRepository
	Doctor      DoctorRepository
	Appointment AppointmentRepository
}
