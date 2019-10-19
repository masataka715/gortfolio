package auth

import "gortfolio/database"

type User struct {
	ID       int
	Email    string
	Password string
}

func UserInsert(user *User) {
	db := database.Open()
	db.Create(&user)
	defer db.Close()
}
