package internal

import discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"

type Guild discord.Guild

func NewGuild(ctx *EventContext, guildID discord.Snowflake) *Guild {
	return &Guild{
		ID: guildID,
	}
}

func (g *Guild) GetMemberNamed(ctx *EventContext, argument string) (member *GuildMember, err error) {
	if len(argument) > 5 && argument[len(argument)-5] == '#' {
		// username, _, discriminator := rpartition(argument, "#")
		// TODO: GRPC to fetch user by username
		// return first result by discriminator
	} else {
		// TODO: GRPC to fetch user
		// find by name or nick
	}

	return nil, nil
}

func (g *Guild) GetMemberById(ctx *EventContext, userID discord.Snowflake) (member *GuildMember, err error) {
	// TODO: GRPC to fetch user by id

	return nil, nil
}

type UnavailableGuild discord.UnavailableGuild

type GuildMember discord.GuildMember

type VoiceState discord.VoiceState
