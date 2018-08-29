package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abidnurulhakim/jarpeace/model"
	"github.com/osteele/liquid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	cron "gopkg.in/robfig/cron.v2"
)

type MongoDB struct {
	Session  *mgo.Session
	Database string
}

// Option holds all necessary options for database.
type MongoDBOption struct {
	User     string
	Password string
	Host     string
	Database string
}

func NewMongoDB(opt *MongoDBOption) (*MongoDB, error) {
	mongoDialInfo := &mgo.DialInfo{
		Addrs:    []string{opt.Host},
		Database: opt.Database,
		Username: opt.User,
		Password: opt.Password,
		Timeout:  60 * time.Second,
	}
	session, err := mgo.DialWithInfo(mongoDialInfo)
	if err != nil {
		panic(err)
	}
	return &MongoDB{Session: session, Database: opt.Database}, err
}

func (db *MongoDB) Collection(collection string) *mgo.Collection {
	c := db.Session.DB(db.Database).C(collection)
	return c
}

func (db *MongoDB) Close() {
	db.Session.Close()
}

func (db *MongoDB) Copy() *MongoDB {
	return &MongoDB{Session: db.Session.Copy(), Database: db.Database}
}

func (db *MongoDB) CreateMessage(message *model.Message) error {
	var err error
	now := time.Now()
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

func (db *MongoDB) ProcessMessage(message *model.Message) (bool, string) {
	data := make(map[string]interface{})
	text := ""
	if message.Content != "" {
		arr := strings.SplitN(message.Content, " ", 2)
		arr2 := strings.SplitN(arr[0], "@"+os.Getenv("TELEGRAM_BOT_USERNAME"), 2)
		command := ""
		if len(arr2) > 1 {
			command = arr2[1]
		} else {
			command = arr[0]
		}
		if len(arr) > 1 {
			text = arr[1]
		}
		if string(command[0]) == "/" {
			arr3 := strings.SplitN(command, "/", 2)
			fmt.Println(arr3)
			if arr3[1] == "reminderadd" {
				if text == "" {
					return false, "Invalid format. Please check `/reminderadd help`"
				}
				params := strings.Split(text, ";")
				if len(params) == 1 {
					if params[0] != "help" {
						return false, "Invalid format. Please check `/reminderadd help`"
					} else {
						return true, "/reminderadd `SCHEDULE;REMINDER_CONTENT;var1;var2;...`\n`SCHEDULE: minute hour day month years` (cron format) (http://www.adminschoice.com/crontab-quick-reference)\n`REMINDER_CONTENT` can fill with liquid template (https://github.com/Shopify/liquid/wiki/Liquid-for-Designers)\n`var1`...`varN` is additional data that can use in `REMINDER CONTENT`"
					}
				}
				reminder := model.Reminder{}
				reminder.ChatID = message.ChatID
				reminder.UserID = message.UserID
				reminder.Username = message.Username
				reminder.Active = true
				reminder.Schedule = "0 " + strings.TrimSpace(params[0])
				reminder.Content = strings.TrimSpace(params[1])
				if len(params) > 2 {
					for i := 2; i < len(params); i++ {
						data["var"+strconv.Itoa(i-1)] = strings.TrimSpace(params[i])
					}
				}
				reminder.Data = data
				err := db.CreateReminder(&reminder)
				if err != nil {
					return false, err.Error()
				}
				return true, "Reminder with ID: `" + reminder.ID.Hex() + "` was added"
			}
			if arr3[1] == "reminderactive" {
				if text == "help" {
					return true, "/reminderactive `REMINDER_ID`\n`REMINDER_ID`: ID reminder was created"
				}
				if !bson.IsObjectIdHex(text) {
					return false, "Invalid format. Please check `/reminderactive help`"
				}
				reminders, err := db.GetReminders(bson.M{"_id": bson.ObjectIdHex(text), "chat_id": message.ChatID})
				if err != nil {
					return false, err.Error()
				}
				if len(reminders) == 0 {
					return false, "Reminder with ID: `" + text + "` not found"
				}
				reminder := reminders[0]
				reminder.Active = true
				err = db.UpdateReminder(&reminder)
				if err != nil {
					return false, err.Error()
				}
				return true, "Reminder with ID: `" + reminder.ID.Hex() + "` active"
			}
			if arr3[1] == "reminderinactive" {
				if text == "help" {
					return true, "/reminderactive `REMINDER_ID`\n`REMINDER_ID`: ID reminder was created"
				}
				if !bson.IsObjectIdHex(text) {
					return false, "Invalid format. Please check `/reminderinactive help`"
				}
				reminders, err := db.GetReminders(bson.M{"_id": bson.ObjectIdHex(text), "chat_id": message.ChatID})
				if err != nil {
					return false, err.Error()
				}
				if len(reminders) == 0 {
					return false, "Reminder with ID: `" + text + "` not found"
				}
				reminder := reminders[0]
				reminder.Active = false
				err = db.UpdateReminder(&reminder)
				if err != nil {
					return false, err.Error()
				}
				return true, "Reminder with ID: `" + reminder.ID.Hex() + "` inactive"
			}
			if arr3[1] == "reminderlist" {
				reminders, err := db.GetReminders(bson.M{"chat_id": message.ChatID})
				if err != nil {
					return false, err.Error()
				}
				if len(reminders) == 0 {
					return false, "There are no reminder."
				}
				str := ""
				for i := 0; i < len(reminders); i++ {
					str += strconv.Itoa(i+1) + ". ID: `" + reminders[i].ID.Hex() + "`\n    ACTIVE: `" + strconv.FormatBool(reminders[i].Active) + "`\n    Schedule: `" + reminders[i].Schedule + "`\n"
				}
				return true, str
			}
			if arr3[1] == "autoreplyadd" {
				if text == "" {
					return false, "Invalid format. Please check `/autoreplyadd help`"
				}
				params := strings.Split(text, ";")
				if len(params) == 1 {
					if params[0] != "help" {
						return false, "Invalid format. Please check `/autoreplyadd help`"
					} else {
						return true, "/autoreplyadd `TEXT;ANSWER;var1;var2;...`\n`TEXT`: text will autoreply by bot\n`ANSWER`: answer that bot will give. It can fill with liquid template (https://github.com/Shopify/liquid/wiki/Liquid-for-Designers)\n`var1`...`varN` is additional data that can use in `ANSWER`"
					}
				}
				autoReply := model.AutoReply{}
				autoReply.ChatID = message.ChatID
				autoReply.UserID = message.UserID
				autoReply.Username = message.Username
				autoReply.Active = true
				autoReply.Text = strings.TrimSpace(params[0])
				autoReply.Answer = strings.TrimSpace(params[1])
				if len(params) > 2 {
					for i := 2; i < len(params); i++ {
						data["var"+strconv.Itoa(i-1)] = strings.TrimSpace(params[i])
					}
				}
				autoReply.Data = data
				err := db.CreateAutoReply(&autoReply)
				if err != nil {
					return false, err.Error()
				}
				return true, "Auto reply with ID: `" + autoReply.ID.Hex() + "` was added"
			}
			if arr3[1] == "autoreplyactive" {
				if text == "help" {
					return true, "/autoreplyactive `AUTOREPLY_ID`\n`AUTOREPLY_ID`: ID auto reply was created"
				}
				if !bson.IsObjectIdHex(text) {
					return false, "Invalid format. Please check `/autoreplyactive help`"
				}
				autoReplies, err := db.GetAutoReplies(bson.M{"_id": bson.ObjectIdHex(text), "chat_id": message.ChatID})
				if err != nil {
					return false, err.Error()
				}
				if len(autoReplies) == 0 {
					return false, "Auto reply with ID: `" + text + "` not found"
				}
				autoReply := autoReplies[0]
				autoReply.Active = true
				err = db.UpdateAutoReply(&autoReply)
				if err != nil {
					return false, err.Error()
				}
				return true, "Auto reply with ID: `" + autoReply.ID.Hex() + "` active"
			}
			if arr3[1] == "autoreplyinactive" {
				if text == "help" {
					return true, "/autoreplyinactive `REMINDER_ID`\n`REMINDER_ID`: ID reminder was created"
				}
				if !bson.IsObjectIdHex(text) {
					return false, "Invalid format. Please check `/autoreplyinactive help`"
				}
				autoReplies, err := db.GetAutoReplies(bson.M{"_id": bson.ObjectIdHex(text), "chat_id": message.ChatID})
				if err != nil {
					return false, err.Error()
				}
				if len(autoReplies) == 0 {
					return false, "Auto reply with ID: `" + text + "` not found"
				}
				autoReply := autoReplies[0]
				autoReply.Active = false
				err = db.UpdateAutoReply(&autoReply)
				if err != nil {
					return false, err.Error()
				}
				return true, "Auto reply with ID: `" + autoReply.ID.Hex() + "` inactive"
			}
			if arr3[1] == "autoreplylist" {
				autoReplies, err := db.GetAutoReplies(bson.M{"chat_id": message.ChatID})
				if err != nil {
					return false, err.Error()
				}
				if len(autoReplies) == 0 {
					return false, "There are no auto reply."
				}
				str := ""
				for i := 0; i < len(autoReplies); i++ {
					str += strconv.Itoa(i+1) + ". ID: `" + autoReplies[i].ID.Hex() + "`\n    ACTIVE: `" + strconv.FormatBool(autoReplies[i].Active) + "`\n    TEXT: `" + autoReplies[i].Text + "`\n    ANSWER: `" + autoReplies[i].Answer + "`\n"
				}
				return true, str
			}
		} else {
			autoReplies, _ := db.GetAutoReplies(bson.M{"active": true, "chat_id": message.ChatID})
			if len(autoReplies) > 0 {
				str := ""
				for _, autoReply := range autoReplies {
					if strings.Contains(strings.ToLower(message.Content), autoReply.Text) {
						engine := liquid.NewEngine()
						template := autoReply.Answer
						bindings := autoReply.Data
						bindings["now"] = time.Now()
						bindings["username"] = message.Username
						out, errLiquid := engine.ParseAndRenderString(template, bindings)
						if errLiquid != nil {
							autoReply.Active = false
							db.UpdateAutoReply(&autoReply)
						} else {
							str += out + "\n"
						}
					}
				}
				return true, str
			}
		}
	}
	return true, ""
}

func (db *MongoDB) CreateReminder(reminder *model.Reminder) error {
	var err error
	now := time.Now()
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
	now := time.Now()
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

func (db *MongoDB) CreateAutoReply(autoReply *model.AutoReply) error {
	var err error
	now := time.Now()
	autoReply.UpdatedAt = now
	autoReply.CreatedAt = now
	autoReply.Text = strings.ToLower(autoReply.Text)
	if autoReply.Text == "" {
		return errors.New("Precondition text auto reply cannot be empty.")
	}
	if autoReply.Answer == "" {
		return errors.New("Answer auto reply cannot be empty.")
	}
	if autoReply.ChatID == 0 {
		return errors.New("Chat ID cannot be empty.")
	}
	c := db.Collection("auto_replies")
	info, err := c.Upsert(bson.M{"chat_id": autoReply.ChatID, "created_at": autoReply.CreatedAt}, autoReply)
	if info != nil && info.UpsertedId != nil {
		autoReply.ID = info.UpsertedId.(bson.ObjectId)
	}
	return err
}

func (db *MongoDB) UpdateAutoReply(autoReply *model.AutoReply) error {
	var err error
	now := time.Now()
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
	if autoReply.ChatID == 0 {
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
