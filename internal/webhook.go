package internal

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
)

var WebhookURLRegex = regexp.MustCompile("discord(?:app)?.com/api/webhooks/(?P<id>[0-9]{17,20})/(?P<token>[A-Za-z0-9.-_]{60,68})")

type Webhook discord_structs.Webhook

type WebhookMessage discord_structs.WebhookMessage

func NewWebhook(id discord.Snowflake, token string, webhookType discord_structs.WebhookType) (w *Webhook) {
	w = &Webhook{
		ID:    id,
		Token: token,
		Type:  webhookType,
	}

	return w
}

// WebhookFromURL attempts to return a partial webhook from a URL.
func WebhookFromURL(url string) (w *Webhook, err error) {
	groups := findAllGroups(WebhookURLRegex, url)

	if len(groups) == 0 {
		return nil, ErrBadWebhookArgument
	}

	webhookID, _ := strconv.ParseInt(groups["id"], 10, 64)

	w = NewWebhook(discord.Snowflake(webhookID), groups["token"], discord_structs.WebhookTypeIncoming)

	return
}

// URL returns the URL for this webhook.
func (w *Webhook) URL() string {
	return discord.EndpointDiscord + discord.EndpointWebhookToken(w.ID.String(), w.Token)
}

// Partial returns a partial webhook which only has an ID and Token.
func (w *Webhook) Partial() *Webhook {
	return &Webhook{
		ID:    w.ID,
		Token: w.Token,
		Type:  w.Type,
	}
}

// Delete deletes a webhook.
func (w *Webhook) Delete(ctx *EventContext, reason string, preferAuth bool) (err error) {
	if w.Token == "" && ctx.Identifier.Token == "" {
		return ErrWebhookMissingToken
	}

	if preferAuth && ctx.Identifier.Token != "" {
		_, err = ctx.HTTPSession.FetchBot(ctx, http.MethodDelete, discord.EndpointWebhook(w.ID.String()), "", nil)
	} else {
		_, err = ctx.HTTPSession.Fetch(ctx, http.MethodDelete, discord.EndpointWebhookToken(w.ID.String(), w.Token), "", nil, "")
	}

	return
}

// Edit edits a webhook.
func (w *Webhook) Edit(ctx *EventContext, reason string, name string, avatar []byte, channelID *discord.Snowflake, preferAuth bool) (webhook *Webhook, err error) {
	if w.Token == "" && ctx.Identifier.Token == "" {
		return nil, ErrWebhookMissingToken
	}

	if channelID != nil && ctx.Identifier.Token == "" {
		return nil, ErrWebhookMissingToken
	}

	avatarData, err := bytesToBase64Data(avatar)
	if err != nil {
		return nil, err
	}

	webhook = &Webhook{
		Name:      name,
		Avatar:    avatarData,
		ChannelID: channelID,
	}

	if preferAuth && ctx.Identifier.Token != "" {
		err = ctx.HTTPSession.FetchJSONBot(ctx, http.MethodPatch, discord.EndpointWebhook(w.ID.String()), "", &webhook)
	} else {
		err = ctx.HTTPSession.FetchJSON(ctx, http.MethodPatch, discord.EndpointWebhookToken(w.ID.String(), w.Token), "", &webhook, "")
	}

	return
}

// Send sends a message.
func (w *Webhook) Send(msg *WebhookMessage) (message *WebhookMessage, err error) {
	// TODO

	return nil, nil
}
