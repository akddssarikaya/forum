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

	tx, err := database.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var rowCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&rowCount)
	if err != nil {
		return 0, err
	}

	if rowCount == 0 {
		resetAIStmt := "DELETE FROM sqlite_sequence WHERE name = 'users';"
		if _, err := tx.Exec(resetAIStmt); err != nil {
			return 0, err
		}
	}

	insertUserSQL := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`
	statement, err := tx.Prepare(insertUserSQL)
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

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	log.Printf("User %s inserted with ID: %d", user.Username, lastID)
	return lastID, nil
}

func InsertOrUpdateProfile(db *sql.DB, userID int64, username, email string) (int64, error) {
	var profileID int64

	err := db.QueryRow("SELECT id FROM profile WHERE user_id = ?", userID).Scan(&profileID)
	if err == sql.ErrNoRows {
		// Insert new profile
		insertProfileSQL := `
            INSERT INTO profile (user_id, username, email)
            VALUES (?, ?, ?)
        `
		res, err := db.Exec(insertProfileSQL, userID, username, email)
		if err != nil {
			return 0, err
		}
		profileID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	} else {
		// Update existing profile
		updateProfileSQL := `
            UPDATE profile
            SET username = ?, email = ?
            WHERE user_id = ?
        `
		_, err := db.Exec(updateProfileSQL, username, email, userID)
		if err != nil {
			return 0, err
		}
	}
	return profileID, nil
}
