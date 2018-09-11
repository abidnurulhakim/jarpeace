package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abidnurulhakim/jarpeace/channel"
	"github.com/abidnurulhakim/jarpeace/database"
	"github.com/abidnurulhakim/jarpeace/handler"
	"github.com/abidnurulhakim/jarpeace/middleware"
	"github.com/abidnurulhakim/jarpeace/model"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/osteele/liquid"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
	cron "gopkg.in/robfig/cron.v2"
)

func main() {
	gotenv.Load()

	fmt.Println("jarpeace service is starting...")
	// DB setup
	fmt.Println("Connecting to", os.Getenv("MONGO_DB"), "on", os.Getenv("MONGO_HOST"))
	option := &database.MongoDBOption{
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Database: os.Getenv("MONGO_DB"),
	}
	db, err := database.NewMongoDB(option)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Using", os.Getenv("MONGO_DB"), "on", os.Getenv("MONGO_HOST"))

	// Routing
	router := httprouter.New()
	telegram := channel.Telegram{Token: os.Getenv("TELEGRAM_TOKEN")}
	handler := handler.Handler{Db: db, Telegram: &telegram}

	fmt.Println("Set webhook")
	telegram.SetWebhook(os.Getenv("BASE_URL") + "/webhook/" + telegram.Token)

	router.GET("/healthz", middleware.HTTP(handler.Index))
	router.POST("/reminders", middleware.HTTP(handler.CreateReminder))
	router.POST("/webhook/"+os.Getenv("TELEGRAM_TOKEN"), middleware.HTTP(handler.Webhook))

	c := cron.New()
	c.AddFunc("TZ=Asia/Jakarta 0 * * * * *", RunReminders)
	c.Start()

	// Start server
	fmt.Println("Listening at port", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func RunReminders() {
	now := time.Now()
	option := &database.MongoDBOption{
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Database: os.Getenv("MONGO_DB"),
	}
	db, err := database.NewMongoDB(option)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	reminders, err := db.GetReminders(bson.M{"active": true})
	if err != nil {
		panic(err)
	}
	validReminders := []model.Reminder{}
	for _, reminder := range reminders {
		schedule, err := cron.Parse(reminder.Schedule)
		since := now.Add(-59 * time.Second)
		nextTime := schedule.Next(since)
		if err != nil || nextTime.IsZero() {
			reminder.Active = false
			db.UpdateReminder(&reminder)
		}
		if now.Year() == nextTime.Year() && now.Month() == nextTime.Month() && now.Day() == nextTime.Day() && now.Hour() == nextTime.Hour() && now.Minute() == nextTime.Minute() {
			validReminders = append(validReminders, reminder)
		}
	}
	for _, validReminder := range validReminders {
		start := time.Now()
		engine := liquid.NewEngine()
		telegram := channel.Telegram{Token: os.Getenv("TELEGRAM_TOKEN")}
		message := channel.TelegramParamMessageText{}
		template := validReminder.Content
		bindings := validReminder.Data
		bindings["now"] = time.Now()
		out, err := engine.ParseAndRenderString(template, bindings)
		if err != nil {
			logger, _ := zap.NewProduction()
			elapsed := time.Since(start).Seconds() * 1000
			elapsedStr := strconv.FormatFloat(elapsed, 'f', -1, 64)
			requestId, _ := uuid.NewRandom()
			logger.Error(err.Error(),
				zap.String("request_id", requestId.String()),
				zap.String("duration", elapsedStr),
				zap.Strings("tags", []string{"render-liquid-template"}),
			)
			message.Text = validReminder.Content
		} else {
			message.Text = out
		}
		message.ChatId = strconv.Itoa(validReminder.ChatID)
		message.ParseMode = "markdown"
		telegram.SendMessage(message)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("BASE_URL")+"/healthz", nil)
	req.Header.Set("Content-Type", "application/json")
	client.Do(req)
}
