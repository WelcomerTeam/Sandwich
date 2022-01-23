package internal

import (
	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
	"golang.org/x/xerrors"
)

type Channel discord.Channel

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

	*c = *channel

	return
}

type ChannelOverwrite discord.ChannelOverwrite

type ThreadMetadata discord.ThreadMetadata

type ThreadMember discord.ThreadMember

type StageInstance discord.StageInstance
