package models

type Notification struct {
	ID          int    `bson:"_id"`
	Subject     string `bson:"subject"`
	FullMessage string `bson:"full_message"`
	TimeCreated int64  `json:"timecreated"`
}
