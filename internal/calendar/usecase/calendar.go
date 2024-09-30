package usecase

import (
	"context"
	"fmt"

	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/models"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/microservice/notification"
	"go.mongodb.org/mongo-driver/bson"
)

func (uc implUseCase) GetFromCalendar(ctx context.Context, input calendar.GetFromCalendarInput) ([]models.Calendar, error) {
	calendarSrv, err := uc.notificationSrv.GetFromCalendar(ctx, notification.GetFromCalendarFilter{
		Year:  input.Year,
		Month: input.Month,
		Day:   input.Day,
	})
	if err != nil {
		return nil, err
	}

	calendarEvents := make([]models.Calendar, 0)

	for _, event := range calendarSrv.Events {
		calendarEvents = append(calendarEvents, models.Calendar{
			ID:            event.ID,
			Name:          event.Name,
			Description:   event.Description,
			FormattedTime: event.FormattedTime,
			CourseName:    event.Course.FullName,
			CourseID:      event.Course.ID,
			URL:           event.URL,
		})
	}

	calendarOutputs := make([]models.Calendar, 0)

	if len(calendarEvents) > 0 {
		collection := uc.db.Collection("calendar_events")

		for _, c := range calendarEvents {
			filter := bson.M{"_id": c.ID}
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				uc.l.Error(ctx, "usecase.GetFromCalendar.CountDocuments", err.Error())
				continue
			}

			if count == 0 {
				_, err := collection.InsertOne(ctx, c)
				if err != nil {
					uc.l.Error(ctx, "usecase.GetFromCalendar.InsertOne", err.Error())
				}

				calendarOutputs = append(calendarOutputs, c)
			} else {
				uc.l.Info(ctx, "usecase.GetFromCalendar.EventExists", fmt.Sprintf("Sự kiện với ID %d đã tồn tại", c.ID))
			}
		}

		msgTexts := uc.createMsgCalendarForTelegram(ctx, calendarOutputs)

		for _, msgText := range msgTexts {
			error := uc.telegramUC.SendMessage(ctx, msgText)
			if error != nil {
				uc.l.Error(ctx, "Failed to send message to telegram", error)
			}
		}

	}

	return calendarOutputs, nil
}

func (uc implUseCase) GetFromNotification(ctx context.Context, input calendar.GetFromNotificationInput) ([]models.Notification, error) {
	notificationSrv, err := uc.notificationSrv.GetFromNotification(ctx, notification.GetFromNotificationFilter{
		Limit:    input.Limit,
		Offset:   input.Offset,
		UserIDTo: input.UserIDTo,
	})
	if err != nil {
		return nil, err
	}

	notificationEvents := make([]models.Notification, 0)

	for _, event := range notificationSrv.Events {
		notificationEvents = append(notificationEvents, models.Notification{
			ID:          event.ID,
			Subject:     event.Subject,
			FullMessage: event.FullMessage,
			TimeCreated: event.TimeCreated,
		})
	}

	for i, j := 0, len(notificationEvents)-1; i < j; i, j = i+1, j-1 {
		notificationEvents[i], notificationEvents[j] = notificationEvents[j], notificationEvents[i]
	}

	notificationOutputs := make([]models.Notification, 0)

	if len(notificationEvents) > 0 {
		collection := uc.db.Collection("notifications")

		for _, n := range notificationEvents {
			filter := bson.M{"_id": n.ID}
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				uc.l.Error(ctx, "usecase.GetFromNotification.CountDocuments", err.Error())
				continue
			}

			if count == 0 {
				_, err := collection.InsertOne(ctx, n)
				if err != nil {
					uc.l.Error(ctx, "usecase.GetFromNotification.InsertOne", err.Error())
				}

				notificationOutputs = append(notificationOutputs, n)
			} else {
				uc.l.Info(ctx, "usecase.GetFromNotification.NotificationExists", fmt.Sprintf("Thông báo với ID %d đã tồn tại", n.ID))
			}
		}

		msgTexts := uc.createMsgNotification(ctx, notificationOutputs)

		for _, msgText := range msgTexts {
			error := uc.telegramUC.SendMessage(ctx, msgText)
			if error != nil {
				uc.l.Error(ctx, "Failed to send message to telegram", error)
			}
		}
	}

	return notificationOutputs, nil

}
