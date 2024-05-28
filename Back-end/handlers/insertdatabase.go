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

func InsertProfile(db *sql.DB, userID int64) (int64, error) {
	var user User
	err := db.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}

	// Check if profile already exists for the user
	var existingID int
	err = db.QueryRow("SELECT id FROM profile WHERE user_id = ?", userID).Scan(&existingID)
	if err == nil {
		// Profile exists, perform an update
		updateProfileSQL := `
            UPDATE profile
            SET email = ?, username = ?
            WHERE user_id = ?
        `
		_, err := db.Exec(updateProfileSQL, user.Email, user.Username, userID)
		if err != nil {
			return 0, err
		}
		return int64(existingID), nil // Return the existing profile ID
	} else if err == sql.ErrNoRows {
		// Profile doesn't exist, insert a new profile for the user
		insertProfileSQL := `
            INSERT INTO profile (email, username, user_id)
            VALUES (?, ?, ?)
        `
		res, err := db.Exec(insertProfileSQL, user.Email, user.Username, userID)
		if err != nil {
			return 0, err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		return lastID, nil
	}
	return 0, err
}
