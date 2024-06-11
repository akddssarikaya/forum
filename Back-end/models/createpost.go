package models

import "net/http"

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["create_post"]
	if !ok {
		http.Error(w, "Could not load create_post template", http.StatusInternalServerError)
		return
	}

	// Kullanıcı giriş yapmış mı diye kontrol edelim
	_, err := r.Cookie("user_id")
	loggedIn := err == nil

	// Şablonu render et
	data := map[string]interface{}{
		"LoggedIn": loggedIn,
	}

	// Kullanıcı giriş yapmamışsa Login ve Register bağlantılarını ekle
	if !loggedIn {
		data["ShowLoginRegister"] = true
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		return
	}
}
