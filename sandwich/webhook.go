package internal

import (
	"github.com/WelcomerTeam/Discord/discord"
	"regexp"
	"strconv"
)

var WebhookURLRegex = regexp.MustCompile("discord(?:app)?.com/api/webhooks/(?P<id>[0-9]{17,20})/(?P<token>[A-Za-z0-9._-]{60,68})")

// NewWebhook creates a partial webhook. Use Fetch() to populate the webhook.
func NewWebhook(id discord.Snowflake, token string, webhookType discord.WebhookType) (w *discord.Webhook) {
	w = &discord.Webhook{
		ID:    id,
		Token: token,
		Type:  webhookType,
	}

	return w
}

// WebhookFromURL attempts to return a partial webhook from a URL.
func WebhookFromURL(url string) (w *discord.Webhook, err error) {
	groups := findAllGroups(WebhookURLRegex, url)

	if len(groups) == 0 {
		return nil, ErrBadWebhookArgument
	}

	webhookID, _ := strconv.ParseInt(groups["id"], 10, 64)

	w = NewWebhook(discord.Snowflake(webhookID), groups["token"], discord.WebhookTypeIncoming)

	return
}
