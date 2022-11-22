package internal

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

var (
	IDRegex             = regexp.MustCompile("([0-9]{15,20})$")
	GenericMentionRegex = regexp.MustCompile("<(?:@(?:!|&)?|#)([0-9]{15,20})>$")
	UserMentionRegex    = regexp.MustCompile("<@!?([0-9]{15,20})>$")
	ChannelMentionRegex = regexp.MustCompile("<#([0-9]{15,20})>")
	RoleMentionRegex    = regexp.MustCompile("<@&([0-9]{15,20})>$")
	EmojiRegex          = regexp.MustCompile("<a?:[a-zA-Z0-9_]{1,32}:([0-9]{15,20})>$")
	PartialEmojiRegex   = regexp.MustCompile("<(a?):([a-zA-Z0-9_]{1,32}):([0-9]{15,20})>$")
)

var (
	BoolTrue = map[string]bool{
		"yes":    true,
		"y":      true,
		"true":   true,
		"t":      true,
		"1":      true,
		"enable": true,
		"on":     true,
	}

	BoolFalse = map[string]bool{
		"no":      true,
		"n":       true,
		"false":   true,
		"f":       true,
		"0":       true,
		"disable": true,
		"off":     true,
	}
)

type ArgumentParameter struct {
	Required     bool
	ArgumentType ArgumentType
	Name         string
	Description  string
}

type Argument struct {
	ArgumentType ArgumentType
	value        interface{}
}

type ArgumentConverterType func(ctx *CommandContext, argument string) (out interface{}, err error)

type Converters struct {
	convertersMu sync.RWMutex
	Converters   map[ArgumentType]*Converter
}

type Converter struct {
	converterType ArgumentConverterType
	data          interface{}
}

// RegisterConverter adds a new converter. If there is already a
// converter registered with its name, it will be overrifden.
func (co *Converters) RegisterConverter(converterName ArgumentType, converter ArgumentConverterType, defaultValue interface{}) {
	co.convertersMu.Lock()
	defer co.convertersMu.Unlock()

	co.Converters[converterName] = &Converter{
		converterType: converter,
		data:          defaultValue,
	}
}

// GetConvert returns a converter based on the converterType passed.
func (co *Converters) GetConverter(converterType ArgumentType) *Converter {
	co.convertersMu.RLock()
	defer co.convertersMu.RUnlock()

	return co.Converters[converterType]
}

// HandleArgumentTypeSnowflake handles converting from a string
// argument into a Snowflake type. Use .Snowflake() within a command
// to get the proper type.
func HandleArgumentTypeSnowflake(ctx *CommandContext, argument string) (out interface{}, err error) {
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

// HandleArgumentTypeMember handles converting from a string
// argument into a Member type. Use .Member() within a command
// to get the proper type.
func HandleArgumentTypeMember(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		matches := UserMentionRegex.FindStringSubmatch(argument)
		if len(matches) > 1 {
			match = matches[1]
		}
	}

	var result *discord.GuildMember

	if match == "" {
		if ctx.GuildID != nil {
			members, err := ctx.EventContext.Sandwich.GRPCInterface.FetchMembersByName(ctx.EventContext, *ctx.GuildID, argument)
			if err != nil {
				return nil, errors.Errorf("Failed to fetch member: %v", err)
			}

			if len(members) > 0 {
				result = members[0]
			}
		}
	} else {
		userID, _ := strconv.ParseInt(match, 10, 64)

		result = NewGuildMember(ctx.EventContext, ctx.GuildID, discord.Snowflake(userID))

		result, err = FetchGuildMember(ctx.EventContext, result)
		if err != nil && !errors.Is(err, ErrMemberNotFound) {
			return nil, err
		}
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
		matches := UserMentionRegex.FindStringSubmatch(argument)
		if len(matches) > 0 {
			match = matches[1]
		}
	}

	var result *discord.User

	if match != "" {
		userID, _ := strconv.ParseInt(match, 10, 64)

		for _, user := range ctx.Mentions {
			if user.ID == discord.Snowflake(userID) {
				result = user

				break
			}
		}

		if result == nil {
			result = NewUser(ctx.EventContext, discord.Snowflake(userID))

			result, err = FetchUser(ctx.EventContext, result, false)
			if err != nil && !errors.Is(err, ErrUserNotFound) {
				return nil, err
			}
		}

		return result, nil
	}

	if result == nil {
		arg := argument

		if arg[0] == '@' {
			arg = arg[1:]
		}

		if len(arg) > 5 && arg[len(arg)-5] == '#' {
			users, err := ctx.EventContext.Sandwich.GRPCInterface.FetchUserByName(ctx.EventContext, arg, false)
			if err != nil {
				return nil, errors.Errorf("Failed to fetch user: %v", err)
			}

			if len(users) > 0 {
				result = users[0]
			}
		}
	}

	if result == nil {
		return nil, ErrUserNotFound
	}

	return result, nil
}

