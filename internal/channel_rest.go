package internal

import (
	"net/http"

	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

// ChannelMessageSend sends a message to a channel.
func ChannelMessageSend(c *Channel, ctx *EventContext, data *discord_structs.MessageParams, files []*discord_structs.File) (message *Message, err error) {
	endpoint := discord.EndpointChannelMessages(c.ID.String())

	if len(data.Files) > 0 {
		var contentType string

		var body []byte

		contentType, body, err = multipartBodyWithJSON(data, data.Files)
		if err != nil {
			return nil, xerrors.Errorf("Failed to create file body: %v", err)
		}

		err = ctx.HTTPSession.FetchBJBot(ctx, http.MethodPost, endpoint, contentType, body, &message)
	} else {
		err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPost, endpoint, data, &message)
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to send message: %v", err)
	}

	return message, nil
}

// TODO: ChannelMessageEdit
// TODO: ChannelMessageDelete
// TODO: ChannelMessagesBulkDelete

// TODO: ChannelMessagePin
// TODO: ChannelMessageUnpin
// TODO: ChannelMessagesPinned

// TODO: ChannelMessageCrosspost
// TODO: ChannelNewsFollow

// TODO: ChannelMessages
// TODO: ChannelMessage

// TODO: FetchChannel
// TODO: ChannelEdit
// TODO: ChannelDelete
// TODO: ChannelTyping

// TODO: ChannelInvites
// TODO: ChannelInviteCreate

// TODO: ChannelPermissionSet
// TODO: ChannelPermissionDelete
