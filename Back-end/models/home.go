package models

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/Back-end/handlers"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["home"]
	if !ok {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Kategorileri çekmek için sorgu
	rows, err := db.Query("SELECT id, name, link FROM categories")
	if err != nil {
		http.Error(w, "Could not retrieve categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []handlers.Category
	for rows.Next() {
		var category handlers.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Link)
		if err != nil {
			http.Error(w, "Error scanning category", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating categories", http.StatusInternalServerError)
		return
	}

	// Gönderileri çekmek için sorgu
	rows1, err := db.Query(`
		SELECT 
			posts.id, posts.user_id, posts.title, posts.content, posts.image, posts.category_id, posts.created_at, posts.total_likes, posts.total_dislikes,
			categories.name, users.username
		FROM posts
		JOIN categories ON posts.category_id = categories.id
		JOIN users ON posts.user_id = users.id`)
	if err != nil {
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows1.Close()

	var posts []handlers.Post
	for rows1.Next() {
		var post handlers.Post
		err := rows1.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.CategoryName, &post.Username)
		if err != nil {
			http.Error(w, "Could not scan post", http.StatusInternalServerError)
			return
		}

		// Görüntü URL'sini oluşturun
		if post.Image != "" {
			post.Image = "/" + post.Image
		}

		// Yorumları çekmek için sorgu
		rows2, err := db.Query(`
			SELECT 
				comments.id, comments.post_id, comments.user_id, comments.content, comments.created_at, users.username,
				IFNULL(SUM(CASE WHEN comment_likes.like_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
				IFNULL(SUM(CASE WHEN comment_likes.like_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
			FROM comments 
			LEFT JOIN comment_likes ON comments.id = comment_likes.comment_id
			JOIN users ON comments.user_id = users.id
			WHERE comments.post_id = ?
			GROUP BY comments.id`, post.ID)
		if err != nil {
			http.Error(w, "Could not retrieve comments", http.StatusInternalServerError)
			return
		}
		defer rows2.Close()

		var comments []handlers.Comment
		for rows2.Next() {
			var comment handlers.Comment
			err := rows2.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username, &comment.Likes, &comment.Dislikes)
			if err != nil {
				http.Error(w, "Could not scan comment", http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}
		post.Comments = comments

		posts = append(posts, post)
	}

	// Kullanıcı giriş yapmış mı diye kontrol edelim
	_, err = r.Cookie("user_id")
	loggedIn := err == nil

	// Şablonu render et
	data := map[string]interface{}{
		"LoggedIn":   loggedIn,
		"Categories": categories,
		"Posts":      posts,
	}

	// Kullanıcı giriş yapmamışsa Login ve Register bağlantılarını ekle
	if !loggedIn {
		data["ShowLoginRegister"] = true
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		return
	}
}

func HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	HandleCommentLikeDislike(w, r, "like")
}

func HandleDislikeComment(w http.ResponseWriter, r *http.Request) {
	HandleCommentLikeDislike(w, r, "dislike")
}

func HandleCommentLikeDislike(w http.ResponseWriter, r *http.Request, likeType string) {
	r.ParseForm()
	commentID, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingLikeType string
	err = db.QueryRow("SELECT like_type FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID.Value, commentID).Scan(&existingLikeType)

	if err == sql.ErrNoRows {
		// Kullanıcı daha önce bu yorumu beğenmemiş veya beğenmemiş
		_, err = db.Exec("INSERT INTO comment_likes (user_id, comment_id, like_type) VALUES (?, ?, ?)", userID.Value, commentID, likeType)
		if err != nil {
			http.Error(w, "Could not insert like/dislike", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Could not query like/dislike", http.StatusInternalServerError)
		return
	} else if existingLikeType == likeType {
		// Kullanıcı zaten bu yorumu beğenmiş veya beğenmeme yapmış, geri al
		_, err = db.Exec("DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID.Value, commentID)
		if err != nil {
			http.Error(w, "Could not remove like/dislike", http.StatusInternalServerError)
			return
		}
	} else {
		// Kullanıcı farklı bir işlem yapmış, güncelle
		_, err = db.Exec("UPDATE comment_likes SET like_type = ? WHERE user_id = ? AND comment_id = ?", likeType, userID.Value, commentID)
		if err != nil {
			http.Error(w, "Could not update like/dislike", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
