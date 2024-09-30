package notification

import "encoding/json"

type NotificationEvent struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	FullMessage string `json:"fullmessage"`
	TimeCreated int64  `json:"timecreated"`
}

type Notification struct {
	Events []NotificationEvent `json:"notifications"`
}

type GetFromNotificationFilter struct {
	Limit    int
	Offset   int
	UserIDTo string
}

type NotificationRequestArgs struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	UserIDTo string `json:"useridto"`
}

type NotificationRequest struct {
	Index      int                     `json:"index"`
	MethodName string                  `json:"methodname"`
	Args       NotificationRequestArgs `json:"args"`
}
type NotificationResponse struct {
	Error bool            `json:"error"`
	Data  json.RawMessage `json:"data"`
}
