package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/telegram"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type implUseCase struct {
	l          log.Logger
	bot        *tgbotapi.BotAPI
	chatID     int64
	calendarUC calendar.UseCase
	db         *mongo.Database
}

func New(
	l log.Logger,
	bot *tgbotapi.BotAPI,
	chatID int64,
	calendarUC calendar.UseCase,
	db *mongo.Database,
) telegram.UseCase {
	return &implUseCase{
		l:          l,
		bot:        bot,
		chatID:     chatID,
		calendarUC: calendarUC,
		db:         db,
	}
}

func (uc *implUseCase) SetCalendarUC(calendarUC calendar.UseCase) {
	uc.calendarUC = calendarUC
}
