package models

import (
	"log"
	"path/filepath"
	"text/template"
)

var tmplCache = make(map[string]*template.Template)

func LoadTemplates() {
	templates := []string{"login", "register", "home", "profile", "panel", "category"}
	for _, tmpl := range templates {
		path := filepath.Join("..", "Front-end", "pages", tmpl+".html")
		tmplCache[tmpl] = template.Must(template.ParseFiles(path))
	}
	log.Println("Templates loaded successfully!")
}
