package internal

import (
	"image/color"
	"strconv"
	"sync"

	discord "github.com/WelcomerTeam/Discord/discord"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
)

type InteractionArgumentConverterType func(ctx *InteractionContext, argument *discord.InteractionDataOption) (out interface{}, err error)

type InteractionConverters struct {
	convertersMu sync.RWMutex
	Converters   map[ArgumentType]*InteractionConverter
}

type InteractionConverter struct {
	converterType InteractionArgumentConverterType
	data          interface{}
}

// RegisterConverter adds a new converter. If there is already a
// converter registered with its name, it will be overrifden.
func (co *InteractionConverters) RegisterConverter(converterName ArgumentType, converter InteractionArgumentConverterType, defaultValue interface{}) {
	co.convertersMu.Lock()
	defer co.convertersMu.Unlock()

	co.Converters[converterName] = &InteractionConverter{
		converterType: converter,
		data:          defaultValue,
	}
}

// GetConverter returns a converter based on the converterType passed.
func (co *InteractionConverters) GetConverter(converterType ArgumentType) *InteractionConverter {
	co.convertersMu.RLock()
	defer co.convertersMu.RUnlock()

	return co.Converters[converterType]
}

// HandleInteractionArgumentTypeSnowflake handles converting from a string
// argument into a Snowflake type. Use .Snowflake() within a command
// to get the proper type.
func HandleInteractionArgumentTypeSnowflake(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	match := IDRegex.FindString(argument)
	if match == "" {
		matches := GenericMentionRegex.FindStringSubmatch(argument)
		if len(matches) > 1 {
			match = matches[1]
		}
	}

	var result *discord.Snowflake

	if match == "" {
		return nil, ErrSnowflakeNotFound
	}

	snowflakeID, _ := strconv.ParseInt(match, 10, 64)
	result = (*discord.Snowflake)(&snowflakeID)

	return result, nil
}

// HandleInteractionArgumentTypeMember handles converting from a string
// argument into a Member type. Use .Member() within a command
// to get the proper type.
func HandleInteractionArgumentTypeMember(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	snowflakeID, _ := strconv.ParseInt(argument, 10, 64)

	result := ctx.Data.Resolved.Members[discord.Snowflake(snowflakeID)]

	if result == nil {
		return nil, ErrMemberNotFound
	}

	return result, nil
}

// HandleInteractionArgumentTypeUser handles converting from a string
// argument into a User type. Use .User() within a command
// to get the proper type.
func HandleInteractionArgumentTypeUser(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	snowflakeID, _ := strconv.ParseInt(argument, 10, 64)

	result := ctx.Data.Resolved.Users[discord.Snowflake(snowflakeID)]

	if result == nil {
		return nil, ErrUserNotFound
	}

	return result, nil
}

// HandleInteractionArgumentTypeGuildChannel handles converting from a string
// argument into a TextChannel type. Use .Channel() within a command
// to get the proper type.
func HandleInteractionArgumentTypeGuildChannel(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	snowflakeID, _ := strconv.ParseInt(argument, 10, 64)

	result := ctx.Data.Resolved.Channels[discord.Snowflake(snowflakeID)]

	if result == nil {
		return nil, ErrChannelNotFound
	}

	return result, nil
}

// HandleInteractionArgumentTypeGuild handles converting from a string
// argument into a Guild type. Use .Guild() within a command
// to get the proper type.
func HandleInteractionArgumentTypeGuild(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	match := IDRegex.FindString(argument)

	var result *discord.Guild

	if match == "" {
		guilds, err := ctx.EventContext.Sandwich.GRPCInterface.FetchGuildsByName(ctx.EventContext, argument)
		if err != nil {
			return nil, xerrors.Errorf("Failed to fetch guild: %v", err)
		}

		if len(guilds) > 0 {
			result = guilds[0]
		}
	} else {
		guildID, _ := strconv.ParseInt(match, 10, 64)

		result = NewGuild(ctx.EventContext, discord.Snowflake(guildID))

		result, err = FetchGuild(ctx.EventContext, result)
		if err != nil && !xerrors.Is(err, ErrGuildNotFound) {
			return nil, err
		}
	}

	if result == nil {
		return nil, ErrGuildNotFound
	}

	return result, nil
}

// HandleInteractionArgumentTypeRole handles converting from a string
// argument into a Role type. Use .Role() within a command
// to get the proper type.
func HandleInteractionArgumentTypeRole(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	snowflakeID, _ := strconv.ParseInt(argument, 10, 64)

	result := ctx.Data.Resolved.Roles[discord.Snowflake(snowflakeID)]

	if result == nil {
		return nil, ErrRoleNotFound
	}

	return result, nil
}

// HandleInteractionArgumentTypeColour handles converting from a string
// argument into a Colour type. Use .Colour() within a command
// to get the proper type.
func HandleInteractionArgumentTypeColour(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	var result *color.RGBA

	switch {
	case argument[0] == '#':
		hexNum, err := parseHexNumber(argument[1:])
		if err != nil {
			return nil, err
		}

		result = intToColour(hexNum)
	case argument[0:2] == "0x":
		hexNum, err := parseHexNumber(argument[2:])
		if err != nil {
			return nil, err
		}

		result = intToColour(hexNum)
	default:
		hexNum, err := parseHexNumber(argument)
		if err == nil {
			result = intToColour(hexNum)
		}
	}

	if result == nil {
		return nil, ErrBadColourArgument
	}

	return result, nil
}

