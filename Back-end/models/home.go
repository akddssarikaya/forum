package models

import "net/http"

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["home"]
	if !ok {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}

	// Kullanıcı oturum bilgisine göre profil ve çıkış bağlantısını ekleyelim
	_, err := r.Cookie("user_id")
	if err == nil {
		tmpl.Execute(w, map[string]interface{}{
			"LoggedIn": true,
		})
	} else {
		tmpl.Execute(w, map[string]interface{}{
			"LoggedIn": false,
		})
	}
}
