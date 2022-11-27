package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
)

func NewChannel(guildID *discord.Snowflake, channelID discord.Snowflake) *discord.Channel {
	return &discord.Channel{
		ID:      channelID,
		GuildID: guildID,
	}
}

func FetchChannel(ctx *GRPCContext, channel *discord.Channel) (*discord.Channel, error) {
	if channel.Name != "" {
		return channel, nil
	}

	if channel.GuildID == nil {
		return channel, ErrFetchMissingGuild
	}

	channel, err := ctx.GRPCInterface.FetchChannelByID(ctx, *channel.GuildID, channel.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel: %w", err)
	}

	if channel == nil {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}
