package handlers

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/home.html"))
	_ = templates.ExecuteTemplate(w, "layout", nil)
}
