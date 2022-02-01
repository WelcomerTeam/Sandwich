package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"

	"golang.org/x/xerrors"
)

func NewChannel(ctx *EventContext, guildID *discord.Snowflake, channelID discord.Snowflake) *discord_structs.Channel {
	return &discord_structs.Channel{
		ID:      channelID,
		GuildID: guildID,
	}
}

func FetchChannel(ctx *EventContext, c *discord_structs.Channel) (channel *discord_structs.Channel, err error) {
	if c.Name != "" {
		return
	}

	if c.GuildID == nil {
		return nil, ErrFetchMissingGuild
	}

	channel, err = ctx.Sandwich.grpcInterface.FetchChannelByID(ctx, *c.GuildID, c.ID)
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch channel: %v", err)
	}

	if channel == nil {
		return nil, ErrChannelNotFound
	}

	return
}
