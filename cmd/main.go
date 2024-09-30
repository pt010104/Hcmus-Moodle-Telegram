package main

import (
	"context"
	"sync"
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

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	notificationSrv := notification.New(l, cfg.HcmusConfig.URL, cfg.HcmusConfig.SessKey, cfg.HcmusConfig.Cookies)

	telegramUC := telegramUC.New(l, bot, cfg.TelegramConfig.ChatID)

	calendarUC := calendarUC.New(l, notificationSrv, db, telegramUC)

	telegramUC.SendMessage(ctx, "Bot started new code at "+util.Now().Format("2006-01-02 15:04:05"))

	for {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var calendarErr error
		var notificationErr error

		wg.Add(2)

		go func() {
			defer wg.Done()
			_, err := calendarUC.GetFromCalendar(ctx)
			if err != nil {
				l.Error(ctx, "main.GetFromCalendar", err.Error())
				mu.Lock()
				calendarErr = err
				mu.Unlock()
			}
		}()

		go func() {
			defer wg.Done()
			notificationInput := calendar.GetFromNotificationInput{
				Limit:    3,
				Offset:   0,
				UserIDTo: "4248",
			}
			_, err := calendarUC.GetFromNotification(ctx, notificationInput)
			if err != nil {
				l.Error(ctx, "main.GetFromNotification", err.Error())
				mu.Lock()
				notificationErr = err
				mu.Unlock()
			}
		}()

		wg.Wait()

		if calendarErr != nil {
			l.Error(ctx, "Error occurred in GetFromCalendar", calendarErr.Error())
		}

		if notificationErr != nil {
			l.Error(ctx, "Error occurred in GetFromNotification", notificationErr.Error())
		}

		time.Sleep(15 * time.Minute)
	}
}
