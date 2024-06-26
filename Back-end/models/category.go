package models

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"forum/Back-end/handlers"
)

var once sync.Once

func HandleCategory(w http.ResponseWriter, r *http.Request) {
	InitializeDatabase() // Veritabanını sadece bir kere başlat

	userIDCookie, err := r.Cookie("user_id")
	loggedIn := err == nil

	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Veritabanından kategorileri çek
	rows, err := db.Query("SELECT id, name, description, link FROM categories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []handlers.Category
	for rows.Next() {
		var category handlers.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Link)
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
	tmplData := map[string]interface{}{
		"Categories": categories,
		"LoggedIn":   loggedIn,
	}

	if loggedIn {
		tmplData["UserID"] = userIDCookie.Value
	}

	if err := tmpl.Execute(w, tmplData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func InitializeDatabase() {
	once.Do(func() {
		db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Tabloyu oluştur
		statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT
			link TEXT NOT NULL
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
			statement, err = db.Prepare("INSERT INTO categories (name, description, link) VALUES (?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Seyahat Önerileri", "En İyi Seyahat Önerileri.", "/seyahatoneri")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Yurt İçi Geziler", "Türkiye İçindeki En İyi Geziler.", "/yurticigezi")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Yurt Dışı Geziler", "Yurt Dışındaki En İyi Geziler.", "/yurtdisigezi")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Vizesiz Geziler", "Vizesiz Gidilebilecek Ülkeler.", "/vizesizgezi")
			if err != nil {
				log.Fatal(err)
			}
			_, err = statement.Exec("Doga ve Kamp", "Doğa ve Kamp Hakkında Bilgiler.", "/dogavekamp")
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Database initialized successfully!")
	})
}
