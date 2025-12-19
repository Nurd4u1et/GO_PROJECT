package handler

import (
	"clinic-cli/internal/middleware"
	"clinic-cli/internal/model"
	"clinic-cli/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	AuthService   *service.AuthService
	ClinicService *service.ClinicService
}

func NewHandler(as *service.AuthService, cs *service.ClinicService) *Handler {
	return &Handler{AuthService: as, ClinicService: cs}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.AuthService.Register(r.Context(), req.Email, req.Password, model.RolePatient)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusCreated, user)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.AuthService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) ListDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := h.ClinicService.ListDoctors(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch doctors")
		return
	}
	jsonResponse(w, http.StatusOK, doctors)
}

type CreateDoctorRequest struct {
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

func (h *Handler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	var req CreateDoctorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	doc, err := h.ClinicService.CreateDoctor(r.Context(), req.Name, req.Specialization)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to create doctor")
		return
	}
	jsonResponse(w, http.StatusCreated, doc)
}

type BookAppointmentRequest struct {
	DoctorID int    `json:"doctor_id"`
	Time     string `json:"time"`
}

func (h *Handler) BookAppointment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	patientID := int(claims["sub"].(float64))

	var req BookAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	app, err := h.ClinicService.BookAppointment(r.Context(), patientID, req.DoctorID, req.Time)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusCreated, app)
}

func (h *Handler) MyAppointments(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	patientID := int(claims["sub"].(float64))

	apps, err := h.ClinicService.MyAppointments(r.Context(), patientID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch appointments")
		return
	}
	jsonResponse(w, http.StatusOK, apps)
}

func (h *Handler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.ClinicService.CancelAppointment(r.Context(), id); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to cancel")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "cancelled"})
}
