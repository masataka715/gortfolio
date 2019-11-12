package blackjack

import (
	"gortfolio/database"
	"net/http"
	"time"
)

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	blackjack := &Blackjack{
		Name:   r.FormValue("name"),
		Result: r.FormValue("result"),
		When:   time.Now().Format("2006年01月02日 15時04分"),
	}
	BlackjackInsert(blackjack)
	w.Header()["Location"] = []string{"/blackjack"}
	w.WriteHeader(http.StatusMovedPermanently)
}

type Blackjack struct {
	ID     int
	Name   string
	Result string
	When   string
}

func BlackjackInsert(blackjack *Blackjack) {
	db := database.Open()
	db.Create(&blackjack)
	defer db.Close()
}

func GetAll() []Blackjack {
	db := database.Open()
	var blackjacks []Blackjack
	db.Order("id desc").Find(&blackjacks)
	db.Close()
	return blackjacks
}
