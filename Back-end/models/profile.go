package models

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/handlers"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user handlers.User
	err = db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if _, err := handlers.InsertOrUpdateProfile(db, userID, user.Username, user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := tmplCache["profile"]
	if !ok {
		http.Error(w, "Could not load profile template", http.StatusInternalServerError)
		return
	}

	// Kullanıcı bilgilerini ve oturum durumunu template'e ekleyelim
	tmpl.Execute(w, map[string]interface{}{
		"User":     user,
		"LoggedIn": true,
	})
}
