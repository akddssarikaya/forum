package models

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/handlers"
)

func HandleViewPost(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["view_post"]
	if !ok {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}

	// Veritabanı bağlantısı açma
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// URL'den id parametresini al
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// idParam'i integer'a çevir
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Gönderi bilgilerini çekme
	var post handlers.Post
	err = db.QueryRow("SELECT id, user_id, title, content, image, category_id, created_at, total_likes, total_dislikes FROM posts WHERE id = ?", postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Gönderi başlığını ve oluşturan kullanıcı adını al
	err = db.QueryRow("SELECT name FROM categories WHERE id = ?", post.Category).Scan(&post.CategoryName)
	if err != nil {
		http.Error(w, "Category not found", http.StatusInternalServerError)
		return
	}
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", post.UserID).Scan(&post.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Görüntü URL'sini oluşturun
	if post.Image != "" {
		post.Image = "/" + post.Image
	}

	// Gönderiye ait yorumları çekme
	rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", post.ID)
	if err != nil {
		http.Error(w, "Could not retrieve comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []handlers.Comment
	for rows.Next() {
		var comment handlers.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			http.Error(w, "Could not scan comment", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", comment.UserID).Scan(&comment.Username)
		if err != nil {
			http.Error(w, "User not found for comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	// Yorumları gönderiye ekle
	post.Comments = comments

	// Kullanıcı giriş yapmış mı kontrol et
	_, err = r.Cookie("user_id")
	loggedIn := err == nil

	// Şablonu render etmek için veri hazırla
	data := map[string]interface{}{
		"LoggedIn":          loggedIn,
		"Post":              post,
		"ShowLoginRegister": !loggedIn, // Kullanıcı giriş yapmamışsa Login ve Register bağlantılarını göster
	}

	// Şablonu execute et
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		return
	}
}
