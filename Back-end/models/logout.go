package models

import "net/http"

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
