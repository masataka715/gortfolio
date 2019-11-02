package chat

import (
	"gortfolio/pkg/auth"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/stretchr/objx"
)

func UploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
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
	authCookie, err := req.Cookie("auth")
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
