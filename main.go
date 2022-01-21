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

	welcomerBot := sandwich.NewBot(sandwich.StaticPrefixCheck("?"))

	welcomerGroup, err := welcomerBot.AddGroup(true, "welcomer")
	if err != nil {
		panic(err)
	}

	command, _ := welcomerBot.AddCommand(func(ctx *sandwich.CommandContext) (err error) {
		println("Arguments test ran!")

		println(ctx.Guild)
		println(ctx.Author)
		println(ctx.Content)

		for k, v := range ctx.Arguments {
			println(k, v.ArgumentType)
		}

		return nil
	}, "argumentstest")
	command.AddArguments(
		sandwich.ArgumentParameter{
			Required:     false,
			ArgumentType: sandwich.ArgumentTypeString,
			Name:         "test",
		},
		sandwich.ArgumentParameter{
			Required:     false,
			ArgumentType: sandwich.ArgumentTypeString,
			Name:         "test2",
		},
	)

	welcomerGroup.SetHandler(func(ctx *sandwich.CommandContext) (err error) {
		if ctx.InvokedSubcommand != nil {
			ctx.EventContext.Logger.Info().Str("subcommand", ctx.InvokedSubcommand.Name).Msg("welcomer group was ran.")
		} else {
			ctx.EventContext.Logger.Info().Msg("welcomer group was ran.")
		}

		return nil
	})

	welcomerGroup.AddCommand(func(ctx *sandwich.CommandContext) (err error) {
		ctx.EventContext.Logger.Info().Msg("welcomer test ran")

		return nil
	}, "test")

	welcomerGroup.AddCommand(func(ctx *sandwich.CommandContext) (err error) {
		ctx.EventContext.Logger.Info().Msg("welcomer add ran")

		return nil
	}, "add")

	welcomerGroup.AddCommand(func(ctx *sandwich.CommandContext) (err error) {
		ctx.EventContext.Logger.Info().Msg("welcomer error ran")

		return xerrors.New("Random error")
	}, "error")

	welcomerBot.AddCommand(func(ctx *sandwich.CommandContext) (err error) {
		ctx.EventContext.Logger.Info().Msg("test ran successfuly")

		return nil
	}, "test", "alias", "testme")

	welcomerBot.RegisterOnMessageCreateEvent(func(ctx *sandwich.EventContext, message sandwich.Message) (err error) {
		err = welcomerBot.ProcessCommands(ctx, message)
		if err != nil {
			ctx.Logger.Warn().Err(err).Str("content", message.Content).Msg("Failed to process command")
		}

		return err
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
