package main

import (
	"flag"
	"gortfolio/config"
	"gortfolio/database"
	"gortfolio/pkg/auth"
	"gortfolio/pkg/chat"
	"gortfolio/pkg/flash"
	"gortfolio/pkg/footprint"
	"gortfolio/pkg/home"
	"gortfolio/pkg/page"
	"gortfolio/pkg/scraping"
	"gortfolio/pkg/shiritori"
	"gortfolio/pkg/todo"
	"gortfolio/trace"
	"gortfolio/utils"
	"time"

	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	footprint.Insert("チャット", when)

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	data["Msg"] = chat.GetMsgAll()
	AuthMessage, _ := flash.Get(w, r, "AuthMessage")
	data["AuthMessage"] = AuthMessage
	_ = t.templ.Execute(w, data)
}

func main() {
	utils.LoggingSettings("go.log")
	// データベース
	db := database.Open()
	db.AutoMigrate(chat.Message{})
	db.AutoMigrate(auth.User{})
	db.AutoMigrate(todo.Todo{})
	db.AutoMigrate(footprint.Footprint{})
	db.AutoMigrate(page.Page{})
	page.Seed()
	defer db.Close()

	var addr = flag.String("addr", ":5002", "アプリケーションのアドレス")
	flag.Parse()
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey(config.Config.GomniauthKey)
	gomniauth.WithProviders(
		google.New(config.Config.GoogleClientID, config.Config.GoogleSecretValue, "http://localhost:5002/auth/callback/google"),
	)

	r := chat.NewRoom()
	r.Tracer = trace.New(os.Stdout)

	http.HandleFunc("/", home.Handler)
	http.Handle("/images/",
		http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/shiritori", shiritori.Handler)
	http.HandleFunc("/scraping", scraping.Handler)
	http.HandleFunc("/footprint", footprint.Handler)
	http.HandleFunc("/download/footprint", footprint.DownloadHandler)
	http.HandleFunc("/todo", todo.Handler)
	http.HandleFunc("/todo/edit/", todo.EditHandler)
	http.HandleFunc("/todo/delete", todo.DeleteHandler)
	http.Handle("/chat", auth.MustAuth(&templateHandler{filename: "chat/chat.html"}))
	http.HandleFunc("/login", auth.LoginScreenHandler)
	http.HandleFunc("/login/form", auth.LoginFormHandler)
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/test/login", auth.TestLoginHandler)
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		flash.Set(w, "AuthMessage", []byte("ログアウトしました"))
		w.Header()["Location"] = []string{"/"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "chat/upload.html"})
	http.HandleFunc("/uploader", chat.UploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", http.FileServer(http.Dir("pkg/chat/avatars"))))
	http.Handle("/room", r)
	go r.Run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
