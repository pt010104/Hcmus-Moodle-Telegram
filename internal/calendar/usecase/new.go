package usecase

import (
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/telegram"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/log"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/microservice/notification"
	"go.mongodb.org/mongo-driver/mongo"
)

type implUseCase struct {
	l               log.Logger
	notificationSrv notification.UseCase
	db              *mongo.Database
	telegramUC      telegram.UseCase
}

func New(
	l log.Logger,
	notificationSrv notification.UseCase,
	db *mongo.Database,
	telegramUC telegram.UseCase,
) calendar.UseCase {
	return &implUseCase{
		l,
		notificationSrv,
		db,
		telegramUC,
	}
}
