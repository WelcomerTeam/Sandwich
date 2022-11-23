package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func NewRole(guildID *discord.Snowflake, roleID discord.Snowflake) *discord.Role {
	return &discord.Role{
		ID:      roleID,
		GuildID: guildID,
	}
}

func FetchRole(ctx *GRPCContext, role *discord.Role) (*discord.Role, error) {
	if role.Name != "" {
		return role, nil
	}

	if role.GuildID == nil {
		return role, ErrFetchMissingGuild
	}

	role, err := ctx.GRPCInterface.FetchRoleByID(ctx, *role.GuildID, role.ID)
	if err != nil {
		return role, errors.Errorf("Failed to fetch role: %v", err)
	}

	if role == nil {
		return role, ErrRoleNotFound
	}

	return role, nil
}
