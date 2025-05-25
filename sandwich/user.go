package internal

import (
	"fmt"

	"github.com/WelcomerTeam/Discord/discord"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
)

func NewUser(userID discord.Snowflake) *discord.User {
	return &discord.User{
		ID: userID,
	}
}

func FetchUser(ctx *GRPCContext, user *discord.User, createDMChannel bool) (*discord.User, error) {
	if user.Username != "" || (createDMChannel && user.DMChannelID != nil) {
		return user, nil
	}

	grpcUsers, err := ctx.SandwichClient.FetchUser(ctx, &sandwich_protobuf.FetchUserRequest{
		UserIds: []int64{int64(user.ID)},
	})
	if err != nil {
		return user, fmt.Errorf("failed to fetch user: %w", err)
	}

	grpcUser, ok := grpcUsers.GetUsers()[int64(user.ID)]
	if !ok {
		return nil, ErrUserNotFound
	}

	if grpcUser.GetID() != 0 {
		user := sandwich_protobuf.PBToUser(grpcUser)

		if !user.ID.IsNil() {
			return user, nil
		}
	}

	// Fetch user from discord, if not found in GRPC cache

	user, err = discord.GetUser(ctx.Context, ctx.Session, user.ID)
	if err != nil {
		return user, ErrUserNotFound
	}

	return user, nil
}
