package auth

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User //型の埋め込み
	uniqueID             string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func LoginScreenHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/login.html"))
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

		m := md5.New()
		io.WriteString(m, strings.ToLower("名前未登録"))
		uniqueID := fmt.Sprintf("%x", m.Sum(nil))
		file, _ := os.Open("pkg/chat/avatars/default.png")
		data, _ := ioutil.ReadAll(file)
		filename := filepath.Join("pkg/chat/avatars", uniqueID+".jpg")
		ioutil.WriteFile(filename, data, 0777)
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     uniqueID,
			"name":       "名前未登録",
			"avatar_url": "/avatars/" + uniqueID + ".jpg",
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusMovedPermanently)
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/auth/register.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

// LoginHandlerはサードパーティーへのログインの処理を受け持ちます。
// パスの形式: /auth/{action}/{provider}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました:", provider, "-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました", provider, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした", provider, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーの取得に失敗しました", provider, "-", err)
		}
		chatUser := &chatUser{User: user}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := Avatars.GetAvatarURL(chatUser)
		log.Println(avatarURL)
		if err != nil {
			log.Fatalln("GetAvatarURLに失敗しました", "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
