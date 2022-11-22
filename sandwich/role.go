package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func NewRole(ctx *EventContext, guildID *discord.Snowflake, roleID discord.Snowflake) *discord.Role {
	return &discord.Role{
		ID:      roleID,
		GuildID: guildID,
	}
}

func FetchRole(ctx *EventContext, r *discord.Role) (role *discord.Role, err error) {
	if r.Name != "" {
		return
	}

	if r.GuildID == nil {
		return r, ErrFetchMissingGuild
	}

	role, err = ctx.Sandwich.GRPCInterface.FetchRoleByID(ctx, *r.GuildID, r.ID)
	if err != nil {
		return r, errors.Errorf("Failed to fetch role: %v", err)
	}

	if role == nil {
		return r, ErrRoleNotFound
	}

	return
}
