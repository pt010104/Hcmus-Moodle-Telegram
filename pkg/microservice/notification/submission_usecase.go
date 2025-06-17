package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/curl"
)

func (uc implUseCase) CheckEventSubmission(ctx context.Context, input EventSubmissionFilter) (EventDetail, error) {
	h := map[string]string{
		"Content-Type": "application/json",
		"Cookie":       uc.cookies,
	}

	baseURL := curl.GetInternalUrl(uc.url, getFromNotificationEndpoint)
	u, err := url.Parse(baseURL)
	if err != nil {
		uc.l.Error(ctx, "notification.usecase.CheckEventSubmission.ParseURL", err.Error())
		return EventDetail{}, err
	}

	query := u.Query()
	query.Set("info", calendarEventByID)
	query.Set("sesskey", uc.sessKey)
	u.RawQuery = query.Encode()

	reqArgs := EventSubmissionArgs{
		EventID: input.EventID,
	}

	reqBody := []EventSubmissionRequest{
		{
			Index:      0,
			MethodName: calendarEventByID,
			Args:       reqArgs,
		},
	}

	response, err := curl.Post(u.String(), h, reqBody)
	if err != nil {
		uc.l.Error(ctx, "notification.usecase.CheckEventSubmission.Post", err.Error())
		return EventDetail{}, err
	}

	var res []EventSubmissionResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		uc.l.Error(ctx, "notification.usecase.CheckEventSubmission.Unmarshal", err.Error())
		return EventDetail{}, err
	}

	// Check if there's an error in the response
	if len(res) == 0 {
		return EventDetail{}, fmt.Errorf("empty response from submission check API")
	}

	if res[0].Error {
		if res[0].Exception != nil {
			return EventDetail{}, fmt.Errorf("submission check API error: %s (code: %s)", res[0].Exception.Message, res[0].Exception.ErrorCode)
		}
		return EventDetail{}, fmt.Errorf("submission check API error: unknown error")
	}

	var data EventSubmissionData
	if err := json.Unmarshal(res[0].Data, &data); err != nil {
		uc.l.Error(ctx, "notification.usecase.CheckEventSubmission.UnmarshalData", err.Error())
		return EventDetail{}, err
	}

	return data.Event, nil
}
