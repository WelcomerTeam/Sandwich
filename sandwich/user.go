package internal

import (
	"fmt"

	"github.com/WelcomerTeam/Discord/discord"
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

	grpcUser, err := ctx.GRPCInterface.FetchUserByID(ctx, ctx.Identifier.Token, user.ID, createDMChannel)
	if err != nil {
		return user, fmt.Errorf("failed to fetch user: %w", err)
	}

	if grpcUser == nil {
		user, err = discord.GetUser(ctx.Session, user.ID)
		if err != nil {
			return user, ErrUserNotFound
		}

		return user, nil
	} else {
		user = grpcUser
	}

	return user, nil
}
