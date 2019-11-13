package provision

import (
	"gortfolio/database"
	"html/template"
)

type Provision struct {
	ID          int
	Number      string
	ContentHTML template.HTML
}

func ProvisionInsert(provision *Provision) {
	db := database.Open()
	db.Create(&provision)
	defer db.Close()
}

func GetOne(id int) Provision {
	db := database.Open()
	var provision Provision
	db.First(&provision, id)
	db.Close()
	return provision
}
