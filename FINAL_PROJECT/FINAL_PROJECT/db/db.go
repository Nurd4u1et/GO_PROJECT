package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./clinic.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	seedDoctors()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password_hash TEXT
	);

	CREATE TABLE IF NOT EXISTS doctors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		specialization TEXT
	);

	CREATE TABLE IF NOT EXISTS appointments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		doctor_id INTEGER,
		datetime TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(doctor_id) REFERENCES doctors(id)
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
}

func seedDoctors() {
	row := DB.QueryRow("SELECT COUNT(*) FROM doctors")
	var count int
	row.Scan(&count)

	if count == 0 {
		doctors := []struct {
			Name           string
			Specialization string
		}{}

		for _, d := range doctors {
			_, err := DB.Exec("INSERT INTO doctors (name, specialization) VALUES (?, ?)", d.Name, d.Specialization)
			if err != nil {
				log.Printf("Error seeding doctor %s: %v", d.Name, err)
			}
		}
		log.Println("Database seeded with initial doctors.")
	}
}
