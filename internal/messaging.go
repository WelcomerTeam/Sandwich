package internal

import (
	"context"

	messaging "github.com/WelcomerTeam/Sandwich/messaging"
	"golang.org/x/xerrors"
)

type MQClient interface {
	String() string
	Channel() string
	Cluster() string

	Connect(ctx context.Context, clientName string, args map[string]interface{}) (err error)
	Subscribe(ctx context.Context, channel string) (err error)
	Unsubscribe()
	Chan() chan []byte
}

func NewMQClient(mqType string) (MQClient, error) {
	switch mqType {
	case "stan":
		return messaging.NewStanMQClient(), nil
	case "kafka":
		return messaging.NewKafkaMQClient(), nil
	case "redis":
		return messaging.NewRedisMQClient(), nil
	default:
		return nil, xerrors.New("No MQ client named " + mqType)
	}
}
