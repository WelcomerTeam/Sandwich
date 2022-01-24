package internal

import (
	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

type User discord_structs.User

func NewUser(ctx *EventContext, userID discord.Snowflake) *User {
	return &User{
		ID: userID,
	}
}

func (u *User) Fetch(ctx *EventContext, createDMChannel bool) (err error) {
	if u.Username != "" || (createDMChannel && u.DMChannelID == nil) {
		return
	}

	user, err := ctx.Sandwich.grpcInterface.FetchUserByID(ctx, u.ID, createDMChannel)
	if err != nil {
		return xerrors.Errorf("Failed to fetch user: %v", err)
	}

	if user != nil {
		*u = *user
	} else {
		return ErrUserNotFound
	}

	return
}
