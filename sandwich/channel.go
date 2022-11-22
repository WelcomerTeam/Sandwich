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

func FetchChannel(ctx *EventContext, c *discord.Channel) (channel *discord.Channel, err error) {
	if c.Name != "" {
		return
	}

	if c.GuildID == nil {
		return nil, ErrFetchMissingGuild
	}

	channel, err = ctx.Sandwich.GRPCInterface.FetchChannelByID(ctx, *c.GuildID, c.ID)
	if err != nil {
		return nil, errors.Errorf("Failed to fetch channel: %v", err)
	}

	if channel == nil {
		return nil, ErrChannelNotFound
	}

	return
}
