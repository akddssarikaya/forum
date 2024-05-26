package handlers

import (
	"database/sql"
	"log"
)

func CreateUserTable(database *sql.DB) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := database.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Users table creation failed: %s", err)
	}
}

func CreatePostTable(database *sql.DB) {
	createPostsTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        image TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
	_, err := database.Exec(createPostsTable)
	if err != nil {
		log.Fatalf("Posts table creation failed: %s", err)
	}
}

func CreateCommentsTable(database *sql.DB) {
	createCommentsTable := `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
	_, err := database.Exec(createCommentsTable)
	if err != nil {
		log.Fatalf("Comments table creation failed: %s", err)
	}
}

func CreateLikesTable(database *sql.DB) {
	createLikesTable := `
    CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER,
        comment_id INTEGER,
        like_count INTEGER DEFAULT 0,
        dislike_count INTEGER DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (comment_id) REFERENCES comments(id)
    );`
	_, err := database.Exec(createLikesTable)
	if err != nil {
		log.Fatalf("Likes table creation failed: %s", err)
	}
}