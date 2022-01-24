package internal

import (
	"net/http"
	"net/url"

	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

// WebhookCreate creates a new webhook for a channel.
func WebhookCreate(c *Channel, ctx *EventContext, name string, avatar string) (webhook *Webhook, err error) {
	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	endpoint := discord.EndpointChannelWebhooks(c.ID.String())

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPost, endpoint, data, &webhook)
	if err != nil {
		return nil, xerrors.Errorf("Failed to create webhook: %v", err)
	}

	return webhook, nil
}

// ChannelWebhooks returns all webhooks for a channel.
func ChannelWebhooks(c *Channel, ctx *EventContext) (webhooks []*Webhook, err error) {
	endpoint := discord.EndpointChannelWebhooks(c.ID.String())

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &webhooks)
	if err != nil {
		return nil, xerrors.Errorf("Failed to create webhook: %v", err)
	}

	return webhooks, nil
}

// GuildWebhooks returns all webhooks for a guild.
func GuildWebhooks(g *Guild, ctx *EventContext) (webhooks []*Webhook, err error) {
	endpoint := discord.EndpointGuildWebhooks(g.ID.String())

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &webhooks)
	if err != nil {
		return nil, xerrors.Errorf("Failed to create webhook: %v", err)
	}

	return webhooks, nil
}

// Webhook returns a webhook from a partial.
func FetchWebhook(w *Webhook, ctx *EventContext, preferBotAuth bool) (webhook *Webhook, err error) {
	if w.Token == "" && ctx.Identifier.Token == "" {
		return nil, ErrWebhookMissingToken
	}

	var endpoint string

	if preferBotAuth && ctx.Identifier.Token != "" {
		endpoint = discord.EndpointWebhook(w.ID.String())
	} else {
		endpoint = discord.EndpointWebhookToken(w.ID.String(), w.Token)
	}

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &webhook)
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch webhook: %v", err)
	}

	return webhook, nil
}

// WebhookDelete deletes a webhook.
func WebhookDelete(w *Webhook, ctx *EventContext, reason string, preferBotAuth bool) (err error) {
	if w.Token == "" && ctx.Identifier.Token == "" {
		return ErrWebhookMissingToken
	}

	var endpoint string

	if preferBotAuth && ctx.Identifier.Token != "" {
		endpoint = discord.EndpointWebhook(w.ID.String())
	} else {
		endpoint = discord.EndpointWebhookToken(w.ID.String(), w.Token)
	}

	_, err = ctx.HTTPSession.FetchBot(ctx, http.MethodDelete, endpoint, "", nil)
	if err != nil {
		return xerrors.Errorf("Failed to delete webhook: %v", err)
	}

	return
}

// WebhookEdit edits a webhook.
func WebhookEdit(w *Webhook, ctx *EventContext, reason string, name string, avatar []byte, channelID *discord.Snowflake, preferBotAuth bool) (webhook *Webhook, err error) {
	if (w.Token == "" && ctx.Identifier.Token == "") || (channelID != nil && ctx.Identifier.Token == "") {
		return nil, ErrWebhookMissingToken
	}

	avatarData, err := bytesToBase64Data(avatar)
	if err != nil {
		return nil, err
	}

	data := struct {
		Name      string `json:"name"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id"`
	}{name, avatarData, channelID.String()}

	var endpoint string

	if preferBotAuth && ctx.Identifier.Token != "" {
		endpoint = discord.EndpointWebhook(w.ID.String())
	} else {
		endpoint = discord.EndpointWebhookToken(w.ID.String(), w.Token)
	}

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPatch, endpoint, data, &webhook)
	if err != nil {
		return nil, xerrors.Errorf("Failed to edit webhook: %v", err)
	}

	return
}

// WebhookExecute sends a webhook message.
func WebhookExecute(w *Webhook, ctx *EventContext, threadID *discord.Snowflake, wait bool, data *discord_structs.WebhookMessageParams, files []*discord_structs.File) (message *WebhookMessage, err error) {
	endpoint := discord.EndpointWebhookToken(w.ID.String(), w.Token)

	values := url.Values{}

	if threadID != nil {
		values.Set("thread_id", threadID.String())
	}

	if wait {
		values.Set("wait", "true")
	}

	if len(values) > 0 {
		endpoint = endpoint + "?" + values.Encode()
	}

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

	// We will not return the message if we are not waiting.
	if !wait {
		return
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to send webhook message: %v", err)
	}

	return message, nil
}

// WebhookMessage fetches a webhook message.
func FetchWebhookMessage(wm *WebhookMessage, ctx *EventContext, token string) (message *WebhookMessage, err error) {
	if token == "" {
		return nil, ErrWebhookMissingToken
	}

	endpoint := discord.EndpointWebhookMessage(wm.WebhookID.String(), token, wm.ID.String())

	err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodGet, endpoint, nil, &message)
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch webhook message: %v", err)
	}

	return message, nil
}

// WebhookMessageEdit edits a webhook message.
func WebhookMessageEdit(wm *WebhookMessage, ctx *EventContext, token string, threadID *discord.Snowflake, data *discord_structs.WebhookMessageParams, files []*discord_structs.File) (message *WebhookMessage, err error) {
	if token == "" {
		return nil, ErrWebhookMissingToken
	}

	endpoint := discord.EndpointWebhookMessage(wm.WebhookID.String(), token, wm.ID.String())

	var values url.Values

	if threadID != nil {
		values.Set("thread_id", threadID.String())
	}

	if len(values) > 0 {
		endpoint = endpoint + "?" + values.Encode()
	}

	if len(data.Files) > 0 {
		var contentType string

		var body []byte

		contentType, body, err = multipartBodyWithJSON(data, data.Files)
		if err != nil {
			return nil, xerrors.Errorf("Failed to create file body: %v", err)
		}

		err = ctx.HTTPSession.FetchBJBot(ctx, http.MethodPatch, endpoint, contentType, body, &message)
	} else {
		err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPatch, endpoint, data, &message)
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to update webhook message: %v", err)
	}

	return message, nil
}

// WebhookMessageDelete deletes a webhook message.
func WebhookMessageDelete(wm *WebhookMessage, ctx *EventContext, token string, threadID *discord.Snowflake) (err error) {
	if token == "" {
		return ErrWebhookMissingToken
	}

	endpoint := discord.EndpointWebhookMessage(wm.WebhookID.String(), token, wm.ID.String())

	if threadID != nil {
		endpoint += "?thread_id=" + threadID.String()
	}

	_, err = ctx.HTTPSession.FetchBot(ctx, http.MethodDelete, endpoint, "", nil)
	if err != nil {
		return xerrors.Errorf("Failed to delete webhook message: %v", err)
	}

	return
}
