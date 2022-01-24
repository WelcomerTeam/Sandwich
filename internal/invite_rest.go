package internal

import (
	"net/http"

	"github.com/WelcomerTeam/Discord/discord"
	"golang.org/x/xerrors"
)

// GuildInvites(guildID string) (st []*Invite, err error) {
// ChannelInvites(channelID string) (st []*Invite, err error) {
// ChannelInviteCreate(channelID string, i Invite) (st *Invite, err error) {

func FetchInvite(i *Invite, ctx *EventContext, withCounts bool, withExpiration bool, guildScheduledEventID *discord.Snowflake) (invite *Invite, err error) {
	url := discord.EndpointInvite(i.Code)

	// TODO: Construct query

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, url, "", &invite)
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch invite: %v", err)
	}

	return
}

func InviteDelete(i *Invite, ctx *EventContext) (err error) {
	url := discord.EndpointInvite(i.Code)

	_, err = ctx.HTTPSession.FetchBot(ctx, http.MethodDelete, url, "", nil)
	if err != nil {
		return xerrors.Errorf("Failed to delete invite: %v", err)
	}

	return
}
