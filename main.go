package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/Back-end/handlers" // Import using module path
	"forum/Back-end/models"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func main() {
	var err error
	database, err = sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	handlers.CreateUserTable(database)
	handlers.CreatePostTable(database)
	handlers.CreateLikesTable(database)
	handlers.CreateCommentsTable(database)
	handlers.CreateProfileTable(database)
	handlers.CreateCommentLikesTable(database)

	models.LoadTemplates()
	log.Println("Tables created successfully!")
	staticFs := http.FileServer(http.Dir("./Front-end/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticFs))
	image := http.FileServer(http.Dir("./Front-end/img"))
	http.Handle("/img/", http.StripPrefix("/img/", image))

	docsFs := http.FileServer(http.Dir("./Front-end/docs"))
	// Handle form submission
	http.Handle("/docs/", http.StripPrefix("/docs/", docsFs))

	imageFs := http.FileServer(http.Dir("./Back-end/uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", imageFs))

	http.HandleFunc("/", models.HandleHome)
	http.HandleFunc("/login", models.HandleLogin)
	http.HandleFunc("/loginSubmit", models.HandleLoginPost)
	http.HandleFunc("/register", models.HandleRegister)
	http.HandleFunc("/registerSubmit", models.HandleRegisterPost)
	http.HandleFunc("/profile", models.HandleProfile)
	http.HandleFunc("/panel", models.HandleAdmin)
	http.HandleFunc("/logout", models.HandleLogout)
	http.HandleFunc("/submit_post", models.HandleSubmitPost)
	http.HandleFunc("/create_post", models.HandleCreatePost)
	http.HandleFunc("/category", models.HandleCategory)
	http.HandleFunc("/like", models.LikePost)
	// Yeni eklenen işlevler için handler fonksiyonlarını ekleyin
	http.HandleFunc("/like_comment", func(w http.ResponseWriter, r *http.Request) {
		models.HandleLikeComment(w, r)
	})

	http.HandleFunc("/dislike_comment", func(w http.ResponseWriter, r *http.Request) {
		models.HandleDislikeComment(w, r)
	})

	http.HandleFunc("/dislike", models.DislikePost)
	http.HandleFunc("/comment", models.CommentPost)
	http.HandleFunc("/view_post", models.HandleViewPost) // Corrected route
	http.HandleFunc("/delete_post", models.HandleDeletePost)
	http.HandleFunc("/delete_comment", models.HandleDeleteComment)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
