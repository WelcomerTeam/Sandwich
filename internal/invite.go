package internal

import (
	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
)

type Invite discord_structs.Invite

func (i *Invite) Fetch(ctx *EventContext, withCounts bool, withExpiration bool, guildScheduledEventID *discord.Snowflake) (err error) {
	invite, err := FetchInvite(i, ctx, withCounts, withExpiration, guildScheduledEventID)
	if err != nil {
		return err
	}

	*i = *invite

	return
}

func (i *Invite) Delete(ctx *EventContext) (err error) {
	return InviteDelete(i, ctx)
}
