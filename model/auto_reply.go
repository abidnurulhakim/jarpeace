package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AutoReply struct {
	ID        bson.ObjectId          `mgorm:"primary_key" bson:"_id,omitempty" json:"id"`
	UpdatedAt time.Time              `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt time.Time              `mgorm:"index" bson:"created_at,omitempty" json:"created_at"`
	Text      string                 `bson:"text" json:"text"`
	Answer    string                 `bson:"answer" json:"answer"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Active    bool                   `bson:"active" json:"active"`
	ChatID    int                    `bson:"chat_id" json:"chat_id"`
	UserID    int                    `bson:"user_id" json:"user_id"`
	Username  string                 `bson:"username" json:"username"`
}

func NewAutoReply() AutoReply {
	autoReply := AutoReply{}
	return autoReply
}
