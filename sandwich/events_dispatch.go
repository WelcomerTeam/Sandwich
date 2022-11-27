package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WelcomerTeam/Discord/discord"
	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
)

type Handlers struct {
	eventHandlersMu sync.RWMutex
	EventHandlers   map[string]*EventHandler
}

// SetupHandler ensures all nullable variables are properly constructed.
func SetupHandler(handler *Handlers) *Handlers {
	if handler == nil {
		handler = &Handlers{
			eventHandlersMu: sync.RWMutex{},
		}
	}

	if handler.EventHandlers == nil {
		handler.EventHandlers = make(map[string]*EventHandler)
	}

	return handler
}

func newDiscordHandlers() *Handlers {
	handler := SetupHandler(nil)

	handler.RegisterEventHandler(discord.DiscordEventReady, OnReady)
	handler.RegisterEventHandler(discord.DiscordEventResumed, OnResumed)
	handler.RegisterEventHandler(discord.DiscordEventApplicationCommandCreate, OnApplicationCommandCreate)
	handler.RegisterEventHandler(discord.DiscordEventApplicationCommandUpdate, OnApplicationCommandUpdate)
	handler.RegisterEventHandler(discord.DiscordEventApplicationCommandDelete, OnApplicationCommandDelete)
	handler.RegisterEventHandler(discord.DiscordEventChannelCreate, OnChannelCreate)
	handler.RegisterEventHandler(discord.DiscordEventChannelUpdate, OnChannelUpdate)
	handler.RegisterEventHandler(discord.DiscordEventChannelDelete, OnChannelDelete)
	handler.RegisterEventHandler(discord.DiscordEventChannelPinsUpdate, OnChannelPinsUpdate)
	handler.RegisterEventHandler(discord.DiscordEventThreadCreate, OnThreadCreate)
	handler.RegisterEventHandler(discord.DiscordEventThreadUpdate, OnThreadUpdate)
	handler.RegisterEventHandler(discord.DiscordEventThreadDelete, OnThreadDelete)
	// h.NewEventHandler(structs.DiscordEventThreadListSync, OnThreadListSync)
	handler.RegisterEventHandler(discord.DiscordEventThreadMemberUpdate, OnThreadMemberUpdate)
	handler.RegisterEventHandler(discord.DiscordEventThreadMembersUpdate, OnThreadMembersUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildCreate, OnGuildCreate)
	handler.RegisterEventHandler(discord.DiscordEventGuildUpdate, OnGuildUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildDelete, OnGuildDelete)
	handler.RegisterEventHandler(discord.DiscordEventGuildBanAdd, OnGuildBanAdd)
	handler.RegisterEventHandler(discord.DiscordEventGuildBanRemove, OnGuildBanRemove)
	handler.RegisterEventHandler(discord.DiscordEventGuildEmojisUpdate, OnGuildEmojisUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildStickersUpdate, OnGuildStickersUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildIntegrationsUpdate, OnGuildIntegrationsUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildMemberAdd, OnGuildMemberAdd)
	handler.RegisterEventHandler(discord.DiscordEventGuildMemberRemove, OnGuildMemberRemove)
	handler.RegisterEventHandler(discord.DiscordEventGuildMemberUpdate, OnGuildMemberUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildRoleCreate, OnGuildRoleCreate)
	handler.RegisterEventHandler(discord.DiscordEventGuildRoleUpdate, OnGuildRoleUpdate)
	handler.RegisterEventHandler(discord.DiscordEventGuildRoleDelete, OnGuildRoleDelete)
	handler.RegisterEventHandler(discord.DiscordEventIntegrationCreate, OnIntegrationCreate)
	handler.RegisterEventHandler(discord.DiscordEventIntegrationUpdate, OnIntegrationUpdate)
	handler.RegisterEventHandler(discord.DiscordEventIntegrationDelete, OnIntegrationDelete)
	handler.RegisterEventHandler(discord.DiscordEventInteractionCreate, OnInteractionCreate)
	handler.RegisterEventHandler(discord.DiscordEventInviteCreate, OnInviteCreate)
	handler.RegisterEventHandler(discord.DiscordEventInviteDelete, OnInviteDelete)
	handler.RegisterEventHandler(discord.DiscordEventMessageCreate, OnMessageCreate)
	handler.RegisterEventHandler(discord.DiscordEventMessageUpdate, OnMessageUpdate)
	handler.RegisterEventHandler(discord.DiscordEventMessageDelete, OnMessageDelete)
	handler.RegisterEventHandler(discord.DiscordEventMessageDeleteBulk, OnMessageDeleteBulk)
	handler.RegisterEventHandler(discord.DiscordEventMessageReactionAdd, OnMessageReactionAdd)
	handler.RegisterEventHandler(discord.DiscordEventMessageReactionRemove, OnMessageReactionRemove)
	handler.RegisterEventHandler(discord.DiscordEventMessageReactionRemoveAll, OnMessageReactionRemoveAll)
	handler.RegisterEventHandler(discord.DiscordEventMessageReactionRemoveEmoji, OnMessageReactionRemoveEmoji)
	handler.RegisterEventHandler(discord.DiscordEventPresenceUpdate, OnPresenceUpdate)
	handler.RegisterEventHandler(discord.DiscordEventStageInstanceCreate, OnStageInstanceCreate)
	handler.RegisterEventHandler(discord.DiscordEventStageInstanceUpdate, OnStageInstanceUpdate)
	handler.RegisterEventHandler(discord.DiscordEventStageInstanceDelete, OnStageInstanceDelete)
	handler.RegisterEventHandler(discord.DiscordEventTypingStart, OnTypingStart)
	handler.RegisterEventHandler(discord.DiscordEventUserUpdate, OnUserUpdate)
	handler.RegisterEventHandler(discord.DiscordEventVoiceStateUpdate, OnVoiceStateUpdate)
	handler.RegisterEventHandler(discord.DiscordEventVoiceServerUpdate, OnVoiceServerUpdate)
	handler.RegisterEventHandler(discord.DiscordEventWebhookUpdate, OnWebhookUpdate)

	// Custom Events.
	handler.RegisterEventHandler(discord.DiscordEventGuildJoin, OnGuildJoin)
	handler.RegisterEventHandler(discord.DiscordEventGuildAvailable, OnGuildAvailable)
	handler.RegisterEventHandler(discord.DiscordEventGuildLeave, OnGuildLeave)
	handler.RegisterEventHandler(discord.DiscordEventGuildUnavailable, OnGuildUnavailable)

	handler.RegisterEventHandler(DiscordEventError, nil)

	return handler
}

