package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"

	"golang.org/x/xerrors"
)

type Guild discord_structs.Guild

// NewGuild creates a new partial guild. Use Fetch() to populate the guild.
func NewGuild(ctx *EventContext, guildID discord.Snowflake) *Guild {
	return &Guild{
		ID: guildID,
	}
}

// Fetch populates the guild.
func (g *Guild) Fetch(ctx *EventContext) (err error) {
	if g.Name != "" {
		return
	}

	guild, err := ctx.Sandwich.grpcInterface.FetchGuildByID(ctx, g.ID)
	if err != nil {
		return xerrors.Errorf("Failed to fetch guild: %v", err)
	}

	if guild != nil {
		*g = *guild
	} else {
		// TODO: Try http

		return ErrGuildNotFound
	}

	return nil
}

// Webhooks returns all webhooks a guild has.
func (g *Guild) Webhooks(ctx *EventContext) (webhooks []*Webhook, err error) {
	return GuildWebhooks(g, ctx)
}

type UnavailableGuild discord_structs.UnavailableGuild

type GuildMember discord_structs.GuildMember

// NewGuildMember creates a new partial guild member. Use Fetch() to populate the member.
func NewGuildMember(ctx *EventContext, guildID *discord.Snowflake, userID discord.Snowflake) *GuildMember {
	return &GuildMember{
		User: &discord_structs.User{
			ID: userID,
		},
		GuildID: guildID,
	}
}

// Fetch populates the guild member.
func (gm *GuildMember) Fetch(ctx *EventContext) (err error) {
	if gm.User.Username != "" {
		return
	}

	if gm.GuildID == nil {
		return ErrFetchMissingGuild
	}

	guildMember, err := ctx.Sandwich.grpcInterface.FetchMemberByID(ctx, *gm.GuildID, gm.User.ID)
	if err != nil {
		return xerrors.Errorf("Failed to fetch member: %v", err)
	}

	if guildMember != nil {
		*gm = *guildMember
	} else {
		// TODO: Try http

		return ErrMemberNotFound
	}

	return
}

type VoiceState discord_structs.VoiceState
