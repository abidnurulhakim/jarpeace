package model

import (
	"github.com/abidnurulhakim/carbon"
	"gopkg.in/mgo.v2/bson"
)

type AutoReply struct {
	ID        bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Text      string                 `bson:"text" json:"text"`
	Answer    string                 `bson:"answer" json:"answer"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Active    bool                   `bson:"active" json:"active"`
	ChatID    int                    `bson:"chat_id" json:"chat_id"`
	UserID    int                    `bson:"user_id" json:"user_id"`
	Username  string                 `bson:"username" json:"username"`
	UpdatedAt *carbon.Carbon         `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt *carbon.Carbon         `bson:"created_at,omitempty" json:"created_at"`
}

func NewAutoReply() AutoReply {
	autoReply := AutoReply{}
	autoReply.Active = true
	return autoReply
}
