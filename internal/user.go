package internal

import (
	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
	"golang.org/x/xerrors"
)

type User discord.User

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

	*u = *user

	return
}
