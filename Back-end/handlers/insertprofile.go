package handlers

import (
	"database/sql"
)

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
