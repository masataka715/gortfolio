package chat

import (
	"gortfolio/database"
)

type Message struct {
	ID        int
	UserID    string
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

func UpdateAvatarURL(userID string, avatarURL string) {
	db := database.Open()
	msg := Message{}
	msg.AvatarURL = avatarURL
	db = db.Where("user_id = ?", userID)
	db.Model(&msg).Update(&msg)
	defer db.Close()
}

func GetMsgAll() []Message {
	db := database.Open()
	var messages []Message
	db.Find(&messages)
	db.Close()
	return messages
}
