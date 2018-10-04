package model

import (
	"strconv"
	"strings"

	"github.com/abidnurulhakim/carbon"
	"github.com/bukalapak/zinc/helper"
	"gopkg.in/mgo.v2/bson"
)

type AutoReply struct {
	ID            bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Text          string                 `bson:"text" json:"text"`
	Answer        string                 `bson:"answer" json:"answer"`
	Active        bool                   `bson:"active" json:"active"`
	ChatIDs       []int                  `bson:"chat_ids" json:"chat_ids"`
	IdenticalText string                 `bson:"indentical_text,omitempty" json:"indentical_text"`
	Data          map[string]interface{} `bson:"data" json:"data"`
	UserID        int                    `bson:"user_id" json:"user_id"`
	Username      string                 `bson:"username,omitempty" json:"username"`
	UpdatedAt     *carbon.Carbon         `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt     *carbon.Carbon         `bson:"created_at,omitempty" json:"created_at"`
	DeletedAt     *carbon.Carbon         `bson:"deleted_at,omitempty" json:"deleted_at"`
}

func NewAutoReply(raw string, chatID int) AutoReply {
	data := make(map[string]interface{})
	autoReply := AutoReply{}
	if raw == "" {
		return autoReply
	}
	params := strings.Split(raw, ";")
	if len(params) > 0 {
		autoReply.Text = strings.TrimSpace(params[0])
	}
	if len(params) > 1 {
		autoReply.Answer = strings.TrimSpace(params[1])
	}
	if len(params) > 2 {
		if helper.Contains([]string{"true", "1"}, strings.ToLower(strings.TrimSpace(params[2]))) {
			autoReply.Active = true
		} else {
			autoReply.Active = false
		}
	}
	autoReply.ChatIDs = []int{chatID}
	if len(params) > 3 {
		chatIDs := strings.Split(raw, ";")
		for i := 0; i < len(chatIDs); i++ {
			if i, err := strconv.Atoi(strings.TrimSpace(chatIDs[i])); err == nil {
				autoReply.ChatIDs = append(autoReply.ChatIDs, i)
			}
		}
	}
	if len(params) > 4 {
		autoReply.IdenticalText = strings.TrimSpace(params[4])
	}
	for i := 5; i < len(params); i++ {
		data["var"+strconv.Itoa(i-2)] = strings.TrimSpace(params[i])
	}
	autoReply.Data = data
	return autoReply
}
