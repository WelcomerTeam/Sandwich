package internal

import (
	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

func NewUser(ctx *EventContext, userID discord.Snowflake) *discord_structs.User {
	return &discord_structs.User{
		ID: userID,
	}
}

func FetchUser(ctx *EventContext, u *discord_structs.User, createDMChannel bool) (user *discord_structs.User, err error) {
	if u.Username != "" || (createDMChannel && u.DMChannelID != nil) {
		return u, nil
	}

	user, err = ctx.Sandwich.grpcInterface.FetchUserByID(ctx, u.ID, createDMChannel)
	if err != nil {
		return u, xerrors.Errorf("Failed to fetch user: %v", err)
	}

	if user == nil {
		user, err = ctx.Session.GetUser(u.ID)
		if err != nil {
			return u, ErrUserNotFound
		}

		return user, nil
	}

	return
}
