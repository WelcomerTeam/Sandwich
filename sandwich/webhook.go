package internal

import (
	"regexp"
	"strconv"

	"github.com/WelcomerTeam/Discord/discord"
)

var WebhookURLRegex = regexp.MustCompile("discord(?:app)?.com/api/webhooks/(?P<id>[0-9]{17,20})/(?P<token>[A-Za-z0-9._-]{60,68})")

// NewWebhook creates a partial webhook. Use Fetch() to populate the webhook.
func NewWebhook(id discord.Snowflake, token string, webhookType discord.WebhookType) *discord.Webhook {
	webhook := &discord.Webhook{
		ID:    id,
		Token: token,
		Type:  webhookType,
	}

	return webhook
}

// WebhookFromURL attempts to return a partial webhook from a URL.
func WebhookFromURL(url string) (*discord.Webhook, error) {
	groups := findAllGroups(WebhookURLRegex, url)

	if len(groups) == 0 {
		return nil, ErrBadWebhookArgument
	}

	webhookID, _ := strconv.ParseInt(groups["id"], 10, 64)

	webhook := NewWebhook(discord.Snowflake(webhookID), groups["token"], discord.WebhookTypeIncoming)

	return webhook, nil
}

func findAllGroups(re *regexp.Regexp, s string) map[string]string {
	matches := re.FindStringSubmatch(s)
	subnames := re.SubexpNames()

	if matches == nil || len(matches) != len(subnames) {
		return nil
	}

	matchMap := map[string]string{}
	for i := 1; i < len(matches); i++ {
		matchMap[subnames[i]] = matches[i]
	}

	return matchMap
}
