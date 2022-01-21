package internal

import "github.com/WelcomerTeam/Sandwich-Daemon/structs"

// Discord Events.

// RegisterOnReadyEvent adds a new event handler for the READY event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnReadyEvent(event OnReadyFuncType) {
	eventName := structs.DiscordEventReady

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnResumedEvent adds a new event handler for the RESUMED event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnResumedEvent(event OnResumedFuncType) {
	eventName := structs.DiscordEventResumed

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnApplicationCommandCreateEvent adds a new event handler for the APPLICATION_COMMAND_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnApplicationCommandCreateEvent(event OnApplicationCommandCreateFuncType) {
	eventName := structs.DiscordEventApplicationCommandCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnApplicationCommandUpdateEvent adds a new event handler for the APPLICATION_COMMAND_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnApplicationCommandUpdateEvent(event OnApplicationCommandUpdateFuncType) {
	eventName := structs.DiscordEventApplicationCommandUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnApplicationCommandDeleteEvent adds a new event handler for the APPLICATION_COMMAND_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnApplicationCommandDeleteEvent(event OnApplicationCommandDeleteFuncType) {
	eventName := structs.DiscordEventApplicationCommandDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnChannelCreateEvent adds a new event handler for the CHANNEL_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelCreateEvent(event OnChannelCreateFuncType) {
	eventName := structs.DiscordEventChannelCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnChannelUpdateEvent adds a new event handler for the CHANNEL_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelUpdateEvent(event OnChannelUpdateFuncType) {
	eventName := structs.DiscordEventChannelUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnChannelDeleteEvent adds a new event handler for the CHANNEL_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelDeleteEvent(event OnChannelDeleteFuncType) {
	eventName := structs.DiscordEventChannelDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnChannelPinsUpdateEvent adds a new event handler for the CHANNEL_PINS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnChannelPinsUpdateEvent(event OnChannelPinsUpdateFuncType) {
	eventName := structs.DiscordEventChannelPinsUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnThreadCreateEvent adds a new event handler for the THREAD_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadCreateEvent(event OnThreadCreateFuncType) {
	eventName := structs.DiscordEventThreadCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnThreadUpdateEvent adds a new event handler for the THREAD_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadUpdateEvent(event OnThreadUpdateFuncType) {
	eventName := structs.DiscordEventThreadUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnThreadDeleteEvent adds a new event handler for the THREAD_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadDeleteEvent(event OnThreadDeleteFuncType) {
	eventName := structs.DiscordEventThreadDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnThreadMemberUpdateEvent adds a new event handler for the THREAD_MEMBER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadMemberUpdateEvent(event OnThreadMemberUpdateFuncType) {
	eventName := structs.DiscordEventThreadMemberUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnThreadMembersUpdateEvent adds a new event handler for the THREAD_MEMBERS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnThreadMembersUpdateEvent(event OnThreadMembersUpdateFuncType) {
	eventName := structs.DiscordEventThreadMembersUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildUpdateEvent adds a new event handler for the GUILD_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildUpdateEvent(event OnGuildUpdateFuncType) {
	eventName := structs.DiscordEventGuildUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildBanAddEvent adds a new event handler for the GUILD_BAN_ADD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildBanAddEvent(event OnGuildBanAddFuncType) {
	eventName := structs.DiscordEventGuildBanAdd

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildBanRemoveEvent adds a new event handler for the GUILD_BAN_REMOVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildBanRemoveEvent(event OnGuildBanRemoveFuncType) {
	eventName := structs.DiscordEventGuildBanRemove

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildEmojisUpdateEvent adds a new event handler for the GUILD_EMOJIS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildEmojisUpdateEvent(event OnGuildEmojisUpdateFuncType) {
	eventName := structs.DiscordEventGuildEmojisUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildStickersUpdateEvent adds a new event handler for the GUILD_STICKERS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildStickersUpdateEvent(event OnGuildStickersUpdateFuncType) {
	eventName := structs.DiscordEventGuildStickersUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildIntegrationsUpdateEvent adds a new event handler for the GUILD_INTEGRATIONS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildIntegrationsUpdateEvent(event OnGuildIntegrationsUpdateFuncType) {
	eventName := structs.DiscordEventGuildIntegrationsUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildMemberAddEvent adds a new event handler for the GUILD_MEMBER_ADD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildMemberAddEvent(event OnGuildMemberAddFuncType) {
	eventName := structs.DiscordEventGuildMemberAdd

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildMemberRemoveEvent adds a new event handler for the GUILD_MEMBER_REMOVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildMemberRemoveEvent(event OnGuildMemberRemoveFuncType) {
	eventName := structs.DiscordEventGuildMemberRemove

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildMemberUpdateEvent adds a new event handler for the GUILD_MEMBER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildMemberUpdateEvent(event OnGuildMemberUpdateFuncType) {
	eventName := structs.DiscordEventGuildMemberUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildRoleCreateEvent adds a new event handler for the GUILD_ROLE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildRoleCreateEvent(event OnGuildRoleCreateFuncType) {
	eventName := structs.DiscordEventGuildRoleCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildRoleUpdateEvent adds a new event handler for the GUILD_ROLE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildRoleUpdateEvent(event OnGuildRoleUpdateFuncType) {
	eventName := structs.DiscordEventGuildRoleUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildRoleDeleteEvent adds a new event handler for the GUILD_ROLE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildRoleDeleteEvent(event OnGuildRoleDeleteFuncType) {
	eventName := structs.DiscordEventGuildRoleDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnIntegrationCreateEvent adds a new event handler for the INTEGRATION_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnIntegrationCreateEvent(event OnIntegrationCreateFuncType) {
	eventName := structs.DiscordEventIntegrationCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnIntegrationUpdateEvent adds a new event handler for the INTEGRATION_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnIntegrationUpdateEvent(event OnIntegrationUpdateFuncType) {
	eventName := structs.DiscordEventIntegrationUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnIntegrationDeleteEvent adds a new event handler for the INTEGRATION_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnIntegrationDeleteEvent(event OnIntegrationDeleteFuncType) {
	eventName := structs.DiscordEventIntegrationDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnInteractionCreateEvent adds a new event handler for the INTERACTION_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnInteractionCreateEvent(event OnInteractionCreateFuncType) {
	eventName := structs.DiscordEventInteractionCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnInviteCreateEvent adds a new event handler for the INVITE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnInviteCreateEvent(event OnInviteCreateFuncType) {
	eventName := structs.DiscordEventInviteCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnInviteDeleteEvent adds a new event handler for the INVITE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnInviteDeleteEvent(event OnInviteDeleteFuncType) {
	eventName := structs.DiscordEventInviteDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageCreateEvent adds a new event handler for the MESSAGE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageCreateEvent(event OnMessageCreateFuncType) {
	eventName := structs.DiscordEventMessageCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageUpdateEvent adds a new event handler for the MESSAGE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageUpdateEvent(event OnMessageUpdateFuncType) {
	eventName := structs.DiscordEventMessageUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageDeleteEvent adds a new event handler for the MESSAGE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageDeleteEvent(event OnMessageDeleteFuncType) {
	eventName := structs.DiscordEventMessageDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageDeleteBulkEvent adds a new event handler for the MESSAGE_DELETE_BULK event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageDeleteBulkEvent(event OnMessageDeleteBulkFuncType) {
	eventName := structs.DiscordEventMessageDeleteBulk

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageReactionAddEvent adds a new event handler for the MESSAGE_REACTION_ADD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionAddEvent(event OnMessageReactionAddFuncType) {
	eventName := structs.DiscordEventMessageReactionAdd

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageReactionRemoveEvent adds a new event handler for the MESSAGE_REACTION_REMOVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionRemoveEvent(event OnMessageReactionRemoveFuncType) {
	eventName := structs.DiscordEventMessageReactionRemove

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageReactionRemoveAllEvent adds a new event handler for the MESSAGE_REACTION_REMOVE_ALL event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionRemoveAllEvent(event OnMessageReactionRemoveAllFuncType) {
	eventName := structs.DiscordEventMessageReactionRemoveAll

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnMessageReactionRemoveEmojiEvent adds a new event handler for the MESSAGE_REACTION_REMOVE_EMOJI event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnMessageReactionRemoveEmojiEvent(event OnMessageReactionRemoveEmojiFuncType) {
	eventName := structs.DiscordEventMessageReactionRemoveEmoji

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnPresenceUpdateEvent adds a new event handler for the PRESENCE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnPresenceUpdateEvent(event OnPresenceUpdateFuncType) {
	eventName := structs.DiscordEventPresenceUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnStageInstanceCreateEvent adds a new event handler for the STAGE_INSTANCE_CREATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnStageInstanceCreateEvent(event OnStageInstanceCreateFuncType) {
	eventName := structs.DiscordEventStageInstanceCreate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnStageInstanceUpdateEvent adds a new event handler for the STAGE_INSTANCE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnStageInstanceUpdateEvent(event OnStageInstanceUpdateFuncType) {
	eventName := structs.DiscordEventStageInstanceUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnStageInstanceDeleteEvent adds a new event handler for the STAGE_INSTANCE_DELETE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnStageInstanceDeleteEvent(event OnStageInstanceDeleteFuncType) {
	eventName := structs.DiscordEventStageInstanceDelete

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnTypingStartEvent adds a new event handler for the TYPING_START event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnTypingStartEvent(event OnTypingStartFuncType) {
	eventName := structs.DiscordEventTypingStart

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnUserUpdateEvent adds a new event handler for the USER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnUserUpdateEvent(event OnUserUpdateFuncType) {
	eventName := structs.DiscordEventUserUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnVoiceStateUpdateEvent adds a new event handler for the VOICE_STATE_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnVoiceStateUpdateEvent(event OnVoiceStateUpdateFuncType) {
	eventName := structs.DiscordEventVoiceStateUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnVoiceServerUpdateEvent adds a new event handler for the VOICE_SERVER_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnVoiceServerUpdateEvent(event OnVoiceServerUpdateFuncType) {
	eventName := structs.DiscordEventVoiceServerUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnWebhookUpdateEvent adds a new event handler for the WEBHOOKS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnWebhookUpdateEvent(event OnWebhookUpdateFuncType) {
	eventName := structs.DiscordEventWebhookUpdate

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildJoinEvent adds a new event handler for the GUILD_JOIN event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildJoinEvent(event OnGuildJoinFuncType) {
	eventName := structs.DiscordEventGuildJoin

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildAvailableEvent adds a new event handler for the GUILD_AVAILABLE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildAvailableEvent(event OnGuildAvailableFuncType) {
	eventName := structs.DiscordEventGuildAvailable

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnGuildLeaveEvent adds a new event handler for the GUILD_LEAVE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildLeaveEvent(event OnGuildLeaveFuncType) {
	eventName := structs.DiscordEventGuildLeave

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)

	eh.eventsMu.Unlock()
}

// RegisterOnGuildUnavailableEvent adds a new event handler for the GUILD_UNAVAILABLE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildUnavailableEvent(event OnGuildUnavailableFuncType) {
	eventName := structs.DiscordEventGuildUnavailable

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// Sandwich Events.

// RegisterOnSandwichConfigurationReload adds a new event handler for the SW_CONFIGURATION_RELOAD event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnSandwichConfigurationReload(event OnSandwichConfigurationReloadFuncType) {
	eventName := "SW_CONFIGURATION_RELOAD"

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnSandwichShardStatusUpdate adds a new event handler for the SW_SHARD_STATUS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnSandwichShardStatusUpdate(event OnSandwichShardStatusUpdateFuncType) {
	eventName := "SW_SHARD_STATUS_UPDATE"

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// RegisterOnSandwichShardGroupStatusUpdate adds a new event handler for the SW_SHARD_GROUP_STATUS_UPDATE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnSandwichShardGroupStatusUpdate(event OnSandwichShardGroupStatusUpdateFuncType) {
	eventName := "SW_SHARD_GROUP_STATUS_UPDATE"

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}

// Generic Events.

// RegisterOnError registers a handler when events raise an error.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnError(event OnErrorFuncType) {
	eventName := "ERROR"

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}