// HandleInteractionArgumentTypeEmoji handles converting from a string
// argument into a Emoji type. Use .Emoji() within a command
// to get the proper type.
func HandleInteractionArgumentTypeEmoji(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	var result *discord.Emoji

	match := IDRegex.FindString(argument)
	if match == "" {
		matches := PartialEmojiRegex.FindStringSubmatch(argument)

		if len(matches) >= 4 {
			animated := matches[1] != ""
			id, _ := strconv.ParseInt(matches[3], 10, 64)

			result = &discord.Emoji{
				Animated: animated,
				Name:     matches[2],
				ID:       discord.Snowflake(id),
			}
		}
	}

	if result == nil {
		if match == "" {
			if ctx.GuildID != nil {
				emojis, err := ctx.EventContext.Sandwich.GRPCInterface.FetchEmojisByName(ctx.EventContext, *ctx.GuildID, argument)
				if err != nil {
					return nil, xerrors.Errorf("Failed to fetch emoji: %v", err)
				}

				if len(emojis) > 0 {
					result = emojis[0]
				}
			}
		} else {
			emojiID, _ := strconv.ParseInt(match, 10, 64)

			result = NewEmoji(ctx.EventContext, ctx.GuildID, discord.Snowflake(emojiID))
		}
	}

	result, err = FetchEmoji(ctx.EventContext, result)
	if err != nil && !xerrors.Is(err, ErrEmojiNotFound) && !xerrors.Is(err, ErrFetchMissingGuild) {
		return nil, err
	}

	return result, nil
}

// HandleInteractionArgumentTypePartialEmoji handles converting from a string
// argument into a Emoji type. Use .Emoji() within a command
// to get the proper type.
func HandleInteractionArgumentTypePartialEmoji(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	matches := PartialEmojiRegex.FindStringSubmatch(argument)

	var result *discord.Emoji

	if len(matches) >= 4 {
		animated := matches[1] != ""
		id, _ := strconv.ParseInt(matches[3], 10, 64)

		result = &discord.Emoji{
			Animated: animated,
			Name:     matches[2],
			ID:       discord.Snowflake(id),
		}
	}

	return result, nil
}

// HandleInteractionArgumentTypeString handles converting from a string
// argument into a String type. Use .String() within a command
// to get the proper type.
func HandleInteractionArgumentTypeString(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	return argument, nil
}

// HandleInteractionArgumentTypeBool handles converting from a string
// argument into a Bool type. Use .Bool() within a command
// to get the proper type.
func HandleInteractionArgumentTypeBool(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument bool

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	return argument, nil
}

// HandleInteractionArgumentTypeInt handles converting from a string
// argument into a Int type. Use .Int64() within a command
// to get the proper type.
func HandleInteractionArgumentTypeInt(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument int64

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	return argument, nil
}

// HandleInteractionArgumentTypeFloat handles converting from a string
// argument into a Float type. Use .Float64() within a command
// to get the proper type.
func HandleInteractionArgumentTypeFloat(ctx *InteractionContext, option *discord.InteractionDataOption) (out interface{}, err error) {
	var argument string

	err = jsoniter.Unmarshal(option.Value, &argument)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal option value: %v", err)
	}

	result, err := strconv.ParseFloat(argument, 64)
	if err != nil {
		return nil, ErrBadFloatArgument
	}

	return result, nil
}

func NewInteractionConverters() (co *InteractionConverters) {
	co = &InteractionConverters{
		convertersMu: sync.RWMutex{},
		Converters:   make(map[ArgumentType]*InteractionConverter),
	}

	co.RegisterConverter(ArgumentTypeSnowflake, HandleInteractionArgumentTypeSnowflake, nil)
	co.RegisterConverter(ArgumentTypeMember, HandleInteractionArgumentTypeMember, nil)
	co.RegisterConverter(ArgumentTypeUser, HandleInteractionArgumentTypeUser, nil)
	co.RegisterConverter(ArgumentTypeTextChannel, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeGuild, HandleInteractionArgumentTypeGuild, nil)
	co.RegisterConverter(ArgumentTypeRole, HandleInteractionArgumentTypeRole, nil)
	co.RegisterConverter(ArgumentTypeColour, HandleInteractionArgumentTypeColour, nil)
	co.RegisterConverter(ArgumentTypeVoiceChannel, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeStageChannel, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeEmoji, HandleInteractionArgumentTypeEmoji, nil)
	co.RegisterConverter(ArgumentTypePartialEmoji, HandleInteractionArgumentTypePartialEmoji, nil)
	co.RegisterConverter(ArgumentTypeCategoryChannel, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeStoreChannel, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeThread, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeGuildChannel, HandleInteractionArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeString, HandleInteractionArgumentTypeString, "")
	co.RegisterConverter(ArgumentTypeBool, HandleInteractionArgumentTypeBool, false)
	co.RegisterConverter(ArgumentTypeInt, HandleInteractionArgumentTypeInt, int64(0))
	co.RegisterConverter(ArgumentTypeFloat, HandleInteractionArgumentTypeFloat, float64(0))
	co.RegisterConverter(ArgumentTypeFill, HandleInteractionArgumentTypeString, "")

	return co
}
