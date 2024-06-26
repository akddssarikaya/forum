package models

import (
	"log"
	"path/filepath"
	"text/template"
)

var tmplCache = make(map[string]*template.Template)

func LoadTemplates() {
	templates := []string{"login", "register", "home", "profile", "panel", "category", "create_post", "view_post"}

	for _, tmpl := range templates {
		path := filepath.Join(".", "Front-end", tmpl+".html")
		tmplContent, err := template.ParseFiles(path)
		if err != nil {
			log.Fatalf("Error parsing template %s: %v", tmpl, err)
		}
		tmplCache[tmpl] = tmplContent
	}

	log.Println("Templates loaded successfully!")
}
