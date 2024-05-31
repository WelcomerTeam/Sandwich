package mqclients

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nuid"
	"github.com/pkg/errors"
)

func init() {
	MQClients = append(MQClients, "jetstream")
}

type JetstreamMQClient struct {
	JetStreamClient jetstream.JetStream `json:"-"`
	ConsumeContext  *jetstream.ConsumeContext

	channel string

	msgChannel chan []byte
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

func (jetstreamMQ *JetstreamMQClient) Connect(ctx context.Context, clientName string, args map[string]interface{}) error {
	var ok bool

	var address string

	if address, ok = GetEntry(args, "Address").(string); !ok {
		return errors.New("jetstreamMQ connect: string type assertion failed for Address")
	}

	var channel string

	if channel, ok = GetEntry(args, "Channel").(string); !ok {
		return errors.New("jetstreamMQ connect: string type assertion failed for Channel")
	}

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
	if jetstreamMQ.ConsumeContext != nil {
		jetstreamMQ.Unsubscribe(ctx)
	}

	handler := func(msg jetstream.Msg) { jetstreamMQ.msgChannel <- msg.Data() }

	var consumer jetstream.Consumer
	var err error

	if v := mustParseBool(os.Getenv("JETSTREAM_USE_INTEREST_POLICY")); v {
		consumer, err = jetstreamMQ.JetStreamClient.OrderedConsumer(
			ctx,
			jetstreamMQ.channel,
			jetstream.OrderedConsumerConfig{
				FilterSubjects: []string{jetstreamMQ.channel + ".*"},
			},
		)
	} else {
		consumer, err = jetstreamMQ.JetStreamClient.CreateOrUpdateConsumer(
			ctx,
			jetstreamMQ.channel,
			jetstream.ConsumerConfig{
				Name:           nuid.Next(),
				DeliverPolicy:  jetstream.DeliverAllPolicy,
				AckPolicy:      jetstream.AckExplicitPolicy,
				FilterSubjects: []string{jetstreamMQ.channel + ".*"},
			},
		)
	}

	if err != nil {
		return fmt.Errorf("failed to create jetstream consumer: %w", err)
	}

	consumeContext, err := consumer.Consume(handler)
	if err != nil {
		return fmt.Errorf("failed to create consume context: %w", err)
	}

	jetstreamMQ.ConsumeContext = &consumeContext

	return nil
}

func (jetstreamMQ *JetstreamMQClient) Unsubscribe(ctx context.Context) {
	if jetstreamMQ.ConsumeContext != nil {
		consumeContext := *jetstreamMQ.ConsumeContext
		consumeContext.Drain()
	}

	jetstreamMQ.ConsumeContext = nil
}

func (jetstreamMQ *JetstreamMQClient) Chan() (ch chan []byte) {
	return jetstreamMQ.msgChannel
}

func mustParseBool(str string) bool {
	boolean, _ := strconv.ParseBool(str)

	return boolean
}
