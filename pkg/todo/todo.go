package todo

import (
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/todo.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
