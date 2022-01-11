package mqclients

import (
	"context"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"golang.org/x/xerrors"
)

func init() {
	MQClients = append(MQClients, "stan")
}

type StanMQClient struct {
	NatsClient *nats.Conn `json:"-"`
	StanClient stan.Conn  `json:"-"`

	channel string
	cluster string

	subscription *stan.Subscription
	msgChannel   chan []byte
}

func NewStanMQClient() (mqC *StanMQClient) {
	mqC = &StanMQClient{
		msgChannel: make(chan []byte, MsgChannelBuffer),
	}

	return
}

func (stanMQ *StanMQClient) String() string {
	return "stan"
}

func (stanMQ *StanMQClient) Channel() string {
	return stanMQ.channel
}

func (stanMQ *StanMQClient) Cluster() string {
	return stanMQ.cluster
}

func (stanMQ *StanMQClient) Connect(ctx context.Context, clientName string, args map[string]interface{}) (err error) {
	var ok bool

	var address string

	if address, ok = GetEntry(args, "Address").(string); !ok {
		return xerrors.New("stanMQ connect: string type assertion failed for Address")
	}

	var cluster string

	if cluster, ok = GetEntry(args, "Cluster").(string); !ok {
		return xerrors.New("stanMQ connect: string type assertion failed for Cluster")
	}

	var channel string

	if channel, ok = GetEntry(args, "Channel").(string); !ok {
		return xerrors.New("stanMQ connect: string type assertion failed for Channel")
	}

	stanMQ.cluster = cluster
	stanMQ.channel = channel

	var useNatsConnection bool

	if useNatsConnectionStr, ok := GetEntry(args, "UseNATSConnection").(string); ok {
		if useNatsConnection, err = strconv.ParseBool(useNatsConnectionStr); err != nil {
			useNatsConnection = true
		}
	} else {
		useNatsConnection = true
	}

	var option stan.Option

	if useNatsConnection {
		stanMQ.NatsClient, err = nats.Connect(address)
		if err != nil {
			return xerrors.Errorf("stanMQ connect nats: %w", err)
		}

		option = stan.NatsConn(stanMQ.NatsClient)
	} else {
		option = stan.NatsURL(address)
	}

	stanMQ.StanClient, err = stan.Connect(
		cluster,
		clientName,
		option,
	)
	if err != nil {
		return xerrors.Errorf("stanMQ connect stan: %w", err)
	}

	return nil
}

func (stanMQ *StanMQClient) Subscribe(ctx context.Context, channelName string) (err error) {
	if stanMQ.subscription != nil {
		stanMQ.Unsubscribe()
	}

	var handler stan.MsgHandler

	handler = func(msg *stan.Msg) { stanMQ.msgChannel <- msg.Data }
	sub, err := stanMQ.StanClient.Subscribe(channelName, handler)

	stanMQ.subscription = &sub

	return
}

func (stanMQ *StanMQClient) Unsubscribe() {
	if stanMQ.subscription != nil {
		sub := *stanMQ.subscription
		sub.Unsubscribe()
	}

	stanMQ.subscription = nil
}

func (stanMQ *StanMQClient) Chan() (ch chan []byte) {
	return stanMQ.msgChannel
}