func newSandwichHandlers() *Handlers {
	handler := SetupHandler(nil)

	handler.RegisterEventHandler(sandwich_structs.SandwichEventConfigurationReload, OnSandwichConfigurationReload)
	handler.RegisterEventHandler(sandwich_structs.SandwichEventShardStatusUpdate, OnSandwichShardStatusUpdate)
	handler.RegisterEventHandler(sandwich_structs.SandwichEventShardGroupStatusUpdate, OnSandwichShardGroupStatusUpdate)

	// Register events that are handled by default.
	handler.RegisterOnSandwichConfigurationReload(func(eventCtx *EventContext) error {
		identifiers, err := eventCtx.Sandwich.GRPCInterface.FetchConsumerConfiguration(eventCtx.ToGRPCContext(), "")
		if err != nil {
			return fmt.Errorf("failed to fetch consumer configuration: %w", err)
		}

		eventCtx.Sandwich.identifiersMu.Lock()
		eventCtx.Sandwich.Identifiers = map[string]*sandwich_structs.ManagerConsumerConfiguration{}

		for k := range identifiers.Identifiers {
			v := identifiers.Identifiers[k]
			eventCtx.Sandwich.Identifiers[k] = &v
		}
		eventCtx.Sandwich.identifiersMu.Unlock()

		return nil
	})

	return handler
}

type EventHandler struct {
	eventName string

	eventsMu sync.RWMutex
	Events   []interface{}

	Parser EventParser

	_handlers *Handlers
}

type EventParser func(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error

// Discord Events.

// RegisterEventHandler adds a new event handler. If there is already
// an event registered with the name, it is overridden.
func (h *Handlers) RegisterEventHandler(eventName string, parser EventParser) *EventHandler {
	h.eventHandlersMu.Lock()
	defer h.eventHandlersMu.Unlock()

	eventHandler := &EventHandler{
		eventName: eventName,
		eventsMu:  sync.RWMutex{},
		Events:    make([]interface{}, 0),
		Parser:    parser,
		_handlers: h,
	}

	h.EventHandlers[eventName] = eventHandler

	return eventHandler
}

// Dispatch dispatches a payload. All dispatched events will be sent through a goroutine, so
// no errors are returned.
func (h *Handlers) Dispatch(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) {
	go h.DispatchType(eventCtx, payload.Type, payload)
}

// DispatchType is similar to Dispatch however a custom event name
// can. be passed, preserving the original payload.
func (h *Handlers) DispatchType(eventCtx *EventContext, eventName string, payload sandwich_structs.SandwichPayload) error {
	if payload.Metadata.Application != "" {
		identifier, ok, err := eventCtx.Sandwich.FetchIdentifier(context.TODO(), payload.Metadata.Application)
		if !ok || err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to fetch identifier for application")

			return err
		}

		eventCtx.Session.Token = "Bot " + identifier.Token
		eventCtx.Identifier = identifier
	}

	if eventHandler, ok := h.EventHandlers[eventName]; ok {
		eventCtx.EventHandler = eventHandler

		defer func() {
			errorValue := recover()
			if errorValue != nil {
				eventCtx.Sandwich.RecoverEventPanic(errorValue, eventCtx, &payload)
			}
		}()

		return eventHandler.Parser(eventCtx, payload)
	}

	eventCtx.Logger.Info().Str("type", payload.Type).Msg("Unknown event handler")

	return ErrUnknownEvent
}

