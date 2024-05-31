package models

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"forum/handlers"

	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["register"]
	if !ok {
		http.Error(w, "Could not load register template", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleRegisterPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		user := handlers.User{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		}
		var userID int64
		err = db.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", email, username).Scan(&userID)
		if err == sql.ErrNoRows {
			userID, err = handlers.InsertUser(db, user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Username or email already exists", http.StatusForbidden)
		}
	}
}
