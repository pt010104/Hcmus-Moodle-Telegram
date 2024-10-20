package models

import "time"

type Calendar struct {
	ID            int       `bson:"_id"`
	Name          string    `bson:"name"`
	Description   string    `bson:"description"`
	FormattedTime string    `bson:"formattedtime"`
	CourseID      int       `bson:"course_id"`
	CourseName    string    `bson:"course_name"`
	URL           string    `bson:"url"`
	Deadline      time.Time `bson:"deadline"`
	TimeRemind    int       `bson:"time_remind"` //0 for 12, 1 for 6, 2 for 3, 3 for 1
}
