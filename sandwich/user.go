package internal

import (
	"github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func NewUser(ctx *EventContext, userID discord.Snowflake) *discord.User {
	return &discord.User{
		ID: userID,
	}
}

func FetchUser(ctx *EventContext, user *discord.User, createDMChannel bool) (*discord.User, error) {
	if user.Username != "" || (createDMChannel && user.DMChannelID != nil) {
		return user, nil
	}

	user, err := ctx.Sandwich.GRPCInterface.FetchUserByID(ctx.ToGRPCContext(), ctx.Identifier.Token, user.ID, createDMChannel)
	if err != nil {
		return user, errors.Errorf("Failed to fetch user: %v", err)
	}

	if user == nil {
		user, err = discord.GetUser(ctx.Session, user.ID)
		if err != nil {
			return user, ErrUserNotFound
		}

		return user, nil
	}

	return user, nil
}
