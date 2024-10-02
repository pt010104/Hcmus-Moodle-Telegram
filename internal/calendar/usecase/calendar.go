package usecase

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/models"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/microservice/notification"
	"github.com/pt010104/Hcmus-Moodle-Telegram/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (uc implUseCase) GetFromCalendar(ctx context.Context) ([]models.Calendar, error) {

	date := util.Now()

	dateEnd := date.AddDate(0, 0, 20)

	var allCalendarOutputs []models.Calendar

	for currentDate := date; !currentDate.After(dateEnd); currentDate = currentDate.AddDate(0, 0, 1) {
		calendarSrv, err := uc.notificationSrv.GetFromCalendar(ctx, notification.GetFromCalendarFilter{
			Year:  currentDate.Format("2006"),
			Month: currentDate.Format("1"),
			Day:   currentDate.Format("2"),
		})
		if err != nil {
			return nil, err
		}

		var calendarEvents []models.Calendar

		for _, event := range calendarSrv.Events {
			eventTime, err := uc.extractEventTime(event.FormattedTime)
			if err != nil {
				uc.l.Error(ctx, "Failed to extract event time", err.Error())
				continue
			}

			calendarEvents = append(calendarEvents, models.Calendar{
				ID:            event.ID,
				Name:          event.Name,
				Description:   event.Description,
				FormattedTime: event.FormattedTime,
				CourseName:    event.Course.FullName,
				CourseID:      event.Course.ID,
				URL:           event.URL,
				Deadline:      eventTime,
			})

			now := time.Now()
			timeDiff := eventTime.Sub(now)
			if timeDiff > 0 && timeDiff <= 2*time.Hour {
				messageText := fmt.Sprintf(
					"<b>Thông báo:</b> có deadline trong 2 tiếng nữa\n"+
						"<b>Môn:</b> %s\n"+
						"<b>Deadline:</b> %s\n"+
						"%s",
					html.EscapeString(event.Course.FullName),
					eventTime.Format("2006-01-02 15:04:05"),
					event.URL,
				)
				err := uc.telegramUC.SendMessage(ctx, messageText)
				if err != nil {
					uc.l.Error(ctx, "Failed to send deadline approaching message to telegram", err)
				}
			}
		}

		if len(calendarEvents) > 0 {
			collection := uc.db.Collection("calendar_events")
			var calendarOutputs []models.Calendar

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
					uc.l.Info(ctx, "usecase.GetFromCalendar.EventExists", fmt.Sprintf(" Event with ID %d already exists", c.ID))
				}
			}

			if len(calendarOutputs) > 0 {
				allCalendarOutputs = append(allCalendarOutputs, calendarOutputs...)

				msgTexts := uc.createMsgCalendarForTelegram(ctx, calendarOutputs)

				for _, msgText := range msgTexts {
					err := uc.telegramUC.SendMessage(ctx, msgText)
					if err != nil {
						uc.l.Error(ctx, "Failed to send message to telegram", err)
					}
				}
			}
		}
	}

	return allCalendarOutputs, nil
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
				uc.l.Info(ctx, "usecase.GetFromNotification.NotificationExists ", n.ID)
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
