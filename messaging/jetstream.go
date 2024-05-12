package mqclients

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

func init() {
	MQClients = append(MQClients, "jetstream")
}

type JetstreamMQClient struct {
	JetStreamClient jetstream.JetStream `json:"-"`

	channel string
	cluster string

	consumeContext *jetstream.ConsumeContext
	msgChannel     chan []byte
}

func NewJetstreamMQClient() (mqC *JetstreamMQClient) {
	mqC = &JetstreamMQClient{
		msgChannel: make(chan []byte, MsgChannelBuffer),
	}

	return
}

func (jetstreamMQ *JetstreamMQClient) String() string {
	return "jetstream"
}

func (jetstreamMQ *JetstreamMQClient) Channel() string {
	return jetstreamMQ.channel
}

func (jetstreamMQ *JetstreamMQClient) Cluster() string {
	return jetstreamMQ.cluster
}

func (jetstreamMQ *JetstreamMQClient) Connect(ctx context.Context, clientName string, args map[string]interface{}) error {
	var ok bool

	var address string

	if address, ok = GetEntry(args, "Address").(string); !ok {
		return errors.New("jetstreamMQ connect: string type assertion failed for Address")
	}

	var cluster string

	if cluster, ok = GetEntry(args, "Cluster").(string); !ok {
		return errors.New("jetstreamMQ connect: string type assertion failed for Cluster")
	}

	var channel string

	if channel, ok = GetEntry(args, "Channel").(string); !ok {
		return errors.New("jetstreamMQ connect: string type assertion failed for Channel")
	}

	jetstreamMQ.cluster = cluster
	jetstreamMQ.channel = channel

	nc, err := nats.Connect(address)
	if err != nil {
		return fmt.Errorf("jetstreamMQ connect nats: %w", err)
	}

	jetstreamMQ.JetStreamClient, err = jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("jetstreamMQ new: %w", err)
	}

	return nil
}

func (jetstreamMQ *JetstreamMQClient) Subscribe(ctx context.Context, channelName string) error {
	if jetstreamMQ.consumeContext != nil {
		jetstreamMQ.Unsubscribe(ctx)
	}

	handler := func(msg jetstream.Msg) { jetstreamMQ.msgChannel <- msg.Data() }
	consumer, err := jetstreamMQ.JetStreamClient.CreateConsumer(ctx, channelName, jetstream.ConsumerConfig{})
	if err != nil {
		return fmt.Errorf("failed to create jetstream consumer: %w", err)
	}

	consumeContext, err := consumer.Consume(handler)
	if err != nil {
		return fmt.Errorf("failed to create consume context: %w", err)
	}

	jetstreamMQ.consumeContext = &consumeContext

	return nil
}

func (jetstreamMQ *JetstreamMQClient) Unsubscribe(ctx context.Context) {
	if jetstreamMQ.consumeContext != nil {
		consumeContext := *jetstreamMQ.consumeContext
		consumeContext.Drain()
	}

	jetstreamMQ.consumeContext = nil
}

func (jetstreamMQ *JetstreamMQClient) Chan() (ch chan []byte) {
	return jetstreamMQ.msgChannel
}
