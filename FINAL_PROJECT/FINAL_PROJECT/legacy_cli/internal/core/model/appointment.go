package model

import "time"

type Appointment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	DoctorID  int64     `json:"doctor_id"`
	DateTime  time.Time `json:"datetime"`
	CreatedAt time.Time `json:"created_at"`

	// Optional: Include relations for response if needed,
	// but purely for domain model, IDs are sufficient usually.
	// We can compose response structs in handlers if needed.
}
