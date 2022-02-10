package main

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"time"

	discord "github.com/WelcomerTeam/Discord/discord"

	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
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

	proxyURL, err := url.Parse("http://localhost:5001/")
	if err != nil {
		panic(err.Error())
	}

	restInterface := discord.NewTwilightProxy(*proxyURL)

	sandwichClient := sandwich.NewSandwich(conn, restInterface, writer)

	bot := sandwich.NewBot(sandwich.StaticPrefixCheck("?"))

	bot.Commands.MustAddCommand(&sandwich.Commandable{
		Name:    "avatar",
		Aliases: []string{"profile"},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Required:     true,
				ArgumentType: sandwich.ArgumentTypeUser,
				Name:         "user",
			},
		},
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			user := ctx.MustGetArgument("user").MustUser()

			avatarURL := discord.EndpointCDN + discord.EndpointUserAvatar(user.ID.String(), user.Avatar)

			ctx.Reply(ctx.EventContext.Session,
				*discord.NewMessage("").
					AddEmbed(
						*discord.NewEmbed(
							discord.EmbedTypeDefault,
						).SetImage(discord.NewEmbedImage(avatarURL)),
					),
			)

			return
		},
	})

	bot.Commands.MustAddCommand(&sandwich.Commandable{
		Name: "filetest",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			ctx.Reply(ctx.EventContext.Session,
				*discord.NewMessage("").
					AddFile(discord.File{
						Name:        "file.txt",
						ContentType: "application/octet-stream",
						Reader:      bytes.NewBufferString("Hello world!"),
					}),
			)

			return
		},
	})

	bot.RegisterOnInteractionCreateEvent(func(ctx *sandwich.EventContext, interaction discord.Interaction) (err error) {
		resp, err := bot.ProcessInteraction(ctx, interaction)
		if err != nil {
			println(err.Error())

			return
		}

		if resp != nil {
			err = interaction.SendResponse(ctx.Session, *resp.Type, *resp.Data.WebhookMessageParams, resp.Data.Choices)
			if err != nil {
				println(err.Error())

				return
			}
		}

		return nil
	})

	bot.RegisterOnMessageCreateEvent(func(ctx *sandwich.EventContext, message discord.Message) (err error) {
		err = bot.ProcessCommands(ctx, message)
		if err != nil {
			ctx.Logger.Warn().Err(err).Str("content", message.Content).Msg("Failed to process command")

			return xerrors.Errorf("Failed to process command: %v", err)
		}

		return nil
	})

	bot.RegisterOnInteractionCreateEvent(func(eventCtx *sandwich.EventContext, interaction discord.Interaction) (err error) {
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

	var p sandwich_structs.SandwichPayload

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
