package main

import (
	"bytes"
	"context"
	"fmt"
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

			_, _ = ctx.Reply(ctx.EventContext.Session,
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

	bot.Commands.MustAddCommand(&sandwich.Commandable{
		Name: "registerCommands",
		Handler: func(ctx *sandwich.CommandContext) (err error) {
			applicationcommands := bot.InteractionCommands.MapApplicationCommands()

			_, err = discord.BulkOverwriteGloblApplicationCommands(ctx.EventContext.Session, ctx.EventContext.Identifier.ID, applicationcommands)

			if err != nil {
				println(err.Error())
			}

			return nil
		},
	})

	bot.InteractionCommands.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "pog",
		Type: sandwich.InteractionCommandableTypeCommand,

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: "<:rock:732274836038221855>📣 pog",
					},
				},
			}, nil
		},
	})

	argumentTestGroup, _ := bot.InteractionCommands.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "argument_test",
		Type: sandwich.InteractionCommandableTypeSubcommandGroup,
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "snowflake",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "snowflake",
				Description:  "Snowflake",
				ArgumentType: sandwich.ArgumentTypeSnowflake,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Snowflake := ctx.MustGetArgument("snowflake").MustSnowflake()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Snowflake.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "member",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "member",
				Description:  "member",
				ArgumentType: sandwich.ArgumentTypeMember,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Member := ctx.MustGetArgument("member").MustMember()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Member.User.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "user",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "user",
				Description:  "user",
				ArgumentType: sandwich.ArgumentTypeUser,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			User := ctx.MustGetArgument("user").MustUser()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", User.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "text_channel",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "text_channel",
				Description:  "text_channel",
				ArgumentType: sandwich.ArgumentTypeTextChannel,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			TextChannel := ctx.MustGetArgument("text_channel").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", TextChannel.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "guild",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "guild",
				Description:  "guild",
				ArgumentType: sandwich.ArgumentTypeGuild,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Guild := ctx.MustGetArgument("guild").MustGuild()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Guild.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "role",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "role",
				Description:  "role",
				ArgumentType: sandwich.ArgumentTypeRole,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Role := ctx.MustGetArgument("role").MustRole()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Role.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "colour",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "colour",
				Description:  "colour",
				ArgumentType: sandwich.ArgumentTypeColour,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Colour := ctx.MustGetArgument("colour").MustColour()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %v", Colour),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "voice_channel",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "voice_channel",
				Description:  "voice_channel",
				ArgumentType: sandwich.ArgumentTypeVoiceChannel,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			VoiceChannel := ctx.MustGetArgument("voice_channel").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", VoiceChannel.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "stage_channel",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "stage_channel",
				Description:  "stage_channel",
				ArgumentType: sandwich.ArgumentTypeStageChannel,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			StageChannel := ctx.MustGetArgument("stage_channel").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", StageChannel.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "emoji",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "emoji",
				Description:  "emoji",
				ArgumentType: sandwich.ArgumentTypeEmoji,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Emoji := ctx.MustGetArgument("emoji").MustEmoji()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Emoji.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "partial_emoji",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "partial_emoji",
				Description:  "partial_emoji",
				ArgumentType: sandwich.ArgumentTypePartialEmoji,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			PartialEmoji := ctx.MustGetArgument("partial_emoji").MustEmoji()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", PartialEmoji.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "category_channel",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "category_channel",
				Description:  "category_channel",
				ArgumentType: sandwich.ArgumentTypeCategoryChannel,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			CategoryChannel := ctx.MustGetArgument("category_channel").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", CategoryChannel.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "store_channel",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "store_channel",
				Description:  "store_channel",
				ArgumentType: sandwich.ArgumentTypeStoreChannel,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			StoreChannel := ctx.MustGetArgument("store_channel").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", StoreChannel.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "thread",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "thread",
				Description:  "thread",
				ArgumentType: sandwich.ArgumentTypeThread,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Thread := ctx.MustGetArgument("thread").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Thread.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "guild_channel",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "guild_channel",
				Description:  "guild_channel",
				ArgumentType: sandwich.ArgumentTypeGuildChannel,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			GuildChannel := ctx.MustGetArgument("guild_channel").MustChannel()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", GuildChannel.ID.String()),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "string",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "string",
				Description:  "string",
				ArgumentType: sandwich.ArgumentTypeString,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			String := ctx.MustGetArgument("string").MustString()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", String),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "bool",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "bool",
				Description:  "bool",
				ArgumentType: sandwich.ArgumentTypeBool,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Bool := ctx.MustGetArgument("bool").MustBool()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %t", Bool),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "int",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "int",
				Description:  "int",
				ArgumentType: sandwich.ArgumentTypeInt,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Int := ctx.MustGetArgument("int").MustInt()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %d", Int),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "float",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "float",
				Description:  "float",
				ArgumentType: sandwich.ArgumentTypeFloat,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Float := ctx.MustGetArgument("float").MustFloat()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %f", Float),
					},
				},
			}, nil
		},
	})

	argumentTestGroup.AddInteractionCommand(&sandwich.InteractionCommandable{
		Name: "fill",
		Type: sandwich.InteractionCommandableTypeSubcommand,

		ArgumentParameter: []sandwich.ArgumentParameter{
			{
				Name:         "fill",
				Description:  "fill",
				ArgumentType: sandwich.ArgumentTypeFill,
			},
		},

		Handler: func(ctx *sandwich.InteractionContext) (resp *sandwich.InteractionResponse, err error) {
			Fill := ctx.MustGetArgument("fill").MustString()

			return &sandwich.InteractionResponse{
				Type: discord.InteractionCallbackTypeChannelMessageSource,
				Data: discord.InteractionCallbackData{
					WebhookMessageParams: discord.WebhookMessageParams{
						Content: fmt.Sprintf("Value is: %s", Fill),
					},
				},
			}, nil
		},
	})

	bot.RegisterOnInteractionCreateEvent(func(ctx *sandwich.EventContext, interaction discord.Interaction) (err error) {
		resp, err := bot.ProcessInteraction(ctx, interaction)
		if err != nil {
			println(err.Error())

			return
		}

		if resp != nil {
			err = interaction.SendResponse(ctx.Session, resp.Type, resp.Data.WebhookMessageParams, resp.Data.Choices)
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
				err = sandwichClient.DispatchSandwichPayload(ctx, p)
				if err != nil {
					println("Fail to dispatch", err.Error())
				}
			} else {
				println(err.Error(), string(m))
			}
		}
	}
}
