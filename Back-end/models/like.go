package models

import (
	"database/sql"
	"net/http"
	"strconv"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
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

	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var likeType string
	err = db.QueryRow("SELECT like_type FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeType)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Could not retrieve like status", http.StatusInternalServerError)
		return
	}

	if likeType == "like" {
		_, err = db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			http.Error(w, "Could not remove like", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Could not update like count", http.StatusInternalServerError)
			return
		}
	} else if likeType == "dislike" {
		_, err = db.Exec("UPDATE likes SET like_type = 'like' WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			http.Error(w, "Could not update like status", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE posts SET total_likes = total_likes + 1, total_dislikes = total_dislikes - 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Could not update like/dislike count", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = db.Exec("INSERT INTO likes (user_id, post_id, like_type) VALUES (?, ?, 'like')", userID, postID)
		if err != nil {
			http.Error(w, "Could not insert like", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE posts SET total_likes = total_likes + 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Could not update like count", http.StatusInternalServerError)
			return
		}
	}

	// Profil tablosunu güncelle
	_, err = db.Exec(`
	  UPDATE profile 
	  SET total_likes = (
		  SELECT COUNT(*) FROM likes 
		  WHERE likes.user_id = ? AND likes.like_type = 'like'
	  ),
	  total_dislikes = (
		  SELECT COUNT(*) FROM likes 
		  WHERE likes.user_id = ? AND likes.like_type = 'dislike'
	  )
	  WHERE user_id = ?
  `, userID, userID, userID)
	if err != nil {
		http.Error(w, "Could not update profile likes/dislikes", http.StatusInternalServerError)
		return
	}

	// Yönlendirme yap
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
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

	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var likeType string
	err = db.QueryRow("SELECT like_type FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeType)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Could not retrieve like status", http.StatusInternalServerError)
		return
	}

	if likeType == "dislike" {
		_, err = db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			http.Error(w, "Could not remove dislike", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE posts SET total_dislikes = total_dislikes - 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Could not update dislike count", http.StatusInternalServerError)
			return
		}
	} else if likeType == "like" {
		_, err = db.Exec("UPDATE likes SET like_type = 'dislike' WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			http.Error(w, "Could not update dislike status", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE posts SET total_likes = total_likes - 1, total_dislikes = total_dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Could not update like/dislike count", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = db.Exec("INSERT INTO likes (user_id, post_id, like_type) VALUES (?, ?, 'dislike')", userID, postID)
		if err != nil {
			http.Error(w, "Could not insert dislike", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE posts SET total_dislikes = total_dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Could not update dislike count", http.StatusInternalServerError)
			return
		}
	}

	// Profil tablosunu güncelle
	_, err = db.Exec(`
  UPDATE profile 
  SET total_likes = (
	  SELECT COUNT(*) FROM likes 
	  WHERE likes.user_id = ? AND likes.like_type = 'like'
  ),
  total_dislikes = (
	  SELECT COUNT(*) FROM likes 
	  WHERE likes.user_id = ? AND likes.like_type = 'dislike'
  )
  WHERE user_id = ?
`, userID, userID, userID)
	if err != nil {
		http.Error(w, "Could not update profile likes/dislikes", http.StatusInternalServerError)
		return
	}

	// Yönlendirme yap
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUserIDFromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
