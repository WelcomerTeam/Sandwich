package internal

import (
	"image/color"
	"regexp"
	"strconv"
	"sync"

	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
)

var (
	IDRegex             = regexp.MustCompile("([-0]{15,20}$")
	GenericMentionRegex = regexp.MustCompile("<(?:@(?:!|&)?|#)([0-9]{15,20})>$")
	UserMentionRegex    = regexp.MustCompile("<@!?([0-9]{15,20})>$")
	ChannelMentionRegex = regexp.MustCompile("<#([0-9]{15,20})>")
	RoleMentionRegex    = regexp.MustCompile("<@&([0-9]{15,20})>$")
	EmojiRegex          = regexp.MustCompile("<a?:[a-zA-Z0-9_]{1,32}:([0-9]{15,20})>$")
)

type ArgumentParameter struct {
	Required     bool
	ArgumentType ArgumentType
	Name         string
}

type Argument struct {
	ArgumentType ArgumentType
	value        interface{}
}

type ArgumentConverterType func(ctx *CommandContext, argument string) (out interface{}, err error)

type Converters struct {
	convertersMu sync.RWMutex
	Converters   map[ArgumentType]ArgumentConverterType
}

// RegisterConverter adds a new converter. If there is already a
// converter registered with its name, it will be overrifden.
func (co *Converters) RegisterConverter(converterName ArgumentType, converter ArgumentConverterType) {
	co.convertersMu.Lock()
	defer co.convertersMu.Unlock()

	co.Converters[converterName] = converter
}

// HandleArgumentTypeSnowflake handles converting from a string
// argument into a Snowflake type. Use .Snowflake() within a command
// to get the proper type.
func HandleArgumentTypeSnowflake(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		match = GenericMentionRegex.FindString(argument)
	}

	var result discord.Snowflake

	if match == "" {
		return nil, ErrSnowflakeNotFound
	}

	snowflakeID, _ := strconv.ParseInt(match, 10, 64)
	result = discord.Snowflake(snowflakeID)

	return result, nil
}

// HandleArgumentTypeMember handles converting from a string
// argument into a Member type. Use .Member() within a command
// to get the proper type.
func HandleArgumentTypeMember(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		match = UserMentionRegex.FindString(argument)
	}

	var result *GuildMember

	if match == "" {
		if ctx.Guild != nil {
			result, err = ctx.Guild.GetMemberNamed(ctx.EventContext, argument)
			if err != nil {
				return nil, err
			}
		}
	} else {
		userID, _ := strconv.ParseInt(match, 10, 64)

		result, err = ctx.Guild.GetMemberById(ctx.EventContext, discord.Snowflake(userID))
		if err != nil {
			return nil, err
		}
	}

	if result == nil {
		// TODO: Do proper member querying, not relying on just builtin cache.
	}

	if result == nil {
		return nil, ErrMemberNotFound
	}

	return result, nil
}

// HandleArgumentTypeUser handles converting from a string
// argument into a User type. Use .User() within a command
// to get the proper type.
func HandleArgumentTypeUser(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		match = UserMentionRegex.FindString(argument)
	}

	var result *User

	if match != "" {
		userID, _ := strconv.ParseInt(match, 10, 64)

		// TODO: Query user via GRPC

		for _, user := range ctx.Mentions {
			if user.ID == discord.Snowflake(userID) {
				result = (*User)(user)

				break
			}
		}

		if result == nil {
			// TODO: Fetch user from discord.

			return nil, ErrUserNotFound
		}

		return result, nil
	}

	arg := argument
	if arg[0] == '@' {
		arg = arg[1:]
	}

	if len(arg) > 5 && arg[len(arg)-5] == '#' {
		// username, _, discriminator := rpartition(arg, "#")

		// TODO: Query user via GRPC
	}

	if result == nil {
		return nil, ErrUserNotFound
	}

	return result, nil
}

// // HandleArgumentTypeMessage handles converting from a string
// // argument into a Message type. Use .Message() within a command
// // to get the proper type.
// func HandleArgumentTypeMessage(ctx *CommandContext, argument string) (out interface{}, err error) {
// 	// TODO: Implement ArgumentTypeMessage converter
// 	return result, nil
// }

// // HandleArgumentTypePartialMessage handles converting from a string
// // argument into a PartialMessage type. Use .PartialMessage() within a command
// // to get the proper type.
// func HandleArgumentTypePartialMessage(ctx *CommandContext, argument string) (out interface{}, err error) {
// 	// TODO: Implement ArgumentTypePartialMessage converter
// 	return result, nil
// }

