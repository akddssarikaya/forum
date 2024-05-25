package handlers

import (
	"database/sql"
	"log"
)

func PrintUsers(database *sql.DB) {
	rows, err := database.Query("SELECT id, email, username, password FROM users ORDER BY id ASC")
	if err != nil {
		log.Fatalf("Failed to fetch users: %s", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
			log.Fatalf("Failed to scan user row: %s", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error during iteration through user rows: %s", err)
	}

	log.Println("Users:")
	for _, user := range users {
		log.Printf("ID: %d, Email: %s, Username: %s, Password: %s", user.ID, user.Email, user.Username, user.Password)
	}
}
