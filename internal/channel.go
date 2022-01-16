package internal

import discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"

type Channel discord.Channel

func NewChannel(ctx *Context, channelID discord.Snowflake, guildID *discord.Snowflake) Channel {
	return Channel{
		ID:      channelID,
		GuildID: guildID,
	}
}

type ChannelOverwrite discord.ChannelOverwrite

type ThreadMetadata discord.ThreadMetadata

type ThreadMember discord.ThreadMember

type StageInstance discord.StageInstance
