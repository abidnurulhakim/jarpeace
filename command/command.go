package command

import (
	"errors"
	"strings"

	"github.com/abidnurulhakim/jarpeace/database"
	"github.com/abidnurulhakim/jarpeace/model"
)

type Command struct {
	Client  *database.MongoDB
	Name    string
	Action  string
	Content string
	Message model.Message
}

type CommandWorker interface {
	Execute() ([]string, error)
}

func Run(db *database.MongoDB, message model.Message) ([]string, error) {
	if message.Content == "" {
		return []string{}, errors.New("🤔")
	}
	if string(message.Content[0]) != "/" || message.Content == "/" {
		return []string{}, errors.New("🤔")
	}
	tmp := strings.SplitN(message.Content, " ", 3)
	commandName := strings.TrimSpace(string(tmp[0][1:]))
	action := "help"
	content := ""
	if len(tmp) > 1 {
		action = strings.ToLower(tmp[1])
	}
	if len(tmp) > 2 {
		content = tmp[2]
	}
	command := Command{}
	command.Name = commandName
	command.Client = db
	command.Action = action
	command.Content = content
	command.Message = message
	return command.Execute()
}

func (cmd *Command) Execute() ([]string, error) {
	switch cmd.Name {
	case "leave":
		return cmd.RunRouteLeave()
	case "reminder":
		return cmd.RunRouteReminder()
	default:
		return []string{}, errors.New("Command not found")
	}
}
