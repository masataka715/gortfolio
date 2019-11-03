package chat

import (
	"gortfolio/pkg/auth"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/stretchr/objx"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	authData := auth.JudgeAuth(w, r)
	if authData == nil {
		return
	}

	data := map[string]interface{}{}
	data["UserData"] = authData

	templates := template.Must(template.ParseFiles("templates/layout.html", "templates/chat/upload.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func UploaderHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join("pkg/chat/avatars", userId+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	authCookie, err := r.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
		return
	}
	// authCookieのavatar_url更新
	userData := objx.MustFromBase64(authCookie.Value)
	userData["avatar_url"] = "/avatars/" + userId + filepath.Ext(header.Filename)
	auth.SetAuthCookie(w, userData["userid"].(string), userData["name"].(string), userData["avatar_url"].(string))
	// データベースのavatar_url更新
	UpdateAvatarURL(userId, "/avatars/"+userId+filepath.Ext(header.Filename))

	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusMovedPermanently)
}
