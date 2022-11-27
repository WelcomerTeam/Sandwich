package mqclients

import (
	"context"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	gotils_strconv "github.com/savsgio/gotils/strconv"
)

func init() {
	MQClients = append(MQClients, "redis")
}

type RedisMQClient struct {
	redisClient *redis.Client

	channel string
	cluster string

	pubsub *redis.PubSub

	msgChannel chan []byte
}

func NewRedisMQClient() (mqC *RedisMQClient) {
	mqC = &RedisMQClient{
		msgChannel: make(chan []byte, MsgChannelBuffer),
	}

	return
}

func (redisMQ *RedisMQClient) String() string {
	return "redis"
}

func (redisMQ *RedisMQClient) Channel() string {
	return redisMQ.channel
}

func (redisMQ *RedisMQClient) Cluster() string {
	return redisMQ.cluster
}

func (redisMQ *RedisMQClient) Connect(ctx context.Context, clientName string, args map[string]interface{}) error {
	var ok bool

	var address string

	if address, ok = GetEntry(args, "Address").(string); !ok {
		return errors.New("redisMQ connect: string type assertion failed for Address")
	}

	var password string

	if password, ok = GetEntry(args, "Password").(string); !ok {
		return errors.New("redisMQ connect: string type assertion failed for Password")
	}

	var db int
	var err error

	if dbStr, ok := GetEntry(args, "DB").(string); !ok {
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return errors.Errorf("redisMQ connect db atoi: %w", err)
		}
	}

	redisMQ.redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	err = redisMQ.redisClient.Ping(ctx).Err()
	if err != nil {
		return errors.Errorf("redisMQ connect ping: %w", err)
	}

	return nil
}

func (redisMQ *RedisMQClient) Subscribe(ctx context.Context, channel string) error {
	if redisMQ.pubsub != nil {
		redisMQ.Unsubscribe()
	}

	redisMQ.pubsub = redisMQ.redisClient.Subscribe(ctx, channel)

	go func(redisMQ *RedisMQClient) {
		channel := redisMQ.pubsub.Channel()

		for {
			msg := <-channel
			redisMQ.msgChannel <- gotils_strconv.S2B(msg.Payload)
		}
	}(redisMQ)

	return nil
}

func (redisMQ *RedisMQClient) Unsubscribe() {
	if redisMQ.pubsub != nil {
		pubsub := redisMQ.pubsub
		pubsub.Close()
	}

	redisMQ.pubsub = nil
}

func (redisMQ *RedisMQClient) Chan() (ch chan []byte) {
	return redisMQ.msgChannel
}
