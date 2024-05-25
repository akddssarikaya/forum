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

	_, err = database.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Failed to clear users table: %s", err)
	}

	log.Println("Tables created successfully!")
	emails := []string{
		"emirtariik@gmail.com",
		"sezerdincer@gmail.com",
		"gulbeyza@gmail.com",
		"sude@gmail.com",
		"akuddusi@gmail.com",
		"bkaan@gmail.com",
	}
	usernames := []string{
		"emirtariik",
		"sezerdincer",
		"gulbeyza",
		"sude",
		"akuddusi",
		"bkaan",
	}
	passwords := []string{
		"12",
		"34",
		"56",
		"78",
		"910",
		"1112",
	}

	for i := 0; i < len(emails); i++ {
		user := handlers.User{
			Email:    emails[i],
			Username: usernames[i],
			Password: passwords[i],
		}
		// Veritabanına kullanıcıyı ekle
		userID, err := handlers.InsertUser(database, user)
		if err != nil {
			fmt.Println("Kullanıcı eklenirken bir hata oluştu:", err)
		} else {
			fmt.Println("Yeni kullanıcı eklendi, kullanıcı ID:", userID)
		}
	}

	handlers.PrintUsers(database)
}
