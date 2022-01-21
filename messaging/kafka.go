package mqclients

import (
	"context"
	"github.com/segmentio/kafka-go"
	"golang.org/x/xerrors"
	"io"
)

func init() {
	MQClients = append(MQClients, "kafka")
}

type KafkaMQClient struct {
	KafkaClient *kafka.Reader

	channel string
	cluster string
	address string

	reader *kafka.Reader

	msgChannel chan []byte
}

func parseKafkaBalancer(balancer string) kafka.Balancer {
	switch balancer {
	case "crc32":
		return &kafka.CRC32Balancer{}
	case "hash":
		return &kafka.Hash{}
	case "murmur2":
		return &kafka.Murmur2Balancer{}
	case "roundrobin":
		return &kafka.RoundRobin{}
	case "leastbytes":
		return &kafka.LeastBytes{}
	default:
		return nil
	}
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

func (kafkaMQ *KafkaMQClient) Cluster() string {
	return kafkaMQ.cluster
}

func (kafkaMQ *KafkaMQClient) Connect(ctx context.Context, clientName string, args map[string]interface{}) (err error) {
	var ok bool

	var address string

	if address, ok = GetEntry(args, "Address").(string); !ok {
		return xerrors.New("kafkaMQ connect: string type assertion failed for Address")
	}

	kafkaMQ.address = address

	return nil
}

func (kafkaMQ *KafkaMQClient) Subscribe(ctx context.Context, channelName string) (err error) {
	if kafkaMQ.reader != nil {
		kafkaMQ.Unsubscribe()
	}

	kafkaMQ.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaMQ.address},
		Topic:   channelName,
	})

	go func(kafkaMQ *KafkaMQClient, ctx context.Context) {
		for {
			msg, err := kafkaMQ.reader.ReadMessage(ctx)
			if err != nil {
				if xerrors.Is(err, io.EOF) {
					return
				}
			}

			kafkaMQ.msgChannel <- msg.Value
		}
	}(kafkaMQ, ctx)

	return nil
}

func (kafkaMQ *KafkaMQClient) Unsubscribe() {
	if kafkaMQ.reader != nil {
		reader := kafkaMQ.reader
		reader.Close()
	}

	kafkaMQ.reader = nil
}

func (kafkaMQ *KafkaMQClient) Chan() (ch chan []byte) {
	return kafkaMQ.msgChannel
}
