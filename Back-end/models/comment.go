package models

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func CommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userID, err := getUserIDFromCookie(r) // Bu fonksiyonun tanımlandığından emin olun
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("comment")
	if comment == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, postID, comment)
	if err != nil {
		http.Error(w, "Could not insert comment", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Comment posted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
