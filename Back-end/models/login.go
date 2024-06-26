package models

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/Back-end/handlers"

	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["login"]
	if !ok {
		http.Error(w, "Could not load login template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Open the login database
		loginDB, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer loginDB.Close()

		var user handlers.User
		err = loginDB.QueryRow("SELECT id, username, password, email FROM users WHERE username = ? OR email = ?", username, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Check if the provided password matches the hashed password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Insert or update the profile data
		_, err = loginDB.Exec(`
            INSERT INTO profile (user_id, username, email, last_login)
            VALUES (?, ?, ?, datetime('now'))
            ON CONFLICT(user_id) DO UPDATE SET 
                last_login = datetime('now')
        `, user.ID, user.Username, user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the user ID in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    fmt.Sprint(user.ID),
			Path:     "/",
			HttpOnly: true,
		})

		// Redirect to profile page after successful login
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
