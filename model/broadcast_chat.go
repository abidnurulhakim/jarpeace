package model

import (
	"github.com/abidnurulhakim/carbon"
	"gopkg.in/mgo.v2/bson"
)

type BroadcastChat struct {
	ID        bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	Type      string         `bson:"type" json:"type"`
	ChatID    int            `bson:"chat_id" json:"chat_id"`
	UserID    int            `bson:"user_id" json:"user_id"`
	Username  string         `bson:"username" json:"username"`
	Active    bool           `bson:"active" json:"active"`
	UpdatedAt *carbon.Carbon `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt *carbon.Carbon `bson:"created_at,omitempty" json:"created_at"`
}

func NewLeaveBroadcastChat(chatID int, userID int, username string, activeStatus bool) BroadcastChat {
	return BroadcastChat{
		Type:     "leave",
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		Active:   activeStatus,
	}
}
