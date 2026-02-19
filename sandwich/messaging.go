package internal

import (
	"context"
)

type MQClient interface {
	String() string
	Channel() string

	Connect(ctx context.Context, clientName string, args map[string]any) error
	Subscribe(ctx context.Context, channel string) error
	Unsubscribe(ctx context.Context)
	Chan() chan []byte
}
