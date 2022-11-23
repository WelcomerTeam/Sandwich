package internal

import (
	"image/color"
	"regexp"
	"strconv"

	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func argumentTypeIs(argumentType ArgumentType, argumentTypes ...ArgumentType) bool {
	for _, aType := range argumentTypes {
		if argumentType == aType {
			return true
		}
	}

	return false
}

func findAllGroups(re *regexp.Regexp, s string) map[string]string {
	matches := re.FindStringSubmatch(s)
	subnames := re.SubexpNames()

	if matches == nil || len(matches) != len(subnames) {
		return nil
	}

	matchMap := map[string]string{}
	for i := 1; i < len(matches); i++ {
		matchMap[subnames[i]] = matches[i]
	}

	return matchMap
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

func findChannel(ctx *CommandContext, argument string, channelTypes ...discord.ChannelType) (out []*discord.Channel, err error) {
	match := IDRegex.FindString(argument)
	if match == "" {
		matches := ChannelMentionRegex.FindStringSubmatch(argument)
		if len(matches) > 1 {
			match = matches[1]
		}
	}

	var results []*discord.Channel

	if match == "" {
		if ctx.GuildID != nil {
			results, err = ctx.EventContext.Sandwich.GRPCInterface.FetchChannelsByName(ctx.EventContext.ToGRPCContext(), *ctx.GuildID, argument)
			if err != nil {
				return nil, errors.Errorf("Failed to fetch channel: %v", err)
			}
		}
	} else {
		channelID, _ := strconv.ParseInt(match, 10, 64)

		result := NewChannel(ctx.GuildID, discord.Snowflake(channelID))

		result, err = FetchChannel(ctx.EventContext.ToGRPCContext(), result)
		if err != nil && !errors.Is(err, ErrChannelNotFound) {
			return nil, err
		}

		results = append(results, result)
	}

	out = make([]*discord.Channel, 0)

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
