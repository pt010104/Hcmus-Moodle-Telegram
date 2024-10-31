package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pt010104/Hcmus-Moodle-Telegram/config"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"
	calendarUC "github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar/usecase"
	telegramUC "github.com/pt010104/Hcmus-Moodle-Telegram/internal/telegram/usecase"
	pkgLog "github.com/pt010104/Hcmus-Moodle-Telegram/pkg/log"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/microservice/notification"
	"github.com/pt010104/Hcmus-Moodle-Telegram/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	util.PrintJson(cfg)

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramConfig.BotToken)
	if err != nil {
		l.Fatal(context.Background(), "Failed to create Telegram bot", err.Error())
	}
	bot.Debug = false

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.Fatal(ctx, "Failed to connect to MongoDB", err.Error())
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			l.Fatal(ctx, "Failed to disconnect MongoDB", err.Error())
		}
	}()

	db := client.Database(cfg.Mongo.Database)

	notificationSrv := notification.New(l, cfg.HcmusConfig.URL, cfg.HcmusConfig.SessKey, cfg.HcmusConfig.Cookies)
	telegramService := telegramUC.New(l, bot, cfg.TelegramConfig.ChatID, nil, db)
	calendarService := calendarUC.New(l, notificationSrv, db, telegramService)

	telegramService.SetCalendarUC(calendarService)

	err = telegramService.SendMessage(ctx, "Bot started at "+util.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		l.Error(ctx, "Failed to send startup message", err.Error())
	}

	c := cron.New(cron.WithSeconds())

	_, err = c.AddFunc("@every 3m", func() {
		l.ResetLogger()
		l.Info(ctx, "Scheduled Task", "Executing GetFromCalendar")
		_, err := calendarService.GetFromCalendar(ctx)
		if err != nil {
			l.Error(ctx, "GetFromCalendar", err.Error())
		} else {
			l.Info(ctx, "GetFromCalendar", "Executed successfully")
		}
	})
	if err != nil {
		l.Fatal(ctx, "Failed to schedule GetFromCalendar", err.Error())
	}

	_, err = c.AddFunc("@every 3m", func() {
		l.ResetLogger()
		l.Info(ctx, "Scheduled Task", "Executing GetFromNotification")
		notificationInput := calendar.GetFromNotificationInput{
			Limit:    3,
			Offset:   0,
			UserIDTo: "4248",
		}
		_, err := calendarService.GetFromNotification(ctx, notificationInput)
		if err != nil {
			l.Error(ctx, "GetFromNotification", err.Error())
		} else {
			l.Info(ctx, "GetFromNotification", "Executed successfully")
		}
	})
	if err != nil {
		l.Fatal(ctx, "Failed to schedule GetFromNotification", err.Error())
	}

	c.Start()
	l.Info(ctx, "Cron Scheduler", "Started successfully")
	defer c.Stop()

	var wg sync.WaitGroup

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for update := range updates {
			if update.Message != nil {
				err := telegramService.CommandHandler(ctx, update.Message)
				if err != nil {
					l.Error(ctx, "Failed to process message", err)
				}
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	l.Info(ctx, "Shutdown", "Received interrupt signal. Shutting down gracefully...")

	c.Stop()
	l.Info(ctx, "Shutdown", "Cron scheduler stopped")

	cancel()

	wg.Wait()

	l.Info(ctx, "Shutdown", "Shutdown complete.")
}
