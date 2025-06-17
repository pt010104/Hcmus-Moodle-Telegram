package notification

import (
	"context"

	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/log"
)

type UseCase interface {
	GetFromCalendar(ctx context.Context, input GetFromCalendarFilter) (Calendar, error)
	GetFromNotification(ctx context.Context, input GetFromNotificationFilter) (Notification, error)
	CheckEventSubmission(ctx context.Context, input EventSubmissionFilter) (EventDetail, error)
}

type implUseCase struct {
	url     string
	l       log.Logger
	sessKey string
	cookies string
}

func New(l log.Logger, url string, sessKey string, cookies string) UseCase {
	return implUseCase{
		url:     url,
		l:       l,
		sessKey: sessKey,
		cookies: cookies,
	}
}
