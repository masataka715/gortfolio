package auth

import (
	"gortfolio/pkg/flash"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	if r.Method == http.MethodPost {
		r.ParseForm()
		register_email := r.FormValue("register_email")
		register_password := r.FormValue("register_password")
		if utf8.RuneCountInString(register_password) <= 8 {
			data["ErrMessage"] = "パスワードの文字数が足りません"
			templates := template.Must(template.ParseFiles("templates/layout.html",
				"templates/auth/register.html"))
			_ = templates.ExecuteTemplate(w, "layout", data)
			return
		}
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
		flash.Set(w, "AuthMessage", []byte("登録されました"))

		cookie := GetRedirectCookie(w, r)
		w.Header()["Location"] = []string{cookie.Value}
		w.WriteHeader(http.StatusMovedPermanently)
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/register.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
