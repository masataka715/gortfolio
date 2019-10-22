package todo

import (
	"gortfolio/pkg/auth"
	"gortfolio/pkg/flash"
	"html/template"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	authData := auth.JudgeAuth(w, r)
	if authData == nil {
		return
	}
	userID := authData["userid"].(string)

	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")
		status := r.FormValue("status")
		if status != "" {
			Insert(userID, text, status)
		}
	}
	data := map[string]interface{}{}
	data["Todos"] = GetAll(userID)
	AuthMessage, _ := flash.Get(w, r, "AuthMessage")
	data["AuthMessage"] = AuthMessage

	templates := template.Must(template.ParseFiles("templates/layout.html", "templates/todo/todo.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	authData := auth.JudgeAuth(w, r)
	if authData == nil {
		return
	}
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
	authData := auth.JudgeAuth(w, r)
	if authData == nil {
		return
	}
	r.ParseForm()
	n := r.FormValue("id")
	id, _ := strconv.Atoi(n)
	Delete(id)
	w.Header()["Location"] = []string{"/todo"}
	w.WriteHeader(http.StatusMovedPermanently)
}
