package main

import (
	"gortfolio/config"
	"gortfolio/database"
	"gortfolio/pkg/auth"
	"gortfolio/pkg/blackjack"
	"gortfolio/pkg/chat"
	"gortfolio/pkg/flash"
	"gortfolio/pkg/footprint"
	"gortfolio/pkg/home"
	"gortfolio/pkg/page"
	"gortfolio/pkg/provision"
	"gortfolio/pkg/scraping"
	"gortfolio/pkg/shiritori"
	"gortfolio/pkg/todo"
	"gortfolio/trace"
	"gortfolio/utils"
	"path/filepath"
	"time"

	"html/template"
	"log"
	"net/http"
	"os"
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
	db.AutoMigrate(provision.Provision{})
	db.AutoMigrate(blackjack.Blackjack{})
	db.AutoMigrate(chat.Message{})
	db.AutoMigrate(auth.User{})
	db.AutoMigrate(todo.Todo{})
	db.AutoMigrate(footprint.Footprint{})
	db.AutoMigrate(page.Page{})
	provision.Seed(db)
	page.Seed(db)
	defer db.Close()

	gomniauth.SetSecurityKey(config.Config.GomniauthKey)
	gomniauth.WithProviders(
		google.New(config.Config.GoogleClientID, config.Config.GoogleSecretValue, config.Config.AppURL+"/auth/callback/google"),
	)

	r := chat.NewRoom()
	r.Tracer = trace.New(os.Stdout)

	http.HandleFunc("/", home.Handler)
	http.Handle("/images/",
		http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/provision", provision.Handler)
	http.HandleFunc("/blackjack/insert", blackjack.InsertHandler)
	http.HandleFunc("/blackjack", blackjack.Handler)
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
	http.HandleFunc("/rename", chat.RenameHandler)
	http.HandleFunc("/upload", chat.UploadHandler)
	http.HandleFunc("/uploader", chat.UploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", http.FileServer(http.Dir("pkg/chat/avatars"))))
	http.Handle("/room", r)
	go r.Run()

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
