package main

import (
	"context"
	"net/url"
	"os"
	"time"

	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
	messaging "github.com/WelcomerTeam/Sandwich/messaging"
	sandwich "github.com/WelcomerTeam/Sandwich/sandwich"
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

	session := sandwich.NewTwilightProxy(*proxyURL)

	sandwichClient := sandwich.NewSandwich(conn, session, writer)

	bot := sandwich.NewBot(sandwich.StaticPrefixCheck("?"))

	argumentTestGroup := bot.MustAddCommand(&sandwich.Commandable{
		Name: "argumentTest",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			if ctx.InvokedSubcommand == nil {
				println("Usage:")

				for _, command := range ctx.Command.GetAllCommands() {
					println(command.Name, command.IsGroup())
				}
			}

			return nil
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Snowflake",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustSnowflake()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeSnowflake,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Member",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustMember()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeMember,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "User",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustUser()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeUser,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "TextChannel",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustChannel()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeTextChannel,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Invite",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustInvite()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeInvite,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Guild",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustGuild()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeGuild,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Role",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustRole()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeRole,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Activity",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustActivity()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeActivity,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Colour",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustColour()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeColour,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "VoiceChannel",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustChannel()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeVoiceChannel,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "StageChannel",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustChannel()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeStageChannel,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Emoji",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustEmoji()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeEmoji,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "PartialEmoji",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustEmoji()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypePartialEmoji,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "CategoryChannel",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustChannel()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeCategoryChannel,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "StoreChannel",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustChannel()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeStoreChannel,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "GuildChannel",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustChannel()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeGuildChannel,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "String",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustString()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeString,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Bool",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustBool()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeBool,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Int",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustInt()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeInt,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Float",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustFloat()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeFloat,
			},
		},
	})

	argumentTestGroup.MustAddCommand(&sandwich.Commandable{
		Name: "Fill",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			value := ctx.MustGetArgument("test").MustString()
			println(value)

			return nil
		},
		ArgumentParameters: []sandwich.ArgumentParameter{
			{
				Name:         "test",
				ArgumentType: sandwich.ArgumentTypeFill,
			},
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
