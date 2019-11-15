package provision

import (
	"gortfolio/pkg/footprint"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	footprint.Insert("条文検索", when)

	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		number := r.FormValue("number")
		slice := strings.Split(number, "")
		kanjiNumber := ""
		switch len(slice) {
		case 1:
			kanjiNumber = GetKanjiNum(slice[0], "一")
		case 2:
			kanjiNumber = GetKanjiNum(slice[0], "十") + GetKanjiNum(slice[1], "一")
		case 3:
			kanjiNumber = GetKanjiNum(slice[0], "百") + GetKanjiNum(slice[1], "十") + GetKanjiNum(slice[2], "一")
		}
		var provision Provision
		if kanjiNumber == "" {
			data["Msg"] = "該当する条文番号はありません"
		} else {
			provision = GetOne(kanjiNumber)
			data["Result"] = template.HTML(provision.Content)
		}
		if provision.ID == 0 {
			data["Msg"] = "該当する条文番号はありません"
		}
	}

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

func GetKanjiNum(number string, place string) string {
	var kanjiNumber string
	switch number {
	case "1":
		if place == "一" {
			kanjiNumber = "一"
		} else {
			kanjiNumber = ""
		}
	case "2":
		kanjiNumber = "二"
	case "3":
		kanjiNumber = "三"
	case "4":
		kanjiNumber = "四"
	case "5":
		kanjiNumber = "五"
	case "6":
		kanjiNumber = "六"
	case "7":
		kanjiNumber = "七"
	case "8":
		kanjiNumber = "八"
	case "9":
		kanjiNumber = "九"
	default:
		kanjiNumber = ""
	}

	if place == "十" {
		kanjiNumber += "十"
	}
	if place == "百" {
		kanjiNumber += "百"
	}

	return kanjiNumber
}
