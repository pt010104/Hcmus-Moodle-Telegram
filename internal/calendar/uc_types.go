package calendar

type GetFromCalendarInput struct {
	Year  string
	Month string
	Day   string
}

type GetFromNotificationInput struct {
	Limit    int
	Offset   int
	UserIDTo string
}
