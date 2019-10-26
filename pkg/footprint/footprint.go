package footprint

import (
	"html/template"
	"net/http"
	"time"

	"github.com/stretchr/objx"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	Insert("あしあと", when)

	data := map[string]interface{}{}
	data["Footprint"] = GetAll()
	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/footprint.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func SetCookie(w http.ResponseWriter, uniqueID string, name string, avatarURL string) {
	FootprintCookie := objx.New(map[string]interface{}{
		"view_page": uniqueID,
		"when":      name,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "footprint",
		Value: FootprintCookie,
		Path:  "/",
	})
}
