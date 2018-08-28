package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Reminder struct {
	ID        bson.ObjectId          `mgorm:"primary_key" bson:"_id,omitempty" json:"id"`
	UpdatedAt time.Time              `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt time.Time              `mgorm:"index" bson:"created_at,omitempty" json:"created_at"`
	Content   string                 `bson:"content" json:"content"`
	Active    bool                   `bson:"active" json:"active"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Schedule  string                 `bson:"schedule" json:"schedule"`
	ChatID    int                    `bson:"chat_id" json:"chat_id"`
	UserID    int                    `bson:"user_id" json:"user_id"`
	Username  string                 `bson:"username" json:"username"`
}

func NewReminder() Reminder {
	var data map[string]interface{}
	reminder := Reminder{}
	reminder.Active = true
	reminder.Schedule = "0 * * * * *"
	reminder.Data = data
	return reminder
}
