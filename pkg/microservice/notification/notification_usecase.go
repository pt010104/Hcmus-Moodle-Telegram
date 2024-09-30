package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/curl"
)

func (uc implUseCase) GetFromNotification(ctx context.Context, input GetFromNotificationFilter) (Notification, error) {
	h := map[string]string{
		"Content-Type": "application/json",
		"Cookie":       uc.cookies,
	}

	baseURL := curl.GetInternalUrl(uc.url, getFromNotificationEndpoint)
	u, err := url.Parse(baseURL)
	if err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.ParseURL", err.Error())
		return Notification{}, err
	}

	query := u.Query()
	query.Set("info", notificationGet)
	query.Set("sesskey", uc.sessKey)
	u.RawQuery = query.Encode()

	reqArgs := NotificationRequestArgs{
		Limit:    input.Limit,
		Offset:   input.Offset,
		UserIDTo: input.UserIDTo,
	}

	reqBody := []NotificationRequest{
		{
			Index:      0,
			MethodName: notificationGet,
			Args:       reqArgs,
		},
	}

	response, err := curl.Post(u.String(), h, reqBody)
	if err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.Post", err.Error())
		return Notification{}, err
	}

	var res []NotificationResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.Unmarshal", err.Error())
		return Notification{}, err
	}

	var data struct {
		Events []NotificationEvent `json:"notifications"`
	}

	if err := json.Unmarshal(res[0].Data, &data); err != nil {
		uc.l.Error(ctx, "notification.usecase.GetFromCalendar.UnmarshalData", err.Error())
		return Notification{}, err
	}

	notification := Notification{
		Events: data.Events,
	}

	fmt.Println(notification)

	return notification, nil
}
