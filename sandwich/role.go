package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
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

	gRoles, err := ctx.SandwichClient.FetchGuildRole(ctx, &sandwich_protobuf.FetchGuildRoleRequest{
		GuildId: int64(*role.GuildID),
		RoleIds: []int64{int64(role.ID)},
	})
	if err != nil {
		return role, fmt.Errorf("failed to fetch role: %w", err)
	}

	gRole, ok := gRoles.GetRoles()[int64(role.ID)]
	if !ok {
		return nil, ErrRoleNotFound
	}

	role = sandwich_protobuf.PBToRole(gRole)

	if role.ID.IsNil() {
		return role, ErrRoleNotFound
	}

	return role, nil
}
