package command

import (
	"fmt"
	"strconv"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/model"
	"gopkg.in/mgo.v2/bson"
)

func (cmd *Command) RunRouteAutoReply() ([]string, error) {
	if cmd.Action == "add" {
		if cmd.Content == "help" {
			return []string{cmd.GetAutoReplyHelpAdd()}, nil
		}
		autoReply := model.NewAutoReply(cmd.Content, cmd.Message.ChatID)
		autoReply.UserID = cmd.Message.UserID
		autoReply.Username = cmd.Message.Username
		err := cmd.Client.CreateAutoReply(&autoReply)
		return []string{"Auto reply with ID: `" + autoReply.ID.Hex() + "` was added"}, err
	} else if cmd.Action == "help" {
		return []string{
			"List available action for /autoreply :\n" + "1. " + cmd.GetAutoReplyHelpAdd(),
			"2. " + cmd.GetAutoReplyHelpList(),
			"3. " + cmd.GetAutoReplyHelpActive(),
			"4. " + cmd.GetAutoReplyHelpInactive(),
			"5. " + cmd.GetAutoReplyHelpRemove(),
		}, nil
	} else if cmd.Action == "list" {
		if cmd.Content == "help" {
			return []string{cmd.GetAutoReplyHelpList()}, nil
		}
		if cmd.Message.ChatID == cmd.Message.UserID {
			autoReplies, err := cmd.Client.GetAutoReplies(bson.M{"user_id": cmd.Message.UserID, "deleted_at": bson.M{"$exists": false}})
			if len(autoReplies) == 0 {
				return []string{"🙏 Sorry, there are no auto reply. Please add new auto reply first."}, nil
			}
			responses := []string{}
			for i := 0; i < len(autoReplies); i++ {
				autoReply := autoReplies[i]
				responses = append(responses, strconv.Itoa(i+1)+". ID: `"+autoReply.ID.Hex()+"`\n    Text: "+autoReply.Text+"\n    Answer: "+autoReply.Answer+"\n    Active: "+strconv.FormatBool(autoReply.Active)+"\n    Chat IDs: "+fmt.Sprint(autoReply.ChatIDs))
			}
			return responses, err
		} else {
			autoReplies, err := cmd.Client.GetAutoReplies(bson.M{"chat_id": bson.M{"$in": []int{cmd.Message.ChatID}}, "deleted_at": bson.M{"$exists": false}})
			if len(autoReplies) == 0 {
				return []string{"🙏 Sorry, there are no auto reply. Please add new auto reply first."}, nil
			}
			responses := []string{}
			for i := 0; i < len(autoReplies); i++ {
				autoReply := autoReplies[i]
				responses = append(responses, strconv.Itoa(i+1)+". Text: "+autoReply.Text+"\n    Answer: "+autoReply.Answer+"\n    Active: "+strconv.FormatBool(autoReply.Active)+"\n    Chat IDs: "+fmt.Sprint(autoReply.ChatIDs)+"\n    Username: "+autoReply.Username)
			}
			return responses, err
		}
	} else if cmd.Action == "active" {
		if cmd.Content == "help" {
			return []string{cmd.GetAutoReplyHelpRemove()}, nil
		}
		if !bson.IsObjectIdHex(cmd.Content) {
			return []string{"Sorry, invalid format. Please check `/autoreply active help`"}, nil
		}
		autoReplies, err := cmd.Client.GetAutoReplies(bson.M{"_id": bson.ObjectIdHex(cmd.Content), "deleted_at": bson.M{"$exists": false}})
		if err != nil {
			return []string{}, err
		}
		if len(autoReplies) == 0 {
			return []string{"Auto reply ID: `" + cmd.Content + "` not found"}, nil
		}
		autoReply := autoReplies[0]
		autoReply.Active = true
		err = cmd.Client.UpdateAutoReply(&autoReply)
		return []string{"Auto reply ID: `" + cmd.Content + "` was removed"}, err
	} else if cmd.Action == "inactive" {
		if cmd.Content == "help" {
			return []string{cmd.GetAutoReplyHelpRemove()}, nil
		}
		if !bson.IsObjectIdHex(cmd.Content) {
			return []string{"Sorry, invalid format. Please check `/autoreply inactive help`"}, nil
		}
		autoReplies, err := cmd.Client.GetAutoReplies(bson.M{"_id": bson.ObjectIdHex(cmd.Content), "deleted_at": bson.M{"$exists": false}})
		if err != nil {
			return []string{}, err
		}
		if len(autoReplies) == 0 {
			return []string{"Auto reply ID: `" + cmd.Content + "` not found"}, nil
		}
		autoReply := autoReplies[0]
		autoReply.Active = false
		err = cmd.Client.UpdateAutoReply(&autoReply)
		return []string{"Auto reply ID: `" + cmd.Content + "` was removed"}, err
	} else if cmd.Action == "remove" {
		if cmd.Content == "help" {
			return []string{cmd.GetAutoReplyHelpRemove()}, nil
		}
		if !bson.IsObjectIdHex(cmd.Content) {
			return []string{"Sorry, invalid format. Please check `/autoreply remove help`"}, nil
		}
		autoReplies, err := cmd.Client.GetAutoReplies(bson.M{"_id": bson.ObjectIdHex(cmd.Content), "deleted_at": bson.M{"$exists": false}})
		if err != nil {
			return []string{}, err
		}
		if len(autoReplies) == 0 {
			return []string{"Auto reply ID: `" + cmd.Content + "` not found"}, nil
		}
		autoReply := autoReplies[0]
		autoReply.DeletedAt = carbon.Now()
		err = cmd.Client.UpdateAutoReply(&autoReply)
		return []string{"Auto reply ID: `" + cmd.Content + "` was removed"}, err
	}
	return []string{"Sorry, invalid format. Please check `/reminder help`"}, nil

}

func (cmd *Command) GetAutoReplyHelpAdd() string {
	return "/autoreply add `TEXT;ANSWER;ACTIVE;CHAT_IDS;var1;var2;`\n`TEXT`: The text that will auto reply\n`ANSWER`: The answer of auto reply\n`ACTIVE`: Status active of auto reply, valid value :['true', 'false']\n`CHAT_IDS`: List of chat ID that can give auto reply. Chat ID will separate by ','\n`var1`...`varN` is additional data that can use in `ANSWER`"
}

func (cmd *Command) GetAutoReplyHelpList() string {
	return "/autoreply list"
}

func (cmd *Command) GetAutoReplyHelpActive() string {
	return "/autoreply active `AUTOREPLY_ID`\n`AUTOREPLY_ID`: ID your auto reply was created"
}

func (cmd *Command) GetAutoReplyHelpInactive() string {
	return "/autoreply inactive `AUTOREPLY_ID`\n`AUTOREPLY_ID`: ID your auto reply was created"
}

func (cmd *Command) GetAutoReplyHelpRemove() string {
	return "/autoreply remove `AUTOREPLY_ID`\n`AUTOREPLY_ID`: ID your auto reply was created"
}
