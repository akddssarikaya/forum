package models

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/Back-end/handlers"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	var user handlers.User
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

	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Kullanıcı gönderilerini çekelim
	rows, err := db.Query("SELECT id, user_id, title, content, image, category_id,  created_at, total_likes, total_dislikes FROM posts WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []handlers.Post
	for rows.Next() {
		var post handlers.Post

		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, "Could not scan post", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT name FROM categories WHERE id = ?", post.Category).Scan(&post.CategoryName)
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

	// Kullanıcının beğendiği gönderileri çekelim
	likedPostsRows, err := db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.image, p.category_id, p.created_at, p.total_likes, p.total_dislikes
		FROM posts p
		INNER JOIN likes l ON p.id = l.post_id
		WHERE l.user_id = ? AND l.like_type = 'like'`, userID)
	if err != nil {
		http.Error(w, "Could not retrieve liked posts", http.StatusInternalServerError)
		return
	}
	defer likedPostsRows.Close()

	var likedPosts []handlers.Post
	for likedPostsRows.Next() {
		var post handlers.Post

		err := likedPostsRows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, "Could not scan liked post", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT name FROM categories WHERE id = ?", post.Category).Scan(&post.CategoryName)
		if err != nil {
			http.Error(w, "Category name not found", http.StatusInternalServerError)
			return
		}

		if post.Image != "" {
			post.Image = "/" + post.Image
		}

		likedPosts = append(likedPosts, post)
	}

	// Kullanıcının beğenmediği gönderileri çekelim
	dislikedPostsRows, err := db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.image, p.category_id, p.created_at, p.total_likes, p.total_dislikes
		FROM posts p
		INNER JOIN likes l ON p.id = l.post_id
		WHERE l.user_id = ? AND l.like_type = 'dislike'`, userID)
	if err != nil {
		http.Error(w, "Could not retrieve disliked posts", http.StatusInternalServerError)
		return
	}
	defer dislikedPostsRows.Close()

	var dislikedPosts []handlers.Post
	for dislikedPostsRows.Next() {
		var post handlers.Post

		err := dislikedPostsRows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, "Could not scan disliked post", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT name FROM categories WHERE id = ?", post.Category).Scan(&post.CategoryName)
		if err != nil {
			http.Error(w, "Category name not found", http.StatusInternalServerError)
			return
		}

		if post.Image != "" {
			post.Image = "/" + post.Image
		}

		dislikedPosts = append(dislikedPosts, post)
	}
	// Kullanıcının yorumlarını çekelim
	commentsRows, err := db.Query(`
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
	FROM comments c
	INNER JOIN users u ON c.user_id = u.id
	WHERE c.user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Could not retrieve comments", http.StatusInternalServerError)
		return
	}
	defer commentsRows.Close()

	var comments []handlers.Comment
	for commentsRows.Next() {
		var comment handlers.Comment

		err := commentsRows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username)
		if err != nil {
			http.Error(w, "Could not scan comment", http.StatusInternalServerError)
			return
		}

		comments = append(comments, comment)
	}

	tmpl, ok := tmplCache["profile"]
	if !ok {
		http.Error(w, "Could not load profile template", http.StatusInternalServerError)
		return
	}

	// Kullanıcı bilgilerini ve gönderilerini template'e ekleyelim
	tmpl.Execute(w, map[string]interface{}{
		"Username":      user.Username,
		"Email":         user.Email,
		"LoggedIn":      true,
		"Posts":         posts,
		"LikedPosts":    likedPosts,
		"DislikedPosts": dislikedPosts,
		"Comments":      comments,
	})
}
