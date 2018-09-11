package model

import (
	"github.com/abidnurulhakim/carbon"
	"gopkg.in/mgo.v2/bson"
)

type CallbackData struct {
	ID         bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Data       map[string]interface{} `bson:"data" json:"data"`
	State      string                 `bson:"state" json:"state"`
	PendingAt  *carbon.Carbon         `bson:"pending_at,omitempty" json:"pending_at"`
	FinishedAt *carbon.Carbon         `bson:"finished_at,omitempty" json:"finished_at"`
	UpdatedAt  *carbon.Carbon         `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt  *carbon.Carbon         `bson:"created_at,omitempty" json:"created_at"`
}

func NewCallbackData() CallbackData {
	var data map[string]interface{}
	callbackData := CallbackData{}
	callbackData.Data = data
	callbackData.State = "pending"
	return callbackData
}
