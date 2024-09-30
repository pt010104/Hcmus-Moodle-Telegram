package main

import (
	"context"
	"time"

	"github.com/pt010104/Hcmus-Moodle-Telegram/config"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"
	calendarUC "github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar/usecase"
	telegramUC "github.com/pt010104/Hcmus-Moodle-Telegram/internal/telegram/usecase"
	pkgLog "github.com/pt010104/Hcmus-Moodle-Telegram/pkg/log"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/microservice/notification"
	"github.com/pt010104/Hcmus-Moodle-Telegram/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramConfig.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(cfg.Mongo.Database)

	// Initialize the logger
	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	notificationSrv := notification.New(l, cfg.HcmusConfig.URL, cfg.HcmusConfig.SessKey, cfg.HcmusConfig.Cookies)

	telegramUC := telegramUC.New(l, bot, cfg.TelegramConfig.ChatID)

	calendarUC := calendarUC.New(l, notificationSrv, db, telegramUC)

	for {

		// calendarInput := calendar.GetFromCalendarInput{
		// 	Year:  "2024",
		// 	Month: "7",
		// 	Day:   "25",
		// }
		today := util.Now()

		calendarInput := calendar.GetFromCalendarInput{
			Year:  today.Format("2006"),
			Month: today.Format("1"),
			Day:   today.Format("2"),
		}
		_, err := calendarUC.GetFromCalendar(ctx, calendarInput)
		if err != nil {
			l.Error(ctx, "main.GetFromCalendar", err.Error())
			return
		}

		notificationInput := calendar.GetFromNotificationInput{
			Limit:    3,
			Offset:   0,
			UserIDTo: "4248",
		}
		_, err = calendarUC.GetFromNotification(ctx, notificationInput)
		if err != nil {
			l.Error(ctx, "main.GetFromNotification", err.Error())
			return
		}

		time.Sleep(1 * time.Minute)
	}
}
