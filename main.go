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
	"golang.org/x/xerrors"
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

	bot := sandwich.NewBot(sandwich.StaticPrefixCheck("?"))

	bot.MustAddCommand(&sandwich.Commandable{
		Name: "argumenttest",
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Required:     false,
				ArgumentType: sandwich.ArgumentTypeTextChannel,
				Name:         "test",
			},
		},
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			val := ctx.MustGetArgument("test").MustChannel()

			if val != nil {
				println("Got channel", val.Name)
			} else {
				println("No channel was found...")
			}

			return nil
		},
	})

	bot.RegisterOnMessageCreateEvent(func(ctx *sandwich.EventContext, message sandwich.Message) (err error) {
		err = bot.ProcessCommands(ctx, message)
		if err != nil {
			ctx.Logger.Warn().Err(err).Str("content", message.Content).Msg("Failed to process command")

			return xerrors.Errorf("Failed to process command: %v", err)
		}

		return nil
	})

	sandwichClient.RegisterBot("welcomer", bot)

	ctx := context.Background()

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
				sandwichClient.DispatchSandwichPayload(ctx, p)
			} else {
				println(err.Error(), string(m))
			}
		}
	}
}
