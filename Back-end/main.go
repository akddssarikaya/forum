package main

import (
	"database/sql"
	"fmt"
	"log"

	"forum/handlers" // Import using module path

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	handlers.CreateUserTable(database) // Use exported function

	log.Println("Tables created successfully!")
	emails := []string{
		"emirtariik@gmail.com",
		"sezerdincer@gmail.com",
		"gulbeyza@gmail.com",
		"sude@gmail.com",
		"akuddusi@gmail.com",
		"bkaan@gmail.com",
		"emirtariik@gmail.com",
	}
	usernames := []string{
		"emirtariik",
		"sezerdincer",
		"gulbeyza",
		"sude",
		"akuddusi",
		"bkaan",
		"emirtariik",
	}
	passwords := []string{
		"12",
		"34",
		"56",
		"78",
		"910",
		"1112",
		"12",
	}

	for i := 0; i < len(emails); i++ {
		user := handlers.User{
			Email:    emails[i],
			Username: usernames[i],
			Password: passwords[i],
		}

		// Kullanıcı zaten mevcut mu kontrol et
		var userID int64
		err := database.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", user.Email, user.Username).Scan(&userID)
		if err == sql.ErrNoRows {
			// Kullanıcı mevcut değil, ekleyelim
			userID, err = handlers.InsertUser(database, user)
			if err != nil {
				fmt.Println("Kullanıcı eklenirken bir hata oluştu:", err)
			} else {
				fmt.Println("Yeni kullanıcı eklendi, kullanıcı ID:", userID)
			}
		} else if err != nil {
			fmt.Println("Kullanıcı kontrol edilirken bir hata oluştu:", err)
		} else {
			fmt.Println("Kullanıcı zaten mevcut, kullanıcı ID:", userID)
			// Burada kullanıcı zaten mevcut olduğu için bir uyarı verebiliriz.
			fmt.Println("Uyarı: Bu kullanıcı zaten mevcut!")
		}
	}

	handlers.PrintUsers(database)
}
