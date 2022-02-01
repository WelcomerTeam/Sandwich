package internal

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"

	"golang.org/x/xerrors"
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
	f ArgumentConverterType
	d interface{}
}

// RegisterConverter adds a new converter. If there is already a
// converter registered with its name, it will be overrifden.
func (co *Converters) RegisterConverter(converterName ArgumentType, converter ArgumentConverterType, defaultValue interface{}) {
	co.convertersMu.Lock()
	defer co.convertersMu.Unlock()

	co.Converters[converterName] = &Converter{
		f: converter,
		d: defaultValue,
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

	var result *discord_structs.GuildMember

	if match == "" {
		if ctx.GuildID != nil {
			members, err := ctx.EventContext.Sandwich.grpcInterface.FetchMembersByName(ctx.EventContext, *ctx.GuildID, argument)
			if err != nil {
				return nil, xerrors.Errorf("Failed to fetch member: %v", err)
			}

			if len(members) > 0 {
				result = members[0]
			}
		}
	} else {
		userID, _ := strconv.ParseInt(match, 10, 64)

		result = NewGuildMember(ctx.EventContext, ctx.GuildID, discord.Snowflake(userID))

		result, err = FetchGuildMember(ctx.EventContext, result)
		if err != nil && !xerrors.Is(err, ErrMemberNotFound) {
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

	var result *discord_structs.User

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
			if err != nil && !xerrors.Is(err, ErrUserNotFound) {
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
			users, err := ctx.EventContext.Sandwich.grpcInterface.FetchUserByName(ctx.EventContext, arg, false)
			if err != nil {
				return nil, xerrors.Errorf("Failed to fetch user: %v", err)
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
	results, err := findChannel(ctx, argument, discord_structs.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrChannelNotFound
	}

	return results[0], nil
}

func findChannel(ctx *CommandContext, argument string, channelTypes ...discord_structs.ChannelType) (out []*discord_structs.Channel, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		matches := ChannelMentionRegex.FindStringSubmatch(argument)
		if len(matches) > 1 {
			match = matches[1]
		}
	}

	var results []*discord_structs.Channel

	if match == "" {
		if ctx.GuildID != nil {
			results, err = ctx.EventContext.Sandwich.grpcInterface.FetchChannelsByName(ctx.EventContext, *ctx.GuildID, argument)
			if err != nil {
				return nil, xerrors.Errorf("Failed to fetch channel: %v", err)
			}
		}
	} else {
		channelID, _ := strconv.ParseInt(match, 10, 64)

		result := NewChannel(ctx.EventContext, ctx.GuildID, discord.Snowflake(channelID))

		result, err = FetchChannel(ctx.EventContext, result)
		if err != nil && !xerrors.Is(err, ErrChannelNotFound) {
			return nil, err
		}

		results = append(results, result)
	}

	out = make([]*discord_structs.Channel, 0)

	for _, result := range results {
		if len(channelTypes) == 0 {
			out = append(out, result)
		} else {
			for _, channelType := range channelTypes {
				if result.Type == channelType {
					out = append(out, result)

					break
				}
			}
		}
	}

	return out, nil
}

// HandleArgumentTypeInvite handles converting from a string
// argument into a Invite type. Use .Invite() within a command
// to get the proper type.
func HandleArgumentTypeInvite(ctx *CommandContext, argument string) (out interface{}, err error) {
	var result *discord_structs.Invite

	result, err = ctx.EventContext.Session.GetInvite(argument, nil, nil, nil)
	if err != nil {
		return nil, xerrors.Errorf("Failed to convert invite: %v", err)
	}

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

	var result *discord_structs.Guild

	if match == "" {
		guilds, err := ctx.EventContext.Sandwich.grpcInterface.FetchGuildsByName(ctx.EventContext, argument)
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

	var result *discord_structs.Role

	if match == "" {
		if ctx.GuildID != nil {
			roles, err := ctx.EventContext.Sandwich.grpcInterface.FetchRolesByName(ctx.EventContext, *ctx.GuildID, argument)
			if err != nil {
				return nil, xerrors.Errorf("Failed to fetch role: %v", err)
			}

			if len(roles) > 0 {
				result = roles[0]
			}
		}
	} else {
		roleID, _ := strconv.ParseInt(match, 10, 64)

		result = NewRole(ctx.EventContext, ctx.GuildID, discord.Snowflake(roleID))

		result, err = FetchRole(ctx.EventContext, result)
		if err != nil && !xerrors.Is(err, ErrRoleNotFound) {
			return nil, err
		}
	}

	if result == nil {
		return nil, ErrRoleNotFound
	}

	return result, nil
}

// HandleArgumentTypeActivity handles converting from a string
// argument into a Activity type. Use .Activity() within a command
// to get the proper type.
func HandleArgumentTypeActivity(ctx *CommandContext, argument string) (out interface{}, err error) {
	result := &discord_structs.Activity{
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
	} else {
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

func parseHexNumber(arg string) (out uint64, err error) {
	return strconv.ParseUint(arg, 16, 64)
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
// argument into a VoiceChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeVoiceChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord_structs.ChannelTypeGuildVoice)
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
	results, err := findChannel(ctx, argument, discord_structs.ChannelTypeGuildStageVoice)
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
	match := IDRegex.FindString(argument)
	if match == "" {
		matches := EmojiRegex.FindStringSubmatch(argument)
		if len(matches) > 0 {
			match = matches[0]
		}
	}

	var result *discord_structs.Emoji

	if match == "" {
		if ctx.GuildID != nil {
			emojis, err := ctx.EventContext.Sandwich.grpcInterface.FetchEmojisByName(ctx.EventContext, *ctx.GuildID, argument)
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

	result, err = FetchEmoji(ctx.EventContext, result)
	if err != nil && !xerrors.Is(err, ErrEmojiNotFound) {
		return nil, err
	}

	return result, nil
}

// HandleArgumentTypePartialEmoji handles converting from a string
// argument into a Emoji type. Use .Emoji() within a command
// to get the proper type.
func HandleArgumentTypePartialEmoji(ctx *CommandContext, argument string) (out interface{}, err error) {
	matches := PartialEmojiRegex.FindStringSubmatch(argument)

	var result *discord_structs.Emoji

	if len(matches) >= 3 {
		animated, _ := strconv.ParseBool(matches[0])
		id, _ := strconv.ParseInt(matches[2], 10, 64)

		result = &discord_structs.Emoji{
			Animated: animated,
			Name:     matches[1],
			ID:       discord.Snowflake(id),
		}
	}

	return result, nil
}

// HandleArgumentTypeCategoryChannel handles converting from a string
// argument into a CategoryChannel type. Use .Channel() within a command
// to get the proper type.
func HandleArgumentTypeCategoryChannel(ctx *CommandContext, argument string) (out interface{}, err error) {
	results, err := findChannel(ctx, argument, discord_structs.ChannelTypeGuildCategory)
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
	results, err := findChannel(ctx, argument, discord_structs.ChannelTypeGuildStore)
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
	results, err := findChannel(ctx, argument, discord_structs.ChannelTypeGuildNewsThread,
		discord_structs.ChannelTypeGuildPrivateThread, discord_structs.ChannelTypeGuildPublicThread)
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
func (a *Argument) Member() (value *discord_structs.GuildMember, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeMember) {
		value, _ = a.value.(*discord_structs.GuildMember)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustMember will attempt to do Member() and will panic if not possible.
func (a *Argument) MustMember() (value *discord_structs.GuildMember) {
	value, err := a.Member()
	if err != nil {
		panic(fmt.Sprintf(`argument: Member(): %v`, err.Error()))
	}

	return value
}

// User returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) User() (value *discord_structs.User, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeUser) {
		value, _ = a.value.(*discord_structs.User)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustUser will attempt to do User() and will panic if not possible.
func (a *Argument) MustUser() (value *discord_structs.User) {
	value, err := a.User()
	if err != nil {
		panic(fmt.Sprintf(`argument: User(): %v`, err.Error()))
	}

	return value
}

// Channel returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Channel() (value *discord_structs.Channel, err error) {
	if argumentTypeIs(a.ArgumentType,
		ArgumentTypeTextChannel, ArgumentTypeVoiceChannel, ArgumentTypeStageChannel,
		ArgumentTypeCategoryChannel, ArgumentTypeStoreChannel, ArgumentTypeGuildChannel) {
		value, _ = a.value.(*discord_structs.Channel)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustTextChannel will attempt to do Channel() and will panic if not possible.
func (a *Argument) MustChannel() (value *discord_structs.Channel) {
	value, err := a.Channel()
	if err != nil {
		panic(fmt.Sprintf(`argument: Channel(): %v`, err.Error()))
	}

	return value
}

// Invite returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Invite() (value *discord_structs.Invite, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeInvite) {
		value, _ = a.value.(*discord_structs.Invite)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustInvite will attempt to do Invite() and will panic if not possible.
func (a *Argument) MustInvite() (value *discord_structs.Invite) {
	value, err := a.Invite()
	if err != nil {
		panic(fmt.Sprintf(`argument: Invite(): %v`, err.Error()))
	}

	return value
}

// Guild returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Guild() (value *discord_structs.Guild, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeGuild) {
		value, _ = a.value.(*discord_structs.Guild)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustGuild will attempt to do Guild() and will panic if not possible.
func (a *Argument) MustGuild() (value *discord_structs.Guild) {
	value, err := a.Guild()
	if err != nil {
		panic(fmt.Sprintf(`argument: Guild(): %v`, err.Error()))
	}

	return value
}

// Role returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Role() (value *discord_structs.Role, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeRole) {
		value, _ = a.value.(*discord_structs.Role)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustRole will attempt to do Role() and will panic if not possible.
func (a *Argument) MustRole() (value *discord_structs.Role) {
	value, err := a.Role()
	if err != nil {
		panic(fmt.Sprintf(`argument: Role(): %v`, err.Error()))
	}

	return value
}

// Activity returns an argument as the specified Type.
// If the argument is not the right type for the converter
// that made the argument, ErrInvalidArgumentType will be returned.
func (a *Argument) Activity() (value *discord_structs.Activity, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeActivity) {
		value, _ = a.value.(*discord_structs.Activity)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustActivity will attempt to do Activity() and will panic if not possible.
func (a *Argument) MustActivity() (value *discord_structs.Activity) {
	value, err := a.Activity()
	if err != nil {
		panic(fmt.Sprintf(`argument: Activity(): %v`, err.Error()))
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
func (a *Argument) Emoji() (value *discord_structs.Emoji, err error) {
	if argumentTypeIs(a.ArgumentType, ArgumentTypeEmoji, ArgumentTypePartialEmoji) {
		value, _ = a.value.(*discord_structs.Emoji)

		return
	}

	return nil, ErrInvalidArgumentType
}

// MustEmoji will attempt to do Emoji() and will panic if not possible.
func (a *Argument) MustEmoji() (value *discord_structs.Emoji) {
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

	return false, ErrInvalidApplication
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
	co.RegisterConverter(ArgumentTypeInvite, HandleArgumentTypeInvite, nil)
	co.RegisterConverter(ArgumentTypeGuild, HandleArgumentTypeGuild, nil)
	co.RegisterConverter(ArgumentTypeRole, HandleArgumentTypeRole, nil)
	co.RegisterConverter(ArgumentTypeActivity, HandleArgumentTypeActivity, nil)
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
