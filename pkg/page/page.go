package page

import "gortfolio/database"

type Page struct {
	PageID   int
	PageName string
}

func Seed() {
	db := database.Open()
	var page Page
	if err := db.Where("page_id = ?", 1).First(&page).Error; err == nil {
		return
	}
	db.Create(&Page{PageID: 1, PageName: "ホーム"})
	db.Create(&Page{PageID: 2, PageName: "チャット"})
	db.Create(&Page{PageID: 3, PageName: "スクレイピング"})
	db.Create(&Page{PageID: 4, PageName: "しりとり"})
	db.Create(&Page{PageID: 5, PageName: "タスク"})
	db.Create(&Page{PageID: 6, PageName: "あしあと"})
	defer db.Close()
}
