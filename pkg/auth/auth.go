package auth

import (
	"gortfolio/pkg/flash"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stretchr/objx"
)

func JudgeAuth(w http.ResponseWriter, r *http.Request) objx.Map {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		http.SetCookie(w, &http.Cookie{
			Name:  "redirectUrl",
			Value: r.URL.Path,
			Path:  "/",
		})
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return nil
	} else if err != nil {
		log.Println(err.Error())
		return nil
	} else {
		authData, _ := objx.FromBase64(cookie.Value)
		return authData
	}
}

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
			flash.Set(w, "AuthMessage", []byte("ログインしました"))

			cookie := GetRedirectCookie(w, r)
			w.Header()["Location"] = []string{cookie.Value}
			w.WriteHeader(http.StatusMovedPermanently)
			return
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
		flash.Set(w, "AuthMessage", []byte("登録されました"))

		cookie := GetRedirectCookie(w, r)
		w.Header()["Location"] = []string{cookie.Value}
		w.WriteHeader(http.StatusMovedPermanently)
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/register.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func TestLoginHandler(w http.ResponseWriter, r *http.Request) {
	uniqueID := GetUniqueID("テストユーザー")
	file, _ := os.Open("pkg/chat/avatars/default.png")
	data, _ := ioutil.ReadAll(file)
	filename := filepath.Join("pkg/chat/avatars", uniqueID+".jpg")
	ioutil.WriteFile(filename, data, 0777)
	SetAuthCookie(w, uniqueID, "テストユーザー", "/avatars/"+uniqueID+".jpg")
	flash.Set(w, "AuthMessage", []byte("テストユーザーでログインしました"))

	cookie := GetRedirectCookie(w, r)
	w.Header()["Location"] = []string{cookie.Value}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GetRedirectCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, _ := r.Cookie("redirectUrl")
	http.SetCookie(w, &http.Cookie{
		Name:  "redirectUrl",
		Value: "",
		Path:  "/",
	})
	return cookie
}