// HandleArgumentTypeTextChannel handles converting from a string
// argument into a TextChannel type. Use .TextChannel() within a command
// to get the proper type.
func HandleArgumentTypeTextChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		match = ChannelMentionRegex.FindString(argument)
	}

	var result *Channel

	if match == "" {
		// TODO: Fetch channel by name
	} else {
		// channelID, _ := strconv.ParseInt(match, 10, 64)

		// TODO: Fetch channel by ID
	}

	if result == nil {
		return nil, ErrChannelNotFound
	}

	return result, nil
}

// HandleArgumentTypeInvite handles converting from a string
// argument into a Invite type. Use .Invite() within a command
// to get the proper type.
func HandleArgumentTypeInvite(ctx *CommandContext, argument string) (out interface{}, err error) {
	var result *Invite

	// TODO: Fetch invite from discord

	if result == nil {
		return nil, ErrBadInviteArgument
	}

	return result, nil
}

// HandleArgumentTypeGuild handles converting from a string
// argument into a Guild type. Use .Guild() within a command
// to get the proper type.
func HandleArgumentTypeGuild(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)

	var result *Guild

	if match == "" {
		// TODO: Fetch guild by name
	} else {
		// guildID, _ := strconv.ParseInt(match, 10, 64)

		// TODO: Fetch guild by ID
	}

	if result == nil {
		return nil, ErrGuildNotFound
	}

	return result, nil
}

// HandleArgumentTypeRole handles converting from a string
// argument into a Role type. Use .Role() within a command
// to get the proper type.
func HandleArgumentTypeRole(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		match = RoleMentionRegex.FindString(argument)
	}

	var result *Role

	if match == "" {
		// TODO: Fetch role by name
	} else {
		// roleID, _ := strconv.ParseInt(match, 10, 64)

		// TODO: Fetch role by ID
	}

	if result == nil {
		return nil, ErrRoleNotFound
	}

	return result, nil
}

// HandleArgumentTypeGame handles converting from a string
// argument into a Game type. Use .Game() within a command
// to get the proper type.
func HandleArgumentTypeGame(ctx *CommandContext, argument string) (out interface{}, err error) {
	result := &Activity{
		Name: argument,
	}

	return result, nil
}

// HandleArgumentTypeColour handles converting from a string
// argument into a Colour type. Use .Colour() within a command
// to get the proper type.
func HandleArgumentTypeColour(ctx *CommandContext, argument string) (out interface{}, err error) {
	var result *color.RGBA

	if argument[0] == '#' {
		hexNum, err := parseHexNumber(argument[1:])
		if err != nil {
			return nil, err
		}

		result = intToColour(hexNum)
	} else if argument[0:2] == "0x" {
		hexNum, err := parseHexNumber(argument[2:])
		if err != nil {
			return nil, err
		}

		result = intToColour(hexNum)
	}

	if result == nil {
		return nil, ErrBadColourArgument
	}

	return result, nil
}

func parseHexNumber(arg string) (out uint64, err error) {
	return strconv.ParseUint(arg, 16, 8)
}

func intToColour(val uint64) (out *color.RGBA) {
	return &color.RGBA{
		R: uint8((val >> 24) & 0xFF),
		G: uint8((val >> 16) & 0xFF),
		B: uint8((val >> 8) & 0xFF),
		A: uint8(val & 0xFF),
	}
}

// HandleArgumentTypeVoiceChannel handles converting from a string
// argument into a VoiceChannel type. Use .VoiceChannel() within a command
// to get the proper type.
func HandleArgumentTypeVoiceChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeVoiceChannel converter
	return argument, nil
}

// HandleArgumentTypeStageChannel handles converting from a string
// argument into a StageChannel type. Use .StageChannel() within a command
// to get the proper type.
func HandleArgumentTypeStageChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeStageChannel converter
	return argument, nil
}

// HandleArgumentTypeEmoji handles converting from a string
// argument into a Emoji type. Use .Emoji() within a command
// to get the proper type.
func HandleArgumentTypeEmoji(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeEmoji converter
	return argument, nil
}

// HandleArgumentTypePartialEmoji handles converting from a string
// argument into a PartialEmoji type. Use .PartialEmoji() within a command
// to get the proper type.
func HandleArgumentTypePartialEmoji(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypePartialEmoji converter
	return argument, nil
}

// HandleArgumentTypeCategoryChannel handles converting from a string
// argument into a CategoryChannel type. Use .CategoryChannel() within a command
// to get the proper type.
func HandleArgumentTypeCategoryChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeCategoryChannel converter
	return argument, nil
}

