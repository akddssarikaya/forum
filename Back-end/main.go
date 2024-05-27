package main

import (
	"database/sql"
	"encoding/json"
	"forum/handlers" // Import using module path
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var tmplCache = make(map[string]*template.Template)
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
	loadTemplates()
	log.Println("Tables created successfully!")
	staticFs := http.FileServer(http.Dir("../Front-end/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticFs))

	docsFs := http.FileServer(http.Dir("../Front-end/docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", docsFs))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/loginSubmit", handleLoginPost)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/registerSubmit", handleRegisterPost)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loadTemplates() {

	tmplCache["login"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "login.html")))
	tmplCache["register"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "register.html")))
	tmplCache["home"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "home.html")))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["home"]
	if !ok {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	tmpl, ok := tmplCache["login"]
	if !ok {
		http.Error(w, "Could not load login template", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func handleLoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var user handlers.User
		err = db.QueryRow("SELECT username, password, email FROM users WHERE username = ? OR email = ?", username, username).Scan(&user.Username, &user.Password, &user.Email)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Check if the provided password matches the hashed password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Redirect to profile page after successful login
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
func handleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["register"]
	if !ok {
		http.Error(w, "Could not load register template", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		user := handlers.User{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		}
		var userID int64
		err = db.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", email, username).Scan(&userID)
		if err == sql.ErrNoRows {
			userID, err = handlers.InsertUser(db, user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Username or email already exists", http.StatusForbidden)
		}
	}
}
