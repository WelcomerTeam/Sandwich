package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	sandwich_daemon "github.com/WelcomerTeam/Sandwich-Daemon"
)

// Discord Events.

// RegisterOnReadyEvent adds a new event handler for the READY event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnReadyEvent(event OnReadyFuncType) {
	eventName := discord.DiscordEventReady

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnResumedEvent adds a new event handler for the RESUMED event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnResumedEvent(event OnResumedFuncType) {
	eventName := discord.DiscordEventResumed

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnApplicationCommandCreateEvent adds a new event handler for the APPLICATION_COMMAND_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnApplicationCommandCreateEvent(event OnApplicationCommandCreateFuncType) {
	eventName := discord.DiscordEventApplicationCommandCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnApplicationCommandUpdateEvent adds a new event handler for the APPLICATION_COMMAND_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnApplicationCommandUpdateEvent(event OnApplicationCommandUpdateFuncType) {
	eventName := discord.DiscordEventApplicationCommandUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnApplicationCommandDeleteEvent adds a new event handler for the APPLICATION_COMMAND_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnApplicationCommandDeleteEvent(event OnApplicationCommandDeleteFuncType) {
	eventName := discord.DiscordEventApplicationCommandDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnChannelCreateEvent adds a new event handler for the CHANNEL_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelCreateEvent(event OnChannelCreateFuncType) {
	eventName := discord.DiscordEventChannelCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnChannelUpdateEvent adds a new event handler for the CHANNEL_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelUpdateEvent(event OnChannelUpdateFuncType) {
	eventName := discord.DiscordEventChannelUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnChannelDeleteEvent adds a new event handler for the CHANNEL_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelDeleteEvent(event OnChannelDeleteFuncType) {
	eventName := discord.DiscordEventChannelDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnChannelPinsUpdateEvent adds a new event handler for the CHANNEL_PINS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelPinsUpdateEvent(event OnChannelPinsUpdateFuncType) {
	eventName := discord.DiscordEventChannelPinsUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnEntitlementCreate adds a new event handler for the ENTITLEMENT_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnEntitlementCreate(event OnEntitlementCreateFuncType) {
	eventName := discord.DiscordEventEntitlementCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnEntitlementUpdate adds a new event handler for the ENTITLEMENT_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnEntitlementUpdate(event OnEntitlementCreateFuncType) {
	eventName := discord.DiscordEventEntitlementUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnEntitlementDelete adds a new event handler for the ENTITLEMENT_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnEntitlementDelete(event OnEntitlementCreateFuncType) {
	eventName := discord.DiscordEventEntitlementDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnThreadCreateEvent adds a new event handler for the THREAD_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadCreateEvent(event OnThreadCreateFuncType) {
	eventName := discord.DiscordEventThreadCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnThreadUpdateEvent adds a new event handler for the THREAD_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadUpdateEvent(event OnThreadUpdateFuncType) {
	eventName := discord.DiscordEventThreadUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnThreadDeleteEvent adds a new event handler for the THREAD_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadDeleteEvent(event OnThreadDeleteFuncType) {
	eventName := discord.DiscordEventThreadDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnThreadMemberUpdateEvent adds a new event handler for the THREAD_MEMBER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadMemberUpdateEvent(event OnThreadMemberUpdateFuncType) {
	eventName := discord.DiscordEventThreadMemberUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnThreadMembersUpdateEvent adds a new event handler for the THREAD_MEMBERS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadMembersUpdateEvent(event OnThreadMembersUpdateFuncType) {
	eventName := discord.DiscordEventThreadMembersUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildUpdateEvent adds a new event handler for the GUILD_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildUpdateEvent(event OnGuildUpdateFuncType) {
	eventName := discord.DiscordEventGuildUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnAuditLogEntryCreateEvent adds a new event handler for the GUILD_AUDIT_LOG_ENTRY_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnAuditGuildAuditLogEntryCreateEvent(event OnGuildAuditLogEntryCreateFuncType) {
	eventName := discord.DiscordEventGuildAuditLogEntryCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildBanAddEvent adds a new event handler for the GUILD_BAN_ADD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildBanAddEvent(event OnGuildBanAddFuncType) {
	eventName := discord.DiscordEventGuildBanAdd

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildBanRemoveEvent adds a new event handler for the GUILD_BAN_REMOVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildBanRemoveEvent(event OnGuildBanRemoveFuncType) {
	eventName := discord.DiscordEventGuildBanRemove

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildEmojisUpdateEvent adds a new event handler for the GUILD_EMOJIS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildEmojisUpdateEvent(event OnGuildEmojisUpdateFuncType) {
	eventName := discord.DiscordEventGuildEmojisUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildStickersUpdateEvent adds a new event handler for the GUILD_STICKERS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildStickersUpdateEvent(event OnGuildStickersUpdateFuncType) {
	eventName := discord.DiscordEventGuildStickersUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildIntegrationsUpdateEvent adds a new event handler for the GUILD_INTEGRATIONS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildIntegrationsUpdateEvent(event OnGuildIntegrationsUpdateFuncType) {
	eventName := discord.DiscordEventGuildIntegrationsUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildMemberAddEvent adds a new event handler for the GUILD_MEMBER_ADD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildMemberAddEvent(event OnGuildMemberAddFuncType) {
	eventName := discord.DiscordEventGuildMemberAdd

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildMemberRemoveEvent adds a new event handler for the GUILD_MEMBER_REMOVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildMemberRemoveEvent(event OnGuildMemberRemoveFuncType) {
	eventName := discord.DiscordEventGuildMemberRemove

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildMemberUpdateEvent adds a new event handler for the GUILD_MEMBER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildMemberUpdateEvent(event OnGuildMemberUpdateFuncType) {
	eventName := discord.DiscordEventGuildMemberUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildRoleCreateEvent adds a new event handler for the GUILD_ROLE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildRoleCreateEvent(event OnGuildRoleCreateFuncType) {
	eventName := discord.DiscordEventGuildRoleCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildRoleUpdateEvent adds a new event handler for the GUILD_ROLE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildRoleUpdateEvent(event OnGuildRoleUpdateFuncType) {
	eventName := discord.DiscordEventGuildRoleUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildRoleDeleteEvent adds a new event handler for the GUILD_ROLE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildRoleDeleteEvent(event OnGuildRoleDeleteFuncType) {
	eventName := discord.DiscordEventGuildRoleDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnIntegrationCreateEvent adds a new event handler for the INTEGRATION_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnIntegrationCreateEvent(event OnIntegrationCreateFuncType) {
	eventName := discord.DiscordEventIntegrationCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnIntegrationUpdateEvent adds a new event handler for the INTEGRATION_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnIntegrationUpdateEvent(event OnIntegrationUpdateFuncType) {
	eventName := discord.DiscordEventIntegrationUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnIntegrationDeleteEvent adds a new event handler for the INTEGRATION_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnIntegrationDeleteEvent(event OnIntegrationDeleteFuncType) {
	eventName := discord.DiscordEventIntegrationDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnInteractionCreateEvent adds a new event handler for the INTERACTION_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnInteractionCreateEvent(event OnInteractionCreateFuncType) {
	eventName := discord.DiscordEventInteractionCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnInviteCreateEvent adds a new event handler for the INVITE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnInviteCreateEvent(event OnInviteCreateFuncType) {
	eventName := discord.DiscordEventInviteCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnInviteDeleteEvent adds a new event handler for the INVITE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnInviteDeleteEvent(event OnInviteDeleteFuncType) {
	eventName := discord.DiscordEventInviteDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageCreateEvent adds a new event handler for the MESSAGE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageCreateEvent(event OnMessageCreateFuncType) {
	eventName := discord.DiscordEventMessageCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageUpdateEvent adds a new event handler for the MESSAGE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageUpdateEvent(event OnMessageUpdateFuncType) {
	eventName := discord.DiscordEventMessageUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageDeleteEvent adds a new event handler for the MESSAGE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageDeleteEvent(event OnMessageDeleteFuncType) {
	eventName := discord.DiscordEventMessageDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageDeleteBulkEvent adds a new event handler for the MESSAGE_DELETE_BULK event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageDeleteBulkEvent(event OnMessageDeleteBulkFuncType) {
	eventName := discord.DiscordEventMessageDeleteBulk

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageReactionAddEvent adds a new event handler for the MESSAGE_REACTION_ADD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionAddEvent(event OnMessageReactionAddFuncType) {
	eventName := discord.DiscordEventMessageReactionAdd

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageReactionRemoveEvent adds a new event handler for the MESSAGE_REACTION_REMOVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionRemoveEvent(event OnMessageReactionRemoveFuncType) {
	eventName := discord.DiscordEventMessageReactionRemove

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageReactionRemoveAllEvent adds a new event handler for the MESSAGE_REACTION_REMOVE_ALL event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionRemoveAllEvent(event OnMessageReactionRemoveAllFuncType) {
	eventName := discord.DiscordEventMessageReactionRemoveAll

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnMessageReactionRemoveEmojiEvent adds a new event handler for the MESSAGE_REACTION_REMOVE_EMOJI event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionRemoveEmojiEvent(event OnMessageReactionRemoveEmojiFuncType) {
	eventName := discord.DiscordEventMessageReactionRemoveEmoji

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnPresenceUpdateEvent adds a new event handler for the PRESENCE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnPresenceUpdateEvent(event OnPresenceUpdateFuncType) {
	eventName := discord.DiscordEventPresenceUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnStageInstanceCreateEvent adds a new event handler for the STAGE_INSTANCE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnStageInstanceCreateEvent(event OnStageInstanceCreateFuncType) {
	eventName := discord.DiscordEventStageInstanceCreate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnStageInstanceUpdateEvent adds a new event handler for the STAGE_INSTANCE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnStageInstanceUpdateEvent(event OnStageInstanceUpdateFuncType) {
	eventName := discord.DiscordEventStageInstanceUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnStageInstanceDeleteEvent adds a new event handler for the STAGE_INSTANCE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnStageInstanceDeleteEvent(event OnStageInstanceDeleteFuncType) {
	eventName := discord.DiscordEventStageInstanceDelete

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnTypingStartEvent adds a new event handler for the TYPING_START event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnTypingStartEvent(event OnTypingStartFuncType) {
	eventName := discord.DiscordEventTypingStart

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnUserUpdateEvent adds a new event handler for the USER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnUserUpdateEvent(event OnUserUpdateFuncType) {
	eventName := discord.DiscordEventUserUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnVoiceStateUpdateEvent adds a new event handler for the VOICE_STATE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnVoiceStateUpdateEvent(event OnVoiceStateUpdateFuncType) {
	eventName := discord.DiscordEventVoiceStateUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnVoiceServerUpdateEvent adds a new event handler for the VOICE_SERVER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnVoiceServerUpdateEvent(event OnVoiceServerUpdateFuncType) {
	eventName := discord.DiscordEventVoiceServerUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnWebhookUpdateEvent adds a new event handler for the WEBHOOKS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnWebhookUpdateEvent(event OnWebhookUpdateFuncType) {
	eventName := discord.DiscordEventWebhookUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildJoinEvent adds a new event handler for the GUILD_JOIN event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildJoinEvent(event OnGuildJoinFuncType) {
	eventName := discord.DiscordEventGuildJoin

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildAvailableEvent adds a new event handler for the GUILD_AVAILABLE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildAvailableEvent(event OnGuildAvailableFuncType) {
	eventName := discord.DiscordEventGuildAvailable

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildLeaveEvent adds a new event handler for the GUILD_LEAVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildLeaveEvent(event OnGuildLeaveFuncType) {
	eventName := discord.DiscordEventGuildLeave

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnGuildUnavailableEvent adds a new event handler for the GUILD_UNAVAILABLE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildUnavailableEvent(event OnGuildUnavailableFuncType) {
	eventName := discord.DiscordEventGuildUnavailable

	h.RegisterEvent(eventName, nil, event)
}

// Sandwich Events.

// RegisterOnSandwichConfigurationReload adds a new event handler for the SW_CONFIGURATION_RELOAD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnSandwichConfigurationReload(event OnSandwichConfigurationReloadFuncType) {
	eventName := sandwich_daemon.SandwichEventConfigUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnSandwichShardStatusUpdate adds a new event handler for the SW_SHARD_STATUS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnSandwichShardStatusUpdate(event OnSandwichShardStatusUpdateFuncType) {
	eventName := sandwich_daemon.SandwichShardStatusUpdate

	h.RegisterEvent(eventName, nil, event)
}

// RegisterOnSandwichApplicationStatusUpdate adds a new event handler for the SW_APPLICATION_STATUS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnSandwichApplicationStatusUpdate(event OnSandwichApplicationStatusUpdateFuncType) {
	eventName := sandwich_daemon.SandwichApplicationStatusUpdate

	h.RegisterEvent(eventName, nil, event)
}

// Generic Events.

// RegisterOnError registers a handler when events raise an error.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnError(event OnErrorFuncType) {
	eventName := "ERROR"

	h.RegisterEvent(eventName, nil, event)
}
