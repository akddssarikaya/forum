package handlers

import (
	"database/sql"
	"log"
)

func CreateUserTable(database *sql.DB) {
	createUsersTable := `
	DROP TABLE IF EXISTS users;
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := database.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Users table creation failed: %s", err)
	}
}
