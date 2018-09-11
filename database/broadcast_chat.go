package database

import (
	"errors"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/helper"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) CreateBroadcastChat(broadcastChat *model.BroadcastChat) error {
	var err error
	now := carbon.Now()
	broadcastChat.UpdatedAt = now
	broadcastChat.CreatedAt = now
	if !helper.Contains([]string{"leave"}, broadcastChat.Type) {
		return errors.New("type must in ('leave')")
	}
	if broadcastChat.ChatID == 0 {
		return errors.New("chat ID is required")
	}
	if broadcastChat.UserID == 0 {
		return errors.New("user ID is required")
	}
	if broadcastChat.Username == "" {
		return errors.New("username is required")
	}
	c := db.Collection("broadcast_chats")
	info, err := c.Upsert(bson.M{"user_id": broadcastChat.UserID, "chat_id": broadcastChat.ChatID}, broadcastChat)
	if info != nil && info.UpsertedId != nil {
		broadcastChat.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) UpdateBroadcastChat(broadcastChat *model.BroadcastChat) error {
	var err error
	now := carbon.Now()
	broadcastChat.UpdatedAt = now
	if helper.Contains([]string{"leave"}, broadcastChat.Type) {
		return errors.New("type must in ('leave')")
	}
	if broadcastChat.ChatID == 0 {
		return errors.New("chat ID is required")
	}
	if broadcastChat.UserID == 0 {
		return errors.New("user ID is required")
	}
	if broadcastChat.Username == "" {
		return errors.New("username is required")
	}
	c := db.Collection("broadcast_chats")
	info, err := c.Upsert(bson.M{"_id": broadcastChat.ID}, broadcastChat)
	if info != nil && info.UpsertedId != nil {
		broadcastChat.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) GetBroadcastChats(filter bson.M) ([]model.BroadcastChat, error) {
	var err error
	broadcastChats := []model.BroadcastChat{}
	c := db.Collection("broadcast_chats")
	err = c.Find(filter).Sort("created_at").All(&broadcastChats)
	return broadcastChats, err
}
