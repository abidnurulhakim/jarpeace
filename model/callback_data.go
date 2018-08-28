package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CallbackData struct {
	ID        bson.ObjectId          `mgorm:"primary_key" bson:"_id,omitempty" json:"id"`
	UpdatedAt time.Time              `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt time.Time              `mgorm:"index" bson:"created_at,omitempty" json:"created_at"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	State     string                 `bson:"state" json:"state"`
	PendingAt time.Time              `bson:"pending_at" json:"pending_at"`
}

func NewCallbackData() CallbackData {
	var data map[string]interface{}
	callbackData := CallbackData{}
	callbackData.Data = data
	callbackData.State = "pending"
	return callbackData
}
