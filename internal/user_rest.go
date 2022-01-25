package internal

import (
	"net/http"

	"github.com/WelcomerTeam/Discord/discord"
	"golang.org/x/xerrors"
)

// FetchUser returns a user from a partial.
func FetchUser(u *User, ctx *EventContext) (user *User, err error) {
	endpoint := discord.EndpointUser(u.ID.String())

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &user)
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch user: %v", err)
	}

	return user, nil
}

// FetchCurrentUser returns the current user.
// If the user you are targeting is not the current running application,
// a token must be passed.
func FetchCurrentUser(ctx *EventContext, token *string) (user *User, err error) {
	endpoint := discord.EndpointUser("@me")

	if token == nil {
		err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &user)
	} else {
		err = ctx.HTTPSession.FetchJJ(ctx, http.MethodGet, endpoint, nil, &user, *token)
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch user: %v", err)
	}

	return user, nil
}

// UserUpdate updates the current user.
// If the user you are targeting is not the current running application,
// a token must be passed.
func UserUpdate(ctx *EventContext, username string, avatar []byte, token *string) (user *User, err error) {
	avatarData, err := bytesToBase64Data(avatar)
	if err != nil {
		return nil, err
	}

	data := struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
	}{username, avatarData}

	endpoint := discord.EndpointUser("@me")

	if token == nil {
		err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPatch, endpoint, data, &user)
	} else {
		err = ctx.HTTPSession.FetchJJ(ctx, http.MethodPatch, endpoint, data, &user, *token)
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to update user: %v", err)
	}

	return user, nil
}

// UserGuilds returns a list of all partial guilds a user is in.
func UserGuilds(u *User, ctx *EventContext) (guilds []*Guild, err error) {
	endpoint := discord.EndpointUserGuilds(u.ID.String())

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &guilds)
	if err != nil {
		return nil, xerrors.Errorf("Failed to get current user guilds: %v", err)
	}

	return guilds, nil
}

// UserGuildMember returns a guild member object for the user.
// If the user you are targeting is not the current running application,
// a token must be passed.
func UserGuildMember(ctx *EventContext, guildID discord.Snowflake, token *string) (guildMember *GuildMember, err error) {
	endpoint := discord.EndpointUserGuild("@me", guildID.String())

	if token == nil {
		err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &guildMember)
	} else {
		err = ctx.HTTPSession.FetchJJ(ctx, http.MethodGet, endpoint, nil, &guildMember, *token)
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to get current user guild member: %v", err)
	}

	return guildMember, nil
}

// UserCreateDM creates a new dm channel for a user.
func UserCreateDM(u *User, ctx *EventContext) (dmChannel *Channel, err error) {
	endpoint := discord.EndpointUserChannels("@me")

	data := struct {
		RecipientID discord.Snowflake `json:"recipient_id"`
	}{u.ID}

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPost, endpoint, data, &dmChannel)
	if err != nil {
		return nil, xerrors.Errorf("Failed to create user dm channel: %v", err)
	}

	return dmChannel, nil
}
