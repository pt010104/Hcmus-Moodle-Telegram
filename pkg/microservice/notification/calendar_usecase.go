package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/curl"
)

func (uc implUseCase) GetFromCalendar(ctx context.Context, input GetFromCalendarFilter) (Calendar, error) {
	h := map[string]string{
		"Content-Type": "application/json",
		"Cookie":       uc.cookies,
	}

	baseURL := curl.GetInternalUrl(uc.url, getFromNotificationEndpoint)
	u, err := url.Parse(baseURL)
	if err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.ParseURL", err.Error())
		return Calendar{}, err
	}

	query := u.Query()
	query.Set("info", calendarViewDay)
	query.Set("sesskey", uc.sessKey)
	u.RawQuery = query.Encode()

	reqArgs := CalendarRequestArgs{
		Year:  input.Year,
		Month: input.Month,
		Day:   input.Day,
	}

	reqBody := []CalendarRequest{
		{
			Index:      0,
			MethodName: calendarViewDay,
			Args:       reqArgs,
		},
	}

	response, err := curl.Post(u.String(), h, reqBody)
	if err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.Post", err.Error())
		return Calendar{}, err
	}

	var res []CalendarResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.Unmarshal", err.Error())
		return Calendar{}, err
	}

	// Check if there's an error in the response
	if len(res) == 0 {
		return Calendar{}, fmt.Errorf("empty response from calendar API")
	}

	if res[0].Error {
		if res[0].Exception != nil {
			return Calendar{}, fmt.Errorf("calendar API error: %s (code: %s)", res[0].Exception.Message, res[0].Exception.ErrorCode)
		}
		return Calendar{}, fmt.Errorf("calendar API error: unknown error")
	}

	var data struct {
		Events []CalendarEvent `json:"events"`
	}

	if err := json.Unmarshal(res[0].Data, &data); err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.UnmarshalData", err.Error())
		return Calendar{}, err
	}

	calendar := Calendar{
		Events: data.Events,
	}

	return calendar, nil
}
