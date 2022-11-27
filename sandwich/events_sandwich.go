package internal

import discord "github.com/WelcomerTeam/Discord/discord"

// Discord Events.

// RegisterOnReadyEvent adds a new event handler for the READY event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnReadyEvent(event OnReadyFuncType) {
	eventName := discord.DiscordEventReady

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventResumed

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventApplicationCommandCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventApplicationCommandUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventApplicationCommandDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventChannelCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventChannelUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventChannelDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventChannelPinsUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventThreadCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventThreadUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventThreadDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventThreadMemberUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventThreadMembersUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildBanAdd

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildBanRemove

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildEmojisUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildStickersUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildIntegrationsUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildMemberAdd

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildMemberRemove

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildMemberUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildRoleCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildRoleUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildRoleDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventIntegrationCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventIntegrationUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventIntegrationDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventInteractionCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventInviteCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventInviteDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageDeleteBulk

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageReactionAdd

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageReactionRemove

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageReactionRemoveAll

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventMessageReactionRemoveEmoji

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventPresenceUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventStageInstanceCreate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventStageInstanceUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventStageInstanceDelete

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventTypingStart

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventUserUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventVoiceStateUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventVoiceServerUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventWebhookUpdate

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildJoin

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildAvailable

	h.ensureEvent(eventName)

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
	eventName := discord.DiscordEventGuildLeave

	h.ensureEvent(eventName)

	h.eventHandlersMu.RLock()
	eventHandler := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eventHandler.eventsMu.Lock()
	eventHandler.Events = append(eventHandler.Events, event)

	eventHandler.eventsMu.Unlock()
}

// RegisterOnGuildUnavailableEvent adds a new event handler for the GUILD_UNAVAILABLE event.
// It does not override a handler and instead will add another handler.
func (h *Handlers) RegisterOnGuildUnavailableEvent(event OnGuildUnavailableFuncType) {
	eventName := discord.DiscordEventGuildUnavailable

	h.ensureEvent(eventName)

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

	h.ensureEvent(eventName)

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

	h.ensureEvent(eventName)

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

	h.ensureEvent(eventName)

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

	h.ensureEvent(eventName)

	h.eventHandlersMu.RLock()
	eh := h.EventHandlers[eventName]
	h.eventHandlersMu.RUnlock()

	eh.eventsMu.Lock()
	eh.Events = append(eh.Events, event)
	eh.eventsMu.Unlock()
}
