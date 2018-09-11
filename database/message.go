package database

import (
	"errors"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) CreateMessage(message *model.Message) error {
	var err error
	now := carbon.Now()
	message.UpdatedAt = now
	message.CreatedAt = now
	if message.ChatID == 0 {
		return errors.New("Chat ID cannot be empty.")
	}
	if message.MessageID == 0 {
		return errors.New("Message ID cannot be empty.")
	}
	if message.Content == "" && len(message.Files) == 0 {
		return errors.New("Message content or files cannot be empty.")
	}
	c := db.Collection("messages")
	info, err := c.Upsert(bson.M{"chat_id": message.ChatID, "message_id": message.MessageID}, message)
	if info != nil && info.UpsertedId != nil {
		message.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}
