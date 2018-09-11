package model

import (
	"github.com/abidnurulhakim/carbon"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID             bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	ChatID         int            `bson:"chat_id" json:"chat_id"`
	UserID         int            `bson:"user_id" json:"user_id"`
	Username       string         `bson:"username" json:"username"`
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
