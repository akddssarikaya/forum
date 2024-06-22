package handlers

import (
	"database/sql"
	"log"
)

func CreateCategoryTable(database *sql.DB) {
	CreateCategoryTable := `
	CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		link TEXT NOT NULL
	);
	
	`
	_, err := database.Exec(CreateCategoryTable)
	if err != nil {
		log.Fatalf("User profile table creation failed: %s", err)
	}
}

func CreateCommentLikesTable(database *sql.DB) {
	createCommentLikesTable := `
	CREATE TABLE IF NOT EXISTS comment_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		comment_id INTEGER NOT NULL,
		like_type TEXT NOT NULL,
		UNIQUE(user_id, comment_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
	);`
	_, err := database.Exec(createCommentLikesTable)
	if err != nil {
		log.Fatalf("CommentLikes table creation failed: %s", err)
	}
}

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

func CreateProfileTable(database *sql.DB) {
	CreateProfileTable := `
	CREATE TABLE IF NOT EXISTS profile (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER UNIQUE NOT NULL,
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	last_login TIMESTAMP,
	total_likes INTEGER DEFAULT 0,
	total_dislikes INTEGER DEFAULT 0,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err := database.Exec(CreateProfileTable)
	if err != nil {
		log.Fatalf("User profile table creation failed: %s", err)
	}
}

func CreatePostTable(database *sql.DB) {
	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		image TEXT,
		category_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		total_likes INTEGER DEFAULT 0,
		total_dislikes INTEGER DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
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
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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
    post_id INTEGER NOT NULL,
    like_type TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
`
	_, err := database.Exec(createLikesTable)
	if err != nil {
		log.Fatalf("Likes table creation failed: %s", err)
	}
}
