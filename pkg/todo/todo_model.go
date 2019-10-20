package todo

import (
	"gortfolio/database"
	"time"
)

type Todo struct {
	ID        int
	Text      string
	Status    string
	CreatedAt time.Time
}

func Insert(text string, status string) {
	db := database.Open()
	db.Create(&Todo{Text: text, Status: status})
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

func GetAll() []Todo {
	db := database.Open()
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
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
