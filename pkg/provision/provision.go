package provision

import (
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	provisions := GetAll()
	provisionsHTML := map[int]template.HTML{}
	for i, v := range provisions {
		provisionsHTML[i] = template.HTML(v.Content)
	}
	data["Provisions"] = provisionsHTML

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/provision.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)

}
