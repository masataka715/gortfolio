package footprint

import (
	"gortfolio/database"
)

type Footprint struct {
	ID       int
	PageName string
	When     string
}

type Count struct {
	PageID   int
	PageName string
	Count    int
}

func Insert(viewPage string, when string) {
	db := database.Open()
	db.Create(&Footprint{PageName: viewPage, When: when})
	defer db.Close()
}

func GetAll() []Footprint {
	db := database.Open()
	var footprints []Footprint
	db.Order("id desc").Limit(20).Find(&footprints)
	db.Close()
	return footprints
}

func GetCount() []Count {
	db := database.Open()
	var count []Count
	db.Table("footprints").Select("page_id, footprints.page_name, count(*) as count").Group("footprints.page_name").Joins("inner join pages on pages.page_name = footprints.page_name").Order("page_id").Find(&count)
	return count
}
