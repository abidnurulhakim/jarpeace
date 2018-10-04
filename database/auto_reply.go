package database

import (
	"errors"
	"strings"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) CreateAutoReply(autoReply *model.AutoReply) error {
	var err error
	now := carbon.Now()
	autoReply.UpdatedAt = now
	autoReply.CreatedAt = now
	autoReply.Text = strings.ToLower(autoReply.Text)
	if autoReply.Text == "" {
		return errors.New("Precondition text auto reply cannot be empty.")
	}
	if autoReply.Answer == "" && autoReply.IdenticalText == "" {
		return errors.New("Answer auto reply cannot be empty.")
	}
	if len(autoReply.ChatIDs) == 0 {
		return errors.New("Chat ID cannot be empty.")
	}
	c := db.Collection("auto_replies")
	info, err := c.Upsert(bson.M{"user_id": autoReply.UserID, "created_at": autoReply.CreatedAt}, autoReply)
	if info != nil && info.UpsertedId != nil {
		autoReply.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) UpdateAutoReply(autoReply *model.AutoReply) error {
	var err error
	now := carbon.Now()
	autoReply.UpdatedAt = now
	if autoReply.ID.Hex() == "" {
		return errors.New("Must exists reminder")
	}
	if autoReply.Text == "" {
		return errors.New("Precondition text auto reply cannot be empty.")
	}
	if autoReply.Answer == "" {
		return errors.New("Answer auto reply cannot be empty.")
	}
	if len(autoReply.ChatIDs) == 0 {
		return errors.New("Chat ID cannot be empty.")
	}
	c := db.Collection("auto_replies")
	info, err := c.Upsert(bson.M{"_id": autoReply.ID}, autoReply)
	if info != nil && info.UpsertedId != nil {
		autoReply.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) GetAutoReplies(filter bson.M) ([]model.AutoReply, error) {
	var err error
	autoReplies := []model.AutoReply{}
	c := db.Collection("auto_replies")
	err = c.Find(filter).All(&autoReplies)
	return autoReplies, err
}
