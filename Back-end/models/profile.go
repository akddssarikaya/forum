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
	// Kullanıcı gönderilerini çekelim
	rows, err := db.Query("SELECT id, user_id, content, image, category_id,  created_at, total_likes, total_dislikes FROM posts WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []handlers.Post
	for rows.Next() {
		var post handlers.Post

		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, "Could not scan post", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT name FROM categories WHERE id = ?", post.Category).Scan(&post.Title)
		if err != nil {
			http.Error(w, "Title not found", http.StatusInternalServerError)
			return
		}

		// Görüntü URL'sini oluşturun
		if post.Image != "" {
			post.Image = "/" + post.Image
		}

		posts = append(posts, post)
	}

	tmpl, ok := tmplCache["profile"]
	if !ok {
		http.Error(w, "Could not load profile template", http.StatusInternalServerError)
		return
	}

	// Kullanıcı bilgilerini ve gönderilerini template'e ekleyelim
	tmpl.Execute(w, map[string]interface{}{
		"Username": user.Username,
		"Email":    user.Email,

		"LoggedIn": true,
		"Posts":    posts,
	})
}
