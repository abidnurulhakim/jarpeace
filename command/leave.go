package command

import (
	"errors"
	"strconv"
	"strings"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/helper"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (cmd *Command) RunRouteLeave() ([]string, error) {
	if cmd.Action == "add" {
		if cmd.Content == "help" {
			return []string{cmd.GetLeaveHelpAdd()}, nil
		}
		leave := model.NewLeave(cmd.Content)
		leave.UserID = cmd.Message.UserID
		leave.Username = cmd.Message.Username
		err := cmd.Client.CreateLeave(&leave)
		return []string{"Leave with ID: `" + leave.ID.Hex() + "` was added"}, err
	} else if cmd.Action == "broadcast2here" {
		if cmd.Content == "help" {
			return []string{cmd.GetLeaveHelpBroadcast2Here()}, nil
		}
		if !helper.Contains([]string{"active", "inactive"}, strings.ToLower(cmd.Content)) {
			return []string{"Sorry, invalid value. Please only use `active` or `inactive` as valid value."}, nil
		}
		leaveBroadcast := model.NewLeaveBroadcastChat(cmd.Message.ChatID, cmd.Message.UserID, cmd.Message.Username, strings.ToLower(cmd.Content) == "active")
		err := cmd.Client.CreateBroadcastChat(&leaveBroadcast)
		return []string{"Leave broadcast was " + strings.ToLower(cmd.Content)}, err
	} else if cmd.Action == "help" {
		return []string{
			"List available action for /leave :\n" + "1. " + cmd.GetLeaveHelpAdd(),
			"2. " + cmd.GetLeaveHelpBroadcast2Here(),
			"3. " + cmd.GetLeaveHelpList(),
			"4. " + cmd.GetLeaveHelpRemove(),
		}, nil
	} else if cmd.Action == "list" {
		if cmd.Content == "help" {
			return []string{cmd.GetLeaveHelpList()}, nil
		}
		startDate := carbon.Now().SubCentury()
		endDate := carbon.Now().AddCentury()
		arr := strings.SplitN(strings.ToLower(strings.TrimSpace(cmd.Content)), ";", 2)
		var err error
		if len(arr) > 0 {
			startDate, err = helper.ParseHumanDatetime(arr[0], startDate)
		}
		if err != nil {
			return []string{}, errors.New("Sorry, invalid date. Please check `/leave list help`")
		}
		if len(arr) > 1 {
			endDate, err = helper.ParseHumanDatetime(arr[1], endDate)
		}
		if err != nil {
			return []string{}, errors.New("Sorry, invalid date. Please check `/leave list help`")
		}
		userIDs := []int{}
		if cmd.Message.ChatID == cmd.Message.UserID {
			userIDs = []int{cmd.Message.UserID}
		} else {
			broadcastChats, _ := cmd.Client.GetBroadcastChats(bson.M{"chat_id": cmd.Message.ChatID})
			for _, group := range broadcastChats {
				userIDs = append(userIDs, group.UserID)
			}
		}
		leaves, err := cmd.Client.GetLeaves(bson.M{"user_id": bson.M{"$in": userIDs}, "start": bson.M{"$gte": startDate.Time}, "end": bson.M{"$lte": endDate.Time}, "deleted_at": bson.M{"$exists": false}})
		if len(leaves) == 0 {
			return []string{"Sorry, there are no leaves. Please add new leave first."}, nil
		}
		str := ""
		for i := 0; i < len(leaves); i++ {
			leave := leaves[i]
			str += strconv.Itoa(i+1) + ". ID: `" + leave.ID.Hex() + "`\n    Username: " + leave.Username + "\n    Type: " + leave.Type + "\n    Start Date: " + leave.Start.Format(carbon.DayDateTimeFormat) + "\n    End Date: " + leave.End.Format(carbon.DayDateTimeFormat) + "\n    Note: "
			if leave.Note == "" {
				str += "-\n"
			} else {
				str += strings.TrimSpace(leave.Note) + "\n"
			}
		}
		return []string{str}, nil
	} else if cmd.Action == "remove" {
		if cmd.Content == "help" {
			return []string{cmd.GetLeaveHelpRemove()}, nil
		}
		if !bson.IsObjectIdHex(cmd.Content) {
			return []string{"Sorry, invalid format. Please check `/leave remove help`"}, nil
		}
		leaves, err := cmd.Client.GetLeaves(bson.M{"_id": bson.ObjectIdHex(cmd.Content)})
		if err != nil {
			return []string{}, err
		}
		if len(leaves) == 0 {
			return []string{"Leave ID: `" + cmd.Content + "` not found"}, nil
		}
		leave := leaves[0]
		leave.DeletedAt = carbon.Now()
		err = cmd.Client.UpdateLeave(&leave)
		return []string{"Leave ID: `" + cmd.Content + "` was removed"}, err
	}
	return []string{"Sorry, invalid format. Please check `/reminder help`"}, nil

}

func (cmd *Command) GetLeaveHelpAdd() string {
	return "/leave add `TYPE;START_DATE;END_DATE;NOTE`\n`TYPE`: ['remote', 'cuti', 'sick']\n`START_DATE: Start date of your leave ex: (today, tomorrow, 23/10/2018)`\n`END_DATE: End date of your leave ex: (today, tomorrow, 23/10/2018)`\nNOTE: Note of your leave"
}

func (cmd *Command) GetLeaveHelpBroadcast2Here() string {
	return "/leave broadcast2here `STATUS`\n`STATUS`: active or inactive"
}

func (cmd *Command) GetLeaveHelpList() string {
	return "/leave list `START_DATE;END_DATE`\n`START_DATE`: Date what you want\n`END_DATE`: Date what you want"
}

func (cmd *Command) GetLeaveHelpRemove() string {
	return "/leave remove `LEAVE_ID`\n`LEAVE_ID`: ID your leave was created"
}
