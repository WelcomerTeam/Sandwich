package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func NewChannel(ctx *EventContext, guildID *discord.Snowflake, channelID discord.Snowflake) *discord.Channel {
	return &discord.Channel{
		ID:      channelID,
		GuildID: guildID,
	}
}

func FetchChannel(ctx *EventContext, channel *discord.Channel) (*discord.Channel, error) {
	if channel.Name != "" {
		return channel, nil
	}

	if channel.GuildID == nil {
		return channel, ErrFetchMissingGuild
	}

	channel, err := ctx.Sandwich.GRPCInterface.FetchChannelByID(ctx.ToGRPCContext(), *channel.GuildID, channel.ID)
	if err != nil {
		return nil, errors.Errorf("Failed to fetch channel: %v", err)
	}

	if channel == nil {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}
