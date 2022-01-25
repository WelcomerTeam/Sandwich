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

// Fetch populates a user.
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
		// TODO: Try http

		return ErrUserNotFound
	}

	return
}

// Update updates the current user. This can only be ran on the running application.
func (u *User) Update(ctx *EventContext, username string, avatar []byte) (err error) {
	if u.ID != ctx.payload.Metadata.ApplicationID {
		return ErrInvalidTarget
	}

	user, err := UserUpdate(ctx, username, avatar, nil)
	if err != nil {
		return err
	}

	if user != nil {
		*u = *user
	} else {
		return ErrUserNotFound
	}

	return
}

// FetchGuildMember returns a guild member for the current user.
func (u *User) FetchGuildMember(ctx *EventContext, guildID discord.Snowflake) (guildMember *GuildMember, err error) {
	if u.ID != ctx.payload.Metadata.ApplicationID {
		return nil, ErrInvalidTarget
	}

	guildMember, err = UserGuildMember(ctx, guildID, nil)
	if err != nil {
		return nil, err
	}

	if guildMember != nil {
		return guildMember, nil
	}

	return nil, ErrMemberNotFound
}

// CreateDM creates a new dm channel for a user.
func (u *User) CreateDM(ctx *EventContext) (dmChannel *Channel, err error) {
	if u.DMChannelID != nil {
		return NewChannel(ctx, nil, *u.DMChannelID), nil
	}

	dmChannel, err = UserCreateDM(u, ctx)
	if err != nil {
		return nil, err
	}

	if dmChannel != nil {
		// TODO: Report to sandwich
		u.DMChannelID = &dmChannel.ID
	}

	return dmChannel, nil
}

// Send sends a message to a user.
func (u *User) Send(ctx *EventContext, data *discord_structs.MessageParams, files []*discord_structs.File) (message *Message, err error) {
	var dmChannel *Channel

	if u.DMChannelID == nil {
		dmChannel, err = u.CreateDM(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		dmChannel = NewChannel(ctx, nil, *u.DMChannelID)
	}

	return dmChannel.Send(ctx, data, files)
}
