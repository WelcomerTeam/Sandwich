package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

// NewGuild creates a new partial guild. Use Fetch() to populate the guild.
func NewGuild(ctx *EventContext, guildID discord.Snowflake) *discord.Guild {
	return &discord.Guild{
		ID: guildID,
	}
}

func FetchGuild(ctx *EventContext, guild *discord.Guild) (*discord.Guild, error) {
	if guild.Name != "" {
		return guild, nil
	}

	guild, err := ctx.Sandwich.GRPCInterface.FetchGuildByID(ctx.ToGRPCContext(), guild.ID)
	if err != nil {
		return guild, errors.Errorf("Failed to fetch guild: %v", err)
	}

	if guild == nil {
		return guild, ErrGuildNotFound
	}

	return guild, nil
}

// NewGuildMember creates a new partial guild member. Use Fetch() to populate the member.
func NewGuildMember(ctx *EventContext, guildID *discord.Snowflake, userID discord.Snowflake) *discord.GuildMember {
	return &discord.GuildMember{
		User: &discord.User{
			ID: userID,
		},
		GuildID: guildID,
	}
}

func FetchGuildMember(ctx *EventContext, guildMember *discord.GuildMember) (*discord.GuildMember, error) {
	if guildMember.User.Username != "" {
		return guildMember, nil
	}

	if guildMember.GuildID == nil {
		return guildMember, ErrFetchMissingGuild
	}

	guildMember, err := ctx.Sandwich.GRPCInterface.FetchMemberByID(ctx.ToGRPCContext(), *guildMember.GuildID, guildMember.User.ID)
	if err != nil {
		return guildMember, errors.Errorf("Failed to fetch member: %v", err)
	}

	if guildMember == nil {
		return guildMember, ErrMemberNotFound
	}

	return guildMember, nil
}
