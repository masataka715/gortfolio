package provision

import (
	"gortfolio/database"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
)

type Provision struct {
	ID      int
	Number  string
	Content string
}

func (provision *Provision) Insert() {
	db := database.Open()
	db.Create(&provision)
	defer db.Close()
}

func GetOne(kanjiNumber string) Provision {
	db := database.Open()
	var provision Provision
	db.Where("number LIKE ?", "%"+kanjiNumber+"%").First(&provision)
	db.Close()
	return provision
}

func GetAll() []Provision {
	db := database.Open()
	var provisions []Provision
	db.Find(&provisions)
	db.Close()
	return provisions
}

func Seed(db *gorm.DB) {
	var p Provision
	if err := db.Where("id = ?", 1).First(&p).Error; err == nil {
		return
	}

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

	content_html := map[int]string{}
	doc.Find("div.Article").EachWithBreak(func(i int, s *goquery.Selection) bool {
		html2, _ := s.Html()
		content_html[i] = html2
		if strings.Contains(html2, "第二百五十九条、第二百六十一条及び前条の罪") {
			return false
		}
		return true
	})

	// データベースへ保存
	for i := 0; i < len(number); i++ {
		provision := &Provision{
			Number:  number[i],
			Content: content_html[i],
		}
		provision.Insert()
	}
}
