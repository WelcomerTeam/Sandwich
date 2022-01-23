package internal

import (
	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
	"golang.org/x/xerrors"
)

type Guild discord.Guild

func NewGuild(ctx *EventContext, guildID discord.Snowflake) *Guild {
	return &Guild{
		ID: guildID,
	}
}

func (g *Guild) Fetch(ctx *EventContext) (err error) {
	if g.Name != "" {
		return
	}

	guild, err := ctx.Sandwich.grpcInterface.FetchGuildByID(ctx, g.ID)
	if err != nil {
		return xerrors.Errorf("Failed to fetch guild: %v", err)
	}

	*g = *guild

	return nil
}

type UnavailableGuild discord.UnavailableGuild

type GuildMember discord.GuildMember

func NewGuildMember(ctx *EventContext, guildID *discord.Snowflake, userID discord.Snowflake) *GuildMember {
	return &GuildMember{
		User: &discord.User{
			ID: userID,
		},
		GuildID: guildID,
	}
}

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

	*gm = *guildMember

	return
}

type VoiceState discord.VoiceState