// WrapFuncType handles the error of a FuncType if it returns an error.
// It will call any ERROR handlers. Errors that occur in the ERROR handler
// will not trigger the ERROR handler.
func (h *Handlers) WrapFuncType(eventCtx *EventContext, funcTypeErr error) error {
	if funcTypeErr != nil {
		if ev, ok := h.EventHandlers["ERROR"]; ok {
			ev.eventsMu.RLock()
			defer ev.eventsMu.RUnlock()

			for _, event := range ev.Events {
				if f, ok := event.(OnErrorFuncType); ok {
					_ = f(eventCtx, funcTypeErr)
				}
			}
		}
	}

	return nil
}

// OnReady.
func OnReady(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var readyPayload discord.Ready
	if err := eventCtx.decodeContent(payload, &readyPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnReadyFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnReadyFuncType func(eventCtx *EventContext) error

// OnResumed.
func OnResumed(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var resumePayload discord.Resume
	if err := eventCtx.decodeContent(payload, &resumePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnResumedFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnResumedFuncType func(eventCtx *EventContext) error

// OnApplicationCommandCreate.
func OnApplicationCommandCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var applicationCommandCreatePayload discord.ApplicationCommandCreate
	if err := eventCtx.decodeContent(payload, &applicationCommandCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if applicationCommandCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*applicationCommandCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *applicationCommandCreatePayload))
		}
	}

	return nil
}

type OnApplicationCommandCreateFuncType func(eventCtx *EventContext, command discord.ApplicationCommand) error

// OnApplicationCommandUpdate.
func OnApplicationCommandUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var applicationCommandUpdatePayload discord.ApplicationCommandUpdate
	if err := eventCtx.decodeContent(payload, &applicationCommandUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if applicationCommandUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*applicationCommandUpdatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *applicationCommandUpdatePayload))
		}
	}

	return nil
}

type OnApplicationCommandUpdateFuncType func(eventCtx *EventContext, command discord.ApplicationCommand) error

// OnApplicationCommandDelete.
func OnApplicationCommandDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var applicationCommandDeletePayload discord.ApplicationCommandDelete
	if err := eventCtx.decodeContent(payload, &applicationCommandDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if applicationCommandDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*applicationCommandDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *applicationCommandDeletePayload))
		}
	}

	return nil
}

type OnApplicationCommandDeleteFuncType func(eventCtx *EventContext, command discord.ApplicationCommand) error

// OnChannelCreate.
func OnChannelCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelCreatePayload discord.ChannelCreate
	if err := eventCtx.decodeContent(payload, &channelCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if channelCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*channelCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *channelCreatePayload))
		}
	}

	return nil
}

type OnChannelCreateFuncType func(eventCtx *EventContext, channel discord.Channel) error

// OnChannelUpdate.
func OnChannelUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelUpdatePayload discord.ChannelUpdate
	if err := eventCtx.decodeContent(payload, &channelUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if channelUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*channelUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeChannel); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeChannel, *channelUpdatePayload))
		}
	}

	return nil
}

type OnChannelUpdateFuncType func(eventCtx *EventContext, before discord.Channel, after discord.Channel) error

// OnChannelDelete.
func OnChannelDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelDeletePayload discord.ChannelDelete
	if err := eventCtx.decodeContent(payload, &channelDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if channelDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*channelDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *channelDeletePayload))
		}
	}

	return nil
}

type OnChannelDeleteFuncType func(eventCtx *EventContext, channel discord.Channel) error

// OnChannelPinsUpdate.
func OnChannelPinsUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelPinsUpdatePayload discord.ChannelPinsUpdate
	if err := eventCtx.decodeContent(payload, &channelPinsUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(channelPinsUpdatePayload.GuildID)

	channel := NewChannel(&channelPinsUpdatePayload.GuildID, channelPinsUpdatePayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelPinsUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, channelPinsUpdatePayload.LastPinTimestamp))
		}
	}

	return nil
}

type OnChannelPinsUpdateFuncType func(eventCtx *EventContext, channel *discord.Channel, lastPinTimestamp time.Time) error

// OnThreadCreate.
func OnThreadCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadCreatePayload discord.ThreadCreate
	if err := eventCtx.decodeContent(payload, &threadCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *threadCreatePayload))
		}
	}

	return nil
}

type OnThreadCreateFuncType func(eventCtx *EventContext, thread discord.Channel) error

// OnThreadUpdate.
func OnThreadUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadUpdatePayload discord.ThreadUpdate
	if err := eventCtx.decodeContent(payload, &threadUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeChannel); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeChannel, *threadUpdatePayload))
		}
	}

	return nil
}

type OnThreadUpdateFuncType func(eventCtx *EventContext, before discord.Channel, after discord.Channel) error

