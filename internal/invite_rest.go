package internal

import (
	"net/http"
	"net/url"

	"github.com/WelcomerTeam/Discord/discord"
	"golang.org/x/xerrors"
)

func FetchInvite(i *Invite, ctx *EventContext, withCounts bool, withExpiration bool, guildScheduledEventID *discord.Snowflake) (invite *Invite, err error) {
	endpoint := discord.EndpointInvite(i.Code)

	values := url.Values{}

	if withCounts {
		values.Set("with_counts", "true")
	}

	if withExpiration {
		values.Set("with_expiration", "true")
	}

	if guildScheduledEventID != nil {
		values.Set("guild_scheduled_event_id", guildScheduledEventID.String())
	}

	if len(values) > 0 {
		endpoint = endpoint + "?" + values.Encode()
	}

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, "", &invite)
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch invite: %v", err)
	}

	return
}

func InviteDelete(i *Invite, ctx *EventContext) (err error) {
	endpoint := discord.EndpointInvite(i.Code)

	_, err = ctx.HTTPSession.FetchBot(ctx, http.MethodDelete, endpoint, "", nil)
	if err != nil {
		return xerrors.Errorf("Failed to delete invite: %v", err)
	}

	return
}
