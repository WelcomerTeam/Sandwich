package internal

import (
	"regexp"
	"strconv"

	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
)

var WebhookURLRegex = regexp.MustCompile("discord(?:app)?.com/api/webhooks/(?P<id>[0-9]{17,20})/(?P<token>[A-Za-z0-9.-_]{60,68})")

type Webhook discord_structs.Webhook

type WebhookMessage discord_structs.Message

// NewWebhook creates a partial webhook. Use Fetch() to populate the webhook.
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

// Fetch populates a webhook.
func (w *Webhook) Fetch(ctx *EventContext, preferBotAuth bool) (err error) {
	webhook, err := FetchWebhook(w, ctx, preferBotAuth)
	if err != nil {
		return
	}

	*w = *webhook

	return
}

// Edit edits a webhook.
func (w *Webhook) Edit(ctx *EventContext, reason string, name string, avatar []byte, channelID *discord.Snowflake) (err error) {
	webhook, err := WebhookEdit(w, ctx, reason, name, avatar, channelID, false)
	if err != nil {
		return
	}

	*w = *webhook

	return
}

// Delete deletes a webhook.
func (w *Webhook) Delete(ctx *EventContext, reason string) (err error) {
	return WebhookDelete(w, ctx, reason, false)
}

// Send executes a webhook message.
func (w *Webhook) Send(ctx *EventContext, data *discord_structs.WebhookMessageParams, files []*discord_structs.File, wait bool, threadID *discord.Snowflake) (message *WebhookMessage, err error) {
	return WebhookExecute(w, ctx, threadID, wait, data, files)
}

// Fetch populates a webhook message.
func (wm *WebhookMessage) Fetch(ctx *EventContext, token string) (message *WebhookMessage, err error) {
	return FetchWebhookMessage(wm, ctx, token)
}

// Edit edits a webhook message.
func (wm *WebhookMessage) Edit(ctx *EventContext, data *discord_structs.WebhookMessageParams, files []*discord_structs.File, threadID *discord.Snowflake, token string) (message *WebhookMessage, err error) {
	return WebhookMessageEdit(wm, ctx, token, threadID, data, files)
}

// Delete deletes a webhook message.
func (wm *WebhookMessage) Delete(ctx *EventContext, threadID *discord.Snowflake, token string) (err error) {
	return WebhookMessageDelete(wm, ctx, token, threadID)
}
