package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *Message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *Message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.UserID = c.userData["userid"].(string)
			msg.Name = c.userData["name"].(string)
			msg.When = time.Now().Format("2006年01月02日 15時04分")
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
			MsgInsert(msg)
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
