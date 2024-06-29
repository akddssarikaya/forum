package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	githubOauthConfig = &oauth2.Config{
		ClientID:     "Ov23liJECUYvrVc2Bxyi",
		ClientSecret: "855297a4c778a144aa97607374162fd084acfd87",
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	googleOauthConfig = &oauth2.Config{
		ClientID:     "YOUR_GOOGLE_CLIENT_ID",
		ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "randomstring"
)

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := githubOauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var user struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	storeUserInDB(user.Login, user.Email, "")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := googleOauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	storeUserInDB(user.Name, user.Email, "")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func storeUserInDB(username, email, password string) {
	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Kullanıcıyı 'users' tablosuna ekle
	insertUserQuery := `
        INSERT INTO users (username, email, password)
        VALUES (?, ?, ?)
        ON CONFLICT(email) DO UPDATE SET
        username = excluded.username,
        password = excluded.password
    `
	_, err = db.Exec(insertUserQuery, username, email, password)
	if err != nil {
		log.Println("Error inserting user:", err)
	}

	// Kullanıcıyı 'profile' tablosuna ekle
	insertProfileQuery := `
        INSERT INTO profile (user_id, username, email)
        VALUES ((SELECT id FROM users WHERE username = ?), ?, ?)
        ON CONFLICT(user_id) DO UPDATE SET
        username = excluded.username,
        email = excluded.email
    `
	_, err = db.Exec(insertProfileQuery, username, username, email)
	if err != nil {
		log.Println("Error inserting profile:", err)
	}
}
