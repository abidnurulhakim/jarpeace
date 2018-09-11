package database

import (
	"errors"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
	cron "gopkg.in/robfig/cron.v2"
)

func (db *MongoDB) CreateReminder(reminder *model.Reminder) error {
	var err error
	now := carbon.Now()
	reminder.UpdatedAt = now
	reminder.CreatedAt = now
	if reminder.Content == "" {
		return errors.New("Content reminder cannot be empty.")
	}
	if _, err = cron.Parse(reminder.Schedule); err != nil {
		return errors.New("Schedule invalid.")
	}
	if reminder.ChatID == 0 {
		return errors.New("Chat ID cannot be empty.")
	}
	c := db.Collection("reminders")
	info, err := c.Upsert(bson.M{"chat_id": reminder.ChatID, "created_at": reminder.CreatedAt}, reminder)
	if info != nil && info.UpsertedId != nil {
		reminder.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) UpdateReminder(reminder *model.Reminder) error {
	var err error
	now := carbon.Now()
	reminder.UpdatedAt = now
	if reminder.ID.Hex() == "" {
		return errors.New("Must exists reminder")
	}
	if reminder.Content == "" {
		return errors.New("Content reminder cannot be empty.")
	}
	if reminder.ChatID == 0 {
		return errors.New("Chat ID cannot be empty.")
	}
	c := db.Collection("reminders")
	info, err := c.Upsert(bson.M{"_id": reminder.ID}, reminder)
	if info != nil && info.UpsertedId != nil {
		reminder.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) GetReminders(filter bson.M) ([]model.Reminder, error) {
	var err error
	reminders := []model.Reminder{}
	c := db.Collection("reminders")
	err = c.Find(filter).All(&reminders)
	return reminders, err
}
