package auth

import (
	"crypto/md5"
	"fmt"
	"gortfolio/pkg/flash"
	"io"
	"log"
	"net/http"
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
		http.SetCookie(w, &http.Cookie{
			Name:  "redirectUrl",
			Value: r.URL.Path,
			Path:  "/",
		})
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		log.Println(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
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
		chatUser.uniqueID = GetUniqueID(user.Name())
		avatarURL, err := Avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("GetAvatarURLに失敗しました", "-", err)
		}
		SetAuthCookie(w, chatUser.uniqueID, user.Name(), avatarURL)
		flash.Set(w, "AuthMessage", []byte("Googleアカウントでログインしました"))

		cookie := GetRedirectCookie(w, r)
		w.Header()["Location"] = []string{cookie.Value}
		w.WriteHeader(http.StatusMovedPermanently)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}

func GetUniqueID(name string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(name))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func SetAuthCookie(w http.ResponseWriter, uniqueID string, name string, avatarURL string) {
	authCookieValue := objx.New(map[string]interface{}{
		"userid":     uniqueID,
		"name":       name,
		"avatar_url": avatarURL,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})
}
