package handlers

import (
	"database/sql"
	"errors"
	"log"
)

func InsertUser(database *sql.DB, user User) (int64, error) {
	var existingID int
	err := database.QueryRow("SELECT id FROM users WHERE username = ?", user.Username).Scan(&existingID)
	if err == nil {
		return 0, errors.New("username already exists")
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	insertUserSQL := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`
	statement, err := database.Prepare(insertUserSQL)
	if err != nil {
		return 0, err
	}
	res, err := statement.Exec(user.Email, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Printf("User %s inserted with ID: %d", user.Username, lastID)
	return lastID, nil
}
