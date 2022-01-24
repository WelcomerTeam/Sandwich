package internal

import (
	"strconv"
)

// RestVersion is the Discord API version used for the REST api.
var RestVersion = "9"

// Endpoints.
var (
	EndpointDiscord = "https://discord.com/"
	EndpointCDN     = "https://cdn.discordapp.com/"
	EndpointAPI     = EndpointDiscord + "api/v" + RestVersion + "/"

	EndpointGuilds     = EndpointAPI + "guilds/"
	EndpointChannels   = EndpointAPI + "channels/"
	EndpointUsers      = EndpointAPI + "users/"
	EndpointGateway    = EndpointAPI + "gateway"
	EndpointGatewayBot = EndpointGateway + "/bot"
	EndpointWebhooks   = EndpointAPI + "webhooks/"

	EndpointCDNAttachments  = EndpointCDN + "attachments/"
	EndpointCDNAvatars      = EndpointCDN + "avatars/"
	EndpointCDNIcons        = EndpointCDN + "icons/"
	EndpointCDNSplashes     = EndpointCDN + "splashes/"
	EndpointCDNChannelIcons = EndpointCDN + "channel-icons/"
	EndpointCDNBanners      = EndpointCDN + "banners/"

	EndpointUser = func(userID string) string {
		return EndpointUsers + userID
	}
	EndpointUserAvatar = func(userID string, avataroleID string) string {
		return EndpointCDNAvatars + userID + "/" + avataroleID + ".png"
	}
	EndpointUserAvatarAnimated = func(userID string, avataroleID string) string {
		return EndpointCDNAvatars + userID + "/" + avataroleID + ".gif"
	}
	EndpointDefaultUserAvatar = func(discriminator string) string {
		discriminatorInt, _ := strconv.Atoi(discriminator)
		return EndpointCDN + "embed/avatars/" + strconv.Itoa(discriminatorInt%5) + ".png"
	}
	EndpointUserGuilds = func(userID string) string {
		return EndpointUsers + userID + "/guilds"
	}
	EndpointUserGuild = func(userID, guildID string) string {
		return EndpointUsers + userID + "/guilds/" + guildID
	}
	EndpointUserGuildMember = func(userID, guildID string) string {
		return EndpointUserGuild(userID, guildID) + "/member"
	}
	EndpointUserChannels = func(userID string) string {
		return EndpointUsers + userID + "/channels"
	}

	EndpointGuild = func(guildID string) string {
		return EndpointGuilds + guildID
	}
	EndpointGuildPreview = func(guildID string) string {
		return EndpointGuilds + guildID + "/preview"
	}
	EndpointGuildChannels = func(guildID string) string {
		return EndpointGuilds + guildID + "/channels"
	}
	EndpointGuildMembers = func(guildID string) string {
		return EndpointGuilds + guildID + "/members"
	}
	EndpointGuildMember = func(guildID, roleID string) string {
		return EndpointGuilds + guildID + "/members/" + roleID
	}
	EndpointGuildMemberRole = func(guildID, userID, roleID string) string {
		return EndpointGuilds + guildID + "/members/" + roleID + "/roles/" + roleID
	}
	EndpointGuildBans = func(guildID string) string {
		return EndpointGuilds + guildID + "/bans"
	}
	EndpointGuildBan = func(guildID, roleID string) string {
		return EndpointGuilds + guildID + "/bans/" + roleID
	}
	EndpointGuildIntegrations = func(guildID string) string {
		return EndpointGuilds + guildID + "/integrations"
	}
	EndpointGuildIntegration = func(guildID, integrationID string) string {
		return EndpointGuilds + guildID + "/integrations/" + integrationID
	}
	EndpointGuildIntegrationSync = func(guildID, integrationID string) string {
		return EndpointGuilds + guildID + "/integrations/" + integrationID + "/sync"
	}
	EndpointGuildRoles = func(guildID string) string {
		return EndpointGuilds + guildID + "/roles"
	}
	EndpointGuildRole = func(guildID, roleID string) string {
		return EndpointGuilds + guildID + "/roles/" + roleID
	}
	EndpointGuildInvites = func(guildID string) string {
		return EndpointGuilds + guildID + "/invites"
	}
	EndpointGuildWidget = func(guildID string) string {
		return EndpointGuilds + guildID + "/widget"
	}

	EndpointGuildEmbed = EndpointGuildWidget

	EndpointGuildPrune = func(guildID string) string {
		return EndpointGuilds + guildID + "/prune"
	}
	EndpointGuildIcon = func(guildID string, hash string) string {
		return EndpointCDNIcons + guildID + "/" + hash + ".png"
	}
	EndpointGuildIconAnimated = func(guildID string, hash string) string {
		return EndpointCDNIcons + guildID + "/" + hash + ".gif"
	}
	EndpointGuildSplash = func(guildID string, hash string) string {
		return EndpointCDNSplashes + guildID + "/" + hash + ".png"
	}
	EndpointGuildWebhooks = func(guildID string) string {
		return EndpointGuilds + guildID + "/webhooks"
	}
	EndpointGuildAuditLogs = func(guildID string) string {
		return EndpointGuilds + guildID + "/audit-logs"
	}
	EndpointGuildEmojis = func(guildID string) string {
		return EndpointGuilds + guildID + "/emojis"
	}
	EndpointGuildEmoji = func(guildID, emojiID string) string {
		return EndpointGuilds + guildID + "/emojis/" + emojiID
	}
	EndpointGuildBanner = func(guildID string, hash string) string {
		return EndpointCDNBanners + guildID + "/" + hash + ".png"
	}

	EndpointGuildScheduledEvents = func(guildID string) string {
		return EndpointGuilds + guildID + "/scheduled-events"
	}
	EndpointGuildScheduledEvent = func(guildID, eventID string) string {
		return EndpointGuilds + guildID + "/scheduled-events/" + eventID
	}
	EndpointGuildScheduledEventUsers = func(guildID, eventID string) string {
		return EndpointGuildScheduledEvent(guildID, eventID) + "/users"
	}

	EndpointChannel = func(channelID string) string {
		return EndpointChannels + channelID
	}
	EndpointChannelPermissions = func(channelID string) string {
		return EndpointChannels + channelID + "/permissions"
	}
	EndpointChannelPermission = func(channelID, overwriteID string) string {
		return EndpointChannels + channelID + "/permissions/" + overwriteID
	}
	EndpointChannelInvites = func(channelID string) string {
		return EndpointChannels + channelID + "/invites"
	}
	EndpointChannelTyping = func(channelID string) string {
		return EndpointChannels + channelID + "/typing"
	}
	EndpointChannelMessages = func(channelID string) string {
		return EndpointChannels + channelID + "/messages"
	}
	EndpointChannelMessage = func(channelID, messageID string) string {
		return EndpointChannels + channelID + "/messages/" + messageID
	}
	EndpointChannelMessageAck = func(channelID, messageID string) string {
		return EndpointChannels + channelID + "/messages/" + messageID + "/ack"
	}
	EndpointChannelMessagesBulkDelete = func(channelID string) string {
		return EndpointChannel(channelID) + "/messages/bulk-delete"
	}
	EndpointChannelMessagesPins = func(channelID string) string {
		return EndpointChannel(channelID) + "/pins"
	}
	EndpointChannelMessagePin = func(channelID, messageID string) string {
		return EndpointChannel(channelID) + "/pins/" + messageID
	}
	EndpointChannelMessageCrosspost = func(channelID, messageID string) string {
		return EndpointChannel(channelID) + "/messages/" + messageID + "/crosspost"
	}
	EndpointChannelFollow = func(channelID string) string {
		return EndpointChannel(channelID) + "/followers"
	}

	EndpointGroupIcon = func(channelID string, hash string) string {
		return EndpointCDNChannelIcons + channelID + "/" + hash + ".png"
	}

	EndpointChannelWebhooks = func(channelID string) string {
		return EndpointChannel(channelID) + "/webhooks"
	}
	EndpointWebhook = func(webhookID string) string {
		return EndpointWebhooks + webhookID
	}
	EndpointWebhookToken = func(webhookID string, token string) string {
		return EndpointWebhooks + webhookID + "/" + token
	}
	EndpointWebhookMessage = func(webhookID string, token string, messageID string) string {
		return EndpointWebhookToken(webhookID, token) + "/messages/" + messageID
	}

	EndpointMessageReactionsAll = func(channelID, messageID string) string {
		return EndpointChannelMessage(channelID, messageID) + "/reactions"
	}
	EndpointMessageReactions = func(channelID, messageID, emojiID string) string {
		return EndpointChannelMessage(channelID, messageID) + "/reactions/" + emojiID
	}
	EndpointMessageReaction = func(channelID, messageID, emojiID, roleID string) string {
		return EndpointMessageReactions(channelID, messageID, emojiID) + "/" + roleID
	}

	EndpointApplicationGlobalCommands = func(applicationID string) string {
		return EndpointApplication(applicationID) + "/commands"
	}
	EndpointApplicationGlobalCommand = func(applicationID, channelID string) string {
		return EndpointApplicationGlobalCommands(applicationID) + "/" + channelID
	}

	EndpointApplicationGuildCommands = func(applicationID, guildID string) string {
		return EndpointApplication(applicationID) + "/guilds/" + guildID + "/commands"
	}
	EndpointApplicationGuildCommand = func(applicationID, guildID, channelID string) string {
		return EndpointApplicationGuildCommands(applicationID, guildID) + "/" + channelID
	}
	EndpointInteraction = func(applicationID string, interactionToken string) string {
		return EndpointAPI + "interactions/" + applicationID + "/" + interactionToken
	}
	EndpointInteractionResponse = func(interactionID string, interactionToken string) string {
		return EndpointInteraction(interactionID, interactionToken) + "/callback"
	}
	EndpointInteractionResponseActions = func(applicationID string, interactionToken string) string {
		return EndpointWebhookMessage(applicationID, interactionToken, "@original")
	}
	EndpointFollowupMessage = func(applicationID string, interactionToken string) string {
		return EndpointWebhookToken(applicationID, interactionToken)
	}
	EndpointFollowupMessageActions = func(applicationID string, interactionToken string, messageID string) string {
		return EndpointWebhookMessage(applicationID, interactionToken, messageID)
	}

	EndpointRelationships = func() string {
		return EndpointUsers + "@me" + "/relationships"
	}
	EndpointRelationship = func(roleID string) string {
		return EndpointRelationships() + "/" + roleID
	}
	EndpointRelationshipsMutual = func(roleID string) string {
		return EndpointUsers + roleID + "/relationships"
	}

	EndpointGuildCreate = EndpointAPI + "guilds"

	EndpointInvite = func(interactionID string) string {
		return EndpointAPI + "invites/" + interactionID
	}

	EndpointIntegrationsJoin = func(interactionID string) string {
		return EndpointAPI + "integrations/" + interactionID + "/join"
	}

	EndpointEmoji = func(emojiID string) string {
		return EndpointCDN + "emojis/" + emojiID + ".png"
	}
	EndpointEmojiAnimated = func(emojiID string) string {
		return EndpointCDN + "emojis/" + emojiID + ".gif"
	}

	EndpointApplications = EndpointAPI + "applications"
	EndpointApplication  = func(applicationID string) string {
		return EndpointApplications + "/" + applicationID
	}

	EndpointOAuth2             = EndpointAPI + "oauth2"
	EndpointOAuth2Applications = EndpointOAuth2 + "/applications"
	EndpointOAuth2Authorize    = EndpointOAuth2 + "/authorize"
	EndpointOAuth2Token        = EndpointOAuth2 + "/token"
	EndpointOAuth2TokenRevoke  = EndpointOAuth2Token + "/revoke"

	EndpointOAuth2Application = func(applicationID string) string {
		return EndpointOAuth2Applications + "/" + applicationID
	}
	EndpointOAuth2ApplicationsBot = func(applicationID string) string {
		return EndpointOAuth2Applications + "/" + applicationID + "/bot"
	}
	EndpointOAuth2ApplicationAssets = func(applicationID string) string {
		return EndpointOAuth2Applications + "/" + applicationID + "/assets"
	}
)
