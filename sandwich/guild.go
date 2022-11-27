package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
)

// NewGuild creates a new partial guild. Use Fetch() to populate the guild.
func NewGuild(guildID discord.Snowflake) *discord.Guild {
	return &discord.Guild{
		ID: guildID,
	}
}

func FetchGuild(ctx *GRPCContext, guild *discord.Guild) (*discord.Guild, error) {
	if guild.Name != "" {
		return guild, nil
	}

	guild, err := ctx.GRPCInterface.FetchGuildByID(ctx, guild.ID)
	if err != nil {
		return guild, fmt.Errorf("failed to fetch guild: %w", err)
	}

	if guild == nil {
		return guild, ErrGuildNotFound
	}

	return guild, nil
}

// NewGuildMember creates a new partial guild member. Use Fetch() to populate the member.
func NewGuildMember(guildID *discord.Snowflake, userID discord.Snowflake) *discord.GuildMember {
	return &discord.GuildMember{
		User: &discord.User{
			ID: userID,
		},
		GuildID: guildID,
	}
}

func FetchGuildMember(ctx *GRPCContext, guildMember *discord.GuildMember) (*discord.GuildMember, error) {
	if guildMember.User.Username != "" {
		return guildMember, nil
	}

	if guildMember.GuildID == nil {
		return guildMember, ErrFetchMissingGuild
	}

	guildMember, err := ctx.GRPCInterface.FetchMemberByID(ctx, *guildMember.GuildID, guildMember.User.ID)
	if err != nil {
		return guildMember, fmt.Errorf("failed to fetch member: %w", err)
	}

	if guildMember == nil {
		return guildMember, ErrMemberNotFound
	}

	return guildMember, nil
}
