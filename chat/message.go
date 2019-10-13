package main

import "go_oreilly_app/chat/database"

type message struct {
	ID        int
	Name      string
	Message   string
	When      string
	AvatarURL string
}

func MsgInsert(message *message) {
	db := database.Open()
	db.Create(&message)
	defer db.Close()
}

func GetMsgAll() []message {
	db := database.Open()
	var messages []message
	db.Find(&messages)
	db.Close()
	return messages
}
