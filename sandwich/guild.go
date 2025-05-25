package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
	sandwich_daemon "github.com/WelcomerTeam/Sandwich-Daemon"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
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

	gGuilds, err := ctx.SandwichClient.FetchGuild(ctx, &sandwich_protobuf.FetchGuildRequest{
		GuildIds: []int64{int64(guild.ID)},
	})
	if err != nil {
		return guild, fmt.Errorf("failed to fetch guild: %w", err)
	}

	gGuild, ok := gGuilds.GetGuilds()[int64(guild.ID)]
	if !ok {
		return guild, ErrGuildNotFound
	}

	guild = sandwich_daemon.PBToGuild(gGuild)

	if guild.ID.IsNil() {
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

	gGuildMembers, err := ctx.SandwichClient.FetchGuildMember(ctx, &sandwich_protobuf.FetchGuildMemberRequest{
		GuildId: int64(*guildMember.GuildID),
		UserIds: []int64{int64(guildMember.User.ID)},
	})
	if err != nil {
		return guildMember, fmt.Errorf("failed to fetch member: %w", err)
	}

	gGuildMember, ok := gGuildMembers.GetGuildMembers()[int64(guildMember.User.ID)]
	if !ok {
		return nil, ErrMemberNotFound
	}

	guildMember = sandwich_daemon.PBToGuildMember(gGuildMember)

	if guildMember.User == nil || guildMember.User.ID.IsNil() {
		return guildMember, ErrMemberNotFound
	}

	return guildMember, nil
}
