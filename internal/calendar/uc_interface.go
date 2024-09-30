package calendar

import (
	"context"

	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/models"
)

type UseCase interface {
	GetFromCalendar(ctx context.Context) ([]models.Calendar, error)
	GetFromNotification(ctx context.Context, input GetFromNotificationInput) ([]models.Notification, error)
}
