package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"

	"golang.org/x/xerrors"
)

type Channel discord_structs.Channel

func NewChannel(ctx *EventContext, guildID *discord.Snowflake, channelID discord.Snowflake) *Channel {
	return &Channel{
		ID:      channelID,
		GuildID: guildID,
	}
}

func (c *Channel) Fetch(ctx *EventContext) (err error) {
	if c.Name != "" {
		return
	}

	if c.GuildID == nil {
		return ErrFetchMissingGuild
	}

	channel, err := ctx.Sandwich.grpcInterface.FetchChannelByID(ctx, *c.GuildID, c.ID)
	if err != nil {
		return xerrors.Errorf("Failed to fetch channel: %v", err)
	}

	if channel != nil {
		*c = *channel
	} else {
		return ErrChannelNotFound
	}

	return
}

type ChannelOverwrite discord_structs.ChannelOverwrite

type ThreadMetadata discord_structs.ThreadMetadata

type ThreadMember discord_structs.ThreadMember

type StageInstance discord_structs.StageInstance
