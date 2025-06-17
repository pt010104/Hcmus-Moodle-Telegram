package models

import "time"

type Calendar struct {
	ID                 int       `bson:"_id"`
	Name               string    `bson:"name"`
	Description        string    `bson:"description"`
	FormattedTime      string    `bson:"formattedtime"`
	CourseID           int       `bson:"course_id"`
	CourseName         string    `bson:"course_name"`
	URL                string    `bson:"url"`
	Deadline           time.Time `bson:"deadline"`
	TimeRemind         int       `bson:"time_remind"` //0 for 24h, 1 for 12h, 2 for 6h, 3 for 3h, 4 for 1h, 5 for 30m
	IsSubmitted        bool      `bson:"is_submitted"`
	SubmissionNotified bool      `bson:"submission_notified"`
}
