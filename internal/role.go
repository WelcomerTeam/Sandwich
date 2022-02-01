package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

func NewRole(ctx *EventContext, guildID *discord.Snowflake, roleID discord.Snowflake) *discord_structs.Role {
	return &discord_structs.Role{
		ID:      roleID,
		GuildID: guildID,
	}
}

func FetchRole(ctx *EventContext, r *discord_structs.Role) (role *discord_structs.Role, err error) {
	if r.Name != "" {
		return
	}

	if r.GuildID == nil {
		return r, ErrFetchMissingGuild
	}

	role, err = ctx.Sandwich.grpcInterface.FetchRoleByID(ctx, *r.GuildID, r.ID)
	if err != nil {
		return r, xerrors.Errorf("Failed to fetch role: %v", err)
	}

	if role == nil {
		return r, ErrRoleNotFound
	}

	return
}
