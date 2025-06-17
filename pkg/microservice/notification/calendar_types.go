package notification

import "encoding/json"

type Course struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
}

type CalendarEvent struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	FormattedTime string `json:"formattedtime"`
	Course        Course `json:"course"`
	URL           string `json:"url"`
}

type Calendar struct {
	Events []CalendarEvent `json:"events"`
}

type GetFromCalendarFilter struct {
	Year  string
	Month string
	Day   string
}

type CalendarRequestArgs struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

type CalendarRequest struct {
	Index      int                 `json:"index"`
	MethodName string              `json:"methodname"`
	Args       CalendarRequestArgs `json:"args"`
}

type CalendarResponse struct {
	Error     bool            `json:"error"`
	Data      json.RawMessage `json:"data,omitempty"`
	Exception *ErrorException `json:"exception,omitempty"`
}

type ErrorException struct {
	Message     string `json:"message"`
	ErrorCode   string `json:"errorcode"`
	Link        string `json:"link,omitempty"`
	MoreInfoURL string `json:"moreinfourl,omitempty"`
}

// Types for checking event submission status
type EventSubmissionFilter struct {
	EventID int
}

type EventSubmissionRequest struct {
	Index      int                 `json:"index"`
	MethodName string              `json:"methodname"`
	Args       EventSubmissionArgs `json:"args"`
}

type EventSubmissionArgs struct {
	EventID int `json:"eventid"`
}

type EventSubmissionResponse struct {
	Error     bool            `json:"error"`
	Data      json.RawMessage `json:"data,omitempty"`
	Exception *ErrorException `json:"exception,omitempty"`
}

type EventSubmissionData struct {
	Event EventDetail `json:"event"`
}

type EventDetail struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Course Course      `json:"course"`
	Action *ActionInfo `json:"action,omitempty"` // nil means submitted, non-nil means not submitted
}

type ActionInfo struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	ItemCount     int    `json:"itemcount"`
	Actionable    bool   `json:"actionable"`
	ShowItemCount bool   `json:"showitemcount"`
}
