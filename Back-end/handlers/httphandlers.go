package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var tmplCache = make(map[string]*template.Template)

func LoadTemplates() {

	tmplCache["login"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "login.html")))
	tmplCache["register"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "register.html")))
	tmplCache["home"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "home.html")))
	tmplCache["profile"] = template.Must(template.ParseFiles(filepath.Join("..", "Front-end", "pages", "profile.html")))
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
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

func HandleLogin(w http.ResponseWriter, r *http.Request) {

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

func HandleProfile(w http.ResponseWriter, r *http.Request) {
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

	// Open the database
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Check if profile already exists for the user
	var existingID int
	err = db.QueryRow("SELECT id FROM profile WHERE user_id = ?", userID).Scan(&existingID)
	if err == nil {
		// Profile already exists, proceed with retrieving and displaying the profile
		var user = User{} // Define an empty User struct
		err := db.QueryRow("SELECT username, email FROM profile WHERE user_id = ?", userID).Scan(&user.Username, &user.Email)
		if err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}

		tmpl, ok := tmplCache["profile"]
		if !ok {
			http.Error(w, "Could not load profile template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err == sql.ErrNoRows {
		// Profile doesn't exist, insert a new profile for the user
		_, err := InsertProfile(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the profile page again to display the newly inserted profile
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	} else {
		// Error occurred while querying for profile existence
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Open the login database
		loginDB, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer loginDB.Close()

		var user User
		err = loginDB.QueryRow("SELECT id, username, password, email FROM users WHERE username = ? OR email = ?", username, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
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

		// Insert or update the profile data
		_, err = loginDB.Exec(`
            INSERT INTO profile (user_id, username, email, last_login)
            VALUES (?, ?, ?, datetime('now'))
            ON CONFLICT(user_id) DO UPDATE SET 
                last_login = datetime('now')
        `, user.ID, user.Username, user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the user ID in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    fmt.Sprint(user.ID),
			Path:     "/",
			HttpOnly: true,
		})

		// Redirect to profile page after successful login
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
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
func HandleRegisterPost(w http.ResponseWriter, r *http.Request) {
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

		user := User{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		}
		var userID int64
		err = db.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", email, username).Scan(&userID)
		if err == sql.ErrNoRows {
			userID, err = InsertUser(db, user)
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
