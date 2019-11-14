package page

import "github.com/jinzhu/gorm"

type Page struct {
	PageID   int
	PageName string
}

func Seed(db *gorm.DB) {
	var page Page
	if err := db.Where("page_id = ?", 1).First(&page).Error; err == nil {
		return
	}
	db.Create(&Page{PageID: 1, PageName: "ホーム"})
	db.Create(&Page{PageID: 2, PageName: "ブラックジャック"})
	db.Create(&Page{PageID: 3, PageName: "チャット"})
	db.Create(&Page{PageID: 4, PageName: "スクレイピング"})
	db.Create(&Page{PageID: 5, PageName: "しりとり"})
	db.Create(&Page{PageID: 6, PageName: "タスク"})
	db.Create(&Page{PageID: 7, PageName: "あしあと"})
}
