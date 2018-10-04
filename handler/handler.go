package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abidnurulhakim/jarpeace/channel"
	"github.com/abidnurulhakim/jarpeace/command"
	"github.com/abidnurulhakim/jarpeace/database"
	"github.com/abidnurulhakim/jarpeace/helper"
	"github.com/abidnurulhakim/jarpeace/model"
	"github.com/abidnurulhakim/jarpeace/response"
	"github.com/julienschmidt/httprouter"
	"github.com/osteele/liquid"
	"gopkg.in/mgo.v2/bson"
)

type Handler struct {
	Db       *database.MongoDB
	Telegram *channel.Telegram
}

func (handler *Handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	object_response(w, struct {
		Message string `json:"message"`
	}{Message: "ok"}, http.StatusOK)
	return nil
}

func (handler *Handler) CreateReminder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	var err error
	var params map[string]interface{}
	db := handler.Db.Copy()
	defer db.Close()
	json.NewDecoder(r.Body).Decode(&params)
	reminder := model.NewReminder("")
	decoder, err := helper.GetDecoder(&reminder)
	if err == nil {
		err = decoder.Decode(params)
	}
	db.CreateReminder(&reminder)
	return success_or_failed(w, reminder, http.StatusOK, err, http.StatusInternalServerError, 500)
}

func (handler *Handler) Webhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	var err error
	var params map[string]interface{}
	db := handler.Db.Copy()
	defer db.Close()
	json.NewDecoder(r.Body).Decode(&params)
	resourceTelegram := channel.TelegramUpdate{}
	decoder, err := helper.GetDecoder(&resourceTelegram)
	if err == nil {
		err = decoder.Decode(params)
	}
	message := model.NewMessage()
	message.ChatID = resourceTelegram.Message.Chat.Id
	message.UserID = resourceTelegram.Message.From.Id
	message.Username = "@" + resourceTelegram.Message.From.Username
	message.FirstName = resourceTelegram.Message.From.FirstName
	message.LastName = resourceTelegram.Message.From.LastName
	message.MessageID = resourceTelegram.Message.MessageId
	message.ReplyMessageID = resourceTelegram.Message.ReplyToMessage.MessageId
	message.Content = resourceTelegram.Message.Text
	document := resourceTelegram.Message.Document
	photos := resourceTelegram.Message.Photo
	if document.FileId != "" {
		message.Files = append(message.Files, "https://api.telegram.org/bot"+handler.Telegram.Token+"/getFile?file_id="+document.FileId)
	}
	if len(photos) > 0 {
		message.Files = append(message.Files, "https://api.telegram.org/bot"+handler.Telegram.Token+"/getFile?file_id="+photos[len(photos)-1].FileId)
	}
	err = db.CreateMessage(&message)
	if err == nil {
		go MessageWebhookProcess(db.Copy(), &message)
	}
	return success_or_failed(w, message, http.StatusOK, err, http.StatusInternalServerError, 500)
}

func object_response(w http.ResponseWriter, o interface{}, httpCode int) {
	successResponse := response.BuildSuccess(o, response.MetaInfo{HTTPStatus: httpCode})
	response.Write(w, successResponse)
}

func error_response(w http.ResponseWriter, err error, code int, httpCode int) {
	ce := response.CustomError{Message: err.Error(), Code: code, HTTPCode: httpCode}
	errorResponse := response.BuildError([]error{ce})
	http.Error(w, http.StatusText(httpCode), httpCode)
	response.Write(w, errorResponse)
}

func success_or_failed(w http.ResponseWriter, o interface{}, sCode int, err error, fCode int, errorCode int) error {
	if err != nil {
		error_response(w, err, errorCode, fCode)
	} else {
		object_response(w, o, sCode)
	}

	return err
}

func FetchUrlQueryInteger(values url.Values, key string, defaultValue int) int {
	result := defaultValue
	if i, err := strconv.Atoi(values.Get(key)); err == nil {
		result = i
	}
	return result
}

func FetchUrlQueryBoolean(values url.Values, key string, defaultValue bool) bool {
	result := defaultValue
	if helper.Contains([]string{"true", "TRUE", "True", "1"}, values.Get(key)) {
		return true
	}
	if helper.Contains([]string{"false", "FALSE", "False", "0"}, values.Get(key)) {
		return false
	}
	return result
}

func MessageWebhookProcess(db *database.MongoDB, message *model.Message) {
	defer db.Close()
	telegram := channel.Telegram{Token: os.Getenv("TELEGRAM_TOKEN")}
	if message.IsCommand() {
		texts, err := command.Run(db, *message)
		if err != nil {
			messageParam := channel.TelegramParamMessageText{}
			messageParam.ChatId = strconv.Itoa(message.ChatID)
			messageParam.ReplyToMessageId = message.MessageID
			messageParam.ParseMode = "markdown"
			messageParam.Text = err.Error()
			telegram.SendMessage(messageParam)
		} else {
			for i := 0; i < len(texts); i++ {
				text := texts[i]
				messageParam := channel.TelegramParamMessageText{}
				messageParam.ChatId = strconv.Itoa(message.ChatID)
				messageParam.ReplyToMessageId = message.MessageID
				messageParam.ParseMode = "markdown"
				messageParam.Text = text
				telegram.SendMessage(messageParam)
			}
		}
	} else if message.Content != "" {
		engine := liquid.NewEngine()
		autoReplies, _ := db.GetAutoReplies(bson.M{"active": true, "deleted_at": bson.M{"$exists": false}, "chat_ids": bson.M{"$in": []int{0, message.ChatID}}})
		for i := 0; i < len(autoReplies); i++ {
			if autoReplies[i].Answer != "" {
				canAutoReply := message.IsMentionBot()
				if message.ChatID == message.UserID {
					canAutoReply = true
				}
				if canAutoReply && strings.Contains(strings.ToLower(message.Content), strings.ToLower(autoReplies[i].Text)) {
					bindings := autoReplies[i].Data
					bindings["now"] = time.Now()
					bindings["username"] = message.Username
					bindings["first_name"] = message.FirstName
					bindings["last_name"] = message.LastName
					bindings["name"] = message.FirstName + " " + message.LastName
					out, _ := engine.ParseAndRenderString(autoReplies[i].Answer, bindings)
					messageParam := channel.TelegramParamMessageText{}
					messageParam.ChatId = strconv.Itoa(message.ChatID)
					messageParam.ReplyToMessageId = message.MessageID
					messageParam.ParseMode = "markdown"
					messageParam.Text = out
					if out != "" {
						telegram.SendMessage(messageParam)
					}
				}
			}
		}
	}
}
