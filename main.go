package main

import (
	"context"
	"os"
	"time"

	"github.com/WelcomerTeam/Sandwich-Daemon/structs"
	sandwich "github.com/WelcomerTeam/Sandwich/internal"
	messaging "github.com/WelcomerTeam/Sandwich/messaging"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:15000", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Stamp,
	}

	log := zerolog.New(writer).With().Timestamp().Logger()
	log.Info().Msg("startup")

	sandwichClient := sandwich.NewSandwich(conn, writer)

	welcomerBot := sandwich.NewBot()

	welcomerBot.RegisterOnMessageCreateEvent(func(ctx *sandwich.Context, message sandwich.Message) (err error) {
		ctx.Logger.Info().Str("author", message.Author.Username+"#"+message.Author.Discriminator).Msg(message.Content)

		return
	})

	sandwichClient.RegisterBot("welcomer", welcomerBot)

	ctx := context.TODO()

	mqC := messaging.NewStanMQClient()
	err = mqC.Connect(ctx, "sdc", map[string]interface{}{
		"Address": "localhost",
		"Cluster": "cluster",
		"Channel": "channel",
	})

	if err != nil {
		panic(err)
	}

	err = mqC.Subscribe(ctx, "sandwich")

	if err != nil {
		panic(err)
	}

	c := mqC.Chan()

	var p structs.SandwichPayload

	for {
		select {
		case m := <-c:
			if err := jsoniter.Unmarshal(m, &p); err == nil {
				err = sandwichClient.DispatchSandwichPayload(p)
				if err != nil {
					println(err.Error(), string(m))
				}
			} else {
				println(err.Error(), string(m))
			}
		}
	}
}
