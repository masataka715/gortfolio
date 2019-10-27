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
	quiita_url    = "https://qiita.com"
	qiita_go_url  = "https://qiita.com/tags/go"
	qiita_gcp_url = "https://qiita.com/tags/gcp"
)

type QiitaTrend struct {
	Text string
	Link string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	footprint.Insert("スクレイピング", when)

	data := map[string]interface{}{}

	sliceGo := GetQiitaTrend(qiita_go_url)
	data["QiitaGo"] = sliceGo

	sliceGCP := GetQiitaTrend(qiita_gcp_url)
	data["QiitaGCP"] = sliceGCP

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/scraping.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func GetQiitaTrend(url string) []QiitaTrend {
	slice := make([]QiitaTrend, 5)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("goqueryのdocumentの取得に失敗しました")
	}
	doc.Find("a.tst-ArticleBody_title").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		slice[i] = QiitaTrend{
			Text: s.Text(),
			Link: quiita_url + href,
		}
	})
	return slice
}