// OnThreadDelete.
func OnThreadDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadDeletePayload discord.ThreadDelete
	if err := eventCtx.decodeContent(payload, &threadDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *threadDeletePayload))
		}
	}

	return nil
}

type OnThreadDeleteFuncType func(eventCtx *EventContext, thread discord.Channel) error

// // OnThreadListSync.
// func. OnThreadListSync(eventCtx *Context, payload structs.SandwichPayload) error {
// 	var threadListSyncPayload discord.ThreadListSync
// 	if err := eventCtx.decodeContent(payload, &threadListSyncPayload); err != nil {
// 		return fmt.Errorf("failed to unmarshal payload: %w", err)
// 	}

// 	for _, event := range eventCtx.EventHandler.Events {
// 		if f, ok := event.(OnThreadListSyncFuncType); ok {
// 			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel()))
// 		}
// 	}
//
// 	return nil
// }

// type. OnThreadListSyncFuncType func(eventCtx *Context, thread Channel) error

// OnThreadMemberUpdate.
func OnThreadMemberUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadMemberUpdatePayload discord.ThreadMemberUpdate
	if err := eventCtx.decodeContent(payload, &threadMemberUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadMemberUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadMemberUpdatePayload.GuildID)
	}

	channel := NewChannel(threadMemberUpdatePayload.GuildID, *threadMemberUpdatePayload.UserID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadMemberUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, *threadMemberUpdatePayload))
		}
	}

	return nil
}

type OnThreadMemberUpdateFuncType func(eventCtx *EventContext, thread *discord.Channel, user discord.ThreadMember) error

// OnThreadMembersUpdate.
func OnThreadMembersUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadMembersUpdatePayload discord.ThreadMembersUpdate
	if err := eventCtx.decodeContent(payload, &threadMembersUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(threadMembersUpdatePayload.GuildID)

	channel := NewChannel(&threadMembersUpdatePayload.GuildID, threadMembersUpdatePayload.ID)

	addedUsers := make([]*discord.User, 0, len(threadMembersUpdatePayload.AddedMembers))
	for _, addedMember := range threadMembersUpdatePayload.AddedMembers {
		addedUsers = append(addedUsers, NewUser(*addedMember.UserID))
	}

	removedUsers := make([]*discord.User, 0, len(threadMembersUpdatePayload.RemovedMemberIDs))
	for _, removedUser := range threadMembersUpdatePayload.RemovedMemberIDs {
		removedUsers = append(removedUsers, NewUser(removedUser))
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadMembersUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, addedUsers, removedUsers))
		}
	}

	return nil
}

type OnThreadMembersUpdateFuncType func(eventCtx *EventContext, thread *discord.Channel, addedUsers []*discord.User, removedUsers []*discord.User) error

// OnGuildCreate.
func OnGuildCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildCreatePayload discord.GuildCreate
	if err := eventCtx.decodeContent(payload, &guildCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildCreatePayload.ID)

	var lazy bool
	if _, err := eventCtx.decodeExtra(payload, "lazy", &lazy); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	// If. true, the guild was previously unavailable.
	var unavailable bool
	if _, err := eventCtx.decodeExtra(payload, "unavailable", &unavailable); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	// Dispatches. either a GUILD_JOIN or GUILD_AVAILABLE event.
	// Guilds. that are lazy loaded are not handled.

	if unavailable {
		return eventCtx.Handlers.DispatchType(eventCtx, "GUILD_AVAILABLE", payload)
	} else if !lazy {
		return eventCtx.Handlers.DispatchType(eventCtx, "GUILD_JOIN", payload)
	}

	return nil
}

// OnGuildUpdate.
func OnGuildUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildUpdatePayload discord.GuildUpdate
	if err := eventCtx.decodeContent(payload, &guildUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	guild := *guildUpdatePayload
	eventCtx.Guild = &guild

	var beforeGuild discord.Guild
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeGuild); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeGuild, guild))
		}
	}

	return nil
}

type OnGuildUpdateFuncType func(eventCtx *EventContext, before discord.Guild, after discord.Guild) error

// OnGuildDelete.
func OnGuildDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildDeletePayload discord.GuildDelete
	if err := eventCtx.decodeContent(payload, &guildDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildDeletePayload.ID)

	if guildDeletePayload.Unavailable {
		return eventCtx.Handlers.DispatchType(eventCtx, "GUILD_UNAVAILABLE", payload)
	}

	return eventCtx.Handlers.DispatchType(eventCtx, "GUILD_REMOVE", payload)
}

// OnGuildBanAdd.
func OnGuildBanAdd(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildBanAddPayload discord.GuildBanAdd
	if err := eventCtx.decodeContent(payload, &guildBanAddPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildBanAddPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildBanAddPayload.GuildID)
	}

	user := *guildBanAddPayload.User

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildBanAddFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, user))
		}
	}

	return nil
}