// HandleArgumentTypeTextChannel handles converting from a string
// argument into a TextChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeTextChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeGuild handles converting  from a string
// argument into a Guild type. Use .Guild() within a command
// to get the proper type.
func HandleArgumentTypeGuild(ctx *CommandContext, argument string) (out interface{}, err error) {
	match := IDRegex.FindString(argument)

	var result *discord.Guild

	if match == "" {
		guilds, err := ctx.EventContext.Sandwich.GRPCInterface.FetchGuildsByName(ctx.EventContext, argument)
		if err != nil {
			return nil, errors.Errorf("Failed to fetch guild: %v", err)
		}

		if len(guilds) > 0 {
			result = guilds[0]
		}
	} else {
		guildID, _ := strconv.ParseInt(match, 10, 64)

		result = NewGuild(ctx.EventContext, discord.Snowflake(guildID))

		result, err = FetchGuild(ctx.EventContext, result)
		if err != nil && !errors.Is(err, ErrGuildNotFound) {
			return nil, err
		}
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
		matches := RoleMentionRegex.FindStringSubmatch(argument)
		if len(matches) > 0 {
			match = matches[0]
		}
	}

	var result *discord.Role

	if match == "" {
		if ctx.GuildID != nil {
			roles, err := ctx.EventContext.Sandwich.GRPCInterface.FetchRolesByName(ctx.EventContext, *ctx.GuildID, argument)
			if err != nil {
				return nil, errors.Errorf("Failed to fetch role: %v", err)
			}

			if len(roles) > 0 {
				result = roles[0]
			}
		}
	} else {
		roleID, _ := strconv.ParseInt(match, 10, 64)

		result = NewRole(ctx.EventContext, ctx.GuildID, discord.Snowflake(roleID))

		result, err = FetchRole(ctx.EventContext, result)
		if err != nil && !errors.Is(err, ErrRoleNotFound) {
			return nil, err
		}
	}

	if result == nil {
		return nil, ErrRoleNotFound
	}

	return result, nil
}

