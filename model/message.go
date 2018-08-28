package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID             bson.ObjectId `mgorm:"primary_key" bson:"_id,omitempty" json:"id"`
	UpdatedAt      time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt      time.Time     `mgorm:"index" bson:"created_at,omitempty" json:"created_at"`
	ChatID         int           `bson:"chat_id" json:"chat_id"`
	UserID         int           `bson:"user_id" json:"user_id"`
	Username       string        `bson:"username" json:"username"`
	MessageID      int           `bson:"message_id" json:"message_id"`
	ReplyMessageID int           `bson:"reply_message_id" json:"reply_message_id"`
	Content        string        `bson:"content" json:"content"`
	Files          []string      `bson:"files" json:"files"`
}

func NewMessage() Message {
	message := Message{}
	return message
}
