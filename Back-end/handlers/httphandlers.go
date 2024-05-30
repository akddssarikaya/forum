package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var tmplCache = make(map[string]*template.Template)

func LoadTemplates() {
	templates := []string{"login", "register", "home", "profile", "panel", "category"}
	for _, tmpl := range templates {
		path := filepath.Join("..", "Front-end", "pages", tmpl+".html")
		tmplCache[tmpl] = template.Must(template.ParseFiles(path))
	}
	log.Println("Templates loaded successfully!")
}

func HandleCategory(w http.ResponseWriter, r *http.Request) {
	InitializeDatabase() // Veritabanını sadece bir kere başlat

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Veritabanından kategorileri çek
	rows, err := db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}

	tmpl, ok := tmplCache["category"]
	if !ok {
		http.Error(w, "Could not load category template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var once sync.Once

func InitializeDatabase() {
	once.Do(func() {
		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Tabloyu oluştur
		statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT
		)`)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()

		// Tablo boş mu diye kontrol et
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		// Eğer tablo boşsa örnek verileri ekle
		if count == 0 {
			statement, err = db.Prepare("INSERT INTO categories (name, description) VALUES (?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Seyehat Önerileri", "En iyi seyehat önerileri")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Yurt İçi Geziler", "Türkiye içindeki en iyi geziler")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Yurt Dışı Geziler", "Yurt dışındaki en iyi geziler")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Vizesiz Yurt Dışı Gezileri", "Vizesiz gidilebilecek ülkeler")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Doğa ve Kamp", "Doğa ve kamp hakkında bilgiler")
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Database initialized successfully!")
	})
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["home"]
	if !ok {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleLogout clears the user cookie and redirects to the login page
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear the user cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["login"]
	if !ok {
		http.Error(w, "Could not load login template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleAdmin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	tmpl, ok := tmplCache["panel"]
	if !ok {
		http.Error(w, "Could not load panel template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		userIdStr := r.FormValue("userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if err := deletetable(db, userId); err != nil {
			http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
			log.Println("Failed to delete user:", err)
			return
		}

		fmt.Fprintf(w, "User with ID %d deleted successfully", userId)
	} else {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func deletetable(database *sql.DB, userId int) error {
	// Transaction başlat
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Kullanıcıyı sil
	deleteStmt := "DELETE FROM users WHERE id = ?;"
	if _, err := tx.Exec(deleteStmt, userId); err != nil {
		return err
	}

	// Tabloda kalan satır sayısını kontrol et
	rowCount := 0
	countStmt := "SELECT COUNT(*) FROM users;"
	if err := tx.QueryRow(countStmt).Scan(&rowCount); err != nil {
		return err
	}

	// Eğer tablo boş ise, otomatik artan değeri sıfırla
	if rowCount == 0 {
		resetAIStmt := "DELETE FROM SQLITE_SEQUENCE WHERE NAME = 'users';"
		if _, err := tx.Exec(resetAIStmt); err != nil {
			return err
		}
	}

	// Transaction'ı commit et
	return tx.Commit()
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

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if _, err := InsertOrUpdateProfile(db, userID, user.Username, user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := tmplCache["profile"]
	if !ok {
		http.Error(w, "Could not load profile template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
