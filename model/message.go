package model

import (
	"os"
	"strings"

	"github.com/abidnurulhakim/carbon"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID             bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	ChatID         int            `bson:"chat_id" json:"chat_id"`
	UserID         int            `bson:"user_id" json:"user_id"`
	Username       string         `bson:"username" json:"username"`
	FirstName      string         `bson:"first_name" json:"first_name"`
	LastName       string         `bson:"last_name" json:"last_name"`
	MessageID      int            `bson:"message_id" json:"message_id"`
	ReplyMessageID int            `bson:"reply_message_id" json:"reply_message_id"`
	Content        string         `bson:"content" json:"content"`
	Files          []string       `bson:"files" json:"files"`
	UpdatedAt      *carbon.Carbon `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt      *carbon.Carbon `bson:"created_at,omitempty" json:"created_at"`
}

func NewMessage() Message {
	message := Message{}
	return message
}

func (message Message) IsCommand() bool {
	if message.Content == "" {
		return false
	}
	if string(message.Content[0]) != "/" || message.Content == "/" {
		return false
	}
	return true
}

func (message Message) IsMentionBot() bool {
	return strings.Contains(message.Content, os.Getenv("TELEGRAM_USERNAME_BOT"))
}
