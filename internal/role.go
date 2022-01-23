package internal

import (
	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
	"golang.org/x/xerrors"
)

type Role discord.Role

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

	*r = *role

	return
}

type RoleTag discord.RoleTag
