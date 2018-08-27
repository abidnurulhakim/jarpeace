package model

import "time"

type CallbackData struct {
	Model
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
