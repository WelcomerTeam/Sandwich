package internal

// RegisterOnReadyEvent adds a new event handler for the READY event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnReadyEvent(event OnReadyFuncType) {
	eventName := "READY"

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
	eventName := "RESUMED"

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
	eventName := "APPLICATION_COMMAND_CREATE"

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
	eventName := "APPLICATION_COMMAND_UPDATE"

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
	eventName := "APPLICATION_COMMAND_DELETE"

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
	eventName := "CHANNEL_CREATE"

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
	eventName := "CHANNEL_UPDATE"

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
	eventName := "CHANNEL_DELETE"

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
	eventName := "CHANNEL_PINS_UPDATE"

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
	eventName := "THREAD_CREATE"

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
	eventName := "THREAD_UPDATE"

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
	eventName := "THREAD_DELETE"

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
	eventName := "THREAD_MEMBER_UPDATE"

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
	eventName := "THREAD_MEMBERS_UPDATE"

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
	eventName := "GUILD_UPDATE"

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
	eventName := "GUILD_BAN_ADD"

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
	eventName := "GUILD_BAN_REMOVE"

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
	eventName := "GUILD_EMOJIS_UPDATE"

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
	eventName := "GUILD_STICKERS_UPDATE"

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
	eventName := "GUILD_INTEGRATIONS_UPDATE"

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
	eventName := "GUILD_MEMBER_ADD"

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
	eventName := "GUILD_MEMBER_REMOVE"

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
	eventName := "GUILD_MEMBER_UPDATE"

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
	eventName := "GUILD_ROLE_CREATE"

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
	eventName := "GUILD_ROLE_UPDATE"

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
	eventName := "GUILD_ROLE_DELETE"

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
	eventName := "INTEGRATION_CREATE"

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
	eventName := "INTEGRATION_UPDATE"

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
	eventName := "INTEGRATION_DELETE"

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
	eventName := "INTERACTION_CREATE"

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
	eventName := "INVITE_CREATE"

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
	eventName := "INVITE_DELETE"

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
	eventName := "MESSAGE_CREATE"

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
	eventName := "MESSAGE_UPDATE"

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
	eventName := "MESSAGE_DELETE"

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
	eventName := "MESSAGE_DELETE_BULK"

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
	eventName := "MESSAGE_REACTION_ADD"

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
	eventName := "MESSAGE_REACTION_REMOVE"

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
	eventName := "MESSAGE_REACTION_REMOVE_ALL"

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
	eventName := "MESSAGE_REACTION_REMOVE_EMOJI"

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
	eventName := "PRESENCE_UPDATE"

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
	eventName := "STAGE_INSTANCE_CREATE"

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
	eventName := "STAGE_INSTANCE_UPDATE"

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
	eventName := "STAGE_INSTANCE_DELETE"

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
	eventName := "TYPING_START"

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
	eventName := "USER_UPDATE"

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
	eventName := "VOICE_STATE_UPDATE"

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
	eventName := "VOICE_SERVER_UPDATE"

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
	eventName := "WEBHOOKS_UPDATE"

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
	eventName := "GUILD_JOIN"

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
	eventName := "GUILD_AVAILABLE"

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
	eventName := "GUILD_LEAVE"

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
	eventName := "}"

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}
