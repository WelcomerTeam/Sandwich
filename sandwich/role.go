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

func FetchRole(ctx *EventContext, role *discord.Role) (*discord.Role, error) {
	if role.Name != "" {
		return role, nil
	}

	if role.GuildID == nil {
		return role, ErrFetchMissingGuild
	}

	role, err := ctx.Sandwich.GRPCInterface.FetchRoleByID(ctx.ToGRPCContext(), *role.GuildID, role.ID)
	if err != nil {
		return role, errors.Errorf("Failed to fetch role: %v", err)
	}

	if role == nil {
		return role, ErrRoleNotFound
	}

	return role, nil
}
