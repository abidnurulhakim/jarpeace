package command

import (
	"strconv"

	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (cmd *Command) RunRouteReminder() ([]string, error) {
	if cmd.Action == "active" {
		if cmd.Content == "help" {
			return []string{cmd.GetReminderHelpActive()}, nil
		}
		if !bson.IsObjectIdHex(cmd.Content) {
			return []string{"Sorry, invalid format. Please check `/reminder active help`"}, nil
		}
		reminders, err := cmd.Client.GetReminders(bson.M{"_id": bson.ObjectIdHex(cmd.Content), "chat_id": cmd.Message.ChatID})
		if err != nil {
			return []string{}, err
		}
		if len(reminders) == 0 {
			return []string{"Reminder with ID: `" + cmd.Content + "` not found"}, nil
		}
		reminder := reminders[0]
		reminder.Active = true
		err = cmd.Client.UpdateReminder(&reminder)
		if err != nil {
			return []string{}, err
		}
		return []string{"Reminder with ID: `" + reminder.ID.Hex() + "` was active"}, nil

	} else if cmd.Action == "add" {
		if cmd.Content == "help" {
			return []string{cmd.GetReminderHelpAdd()}, nil
		}
		reminder := model.NewReminder(cmd.Content)
		reminder.ChatID = cmd.Message.ChatID
		reminder.UserID = cmd.Message.UserID
		reminder.Username = cmd.Message.Username
		if reminder.Title == "" || reminder.Schedule == "0 " || reminder.Content == "" {
			return []string{"Sorry, invalid format. Please check `/reminder add help`"}, nil
		}
		reminders, err := cmd.Client.GetReminders(bson.M{"title": reminder.Title, "chat_id": cmd.Message.ChatID})
		if err != nil {
			return []string{}, err
		}
		if len(reminders) > 0 {
			return []string{"Sorry, title was taken. Please use different title"}, err
		}
		reminder.Active = true
		err = cmd.Client.CreateReminder(&reminder)
		if err != nil {
			return []string{}, err
		}
		return []string{"Reminder with ID: `" + reminder.ID.Hex() + "` was added"}, nil
	} else if cmd.Action == "help" {
		return []string{
			"List available action for /reminder :\n" + "1. " + cmd.GetReminderHelpAdd(),
			"2. " + cmd.GetReminderHelpActive(),
			"3. " + cmd.GetReminderHelpInactive(),
			"4. " + cmd.GetReminderHelpList(),
		}, nil
	} else if cmd.Action == "inactive" {
		if cmd.Content == "help" {
			return []string{cmd.GetReminderHelpInactive()}, nil
		}
		if !bson.IsObjectIdHex(cmd.Content) {
			return []string{"Sorry, invalid format. Please check `/reminder inactive help`"}, nil
		}
		reminders, err := cmd.Client.GetReminders(bson.M{"_id": bson.ObjectIdHex(cmd.Content), "chat_id": cmd.Message.ChatID})
		if err != nil {
			return []string{}, err
		}
		if len(reminders) == 0 {
			return []string{"Reminder with ID: `" + cmd.Content + "` not found"}, nil
		}
		reminder := reminders[0]
		reminder.Active = false
		err = cmd.Client.UpdateReminder(&reminder)
		if err != nil {
			return []string{}, err
		}
		return []string{"Reminder with ID: `" + reminder.ID.Hex() + "` was inactive"}, nil
	} else if cmd.Action == "list" {
		if cmd.Content == "help" {
			return []string{cmd.GetReminderHelpList()}, nil
		}
		reminders, err := cmd.Client.GetReminders(bson.M{"chat_id": cmd.Message.ChatID})
		if err != nil {
			return []string{}, err
		}
		if len(reminders) == 0 {
			return []string{"Sorry, there are no reminders. Please add new reminder first."}, nil
		}
		str := ""
		for i := 0; i < len(reminders); i++ {
			str += strconv.Itoa(i+1) + ". ID: `" + reminders[i].ID.Hex() + "`\n    TITLE: `" + reminders[i].Title + "`\n    ACTIVE: `" + strconv.FormatBool(reminders[i].Active) + "`\n    Schedule: `" + reminders[i].Schedule + "`\n"
		}
		return []string{str}, nil
	}
	return []string{"Sorry, invalid format. Please check `/reminder help`"}, nil

}

func (cmd *Command) GetReminderHelpAdd() string {
	return "/reminder add `title;SCHEDULE;REMINDER_CONTENT;var1;var2;...`\n`SCHEDULE: minute hour day month day_of_week` (cron format) (http://www.adminschoice.com/crontab-quick-reference)\n`REMINDER_CONTENT` can fill with liquid template (https://github.com/Shopify/liquid/wiki/Liquid-for-Designers)\n`var1`...`varN` is additional data that can use in `REMINDER CONTENT`"
}

func (cmd *Command) GetReminderHelpActive() string {
	return "/reminder active `REMINDER_ID`\n`REMINDER_ID`: ID reminder was created"
}

func (cmd *Command) GetReminderHelpInactive() string {
	return "/reminder inactive `REMINDER_ID`\n`REMINDER_ID`: ID reminder was created"
}

func (cmd *Command) GetReminderHelpList() string {
	return "/reminder list"
}
