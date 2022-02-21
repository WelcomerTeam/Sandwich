package internal

import (
	"github.com/WelcomerTeam/Discord/discord"
	"golang.org/x/xerrors"
)

func NewUser(ctx *EventContext, userID discord.Snowflake) *discord.User {
	return &discord.User{
		ID: userID,
	}
}

func FetchUser(ctx *EventContext, u *discord.User, createDMChannel bool) (user *discord.User, err error) {
	if u.Username != "" || (createDMChannel && u.DMChannelID != nil) {
		return u, nil
	}

	user, err = ctx.Sandwich.GRPCInterface.FetchUserByID(ctx, u.ID, createDMChannel)
	if err != nil {
		return u, xerrors.Errorf("Failed to fetch user: %v", err)
	}

	if user == nil {
		user, err = discord.GetUser(ctx.Session, u.ID)
		if err != nil {
			return u, ErrUserNotFound
		}

		return user, nil
	}

	return
}
