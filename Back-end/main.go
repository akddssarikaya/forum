package main

import (
	"database/sql"
	"fmt"
	"forum/handlers" // Import using module path
	"log"
	"net/http"
	"strconv"

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
	// Handle form submission
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Parse user ID from form data
		userIdStr := r.FormValue("userId")
		userId, err := strconv.Atoi(userIdStr)

		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Call the deletetable function
		if err := deletetable(database, userId); err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		// Respond with success message
		fmt.Fprintf(w, "User with ID %d deleted successfully", userId)
	})

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func deletetable(database *sql.DB, userId int) error {
	// Begin transaction
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Define SQL statements as strings
	deleteStmt := "DELETE FROM users WHERE id = ?;"
	updateStmt := "UPDATE users SET id = id - 1 WHERE id > ?;"

	// Execute the delete statement
	_, err = tx.Exec(deleteStmt, userId)
	if err != nil {
		return err
	}

	// Execute the update statement
	_, err = tx.Exec(updateStmt, userId)
	if err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
