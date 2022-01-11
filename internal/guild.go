package internal

import discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"

type Guild discord.Guild

func NewGuild(ctx *Context, guildID discord.Snowflake) *Guild {
	return &Guild{
		ID: guildID,
	}
}

type UnavailableGuild discord.UnavailableGuild

type GuildMember discord.GuildMember

type VoiceState discord.VoiceState