// HandleArgumentTypeColour handles converting from a string
// argument into a Colour type. Use .Colour() within a command
// to get the proper type.
func HandleArgumentTypeColour(ctx *CommandContext, argument string) (out interface{}, err error) {
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

// HandleArgumentTypeVoiceChannel handles converting from a string
// argument into a VoiceChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeVoiceChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord.ChannelTypeGuildVoice)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeStageChannel handles converting from a string
// argument into a StageChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeStageChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord.ChannelTypeGuildStageVoice)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeEmoji handles converting from a string
// argument into a Emoji type. Use .Emoji() within a command
// to get the proper type.
func HandleArgumentTypeEmoji(ctx *CommandContext, argument string) (out interface{}, err error) {
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
					return nil, errors.Errorf("Failed to fetch emoji: %v", err)
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
	if err != nil && !errors.Is(err, ErrEmojiNotFound) && !errors.Is(err, ErrFetchMissingGuild) {
		return nil, err
	}

	return result, nil
}

// HandleArgumentTypePartialEmoji handles converting from a string
// argument into a Emoji type. Use .Emoji() within a command
// to get the proper type.
func HandleArgumentTypePartialEmoji(ctx *CommandContext, argument string) (out interface{}, err error) {
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

// HandleArgumentTypeCategoryChannel handles converting from a string
// argument into a CategoryChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeCategoryChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord.ChannelTypeGuildCategory)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeStoreChannel handles converting from a string
// argument into a StoreChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeStoreChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord.ChannelTypeGuildStore)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeThread handles converting from a string
// argument into a Thread type. Use .Thread() within a command
// to get the proper type.
func HandleArgumentTypeThread(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord.ChannelTypeGuildNewsThread,
		discord.ChannelTypeGuildPrivateThread, discord.ChannelTypeGuildPublicThread)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeGuildChannel handles converting from a string
// argument into a GuildChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeGuildChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

// HandleArgumentTypeString handles converting from a string
// argument into a String type. Use .String() within a command
// to get the proper type.
func HandleArgumentTypeString(ctx *CommandContext, argument string) (out interface{}, err error) {
	return argument, nil
}

// HandleArgumentTypeBool handles converting from a string
// argument into a Bool type. Use .Bool() within a command
// to get the proper type.
func HandleArgumentTypeBool(ctx *CommandContext, argument string) (out interface{}, err error) {
	argument = strings.ToLower(argument)

	if _, ok := BoolTrue[argument]; ok {
		return true, nil
	}

	if _, ok := BoolFalse[argument]; ok {
		return false, nil
	}

	return argument, ErrBadBoolArgument
}

// HandleArgumentTypeInt handles converting from a string
// argument into a Int type. Use .Int64() within a command
// to get the proper type.
func HandleArgumentTypeInt(ctx *CommandContext, argument string) (out interface{}, err error) {
	result, err := strconv.ParseInt(argument, 10, 64)
	if err != nil {
		return nil, ErrBadIntArgument
	}

	return result, nil
}

// HandleArgumentTypeFloat handles converting from a string
// argument into a Float type. Use .Float64() within a command
// to get the proper type.
func HandleArgumentTypeFloat(ctx *CommandContext, argument string) (out interface{}, err error) {
	result, err := strconv.ParseFloat(argument, 64)
	if err != nil {
		return nil, ErrBadFloatArgument
	}

	return result, nil
}

// Argument fetchers

// Snowflake returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Snowflake() (value *discord.Snowflake, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeSnowflake) {
		value, _ = a.value.(*discord.Snowflake)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustSnowflake will attempt to do Snowflake() and will panic if not possible.
func (a *Argument) MustSnowflake() (value *discord.Snowflake) {
	value, err := a.Snowflake()
	if err != nil {
		panic(fmt.Sprintf(`argument: Snowflake(): %v`, err.Error()))
	}

	return value
}

// Member returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Member() (value *discord.GuildMember, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeMember) {
		value, _ = a.value.(*discord.GuildMember)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustMember will attempt to do Member() and will panic if not possible.
func (a *Argument) MustMember() (value *discord.GuildMember) {
	value, err := a.Member()
	if err != nil {
		panic(fmt.Sprintf(`argument: Member(): %v`, err.Error()))
	}

	return value
}

// User returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) User() (value *discord.User, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeUser) {
		value, _ = a.value.(*discord.User)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustUser will attempt to do User() and will panic if not possible.
func (a *Argument) MustUser() (value *discord.User) {
	value, err := a.User()
	if err != nil {
		panic(fmt.Sprintf(`argument: User(): %v`, err.Error()))
	}

	return value
}

// Channel returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Channel() (value *discord.Channel, err error) {
	if argumentTypeIs(a.ArgumentType,
		ArgumentTypeTextChannel, ArgumentTypeVoiceChannel, ArgumentTypeStageChannel,
		ArgumentTypeCategoryChannel, ArgumentTypeStoreChannel, ArgumentTypeGuildChannel) {
		value, _ = a.value.(*discord.Channel)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustTextChannel will attempt to do Channel() and will panic if not possible.
func (a *Argument) MustChannel() (value *discord.Channel) {
	value, err := a.Channel()
	if err != nil {
		panic(fmt.Sprintf(`argument: Channel(): %v`, err.Error()))
	}

	return value
}

// Guild returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Guild() (value *discord.Guild, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeGuild) {
		value, _ = a.value.(*discord.Guild)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustGuild will attempt to do Guild() and will panic if not possible.
func (a *Argument) MustGuild() (value *discord.Guild) {
	value, err := a.Guild()
	if err != nil {
		panic(fmt.Sprintf(`argument: Guild(): %v`, err.Error()))
	}

	return value
}

// Role returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Role() (value *discord.Role, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeRole) {
		value, _ = a.value.(*discord.Role)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustRole will attempt to do Role() and will panic if not possible.
func (a *Argument) MustRole() (value *discord.Role) {
	value, err := a.Role()
	if err != nil {
		panic(fmt.Sprintf(`argument: Role(): %v`, err.Error()))
	}

	return value
}

// Colour returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Colour() (value *color.RGBA, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeColour) {
		value, _ = a.value.(*color.RGBA)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustColour will attempt to do Colour() and will panic if not possible.
func (a *Argument) MustColour() (value *color.RGBA) {
	value, err := a.Colour()
	if err != nil {
		panic(fmt.Sprintf(`argument: Colour(): %v`, err.Error()))
	}

	return value
}

// Emoji returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Emoji() (value *discord.Emoji, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeEmoji, ArgumentTypePartialEmoji) {
		value, _ = a.value.(*discord.Emoji)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustEmoji will attempt to do Emoji() and will panic if not possible.
func (a *Argument) MustEmoji() (value *discord.Emoji) {
	value, err := a.Emoji()
	if err != nil {
		panic(fmt.Sprintf(`argument: Emoji(): %v`, err.Error()))
	}

	return value
}

// String returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) String() (value string, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeString, ArgumentTypeFill) {
		value, _ = a.value.(string)

		return
	}

	return "", ErrInvalidArgumentType
}

// MustString will attempt to do String() and will panic if not possible.
func (a *Argument) MustString() (value string) {
	value, err := a.String()
	if err != nil {
		panic(fmt.Sprintf(`argument: String(): %v`, err.Error()))
	}

	return value
}

// Bool returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Bool() (value bool, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeBool) {
		value, _ = a.value.(bool)

		return
	}

	return false, ErrInvalidArgumentType
}

