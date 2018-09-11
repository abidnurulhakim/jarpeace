package model

import (
	"strconv"
	"strings"

	"github.com/abidnurulhakim/carbon"
	"gopkg.in/mgo.v2/bson"
)

type Reminder struct {
	ID        bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Title     string                 `bson:"title" json:"title"`
	Content   string                 `bson:"content" json:"content"`
	Active    bool                   `bson:"active" json:"active"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Schedule  string                 `bson:"schedule" json:"schedule"`
	ChatID    int                    `bson:"chat_id" json:"chat_id"`
	UserID    int                    `bson:"user_id" json:"user_id"`
	Username  string                 `bson:"username" json:"username"`
	UpdatedAt *carbon.Carbon         `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt *carbon.Carbon         `bson:"created_at,omitempty" json:"created_at"`
}

func NewReminder(raw string) Reminder {
	data := make(map[string]interface{})
	reminder := Reminder{}
	if raw == "" {
		return reminder
	}
	params := strings.Split(raw, ";")
	if len(params) > 0 {
		reminder.Title = strings.TrimSpace(params[0])
	}
	if len(params) > 1 {
		reminder.Schedule = "0 " + strings.TrimSpace(params[1])
	}
	if len(params) > 2 {
		reminder.Content = strings.TrimSpace(params[2])
	}
	if len(params) > 3 {
		for i := 3; i < len(params); i++ {
			data["var"+strconv.Itoa(i-2)] = strings.TrimSpace(params[i])
		}
	}
	reminder.Data = data
	return reminder
}
