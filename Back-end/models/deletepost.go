package models

import (
	"database/sql"
	"net/http"
)

// Gönderi silme işlemi
func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID missing", http.StatusBadRequest)
		return
	}

	// Kullanıcı kimliği alınır
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}
	userID := cookie.Value

	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	// Gönderinin kullanıcıya ait olup olmadığını kontrol et
	var ownerID string
	err = db.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Transaction başlat
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Could not begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Transaction işlemi başarısız olursa geri alma

	// Gönderiye ait yorumları sil
	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
	if err != nil {
		http.Error(w, "Could not delete comments", http.StatusInternalServerError)
		return
	}

	// Gönderiyi sil
	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		http.Error(w, "Could not delete post", http.StatusInternalServerError)
		return
	}

	// Transaction commit et
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Could not commit transaction", http.StatusInternalServerError)
		return
	}

	// Başarılı bir şekilde silindiğini belirt
	w.Write([]byte("Post and its comments deleted successfully"))
}
