package main

import (
	"database/sql"
	"forum/handlers" // Import using module path
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
)

func main() {
	var err error
	database, err = sql.Open("sqlite3", "./database/forum.db")
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
	handlers.CreateUserProfileTable(database)
	handlers.LoadTemplates()
	log.Println("Tables created successfully!")
	staticFs := http.FileServer(http.Dir("../Front-end/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticFs))

	docsFs := http.FileServer(http.Dir("../Front-end/docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", docsFs))

	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/loginSubmit", handlers.HandleLoginPost)
	http.HandleFunc("/register", handlers.HandleRegister)
	http.HandleFunc("/registerSubmit", handlers.HandleRegisterPost)
	http.HandleFunc("/profile", handlers.HandleProfile)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
