package models

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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

	userID, err := storeUserInDB(user.Login, user.Email, "")
	if err != nil {
		http.Error(w, "Could not store user in database", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.FormatInt(userID, 10),
		Path:  "/",
	})

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

	userID, err := storeUserInDB(user.Name, user.Email, "")
	if err != nil {
		http.Error(w, "Could not store user in database", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.FormatInt(userID, 10),
		Path:  "/",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func storeUserInDB(username, email, password string) (int64, error) {
	db, err := sql.Open("sqlite3", "./Back-end/database/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if userID == 0 {
		result, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			return 0, err
		}
		userID, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		_, err := db.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", username, password, userID)
		if err != nil {
			return 0, err
		}
	}

	_, err = db.Exec("INSERT INTO profile (user_id, username, email) VALUES (?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET username = excluded.username, email = excluded.email", userID, username, email)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
