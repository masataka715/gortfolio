package scraping

import (
	"gortfolio/pkg/footprint"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	quiita_url   = "https://qiita.com"
	qiita_go_url = "https://qiita.com/tags/go"
)

type Qiita struct {
	Text string
	Link string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	footprint.Insert("スクレイピング", when)

	data := map[string]interface{}{}
	slice := make([]Qiita, 5)

	doc, err := goquery.NewDocument(qiita_go_url)
	if err != nil {
		log.Println("goqueryのdocumentの取得に失敗しました")
	}
	doc.Find("a.tst-ArticleBody_title").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		slice[i] = Qiita{
			Text: s.Text(),
			Link: quiita_url + href,
		}
	})
	data["Qiita"] = slice

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/scraping.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}
