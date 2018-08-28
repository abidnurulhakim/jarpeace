package channel

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abidnurulhakim/jarpeace/helper"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Telegram struct {
	Token string
}

// Representative type response telegram
type TelegramCallbackQuery struct {
	MessageId       int    `json:"message_id"`
	Message         string `json:"message"`
	EditedMessage   string `json:"edited_message"`
	ChannelPost     string `json:"channel_post"`
	EditChannelPost string `json:"edit_channel_post"`
	CallbackQuery   string `json:"callback_query"`
}

type TelegramChat struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type TelegramDocument struct {
	FileId   string        `json:"file_id"`
	Thumb    TelegramPhoto `json:"thumb"`
	FileName string        `json:"file_name"`
	MimeType string        `json:"mime_type"`
	FileSize int           `json:"file_size"`
}

type TelegramFile struct {
	FileId   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FilePath string `json:"file_path"`
}

type TelegramMessage struct {
	MessageId        int                  `json:"message_id"`
	From             TelegramUser         `json:"from"`
	Date             int                  `json:"date"`
	Chat             TelegramChat         `json:"chat"`
	ForwardFrom      TelegramUser         `json:"forward_from"`
	ForwardFromChat  TelegramChat         `json:"forward_from_chat"`
	ForwardMessageId int                  `json:"forward_message_id"`
	ForwardDate      int                  `json:"forward_date"`
	ReplyToMessage   TelegramMessageReply `json:"reply_to_message"`
	Text             string               `json:"text"`
	Document         TelegramDocument     `json:"document"`
	Photo            []TelegramPhoto      `json:"photo"`
}

type TelegramMessageReply struct {
	MessageId        int          `json:"message_id"`
	From             TelegramUser `json:"from"`
	Date             int          `json:"date"`
	Chat             TelegramChat `json:"chat"`
	ForwardFromChat  TelegramChat `json:"forward_from_chat"`
	ForwardFrom      TelegramUser `json:"forward_from"`
	ForwardMessageId int          `json:"forward_message_id"`
	ForwardDate      int          `json:"forward_date"`
	Text             string       `json:"text"`
}

type TelegramPhoto struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type TelegramUpdate struct {
	UpdateId        int             `json:"update_id"`
	Message         TelegramMessage `json:"message"`
	EditedMessage   TelegramMessage `json:"edited_message"`
	ChannelPost     TelegramMessage `json:"channel_post"`
	EditChannelPost TelegramMessage `json:"edit_channel_post"`
	CallbackQuery   string          `json:"callback_query"`
}

type TelegramUser struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Representative post parameter request to telegram
type TelegramParamGetFile struct {
	FileId string `json:"file_id"`
}

type TelegramParamMessageText struct {
	ChatId                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessageId      int    `json:"reply_to_message_id"`
}

type TelegramParamMessageDocument struct {
	ChatId                string `json:"chat_id"`
	Document              string `json:"document"`
	Caption               string `json:"caption"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessageId      int    `json:"reply_to_message_id"`
}

type TelegramParamWebhook struct {
	Url            string   `json:"url"`
	MaxConnections int      `json:"max_connections"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type TelegramResponse struct {
	Description string                 `json:"description"`
	ErrorCode   int                    `json:"error_code"`
	Ok          bool                   `json:"ok"`
	Result      map[string]interface{} `json:"result"`
}

// Telegram Request API
func (telegram *Telegram) Request(method string, endpoint string, data string, endpointName string) (string, error) {
	start := time.Now()
	client := &http.Client{}
	var jsonStr = []byte(data)
	req, err := http.NewRequest(method, "https://api.telegram.org/bot"+telegram.Token+"/"+endpoint, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	bodyText, err := ioutil.ReadAll(resp.Body)
	text := string(bodyText)
	if resp.StatusCode < 200 || resp.StatusCode > 302 {
		err = errors.New(text)
	}
	metricEndpoint := strings.Split(strings.Split(endpointName, "?")[0], "/")[0]
	if err != nil {
		logger, _ := zap.NewProduction()
		elapsed := time.Since(start).Seconds() * 1000
		elapsedStr := strconv.FormatFloat(elapsed, 'f', -1, 64)
		requestId, _ := uuid.NewRandom()
		logger.Error(err.Error(),
			zap.String("request_id", requestId.String()),
			zap.String("duration", elapsedStr),
			zap.Strings("tags", []string{metricEndpoint, "telegram-api"}),
		)
	}
	return text, err
}

func (telegram *Telegram) SetWebhook(urlWebhook string) error {
	data, err := json.Marshal(&TelegramParamWebhook{Url: urlWebhook, MaxConnections: 100})
	if err == nil {
		response, err := telegram.Request("POST", "setWebhook", string(data), "set-webhook")
		if err == nil {
			GetResultOfTelegramResponse(response)
		}
	}
	return err
}

func (telegram *Telegram) SendMessage(message TelegramParamMessageText) (TelegramMessage, error) {
	telegramMessage := TelegramMessage{}
	message.ParseMode = strings.ToLower(message.ParseMode)
	if !helper.Contains([]string{"markdown", "html"}, strings.ToLower(message.ParseMode)) {
		message.ParseMode = "markdown"
	}
	data, err := json.Marshal(&message)
	if err == nil {
		response, err := telegram.Request("POST", "sendMessage", string(data), "send-message")
		if err == nil {
			result, err := GetResultOfTelegramResponse(response)
			if err == nil {
				decoder, err := helper.GetDecoder(&telegramMessage)
				if err != nil {
					return telegramMessage, err
				}
				err = decoder.Decode(result)
				return telegramMessage, err
			}
		}
	}
	return telegramMessage, err
}

func (telegram *Telegram) SendDocument(message TelegramParamMessageDocument) (TelegramMessage, error) {
	telegramMessage := TelegramMessage{}
	message.ParseMode = strings.ToLower(message.ParseMode)
	if !helper.Contains([]string{"markdown", "html"}, strings.ToLower(message.ParseMode)) {
		message.ParseMode = "markdown"
	}
	data, err := json.Marshal(&message)
	if err == nil {
		response, err := telegram.Request("POST", "sendDocument", string(data), "send-document")
		if err == nil {
			result, err := GetResultOfTelegramResponse(response)
			if err == nil {
				decoder, err := helper.GetDecoder(&telegramMessage)
				if err != nil {
					return telegramMessage, err
				}
				err = decoder.Decode(result)
				return telegramMessage, err
			}
		}
	}
	return telegramMessage, err
}

func (telegram *Telegram) GetFile(fileId string) (TelegramFile, error) {
	telegramFile := TelegramFile{}
	data, err := json.Marshal(&TelegramParamGetFile{FileId: fileId})
	if err == nil {
		response, err := telegram.Request("GET", "getFile", string(data), "get-file")
		if err == nil {
			result, err := GetResultOfTelegramResponse(response)
			if err == nil {
				decoder, err := helper.GetDecoder(&telegramFile)
				if err != nil {
					return telegramFile, err
				}
				err = decoder.Decode(result)
				return telegramFile, err
			}
		}
	}
	return telegramFile, err
}

func GetResultOfTelegramResponse(rawResponse string) (map[string]interface{}, error) {
	var result map[string]interface{}
	telegramResponse := TelegramResponse{}
	err := json.Unmarshal([]byte(rawResponse), &telegramResponse)
	if err != nil {
		return result, err
	}
	if telegramResponse.Ok {
		return telegramResponse.Result, err
	} else {
		return result, errors.New(telegramResponse.Description)
	}
}
