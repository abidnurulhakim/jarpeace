package database

import (
	"errors"
	"strings"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/helper"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) CreateLeave(leave *model.Leave) error {
	now := carbon.Now()
	leave.UpdatedAt = now
	leave.CreatedAt = now

	if leave.UserID == 0 {
		return errors.New("User ID must exists")
	}
	leave.Type = strings.ToLower(leave.Type)
	if !helper.Contains([]string{"remote", "cuti", "sick"}, leave.Type) {
		return errors.New("Invalid leave type. Please insert valid value (remote, cuti, sick)")
	}
	leave.Start = leave.Start.StartOfDay()
	if leave.Start.Before(carbon.Now().StartOfDay().Time) {
		return errors.New("Start time must after today")
	}
	leave.End = leave.End.EndOfDay()
	if leave.End.Before(leave.Start.Time) {
		return errors.New("End time must after start time")
	}
	c := db.Collection("leaves")
	oldLeaves, _ := db.GetLeaves(bson.M{"user_id": leave.UserID, "start": bson.M{"$lte": leave.Start}, "end": bson.M{"$gte": leave.Start}, "deleted_at": bson.M{"$exists": true}})
	if len(oldLeaves) > 0 {
		return errors.New("There are collision in your new leave")
	}
	info, err := c.Upsert(bson.M{"user_id": leave.UserID, "start": leave.Start}, leave)
	if info != nil && info.UpsertedId != nil {
		leave.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) UpdateLeave(leave *model.Leave) error {
	leave.UpdatedAt = carbon.Now()
	if leave.UserID == 0 {
		return errors.New("User ID must exists")
	}
	leave.Type = strings.ToLower(leave.Type)
	if !helper.Contains([]string{"remote", "cuti", "sick"}, leave.Type) {
		return errors.New("Invalid leave type. Please insert valid value (remote, cuti, sick)")
	}
	leave.Start = leave.Start.StartOfDay()
	if leave.Start.Before(carbon.Now().StartOfDay().Time) {
		return errors.New("Start time must after today")
	}
	leave.End = leave.End.EndOfDay()
	if leave.End.Before(leave.Start.Time) {
		return errors.New("End time must after start time")
	}
	c := db.Collection("leaves")
	oldLeaves, _ := db.GetLeaves(bson.M{"user_id": leave.UserID, "start": bson.M{"$lte": leave.Start}, "end": bson.M{"$gte": leave.Start}, "deleted_at": bson.M{"$exists": true}, "_id": bson.M{"$ne": leave.Id}})
	if len(oldLeaves) > 0 {
		return errors.New("There are collision in your new leave")
	}
	info, err := c.Upsert(bson.M{"_id": leave.ID}, leave)
	if info != nil && info.UpsertedId != nil {
		leave.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) GetLeaves(filter bson.M) ([]model.Leave, error) {
	var err error
	leaves := []model.Leave{}
	c := db.Collection("leaves")
	err = c.Find(filter).All(&leaves)
	return leaves, err
}
