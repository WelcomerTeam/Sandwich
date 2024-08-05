package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
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

	gRole, err := ctx.GRPCInterface.FetchRoleByID(ctx, *role.GuildID, role.ID)
	if err != nil {
		return role, fmt.Errorf("failed to fetch role: %w", err)
	}

	role = &gRole

	if role.ID.IsNil() {
		return role, ErrRoleNotFound
	}

	return role, nil
}
