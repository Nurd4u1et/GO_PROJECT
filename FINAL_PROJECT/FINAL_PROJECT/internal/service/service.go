package service

import (
	"clinic-cli/internal/model"
	"clinic-cli/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (s *AuthService) Register(ctx context.Context, email, password string, role model.Role) (*model.User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if role == "" {
		role = model.RolePatient
	}

	return s.repo.Create(ctx, email, string(hashedBytes), role)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type ClinicService struct {
	doctorRepo      repository.DoctorRepository
	appointmentRepo repository.AppointmentRepository
	notifyChan      chan<- model.Appointment
}

func NewClinicService(dr repository.DoctorRepository, ar repository.AppointmentRepository, notifyChan chan<- model.Appointment) *ClinicService {
	return &ClinicService{
		doctorRepo:      dr,
		appointmentRepo: ar,
		notifyChan:      notifyChan,
	}
}

func (s *ClinicService) ListDoctors(ctx context.Context) ([]model.Doctor, error) {
	return s.doctorRepo.GetAll(ctx)
}

func (s *ClinicService) CreateDoctor(ctx context.Context, name, spec string) (*model.Doctor, error) {
	return s.doctorRepo.Create(ctx, name, spec)
}

func (s *ClinicService) BookAppointment(ctx context.Context, patientID, doctorID int, timeStr string) (*model.Appointment, error) {
	app, err := s.appointmentRepo.Create(ctx, patientID, doctorID, timeStr)
	if err != nil {
		return nil, err
	}

	select {
	case s.notifyChan <- *app:
	default:

	}

	return app, nil
}

func (s *ClinicService) MyAppointments(ctx context.Context, patientID int) ([]model.Appointment, error) {
	return s.appointmentRepo.GetByPatientID(ctx, patientID)
}

func (s *ClinicService) CancelAppointment(ctx context.Context, appID int) error {
	return s.appointmentRepo.Cancel(ctx, appID)
}
