package core

import (
	"clinic-cli/auth"
	"clinic-cli/db"
	"clinic-cli/models"
	"fmt"
)

func ListDoctors() ([]models.Doctor, error) {
	rows, err := db.DB.Query("SELECT id, name, specialization FROM doctors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doc models.Doctor
		if err := rows.Scan(&doc.ID, &doc.Name, &doc.Specialization); err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}
	return doctors, nil
}

func BookAppointment(doctorID int, dateTime string) error {
	if auth.CurrentUser == nil {
		return fmt.Errorf("not logged in")
	}

	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM appointments WHERE doctor_id = ? AND datetime = ?", doctorID, dateTime).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("this slot is already booked")
	}

	_, err = db.DB.Exec("INSERT INTO appointments (user_id, doctor_id, datetime) VALUES (?, ?, ?)",
		auth.CurrentUser.ID, doctorID, dateTime)
	return err
}

func ListMyAppointments() ([]struct {
	ID             int
	DoctorName     string
	Specialization string
	DateTime       string
}, error) {
	if auth.CurrentUser == nil {
		return nil, fmt.Errorf("not logged in")
	}

	query := `
		SELECT a.id, d.name, d.specialization, a.datetime
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		WHERE a.user_id = ?
		ORDER BY a.datetime
	`
	rows, err := db.DB.Query(query, auth.CurrentUser.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []struct {
		ID             int
		DoctorName     string
		Specialization string
		DateTime       string
	}

	for rows.Next() {
		var item struct {
			ID             int
			DoctorName     string
			Specialization string
			DateTime       string
		}
		if err := rows.Scan(&item.ID, &item.DoctorName, &item.Specialization, &item.DateTime); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}
