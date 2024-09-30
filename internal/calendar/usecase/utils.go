package usecase

import (
	"context"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"time"

	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/models"
	"github.com/pt010104/Hcmus-Moodle-Telegram/util"
)

func (uc implUseCase) createMsgCalendarForTelegram(ctx context.Context, events []models.Calendar) []string {

	msgTexts := make([]string, 0)

	for _, event := range events {
		eventName := "N/A"
		if event.Name != "" {
			eventName = event.Name
		}

		eventCourseName := "N/A"
		if event.CourseName != "" {
			eventCourseName = html.UnescapeString(event.CourseName)
		}

		formattedTime := "N/A"
		if event.FormattedTime != "" {
			formattedTime = util.ExtractTextFromHTML(event.FormattedTime)
		}

		link := "N/A"
		if event.URL != "" {
			link = event.URL
		}

		messageText := fmt.Sprintf(
			"<b>🟢 New:</b> %s\n"+
				"<b>🟢 Course:</b> %s\n"+
				"<b>🟢 Deadline:</b> %s\n"+
				"%s",
			eventName,
			eventCourseName,
			html.EscapeString(formattedTime),
			link,
		)

		msgTexts = append(msgTexts, messageText)

	}

	return msgTexts

}

// type Notification struct {
// 	ID          int    `bson:"_id"`
// 	Subject     string `bson:"subject"`
// 	FullMessage string `bson:"full_message"`
// 	TimeCreated string `json:"timecreated"`
// }

func (uc implUseCase) createMsgNotification(ctx context.Context, notifications []models.Notification) []string {

	msgTexts := make([]string, 0)

	for _, notification := range notifications {
		subject := "N/A"
		if notification.Subject != "" {
			subject = notification.Subject
		}

		fullMessage := "N/A"
		if notification.FullMessage != "" {
			fullMessage = util.ExtractTextFromHTML(notification.FullMessage)
		}

		timeCreated := time.Unix(notification.TimeCreated, 0).Format("2006-01-02 15:04:05")

		messageText := fmt.Sprintf(
			"<b>🟢 Subject:\n</b> %s\n"+
				"<b>🟢 Message:\n</b> %s\n"+
				"<b>🟢 Created:\n</b> %s",
			subject,
			fullMessage,
			timeCreated,
		)

		msgTexts = append(msgTexts, messageText)

	}

	return msgTexts

}

func (uc implUseCase) extractEventTime(formattedTime string) (time.Time, error) {
	re := regexp.MustCompile(`time=(\d+)`)
	matches := re.FindStringSubmatch(formattedTime)
	if len(matches) < 2 {
		return time.Time{}, errors.New("could not find time parameter in formattedTime")
	}
	timestampStr := matches[1]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	eventTime := time.Unix(timestamp, 0)
	return eventTime, nil
}
