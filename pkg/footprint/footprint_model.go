package footprint

import "gortfolio/database"

type Footprint struct {
	ID       int
	ViewPage string
	When     string
}

func Insert(viewPage string, when string) {
	db := database.Open()
	db.Create(&Footprint{ViewPage: viewPage, When: when})
	defer db.Close()
}

func GetAll() []Footprint {
	db := database.Open()
	var footprints []Footprint
	db.Order("id desc").Limit(20).Find(&footprints)
	db.Close()
	return footprints
}
