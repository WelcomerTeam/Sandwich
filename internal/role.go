package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

type Role discord_structs.Role

func NewRole(ctx *EventContext, guildID *discord.Snowflake, roleID discord.Snowflake) *Role {
	return &Role{
		ID:      roleID,
		GuildID: guildID,
	}
}

func (r *Role) Fetch(ctx *EventContext) (err error) {
	if r.Name != "" {
		return
	}

	if r.GuildID == nil {
		return ErrFetchMissingGuild
	}

	role, err := ctx.Sandwich.grpcInterface.FetchRoleByID(ctx, *r.GuildID, r.ID)
	if err != nil {
		return xerrors.Errorf("Failed to fetch role: %v", err)
	}

	if role != nil {
		*r = *role
	} else {
		// TODO: Try http

		return ErrRoleNotFound
	}

	return
}

type RoleTag discord_structs.RoleTag
