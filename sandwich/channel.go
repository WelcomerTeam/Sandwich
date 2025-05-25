package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
	sandwich_daemon "github.com/WelcomerTeam/Sandwich-Daemon"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
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

	gChannels, err := ctx.SandwichClient.FetchGuildChannel(ctx, &sandwich_protobuf.FetchGuildChannelRequest{
		GuildId:    int64(*channel.GuildID),
		ChannelIds: []int64{int64(channel.ID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel: %w", err)
	}

	gChannel, ok := gChannels.GetChannels()[int64(channel.ID)]
	if !ok {
		return nil, ErrChannelNotFound
	}

	channel = sandwich_daemon.PBToChannel(gChannel)

	if channel.ID.IsNil() {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}
