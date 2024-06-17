package models

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"forum/handlers"
)

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["create_post"]
	if !ok {
		http.Error(w, "Could not load create_post template", http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Kategorileri çekmek için sorgu
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Could not retrieve categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []handlers.Category
	for rows.Next() {
		var category handlers.Category
		err := rows.Scan(&category.ID, &category.Name)
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

	// Kullanıcı giriş yapmış mı diye kontrol edelim
	_, err = r.Cookie("user_id")
	loggedIn := err == nil

	// Şablonu render et
	data := map[string]interface{}{
		"LoggedIn":   loggedIn,
		"Categories": categories,
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

func HandleSubmitPost(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get form values

	content := r.FormValue("content")
	title:=r.FormValue("title")
	categoryID := r.FormValue("category")

	// Get user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value

	// Handle file upload
	var imagePath string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		// Ensure the uploads directory exists
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			os.Mkdir("uploads", os.ModePerm)
		}

		// Create file
		imagePath = filepath.Join("uploads", handler.Filename)
		out, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// Copy the file content
		if _, err := io.Copy(out, file); err != nil {
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	// Insert post into the database
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts (user_id, title, content, image, category_id, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	createdAt := time.Now()
	_, err = stmt.Exec(userID, title, content, imagePath, categoryID, createdAt)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	// Redirect to a confirmation page or back to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
