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

func FetchGuild(ctx *EventContext, g *discord.Guild) (guild *discord.Guild, err error) {
	if g.Name != "" {
		return g, nil
	}

	guild, err = ctx.Sandwich.GRPCInterface.FetchGuildByID(ctx.ToGRPCContext(), g.ID)
	if err != nil {
		return g, errors.Errorf("Failed to fetch guild: %v", err)
	}

	if guild == nil {
		return g, ErrGuildNotFound
	}

	return
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

func FetchGuildMember(ctx *EventContext, gm *discord.GuildMember) (guildMember *discord.GuildMember, err error) {
	if gm.User.Username != "" {
		return gm, nil
	}

	if gm.GuildID == nil {
		return gm, ErrFetchMissingGuild
	}

	guildMember, err = ctx.Sandwich.GRPCInterface.FetchMemberByID(ctx.ToGRPCContext(), *gm.GuildID, gm.User.ID)
	if err != nil {
		return gm, errors.Errorf("Failed to fetch member: %v", err)
	}

	if guildMember == nil {
		return gm, ErrMemberNotFound
	}

	return
}
