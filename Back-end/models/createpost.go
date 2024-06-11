package models

import "net/http"

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := tmplCache["create_post"]
	if !ok {
		http.Error(w, "Could not load create_post template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
