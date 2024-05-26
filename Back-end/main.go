package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"forum/handlers" // Import using module path

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
)

func main() {
	var err error
	database, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	handlers.CreateUserTable(database)
	handlers.CreatePostTable(database)
	handlers.CreateLikesTable(database)
	handlers.CreateCommentsTable(database)

	log.Println("Tables created successfully!")

	// Serve static files
	fs := http.FileServer(http.Dir("../Front-end"))
	http.Handle("/Front-end/", http.StripPrefix("/Front-end/", fs))

	http.HandleFunc("/register", registerHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		user := handlers.User{
			Email:    email,
			Username: username,
			Password: password,
		}

		// Check if the user already exists in the database
		var userID int64
		err := database.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", user.Email, user.Username).Scan(&userID)
		if err == sql.ErrNoRows {
			// If user doesn't exist, insert into the database
			userID, err = handlers.InsertUser(database, user)
			if err != nil {
				http.Error(w, "Error adding user", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
		} else if err != nil {
			http.Error(w, "Error checking user", http.StatusInternalServerError)
		} else {
			// If user already exists, return a conflict error
			http.Error(w, "User already exists", http.StatusConflict)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