// HandleArgumentTypeStoreChannel handles converting from a string
// argument into a StoreChannel type. Use .StoreChannel() within a command
// to get the proper type.
func HandleArgumentTypeStoreChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeStoreChannel converter
	return argument, nil
}

// HandleArgumentTypeThread handles converting from a string
// argument into a Thread type. Use .Thread() within a command
// to get the proper type.
func HandleArgumentTypeThread(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeThread converter
	return argument, nil
}

// HandleArgumentTypeGuildChannel handles converting from a string
// argument into a GuildChannel type. Use .GuildChannel() within a command
// to get the proper type.
func HandleArgumentTypeGuildChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeGuildChannel converter
	return argument, nil
}

// HandleArgumentTypeGuildSticker handles converting from a string
// argument into a GuildSticker type. Use .GuildSticker() within a command
// to get the proper type.
func HandleArgumentTypeGuildSticker(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeGuildSticker converter
	return argument, nil
}

// HandleArgumentTypeString handles converting from a string
// argument into a String type. Use .String() within a command
// to get the proper type.
func HandleArgumentTypeString(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeString converter
	return argument, nil
}

// HandleArgumentTypeBool handles converting from a string
// argument into a Bool type. Use .Bool() within a command
// to get the proper type.
func HandleArgumentTypeBool(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeBool converter
	return argument, nil
}

// HandleArgumentTypeInt handles converting from a string
// argument into a Int type. Use .Int() within a command
// to get the proper type.
func HandleArgumentTypeInt(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeInt converter
	return argument, nil
}

// HandleArgumentTypeFloat handles converting from a string
// argument into a Float type. Use .Float() within a command
// to get the proper type.
func HandleArgumentTypeFloat(ctx *CommandContext, argument string) (out interface{}, err error) {
	// TODO: Implement ArgumentTypeFloat converter
	return argument, nil
}

func NewDefaultConverters() (co *Converters) {
	co = &Converters{
		convertersMu: sync.RWMutex{},
		Converters:   make(map[ArgumentType]ArgumentConverterType),
	}

	co.RegisterConverter(ArgumentTypeSnowflake, HandleArgumentTypeSnowflake)
	co.RegisterConverter(ArgumentTypeMember, HandleArgumentTypeMember)
	co.RegisterConverter(ArgumentTypeUser, HandleArgumentTypeUser)
	// co.RegisterConverter(ArgumentTypeMessage, HandleArgumentTypeMessage)
	// co.RegisterConverter(ArgumentTypePartialMessage, HandleArgumentTypePartialMessage)
	// co.RegisterConverter(ArgumentTypeTextchannel, HandleArgumentTypeTextchannel)
	co.RegisterConverter(ArgumentTypeInvite, HandleArgumentTypeInvite)
	co.RegisterConverter(ArgumentTypeGuild, HandleArgumentTypeGuild)
	co.RegisterConverter(ArgumentTypeRole, HandleArgumentTypeRole)
	co.RegisterConverter(ArgumentTypeGame, HandleArgumentTypeGame)
	co.RegisterConverter(ArgumentTypeColour, HandleArgumentTypeColour)
	co.RegisterConverter(ArgumentTypeVoiceChannel, HandleArgumentTypeVoiceChannel)
	co.RegisterConverter(ArgumentTypeStageChannel, HandleArgumentTypeStageChannel)
	co.RegisterConverter(ArgumentTypeEmoji, HandleArgumentTypeEmoji)
	co.RegisterConverter(ArgumentTypePartialEmoji, HandleArgumentTypePartialEmoji)
	co.RegisterConverter(ArgumentTypeCategoryChannel, HandleArgumentTypeCategoryChannel)
	co.RegisterConverter(ArgumentTypeStoreChannel, HandleArgumentTypeStoreChannel)
	co.RegisterConverter(ArgumentTypeThread, HandleArgumentTypeThread)
	co.RegisterConverter(ArgumentTypeGuildChannel, HandleArgumentTypeGuildChannel)
	co.RegisterConverter(ArgumentTypeGuildSticker, HandleArgumentTypeGuildSticker)
	co.RegisterConverter(ArgumentTypeString, HandleArgumentTypeString)
	co.RegisterConverter(ArgumentTypeBool, HandleArgumentTypeBool)
	co.RegisterConverter(ArgumentTypeInt, HandleArgumentTypeInt)
	co.RegisterConverter(ArgumentTypeFloat, HandleArgumentTypeFloat)
	co.RegisterConverter(ArgumentTypeFill, HandleArgumentTypeString)

	return co
}
