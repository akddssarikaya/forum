package models

import (
	"database/sql"
	"net/http"

	"forum/handlers"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["home"]
	if !ok {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
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
	rows1, err := db.Query("SELECT id, user_id, title, content, image, category_id, created_at, total_likes, total_dislikes FROM posts")
	if err != nil {
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows1.Close()

	var posts []handlers.Post
	for rows1.Next() {
		var post handlers.Post

		err := rows1.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, "Could not scan post", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT name FROM categories WHERE id = ?", post.Category).Scan(&post.CategoryName)
		if err != nil {
			http.Error(w, "Title not found", http.StatusInternalServerError)
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

		// Yorumları çekmek için sorgu
		rows2, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", post.ID)
		if err != nil {
			http.Error(w, "Could not retrieve comments", http.StatusInternalServerError)
			return
		}
		defer rows2.Close()

		var comments []handlers.Comment
		for rows2.Next() {
			var comment handlers.Comment
			err := rows2.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
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

		post.Comments = comments // Post yapısına yorumları ekleyin

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
