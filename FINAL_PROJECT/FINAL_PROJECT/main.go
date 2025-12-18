package main

import (
	"bufio"
	"clinic-cli/auth"
	"clinic-cli/core"
	"clinic-cli/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/term"
)

func main() {
	go startHTTPServer()
	db.InitDB()
	fmt.Println("Welcome ,please choose")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if auth.CurrentUser == nil {
			showGuestMenu()
			if !scanner.Scan() {
				break
			}
			choice := strings.TrimSpace(scanner.Text())
			handleGuestChoice(choice, scanner)
		} else {
			showUserMenu()
			if !scanner.Scan() {
				break
			}
			choice := strings.TrimSpace(scanner.Text())
			handleUserChoice(choice, scanner)
		}
		fmt.Println()
	}
}

func showGuestMenu() {
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("3. Exit")
	fmt.Print("Enter choice: ")
}

func showUserMenu() {
	fmt.Printf("Welcome, %s!\n", auth.CurrentUser.Username)
	fmt.Println("1. List Doctors")
	fmt.Println("2. Book Appointment")
	fmt.Println("3. My Appointments")
	fmt.Println("4. Logout")
	fmt.Print("Enter choice: ")
}

func handleGuestChoice(choice string, scanner *bufio.Scanner) {
	switch choice {
	case "1":
		login(scanner)
	case "2":
		register(scanner)
	case "3":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice")
	}
}

func handleUserChoice(choice string, scanner *bufio.Scanner) {
	switch choice {
	case "1":
		listDoctors()
	case "2":
		bookAppointment(scanner)
	case "3":
		listMyAppointments()
	case "4":
		auth.Logout()
		fmt.Println("Logged out successfully.")
	default:
		fmt.Println("Invalid choice")
	}
}

func login(scanner *bufio.Scanner) {
	fmt.Print("Username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password")
		return
	}
	password := string(bytePassword)
	fmt.Println()

	err = auth.Login(username, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
	} else {
		fmt.Println("Login successful!")
	}
}

func register(scanner *bufio.Scanner) {
	fmt.Print("New Username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Print("New Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password")
		return
	}
	password := string(bytePassword)
	fmt.Println()

	err = auth.Register(username, password)
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
	} else {
		fmt.Println("Registration successful! You can now login.")
	}
}

func listDoctors() {
	doctors, err := core.ListDoctors()
	if err != nil {
		fmt.Printf("Error fetching doctors: %v\n", err)
		return
	}
	fmt.Println("\n--- Available Doctors ---")
	for _, d := range doctors {
		fmt.Printf("[%d] %s (%s)\n", d.ID, d.Name, d.Specialization)
	}
}

func bookAppointment(scanner *bufio.Scanner) {
	listDoctors()
	fmt.Print("\nEnter Doctor ID to book: ")
	scanner.Scan()
	var docID int
	_, err := fmt.Sscan(scanner.Text(), &docID)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	fmt.Print("Enter Date and Time (e.g., 2023-12-01 14:00): ")
	scanner.Scan()
	dateTime := strings.TrimSpace(scanner.Text())

	err = core.BookAppointment(docID, dateTime)
	if err != nil {
		fmt.Printf("Booking failed: %v\n", err)
	} else {
		fmt.Println("Appointment booked successfully!")
	}
}

func getUsersFromDB() ([]string, error) {
	dbConn, err := sql.Open("sqlite3", "clinic.db")
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}

	return users, nil
}

func listMyAppointments() {
	apps, err := core.ListMyAppointments()
	if err != nil {
		fmt.Printf("Error fetching appointments: %v\n", err)
		return
	}
	if len(apps) == 0 {
		fmt.Println("\nYou have no upcoming appointments.")
		return
	}

	fmt.Println("\n--- My Appointments ---")
	for _, a := range apps {
		fmt.Printf("- %s with %s (%s)\n", a.DateTime, a.DoctorName, a.Specialization)
	}
}

func startHTTPServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Doctor Appointment API is running"))
	})

	http.HandleFunc("/appointments", func(w http.ResponseWriter, r *http.Request) {
		data, err := getAppointmentsDetailed()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(data)
	})

	http.HandleFunc("/doctors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`["Dr. Smith", "Dr. Brown"]`))
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := getUsersFromDB()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(users)
	})

	http.ListenAndServe(":8080", nil)
}

type AppointmentResponse struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Doctor   string `json:"doctor"`
	DateTime string `json:"datetime"`
}

func getAppointmentsDetailed() ([]AppointmentResponse, error) {
	dbConn, err := sql.Open("sqlite3", "clinic.db")
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	query := `
	SELECT 
		a.id,
		u.username,
		d.name,
		a.datetime
	FROM appointments a
	JOIN users u ON a.user_id = u.id
	JOIN doctors d ON a.doctor_id = d.id
	`

	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []AppointmentResponse{}

	for rows.Next() {
		var ap AppointmentResponse
		if err := rows.Scan(&ap.ID, &ap.User, &ap.Doctor, &ap.DateTime); err != nil {
			return nil, err
		}
		result = append(result, ap)
	}

	return result, nil
}