// MustBool will attempt to do Bool() and will panic if not possible.
func (a *Argument) MustBool() (value bool) {
	value, err := a.Bool()
	if err != nil {
		panic(fmt.Sprintf(`argument: Bool(): %v`, err.Error()))
	}

	return value
}

// Int returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Int() (value int64, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeInt) {
		value, _ = a.value.(int64)

		return
	}

	return 0, ErrInvalidArgumentType
}

// MustInt will attempt to do Int() and will panic if not possible.
func (a *Argument) MustInt() (value int64) {
	value, err := a.Int()
	if err != nil {
		panic(fmt.Sprintf(`argument: Int(): %v`, err.Error()))
	}

	return value
}

// Float returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Float() (value float64, err error) {
	v, ok := a.value.(float64)
	if !ok {
		return v, ErrInvalidArgumentType
	}

	return v, nil
}

// MustFloat will attempt to do Float() and will panic if not possible.
func (a *Argument) MustFloat() (value float64) {
	value, err := a.Float()
	if err != nil {
		panic(fmt.Sprintf(`argument: Float(): %v`, err.Error()))
	}

	return value
}

func NewDefaultConverters() (co *Converters) {
	co = &Converters{
		convertersMu: sync.RWMutex{},
		Converters:   make(map[ArgumentType]*Converter),
	}

	co.RegisterConverter(ArgumentTypeSnowflake, HandleArgumentTypeSnowflake, nil)
	co.RegisterConverter(ArgumentTypeMember, HandleArgumentTypeMember, nil)
	co.RegisterConverter(ArgumentTypeUser, HandleArgumentTypeUser, nil)
	co.RegisterConverter(ArgumentTypeTextChannel, HandleArgumentTypeTextChannel, nil)
	co.RegisterConverter(ArgumentTypeGuild, HandleArgumentTypeGuild, nil)
	co.RegisterConverter(ArgumentTypeRole, HandleArgumentTypeRole, nil)
	co.RegisterConverter(ArgumentTypeColour, HandleArgumentTypeColour, nil)
	co.RegisterConverter(ArgumentTypeVoiceChannel, HandleArgumentTypeVoiceChannel, nil)
	co.RegisterConverter(ArgumentTypeStageChannel, HandleArgumentTypeStageChannel, nil)
	co.RegisterConverter(ArgumentTypeEmoji, HandleArgumentTypeEmoji, nil)
	co.RegisterConverter(ArgumentTypePartialEmoji, HandleArgumentTypePartialEmoji, nil)
	co.RegisterConverter(ArgumentTypeCategoryChannel, HandleArgumentTypeCategoryChannel, nil)
	co.RegisterConverter(ArgumentTypeStoreChannel, HandleArgumentTypeStoreChannel, nil)
	co.RegisterConverter(ArgumentTypeThread, HandleArgumentTypeThread, nil)
	co.RegisterConverter(ArgumentTypeGuildChannel, HandleArgumentTypeGuildChannel, nil)
	co.RegisterConverter(ArgumentTypeString, HandleArgumentTypeString, "")
	co.RegisterConverter(ArgumentTypeBool, HandleArgumentTypeBool, false)
	co.RegisterConverter(ArgumentTypeInt, HandleArgumentTypeInt, int64(0))
	co.RegisterConverter(ArgumentTypeFloat, HandleArgumentTypeFloat, float64(0))
	co.RegisterConverter(ArgumentTypeFill, HandleArgumentTypeString, "")

	return co
}
