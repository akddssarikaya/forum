package handlers

import "time"

type User struct {
	ID       int64
	Email    string
	Username string
	Password string
}
type Category struct {
	ID          int
	Name        string
	Description string
	Link        string
}
type Post struct {
	ID        int
	UserID    int
	Username  string
	Title     string
	Content   string
	Image     string
	Category  int
	CreatedAt *time.Time
	Likes     int
	Dislikes  int
}
type Profile struct {
	ID             int
	UserID         int
	Username       string
	Email          string
	Last_login     *time.Location
	Total_likes    int
	Total_dislikes int
}
