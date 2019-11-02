package home

import (
	"gortfolio/pkg/flash"
	"gortfolio/pkg/footprint"
	"html/template"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		when := time.Now().Format("2006年01月02日 15時04分")
		footprint.Insert("ホーム", when)
		MakeAccessGraph()
		MakeQRcode()
	}

	data := map[string]interface{}{}
	data["Weather"] = GetWeather()
	AuthMessage, _ := flash.Get(w, r, "AuthMessage")
	data["AuthMessage"] = AuthMessage

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/home.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
