package repository

import (
	"clinic-cli/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, email, passwordHash string, role model.Role) (*model.User, error) {
	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, created_at`
	user := &model.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}
	err := r.pool.QueryRow(ctx, query, email, passwordHash, role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`
	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

type PostgresDoctorRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresDoctorRepository(pool *pgxpool.Pool) *PostgresDoctorRepository {
	return &PostgresDoctorRepository{pool: pool}
}

func (r *PostgresDoctorRepository) Create(ctx context.Context, name, specialization string) (*model.Doctor, error) {
	query := `INSERT INTO doctors (name, specialization) VALUES ($1, $2) RETURNING id, created_at`
	doc := &model.Doctor{Name: name, Specialization: specialization}
	err := r.pool.QueryRow(ctx, query, name, specialization).Scan(&doc.ID, &doc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *PostgresDoctorRepository) GetAll(ctx context.Context) ([]model.Doctor, error) {
	query := `SELECT id, name, specialization, created_at FROM doctors`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []model.Doctor
	for rows.Next() {
		var d model.Doctor
		if err := rows.Scan(&d.ID, &d.Name, &d.Specialization, &d.CreatedAt); err != nil {
			return nil, err
		}
		doctors = append(doctors, d)
	}
	return doctors, nil
}

func (r *PostgresDoctorRepository) GetByID(ctx context.Context, id int) (*model.Doctor, error) {
	query := `SELECT id, name, specialization, created_at FROM doctors WHERE id = $1`
	doc := &model.Doctor{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&doc.ID, &doc.Name, &doc.Specialization, &doc.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return doc, nil
}

type PostgresAppointmentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAppointmentRepository(pool *pgxpool.Pool) *PostgresAppointmentRepository {
	return &PostgresAppointmentRepository{pool: pool}
}

func (r *PostgresAppointmentRepository) Create(ctx context.Context, patientID, doctorID int, timeStr string) (*model.Appointment, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04", timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}

	query := `INSERT INTO appointments (patient_id, doctor_id, time, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	app := &model.Appointment{
		PatientID: patientID,
		DoctorID:  doctorID,
		Time:      parsedTime,
		Status:    model.StatusScheduled,
	}
	err = r.pool.QueryRow(ctx, query, patientID, doctorID, parsedTime, model.StatusScheduled).Scan(&app.ID, &app.CreatedAt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (r *PostgresAppointmentRepository) GetByPatientID(ctx context.Context, patientID int) ([]model.Appointment, error) {
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.time, a.status, a.created_at, d.name, d.specialization
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		WHERE a.patient_id = $1
		ORDER BY a.time ASC
	`
	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(&a.ID, &a.PatientID, &a.DoctorID, &a.Time, &a.Status, &a.CreatedAt, &a.DoctorName, &a.Specialization); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *PostgresAppointmentRepository) GetByDoctorID(ctx context.Context, doctorID int) ([]model.Appointment, error) {
	return nil, nil
}

func (r *PostgresAppointmentRepository) Cancel(ctx context.Context, id int) error {
	query := `UPDATE appointments SET status = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, model.StatusCancelled, id)
	return err
}
