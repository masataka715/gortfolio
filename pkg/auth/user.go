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

func GetMatchingUser(user User) User {
	db := database.Open()
	newUser := User{}
	db.Where("email = ?", user.Email).Where("password = ?", user.Password).First(&newUser)
	db.Close()
	return newUser
}
