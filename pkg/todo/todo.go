package todo

import (
	"html/template"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")
		status := r.FormValue("status")
		Insert(text, status)
	}

	data := map[string]interface{}{}
	data["Todos"] = GetAll()

	templates := template.Must(template.ParseFiles("templates/layout.html", "templates/todo/todo.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	n := r.URL.Path[len("/todo/edit/"):]
	id, _ := strconv.Atoi(n)

	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")
		status := r.FormValue("status")
		Update(id, text, status)
		w.Header()["Location"] = []string{"/todo"}
		w.WriteHeader(http.StatusMovedPermanently)
	}

	data := map[string]interface{}{}
	data["Todo"] = GetOne(id)

	templates := template.Must(template.ParseFiles("templates/layout.html", "templates/todo/edit.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	n := r.FormValue("id")
	id, _ := strconv.Atoi(n)
	Delete(id)
	w.Header()["Location"] = []string{"/todo"}
	w.WriteHeader(http.StatusMovedPermanently)
}
