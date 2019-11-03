package chat

import (
	"gortfolio/pkg/auth"
	"net/http"
	"text/template"
)

func RenameHandler(w http.ResponseWriter, r *http.Request) {
	authData := auth.JudgeAuth(w, r)
	if authData == nil {
		return
	}
	userID := authData["userid"].(string)

	if r.Method == http.MethodPost {
		r.ParseForm()
		newName := r.FormValue("newName")
		testUserID := auth.GetUniqueID("testUserAWByyToQBh")
		if userID != testUserID {
			// authCookieのname更新
			auth.SetAuthCookie(w, userID, newName, authData["avatar_url"].(string))
			// データベースのname更新
			UpdateName(userID, newName)
			w.Header()["Location"] = []string{"/chat"}
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}

	data := map[string]interface{}{}
	data["UserData"] = authData

	templates := template.Must(template.ParseFiles("templates/layout.html", "templates/chat/rename.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
