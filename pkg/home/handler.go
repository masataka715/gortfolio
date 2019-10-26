package home

import (
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	data["Weather"] = GetWeather()

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/home.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