type OnGuildBanAddFuncType func(eventCtx *EventContext, user discord.User) error

// OnGuildBanRemove.
func OnGuildBanRemove(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildBanRemovePayload discord.GuildBanRemove
	if err := eventCtx.decodeContent(payload, &guildBanRemovePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildBanRemovePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildBanRemovePayload.GuildID)
	}

	user := *guildBanRemovePayload.User

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildBanRemoveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, user))
		}
	}

	return nil
}

type OnGuildBanRemoveFuncType func(eventCtx *EventContext, user discord.User) error

// OnGuildEmojisUpdate.
func OnGuildEmojisUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildEmojisUpdatePayload discord.GuildEmojisUpdate
	if err := eventCtx.decodeContent(payload, &guildEmojisUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildEmojisUpdatePayload.GuildID)

	var before []discord.Emoji
	if _, err := eventCtx.decodeExtra(payload, "before", &before); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	after := make([]discord.Emoji, 0, len(guildEmojisUpdatePayload.Emojis))
	for _, emoji := range guildEmojisUpdatePayload.Emojis {
		after = append(after, *emoji)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildEmojisUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, before, after))
		}
	}

	return nil
}

type OnGuildEmojisUpdateFuncType func(eventCtx *EventContext, before []discord.Emoji, after []discord.Emoji) error

// OnGuildStickersUpdate.
func OnGuildStickersUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildStickersUpdatePayload discord.GuildStickersUpdate
	if err := eventCtx.decodeContent(payload, &guildStickersUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildStickersUpdatePayload.GuildID)

	var before []discord.Sticker
	if _, err := eventCtx.decodeExtra(payload, "before", &before); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	after := make([]discord.Sticker, 0, len(guildStickersUpdatePayload.Stickers))
	for _, sticker := range guildStickersUpdatePayload.Stickers {
		after = append(after, *sticker)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildStickersUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, before, after))
		}
	}

	return nil
}

type OnGuildStickersUpdateFuncType func(eventCtx *EventContext, before []discord.Sticker, after []discord.Sticker) error

// OnGuildIntegrationsUpdate.
func OnGuildIntegrationsUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildIntegrationsUpdatePayload discord.GuildIntegrationsUpdate
	if err := eventCtx.decodeContent(payload, &guildIntegrationsUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildIntegrationsUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildIntegrationsUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnGuildIntegrationsUpdateFuncType func(eventCtx *EventContext) error

// OnGuildMemberAdd.
func OnGuildMemberAdd(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildMemberAddPayload discord.GuildMemberAdd
	if err := eventCtx.decodeContent(payload, &guildMemberAddPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(*guildMemberAddPayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberAddFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *guildMemberAddPayload))
		}
	}

	return nil
}

type OnGuildMemberAddFuncType func(eventCtx *EventContext, member discord.GuildMember) error

// OnGuildMemberRemove.
func OnGuildMemberRemove(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildMemberRemovePayload discord.GuildMemberRemove
	if err := eventCtx.decodeContent(payload, &guildMemberRemovePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildMemberRemovePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberRemoveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *guildMemberRemovePayload.User))
		}
	}

	return nil
}

type OnGuildMemberRemoveFuncType func(eventCtx *EventContext, member discord.User) error

// OnGuildMemberUpdate.
func OnGuildMemberUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildMemberUpdatePayload discord.GuildMemberUpdate
	if err := eventCtx.decodeContent(payload, &guildMemberUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(*guildMemberUpdatePayload.GuildID)

	var beforeGuildMember discord.GuildMember
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeGuildMember); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeGuildMember, *guildMemberUpdatePayload.GuildMember))
		}
	}

	return nil
}

type OnGuildMemberUpdateFuncType func(eventCtx *EventContext, before discord.GuildMember, after discord.GuildMember) error

// OnGuildRoleCreate.
func OnGuildRoleCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildRoleCreatePayload discord.GuildRoleCreate
	if err := eventCtx.decodeContent(payload, &guildRoleCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildRoleCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildRoleCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *guildRoleCreatePayload))
		}
	}

	return nil
}

type OnGuildRoleCreateFuncType func(eventCtx *EventContext, role discord.Role) error

// OnGuildRoleUpdate.
func OnGuildRoleUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildRoleUpdatePayload discord.GuildRoleUpdate
	if err := eventCtx.decodeContent(payload, &guildRoleUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildRoleUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildRoleUpdatePayload.GuildID)
	}

	var beforeRole discord.Role
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeRole); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeRole, *guildRoleUpdatePayload))
		}
	}

	return nil
}

type OnGuildRoleUpdateFuncType func(eventCtx *EventContext, before discord.Role, after discord.Role) error

// OnGuildRoleDelete.
func OnGuildRoleDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildRoleDeletePayload discord.GuildRoleDelete
	if err := eventCtx.decodeContent(payload, &guildRoleDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildRoleDeletePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildRoleDeletePayload.RoleID))
		}
	}

	return nil
}

