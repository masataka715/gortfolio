package provision

import (
	"gortfolio/database"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	// data["Provisions"] = GetProvisions()

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/provision.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)

}

func GetProvisions() {
	url := "https://elaws.e-gov.go.jp/search/elawsSearch/elaws_search/lsg0500/detail?lawId=140AC0000000045"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("goqueryのdocumentの取得に失敗しました")
	}

	number := map[int]string{}
	doc.Find("span.ArticleTitle").EachWithBreak(func(i int, s *goquery.Selection) bool {
		html, _ := s.Html()
		number[i] = html
		if strings.Contains(html, "第二百六十四条") {
			return false
		}
		return true
	})

	content_html := map[int]template.HTML{}
	doc.Find("div.Article").EachWithBreak(func(i int, s *goquery.Selection) bool {
		html2, _ := s.Html()
		content_html[i] = template.HTML(html2)
		if strings.Contains(html2, "第二百五十九条、第二百六十一条及び前条の罪") {
			return false
		}
		return true
	})

	// データベースへ保存
	var p Provision
	db := database.Open()
	defer db.Close()
	if err := db.Where("id = ?", 1).First(&p).Error; err == nil {
		return
	}
	for i := 0; i < len(number); i++ {
		provision := &Provision{
			Number:      number[i],
			ContentHTML: content_html[i],
		}
		ProvisionInsert(provision)
	}
}
