package auth

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func LoginScreenHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/login.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	if r.Method == http.MethodPost {
		r.ParseForm()
		user := User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		user = GetMatchingUser(user)
		if user.ID != 0 {
			uniqueID := GetUniqueID("名前未登録")
			SetAuthCookie(w, uniqueID, "名前未登録", "/avatars/default.png")

			w.Header()["Location"] = []string{"/chat"}
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/login_form.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	if r.Method == http.MethodPost {
		// データベース保存
		r.ParseForm()
		register_email := r.FormValue("register_email")
		register_password := r.FormValue("register_password")
		user := &User{
			Email:    register_email,
			Password: register_password,
		}
		UserInsert(user)

		uniqueID := GetUniqueID("名前未登録")
		file, _ := os.Open("pkg/chat/avatars/default.png")
		data, _ := ioutil.ReadAll(file)
		filename := filepath.Join("pkg/chat/avatars", uniqueID+".jpg")
		ioutil.WriteFile(filename, data, 0777)
		SetAuthCookie(w, uniqueID, "名前未登録", "/avatars/"+uniqueID+".jpg")
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusMovedPermanently)
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/register.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
