package auth

import (
	"clinic-cli/db"
	"clinic-cli/models"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

var CurrentUser *models.User

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func Register(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	hashedPwd := hashPassword(password)

	_, err := db.DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hashedPwd)
	if err != nil {
		return fmt.Errorf("could not register user (might already exist): %v", err)
	}
	return nil
}

func Login(username, password string) error {
	hashedPwd := hashPassword(password)

	row := db.DB.QueryRow("SELECT id, username FROM users WHERE username = ? AND password_hash = ?", username, hashedPwd)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		return errors.New("invalid username or password")
	} else if err != nil {
		return err
	}

	CurrentUser = user
	return nil
}

func Logout() {
	CurrentUser = nil
}
