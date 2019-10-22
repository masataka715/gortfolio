package todo

import (
	"gortfolio/database"
	"time"
)

type Todo struct {
	ID        int
	UserID    string
	Text      string
	Status    string
	CreatedAt time.Time
}

func Insert(userID string, text string, status string) {
	db := database.Open()
	db.Create(&Todo{UserID: userID, Text: text, Status: status})
	defer db.Close()
}

func Update(id int, text string, status string) {
	db := database.Open()
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func Delete(id int) {
	db := database.Open()
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func GetAll(userID string) []Todo {
	db := database.Open()
	var todos []Todo
	db.Where("user_id = ?", userID).Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func GetOne(id int) Todo {
	db := database.Open()
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}
