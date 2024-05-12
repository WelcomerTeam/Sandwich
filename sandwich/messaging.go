package internal

import (
	"context"
	"fmt"

	messaging "github.com/WelcomerTeam/Sandwich/messaging"
)

type MQClient interface {
	String() string
	Channel() string

	Connect(ctx context.Context, clientName string, args map[string]interface{}) error
	Subscribe(ctx context.Context, channel string) error
	Unsubscribe(ctx context.Context)
	Chan() chan []byte
}

func NewMQClient(mqType string) (MQClient, error) {
	switch mqType {
	case "stan":
		return messaging.NewJetstreamMQClient(), nil
	case "kafka":
		return messaging.NewKafkaMQClient(), nil
	case "redis":
		return messaging.NewRedisMQClient(), nil
	default:
		panic(fmt.Sprintf(`NewMQClient(%s): no mq with this name`, mqType))
	}
}
