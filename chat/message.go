package chat

import "go_chat/database"

type Message struct {
	ID        int
	Name      string
	Message   string
	When      string
	AvatarURL string
}

func MsgInsert(message *Message) {
	db := database.Open()
	db.Create(&message)
	defer db.Close()
}

func GetMsgAll() []Message {
	db := database.Open()
	var messages []Message
	db.Find(&messages)
	db.Close()
	return messages
}
