package model

type Reminder struct {
	Model
	Name     string                 `bson:"name" json:"name"`
	Content  string                 `bson:"content" json:"content"`
	Active   bool                   `bson:"active" json:"active"`
	Data     map[string]interface{} `bson:"data" json:"data"`
	Schedule string                 `bson:"schedule" json:"schedule"`
	ChatID   int                    `bson:"chat_id" json:"chat_id"`
	UserID   int                    `bson:"user_id" json:"user_id"`
}

func NewReminder() Reminder {
	var data map[string]interface{}
	reminder := Reminder{}
	reminder.Active = true
	reminder.Schedule = "0 * * * * *"
	reminder.Data = data
	return reminder
}
