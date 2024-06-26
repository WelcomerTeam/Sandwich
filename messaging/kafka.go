package mqclients

import (
	"context"
	"errors"
	"io"

	"github.com/segmentio/kafka-go"
)

func init() {
	MQClients = append(MQClients, "kafka")
}

type KafkaMQClient struct {
	KafkaClient *kafka.Reader

	address string
	channel string

	reader *kafka.Reader

	msgChannel chan []byte
}

func NewKafkaMQClient() (mqC *KafkaMQClient) {
	mqC = &KafkaMQClient{
		msgChannel: make(chan []byte, MsgChannelBuffer),
	}

	return
}

func (kafkaMQ *KafkaMQClient) String() string {
	return "kafka"
}

func (kafkaMQ *KafkaMQClient) Channel() string {
	return kafkaMQ.channel
}

func (kafkaMQ *KafkaMQClient) Connect(ctx context.Context, clientName string, args map[string]interface{}) error {
	var ok bool

	var address string

	if address, ok = GetEntry(args, "Address").(string); !ok {
		return errors.New("kafkaMQ connect: string type assertion failed for Address")
	}

	kafkaMQ.address = address

	return nil
}

func (kafkaMQ *KafkaMQClient) Subscribe(ctx context.Context, channelName string) error {
	if kafkaMQ.reader != nil {
		kafkaMQ.Unsubscribe(ctx)
	}

	kafkaMQ.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaMQ.address},
		Topic:   channelName,
	})

	go func(kafkaMQ *KafkaMQClient, ctx context.Context) {
		for {
			msg, err := kafkaMQ.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
			}

			kafkaMQ.msgChannel <- msg.Value
		}
	}(kafkaMQ, ctx)

	return nil
}

func (kafkaMQ *KafkaMQClient) Unsubscribe(ctx context.Context) {
	if kafkaMQ.reader != nil {
		reader := kafkaMQ.reader
		reader.Close()
	}

	kafkaMQ.reader = nil
}

func (kafkaMQ *KafkaMQClient) Chan() (ch chan []byte) {
	return kafkaMQ.msgChannel
}
