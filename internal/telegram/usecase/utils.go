package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/models"
	"go.mongodb.org/mongo-driver/mongo"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (uc implUseCase) sendTextMessage(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := uc.bot.Send(msg)
	if err != nil {
		uc.l.Error(ctx, "Failed to send message", err.Error())
	}
	return err
}

func (uc implUseCase) handleListCourses(ctx context.Context, message *tgbotapi.Message) error {
	collection := uc.db.Collection("calendar_events")

	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", "$course_id"},
			{"course_name", bson.D{{"$first", "$course_name"}}},
		}}},
		{{"$sort", bson.D{{"course_name", 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		uc.l.Error(ctx, "Failed to aggregate courses", err.Error())
		return uc.sendTextMessage(ctx, message.Chat.ID, "Failed to retrieve courses")
	}
	defer cursor.Close(ctx)

	var results []struct {
		CourseID   int64  `bson:"_id"`
		CourseName string `bson:"course_name"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		uc.l.Error(ctx, "Failed to decode results", err.Error())
		return uc.sendTextMessage(ctx, message.Chat.ID, "Failed to retrieve courses")
	}

	if len(results) == 0 {
		return uc.sendTextMessage(ctx, message.Chat.ID, "No courses found")
	}

	var sb strings.Builder
	sb.WriteString("List of courses:\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("- %s (%d)\n", r.CourseName, r.CourseID))
	}

	return uc.sendTextMessage(ctx, message.Chat.ID, sb.String())
}

func (uc implUseCase) handleListDeadlines(ctx context.Context, message *tgbotapi.Message) error {
	collection := uc.db.Collection("calendar_events")

	now := time.Now()

	filter := bson.M{
		"deadline": bson.M{
			"$gte": now,
		},
	}

	opts := options.Find()
	opts.SetSort(bson.D{{"deadline", 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		uc.l.Error(ctx, "Failed to find deadlines", err.Error())
		return uc.sendTextMessage(ctx, message.Chat.ID, "Failed to retrieve deadlines")
	}
	defer cursor.Close(ctx)

	var results []models.Calendar

	if err = cursor.All(ctx, &results); err != nil {
		uc.l.Error(ctx, "Failed to decode results", err.Error())
		return uc.sendTextMessage(ctx, message.Chat.ID, "Failed to retrieve deadlines")
	}

	if len(results) == 0 {
		return uc.sendTextMessage(ctx, message.Chat.ID, "No upcoming deadlines found")
	}

	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	var sb strings.Builder
	for _, r := range results {
		deadlineLocal := r.Deadline.In(location)
		timeDiff := r.Deadline.Sub(now)
		diffString := formatTimeDifference(timeDiff)

		sb.WriteString(fmt.Sprintf("🟢 %s\n", r.CourseName))
		sb.WriteString(fmt.Sprintf("    + %s\n", r.Name))
		sb.WriteString(fmt.Sprintf("    + %s - %s\n", deadlineLocal.Format("2006-01-02T15:04:05-07:00"), diffString))
		sb.WriteString(fmt.Sprintf("    + %s\n", r.URL))
		sb.WriteString("\n")
	}

	return uc.sendTextMessage(ctx, message.Chat.ID, sb.String())
}

func (uc implUseCase) handleCourseDeadlines(ctx context.Context, message *tgbotapi.Message) error {
	parts := strings.Fields(message.Text)
	if len(parts) < 2 {
		return uc.sendTextMessage(ctx, message.Chat.ID, "Please provide a course ID. Usage: /cd <course_id>")
	}
	courseIDStr := parts[1]
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		return uc.sendTextMessage(ctx, message.Chat.ID, "Invalid course ID. Usage: /cd <course_id>")
	}

	collection := uc.db.Collection("calendar_events")

	now := time.Now()

	filter := bson.M{
		"course_id": courseID,
		"deadline": bson.M{
			"$gte": now,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		uc.l.Error(ctx, "Failed to find deadlines for course", err.Error())
		return uc.sendTextMessage(ctx, message.Chat.ID, "Failed to retrieve deadlines for the course")
	}
	defer cursor.Close(ctx)

	var results []models.Calendar

	if err = cursor.All(ctx, &results); err != nil {
		uc.l.Error(ctx, "Failed to decode results", err.Error())
		return uc.sendTextMessage(ctx, message.Chat.ID, "Failed to retrieve deadlines for the course")
	}

	if len(results) == 0 {
		return uc.sendTextMessage(ctx, message.Chat.ID, "No deadlines found for the course")
	}

	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	var sb strings.Builder
	for _, r := range results {
		timeDiff := r.Deadline.Sub(now)
		diffString := formatTimeDifference(timeDiff)
		dlLocal := r.Deadline.In(location)

		sb.WriteString(fmt.Sprintf("🟢 %s\n", r.Name))
		sb.WriteString(fmt.Sprintf("    + %s - %s\n", dlLocal.Format("2006-01-02T15:04:05-07:00"), diffString))
		sb.WriteString(fmt.Sprintf("    + %s\n", r.URL))
		sb.WriteString("\n")
	}

	return uc.sendTextMessage(ctx, message.Chat.ID, sb.String())
}

func formatTimeDifference(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d ngày", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d giờ", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d phút", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d giây", seconds))
	}

	if len(parts) == 0 {
		return "0 giây"
	}

	return strings.Join(parts, " ")
}