type OnGuildRoleDeleteFuncType func(eventCtx *EventContext, roleID discord.Snowflake) error

// OnIntegrationCreate.
func OnIntegrationCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var integrationCreatePayload discord.IntegrationCreate
	if err := eventCtx.decodeContent(payload, &integrationCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if integrationCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*integrationCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *integrationCreatePayload))
		}
	}

	return nil
}

type OnIntegrationCreateFuncType func(eventCtx *EventContext, integration discord.Integration) error

// OnIntegrationUpdate.
func OnIntegrationUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var integrationUpdatePayload discord.IntegrationUpdate
	if err := eventCtx.decodeContent(payload, &integrationUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if integrationUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*integrationUpdatePayload.GuildID)
	}

	var beforeIntegration discord.Integration
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeIntegration); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeIntegration, *integrationUpdatePayload))
		}
	}

	return nil
}

type OnIntegrationUpdateFuncType func(eventCtx *EventContext, before discord.Integration, after discord.Integration) error

// OnIntegrationDelete.
func OnIntegrationDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var integrationDeletePayload discord.IntegrationDelete
	if err := eventCtx.decodeContent(payload, &integrationDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(integrationDeletePayload.GuildID)

	var applicationID discord.Snowflake
	if integrationDeletePayload.ApplicationID != nil {
		applicationID = *integrationDeletePayload.ApplicationID
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, integrationDeletePayload.ID, applicationID))
		}
	}

	return nil
}

type OnIntegrationDeleteFuncType func(eventCtx *EventContext, integrationID discord.Snowflake, applicationID discord.Snowflake) error

// OnInteractionCreate.
func OnInteractionCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var interactionCreatePayload discord.InteractionCreate
	if err := eventCtx.decodeContent(payload, &interactionCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if interactionCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*interactionCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInteractionCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *interactionCreatePayload))
		}
	}

	return nil
}

type OnInteractionCreateFuncType func(eventCtx *EventContext, interaction discord.Interaction) error

// OnInviteCreate.
func OnInviteCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var inviteCreatePayload discord.InviteCreate
	if err := eventCtx.decodeContent(payload, &inviteCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if inviteCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*inviteCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInviteCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *inviteCreatePayload))
		}
	}

	return nil
}

type OnInviteCreateFuncType func(eventCtx *EventContext, invite discord.Invite) error

// OnInviteDelete.
func OnInviteDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var inviteDeletePayload discord.InviteDelete
	if err := eventCtx.decodeContent(payload, &inviteDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if inviteDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*inviteDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInviteDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *inviteDeletePayload))
		}
	}

	return nil
}

type OnInviteDeleteFuncType func(eventCtx *EventContext, invite discord.Invite) error

// OnMessageCreate.
func OnMessageCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageCreatePayload discord.MessageCreate
	if err := eventCtx.decodeContent(payload, &messageCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *messageCreatePayload))
		}
	}

	return nil
}

type OnMessageCreateFuncType func(eventCtx *EventContext, message discord.Message) error

// OnMessageUpdate.
func OnMessageUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageUpdatePayload discord.MessageUpdate
	if err := eventCtx.decodeContent(payload, &messageUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageUpdatePayload.GuildID)
	}

	var beforeMessage discord.Message
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeMessage); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeMessage, *messageUpdatePayload))
		}
	}

	return nil
}

type OnMessageUpdateFuncType func(eventCtx *EventContext, before discord.Message, after discord.Message) error

// OnMessageDelete.
func OnMessageDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageDeletePayload discord.MessageDelete
	if err := eventCtx.decodeContent(payload, &messageDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageDeletePayload.GuildID)
	}

	channel := NewChannel(messageDeletePayload.GuildID, messageDeletePayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageDeletePayload.ID))
		}
	}

	return nil
}

type OnMessageDeleteFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake) error

// OnMessageDeleteBulk.
func OnMessageDeleteBulk(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageDeleteBulkPayload discord.MessageDeleteBulk
	if err := eventCtx.decodeContent(payload, &messageDeleteBulkPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageDeleteBulkPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageDeleteBulkPayload.GuildID)
	}

	channel := NewChannel(messageDeleteBulkPayload.GuildID, messageDeleteBulkPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteBulkFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageDeleteBulkPayload.IDs))
		}
	}

	return nil
}

type OnMessageDeleteBulkFuncType func(eventCtx *EventContext, channel *discord.Channel, messageIDs []discord.Snowflake) error

// OnMessageReactionAdd.
func OnMessageReactionAdd(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionAddPayload discord.MessageReactionAdd
	if err := eventCtx.decodeContent(payload, &messageReactionAddPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(messageReactionAddPayload.GuildID)

	channel := NewChannel(&messageReactionAddPayload.GuildID, messageReactionAddPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionAddFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionAddPayload.MessageID, *messageReactionAddPayload.Emoji, *messageReactionAddPayload.Member))
		}
	}

	return nil
}

type OnMessageReactionAddFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake, emoji discord.Emoji, guildMember discord.GuildMember) error

// OnMessageReactionRemove.
func OnMessageReactionRemove(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionRemovePayload discord.MessageReactionRemove
	if err := eventCtx.decodeContent(payload, &messageReactionRemovePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageReactionRemovePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageReactionRemovePayload.GuildID)
	}

	channel := NewChannel(messageReactionRemovePayload.GuildID, messageReactionRemovePayload.ChannelID)
	user := NewUser(messageReactionRemovePayload.UserID)
	emoji := *messageReactionRemovePayload.Emoji

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemovePayload.MessageID, emoji, user))
		}
	}

	return nil
}

type OnMessageReactionRemoveFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake, emoji discord.Emoji, user *discord.User) error

// OnMessageReactionRemoveAll.
func OnMessageReactionRemoveAll(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionRemoveAllPayload discord.MessageReactionRemoveAll
	if err := eventCtx.decodeContent(payload, &messageReactionRemoveAllPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(messageReactionRemoveAllPayload.GuildID)

	channel := NewChannel(&messageReactionRemoveAllPayload.GuildID, messageReactionRemoveAllPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveAllFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemoveAllPayload.MessageID))
		}
	}

	return nil
}

type OnMessageReactionRemoveAllFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake) error

// OnMessageReactionRemoveEmoji.
func OnMessageReactionRemoveEmoji(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionRemoveEmojiPayload discord.MessageReactionRemoveEmoji
	if err := eventCtx.decodeContent(payload, &messageReactionRemoveEmojiPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageReactionRemoveEmojiPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageReactionRemoveEmojiPayload.GuildID)
	}

	channel := NewChannel(messageReactionRemoveEmojiPayload.GuildID, messageReactionRemoveEmojiPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveEmojiFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemoveEmojiPayload.MessageID, *messageReactionRemoveEmojiPayload.Emoji))
		}
	}

	return nil
}

type OnMessageReactionRemoveEmojiFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake, emoji discord.Emoji) error

// OnPresenceUpdate.
func OnPresenceUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var presenceUpdatePayload discord.PresenceUpdate
	if err := eventCtx.decodeContent(payload, &presenceUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(presenceUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnPresenceUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *presenceUpdatePayload.User, presenceUpdatePayload))
		}
	}

	return nil
}

type OnPresenceUpdateFuncType func(eventCtx *EventContext, user discord.User, payload discord.PresenceUpdate) error

// OnStageInstanceCreate.
func OnStageInstanceCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var stageInstanceCreatePayload discord.StageInstanceCreate
	if err := eventCtx.decodeContent(payload, &stageInstanceCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(stageInstanceCreatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *stageInstanceCreatePayload))
		}
	}

	return nil
}

type OnStageInstanceCreateFuncType func(eventCtx *EventContext, stage discord.StageInstance) error

// OnStageInstanceUpdate.
func OnStageInstanceUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var stageInstanceUpdatePayload discord.StageInstanceUpdate
	if err := eventCtx.decodeContent(payload, &stageInstanceUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(stageInstanceUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *stageInstanceUpdatePayload))
		}
	}

	return nil
}

type OnStageInstanceUpdateFuncType func(eventCtx *EventContext, stage discord.StageInstance) error

// OnStageInstanceDelete.
func OnStageInstanceDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var stageInstanceDeletePayload discord.StageInstanceDelete
	if err := eventCtx.decodeContent(payload, &stageInstanceDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *stageInstanceDeletePayload))
		}
	}

	return nil
}

type OnStageInstanceDeleteFuncType func(eventCtx *EventContext, stage discord.StageInstance) error

// OnTypingStart.
func OnTypingStart(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var typingStartPayload discord.TypingStart
	if err := eventCtx.decodeContent(payload, &typingStartPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if typingStartPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*typingStartPayload.GuildID)
	}

	timestamp := time.Unix(int64(typingStartPayload.Timestamp), 0)
	channel := NewChannel(typingStartPayload.GuildID, typingStartPayload.ChannelID)

	var user *discord.User

	var member *discord.GuildMember

	if typingStartPayload.Member != nil {
		member = typingStartPayload.Member
		user = typingStartPayload.Member.User
	} else {
		user = NewUser(typingStartPayload.UserID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnTypingStartFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, member, user, timestamp))
		}
	}

	return nil
}

type OnTypingStartFuncType func(eventCtx *EventContext, channel *discord.Channel, member *discord.GuildMember, user *discord.User, timestamp time.Time) error

// OnUserUpdate.
func OnUserUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var userUpdatePayload discord.UserUpdate
	if err := eventCtx.decodeContent(payload, &userUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	var beforeUser discord.User
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeUser); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnUserUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeUser, *userUpdatePayload))
		}
	}

	return nil
}

