package telegram

import "context"

type UseCase interface {
	SendMessage(ctx context.Context, message string) error
}