type OnUserUpdateFuncType func(eventCtx *EventContext, before discord.User, after discord.User) error

// OnVoiceStateUpdate.
func OnVoiceStateUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var voiceStateUpdatePayload discord.VoiceStateUpdate
	if err := eventCtx.decodeContent(payload, &voiceStateUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if voiceStateUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*voiceStateUpdatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnVoiceStateUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *voiceStateUpdatePayload.Member, *voiceStateUpdatePayload))
		}
	}

	return nil
}

type OnVoiceStateUpdateFuncType func(eventCtx *EventContext, member discord.GuildMember, voice discord.VoiceState) error

// OnVoiceServerUpdate.
func OnVoiceServerUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var voiceServerUpdatePayload discord.VoiceServerUpdate
	if err := eventCtx.decodeContent(payload, &voiceServerUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(voiceServerUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnVoiceServerUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, voiceServerUpdatePayload))
		}
	}

	return nil
}

type OnVoiceServerUpdateFuncType func(eventCtx *EventContext, payload discord.VoiceServerUpdate) error

// OnWebhookUpdate.
func OnWebhookUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var webhookUpdatePayload discord.WebhookUpdate
	if err := eventCtx.decodeContent(payload, &webhookUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(webhookUpdatePayload.GuildID)

	channel := NewChannel(&webhookUpdatePayload.GuildID, webhookUpdatePayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnWebhookUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel))
		}
	}

	return nil
}

type OnWebhookUpdateFuncType func(eventCtx *EventContext, channel *discord.Channel) error

// OnGuildJoin.
func OnGuildJoin(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildCreatePayload discord.GuildCreate
	if err := eventCtx.decodeContent(payload, &guildCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	guild := *guildCreatePayload
	eventCtx.Guild = &guild

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildJoinFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guild))
		}
	}

	return nil
}

type OnGuildJoinFuncType func(eventCtx *EventContext, guild discord.Guild) error

// OnGuildAvailable.
func OnGuildAvailable(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildCreatePayload discord.GuildCreate
	if err := eventCtx.decodeContent(payload, &guildCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	guild := *guildCreatePayload
	eventCtx.Guild = &guild

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildJoinFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guild))
		}
	}

	return nil
}

type OnGuildAvailableFuncType func(eventCtx *EventContext, guild discord.Guild) error

// OnGuildLeave.
func OnGuildLeave(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildDeletePayload discord.GuildDelete
	if err := eventCtx.decodeContent(payload, &guildDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildDeletePayload.ID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildLeaveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *guildDeletePayload))
		}
	}

	return nil
}

type OnGuildLeaveFuncType func(eventCtx *EventContext, unavailableGuild discord.UnavailableGuild) error

// OnGuildUnavailable.
func OnGuildUnavailable(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildDeletePayload discord.GuildDelete
	if err := eventCtx.decodeContent(payload, &guildDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildDeletePayload.ID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildUnavailableFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, *guildDeletePayload))
		}
	}

	return nil
}

type OnGuildUnavailableFuncType func(eventCtx *EventContext, unavailableGuild discord.UnavailableGuild) error

// Sandwich Events.

// OnSandwichConfigurationReload.
func OnSandwichConfigurationReload(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichConfigurationReloadFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnSandwichConfigurationReloadFuncType func(eventCtx *EventContext) error

// OnSandwichShardStatusUpdate.
func OnSandwichShardStatusUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var shardStatusUpdatePayload sandwich_structs.ShardStatusUpdate
	if err := eventCtx.decodeContent(payload, &shardStatusUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichShardStatusUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(
				eventCtx,
				shardStatusUpdatePayload.Manager,
				shardStatusUpdatePayload.ShardGroup,
				shardStatusUpdatePayload.Shard,
				shardStatusUpdatePayload.Status))
		}
	}

	return nil
}

type OnSandwichShardStatusUpdateFuncType func(eventCtx *EventContext, manager string, shardGroup int32, shard int32, status sandwich_structs.ShardStatus) error

// OnSandwichShardGroupStatusUpdate.
func OnSandwichShardGroupStatusUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var shardGroupStatusUpdatePayload sandwich_structs.ShardGroupStatusUpdate
	if err := eventCtx.decodeContent(payload, &shardGroupStatusUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichShardGroupStatusUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(
				eventCtx,
				shardGroupStatusUpdatePayload.Manager,
				shardGroupStatusUpdatePayload.ShardGroup,
				shardGroupStatusUpdatePayload.Status))
		}
	}

	return nil
}

type OnSandwichShardGroupStatusUpdateFuncType func(eventCtx *EventContext, manager string, shardGroup int32, status sandwich_structs.ShardGroupStatus) error

// Generic Events.

type OnErrorFuncType func(eventCtx *EventContext, eventErr error) error
