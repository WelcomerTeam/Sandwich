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

	EventsMu sync.RWMutex
	Events   []interface{}

	Parser EventParser

	_handlers *Handlers
}

type EventParser func(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error

// Discord Events.

func (h *Handlers) RegisterEvent(eventName string, parser EventParser, event interface{}) *EventHandler {
	h.eventHandlersMu.Lock()
	defer h.eventHandlersMu.Unlock()

	_, ok := h.EventHandlers[eventName]
	if !ok {
		eventHandler := &EventHandler{
			eventName: eventName,
			EventsMu:  sync.RWMutex{},
			Events:    make([]interface{}, 0),
			Parser:    parser,
			_handlers: h,
		}

		h.EventHandlers[eventName] = eventHandler
	}

	eventHandler := h.EventHandlers[eventName]

	if event != nil {
		eventHandler.EventsMu.Lock()
		eventHandler.Events = append(eventHandler.Events, event)
		eventHandler.EventsMu.Unlock()
	}

	return eventHandler
}

// RegisterEventHandler adds a new event handler. If there is already
// an event registered with the name, it is overridden.
func (h *Handlers) RegisterEventHandler(eventName string, parser EventParser) *EventHandler {
	return h.RegisterEvent(eventName, parser, nil)
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

		if identifier != nil {
			eventCtx.Session.Token = "Bot " + identifier.Token
			eventCtx.Identifier = identifier
		}
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
			ev.EventsMu.RLock()
			defer ev.EventsMu.RUnlock()

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
	if err := eventCtx.DecodeContent(payload, &readyPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnReadyFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnReadyFuncType func(eventCtx *EventContext) error

// OnResumed.
func OnResumed(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var resumePayload discord.Resume
	if err := eventCtx.DecodeContent(payload, &resumePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnResumedFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnResumedFuncType func(eventCtx *EventContext) error

// OnApplicationCommandCreate.
func OnApplicationCommandCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var applicationCommandCreatePayload discord.ApplicationCommandCreate
	if err := eventCtx.DecodeContent(payload, &applicationCommandCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if applicationCommandCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*applicationCommandCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.ApplicationCommand(applicationCommandCreatePayload)))
		}
	}

	return nil
}

type OnApplicationCommandCreateFuncType func(eventCtx *EventContext, command discord.ApplicationCommand) error

// OnApplicationCommandUpdate.
func OnApplicationCommandUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var applicationCommandUpdatePayload discord.ApplicationCommandUpdate
	if err := eventCtx.DecodeContent(payload, &applicationCommandUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if applicationCommandUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*applicationCommandUpdatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.ApplicationCommand(applicationCommandUpdatePayload)))
		}
	}

	return nil
}

type OnApplicationCommandUpdateFuncType func(eventCtx *EventContext, command discord.ApplicationCommand) error

// OnApplicationCommandDelete.
func OnApplicationCommandDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var applicationCommandDeletePayload discord.ApplicationCommandDelete
	if err := eventCtx.DecodeContent(payload, &applicationCommandDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if applicationCommandDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*applicationCommandDeletePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.ApplicationCommand(applicationCommandDeletePayload)))
		}
	}

	return nil
}

type OnApplicationCommandDeleteFuncType func(eventCtx *EventContext, command discord.ApplicationCommand) error

// OnChannelCreate.
func OnChannelCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelCreatePayload discord.ChannelCreate
	if err := eventCtx.DecodeContent(payload, &channelCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if channelCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*channelCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Channel(channelCreatePayload)))
		}
	}

	return nil
}

type OnChannelCreateFuncType func(eventCtx *EventContext, channel discord.Channel) error

// OnChannelUpdate.
func OnChannelUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelUpdatePayload discord.ChannelUpdate
	if err := eventCtx.DecodeContent(payload, &channelUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if channelUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*channelUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeChannel); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeChannel, discord.Channel(channelUpdatePayload)))
		}
	}

	return nil
}

type OnChannelUpdateFuncType func(eventCtx *EventContext, before discord.Channel, after discord.Channel) error

// OnChannelDelete.
func OnChannelDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelDeletePayload discord.ChannelDelete
	if err := eventCtx.DecodeContent(payload, &channelDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if channelDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*channelDeletePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Channel(channelDeletePayload)))
		}
	}

	return nil
}

type OnChannelDeleteFuncType func(eventCtx *EventContext, channel discord.Channel) error

// OnChannelPinsUpdate.
func OnChannelPinsUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var channelPinsUpdatePayload discord.ChannelPinsUpdate
	if err := eventCtx.DecodeContent(payload, &channelPinsUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(channelPinsUpdatePayload.GuildID)

	channel := NewChannel(&channelPinsUpdatePayload.GuildID, channelPinsUpdatePayload.ChannelID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelPinsUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, channelPinsUpdatePayload.LastPinTimestamp))
		}
	}

	return nil
}

type OnChannelPinsUpdateFuncType func(eventCtx *EventContext, channel *discord.Channel, lastPinTimestamp time.Time) error

// OnThreadCreate.
func OnThreadCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadCreatePayload discord.ThreadCreate
	if err := eventCtx.DecodeContent(payload, &threadCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Channel(threadCreatePayload)))
		}
	}

	return nil
}

type OnThreadCreateFuncType func(eventCtx *EventContext, thread discord.Channel) error

// OnThreadUpdate.
func OnThreadUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadUpdatePayload discord.ThreadUpdate
	if err := eventCtx.DecodeContent(payload, &threadUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeChannel); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeChannel, discord.Channel(threadUpdatePayload)))
		}
	}

	return nil
}

type OnThreadUpdateFuncType func(eventCtx *EventContext, before discord.Channel, after discord.Channel) error

// OnThreadDelete.
func OnThreadDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadDeletePayload discord.ThreadDelete
	if err := eventCtx.DecodeContent(payload, &threadDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadDeletePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Channel(threadDeletePayload)))
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
// 			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel()))
// 		}
// 	}
//
// 	return nil
// }

// type. OnThreadListSyncFuncType func(eventCtx *Context, thread Channel) error

// OnThreadMemberUpdate.
func OnThreadMemberUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadMemberUpdatePayload discord.ThreadMemberUpdate
	if err := eventCtx.DecodeContent(payload, &threadMemberUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if threadMemberUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*threadMemberUpdatePayload.GuildID)
	}

	channel := NewChannel(threadMemberUpdatePayload.GuildID, *threadMemberUpdatePayload.UserID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadMemberUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, discord.ThreadMember(threadMemberUpdatePayload)))
		}
	}

	return nil
}

type OnThreadMemberUpdateFuncType func(eventCtx *EventContext, thread *discord.Channel, user discord.ThreadMember) error

// OnThreadMembersUpdate.
func OnThreadMembersUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var threadMembersUpdatePayload discord.ThreadMembersUpdate
	if err := eventCtx.DecodeContent(payload, &threadMembersUpdatePayload); err != nil {
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

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadMembersUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, addedUsers, removedUsers))
		}
	}

	return nil
}

type OnThreadMembersUpdateFuncType func(eventCtx *EventContext, thread *discord.Channel, addedUsers []*discord.User, removedUsers []*discord.User) error

// OnGuildCreate.
func OnGuildCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildCreatePayload discord.GuildCreate
	if err := eventCtx.DecodeContent(payload, &guildCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildCreatePayload.ID)

	var lazy bool
	if _, err := eventCtx.DecodeExtra(payload, "lazy", &lazy); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	// If. true, the guild was previously unavailable.
	var unavailable bool
	if _, err := eventCtx.DecodeExtra(payload, "unavailable", &unavailable); err != nil {
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
	if err := eventCtx.DecodeContent(payload, &guildUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	guild := discord.Guild(guildUpdatePayload)
	eventCtx.Guild = &guild

	var beforeGuild discord.Guild
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeGuild); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeGuild, guild))
		}
	}

	return nil
}

type OnGuildUpdateFuncType func(eventCtx *EventContext, before discord.Guild, after discord.Guild) error

// OnGuildDelete.
func OnGuildDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildDeletePayload discord.GuildDelete
	if err := eventCtx.DecodeContent(payload, &guildDeletePayload); err != nil {
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
	if err := eventCtx.DecodeContent(payload, &guildBanAddPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildBanAddPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildBanAddPayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildBanAddFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildBanAddPayload.User))
		}
	}

	return nil
}

type OnGuildBanAddFuncType func(eventCtx *EventContext, user discord.User) error

// OnGuildBanRemove.
func OnGuildBanRemove(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildBanRemovePayload discord.GuildBanRemove
	if err := eventCtx.DecodeContent(payload, &guildBanRemovePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildBanRemovePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildBanRemovePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildBanRemoveFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildBanRemovePayload.User))
		}
	}

	return nil
}

type OnGuildBanRemoveFuncType func(eventCtx *EventContext, user discord.User) error

// OnGuildEmojisUpdate.
func OnGuildEmojisUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildEmojisUpdatePayload discord.GuildEmojisUpdate
	if err := eventCtx.DecodeContent(payload, &guildEmojisUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildEmojisUpdatePayload.GuildID)

	var before []discord.Emoji
	if _, err := eventCtx.DecodeExtra(payload, "before", &before); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	after := make([]discord.Emoji, 0, len(guildEmojisUpdatePayload.Emojis))
	after = append(after, guildEmojisUpdatePayload.Emojis...)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildEmojisUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, before, after))
		}
	}

	return nil
}

type OnGuildEmojisUpdateFuncType func(eventCtx *EventContext, before []discord.Emoji, after []discord.Emoji) error

// OnGuildStickersUpdate.
func OnGuildStickersUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildStickersUpdatePayload discord.GuildStickersUpdate
	if err := eventCtx.DecodeContent(payload, &guildStickersUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildStickersUpdatePayload.GuildID)

	var before []discord.Sticker
	if _, err := eventCtx.DecodeExtra(payload, "before", &before); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	after := make([]discord.Sticker, 0, len(guildStickersUpdatePayload.Stickers))
	after = append(after, guildStickersUpdatePayload.Stickers...)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildStickersUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, before, after))
		}
	}

	return nil
}

type OnGuildStickersUpdateFuncType func(eventCtx *EventContext, before []discord.Sticker, after []discord.Sticker) error

// OnGuildIntegrationsUpdate.
func OnGuildIntegrationsUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildIntegrationsUpdatePayload discord.GuildIntegrationsUpdate
	if err := eventCtx.DecodeContent(payload, &guildIntegrationsUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildIntegrationsUpdatePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildIntegrationsUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnGuildIntegrationsUpdateFuncType func(eventCtx *EventContext) error

// OnGuildMemberAdd.
func OnGuildMemberAdd(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildMemberAddPayload discord.GuildMemberAdd
	if err := eventCtx.DecodeContent(payload, &guildMemberAddPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(*guildMemberAddPayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberAddFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.GuildMember(guildMemberAddPayload)))
		}
	}

	return nil
}

type OnGuildMemberAddFuncType func(eventCtx *EventContext, member discord.GuildMember) error

// OnGuildMemberRemove.
func OnGuildMemberRemove(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildMemberRemovePayload discord.GuildMemberRemove
	if err := eventCtx.DecodeContent(payload, &guildMemberRemovePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildMemberRemovePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberRemoveFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildMemberRemovePayload.User))
		}
	}

	return nil
}

type OnGuildMemberRemoveFuncType func(eventCtx *EventContext, member discord.User) error

// OnGuildMemberUpdate.
func OnGuildMemberUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildMemberUpdatePayload discord.GuildMemberUpdate
	if err := eventCtx.DecodeContent(payload, &guildMemberUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(*guildMemberUpdatePayload.GuildID)

	var beforeGuildMember discord.GuildMember
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeGuildMember); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeGuildMember, discord.GuildMember(guildMemberUpdatePayload)))
		}
	}

	return nil
}

type OnGuildMemberUpdateFuncType func(eventCtx *EventContext, before discord.GuildMember, after discord.GuildMember) error

// OnGuildRoleCreate.
func OnGuildRoleCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildRoleCreatePayload discord.GuildRoleCreate
	if err := eventCtx.DecodeContent(payload, &guildRoleCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if guildRoleCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*guildRoleCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Role(guildRoleCreatePayload)))
		}
	}

	return nil
}

type OnGuildRoleCreateFuncType func(eventCtx *EventContext, role discord.Role) error

// OnGuildRoleUpdate.
func OnGuildRoleUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildRoleUpdatePayload discord.GuildRoleUpdate
	if err := eventCtx.DecodeContent(payload, &guildRoleUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildRoleUpdatePayload.GuildID)

	var beforeRole discord.Role
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeRole); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeRole, guildRoleUpdatePayload.Role))
		}
	}

	return nil
}

type OnGuildRoleUpdateFuncType func(eventCtx *EventContext, before discord.Role, after discord.Role) error

// OnGuildRoleDelete.
func OnGuildRoleDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildRoleDeletePayload discord.GuildRoleDelete
	if err := eventCtx.DecodeContent(payload, &guildRoleDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildRoleDeletePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildRoleDeletePayload.RoleID))
		}
	}

	return nil
}

type OnGuildRoleDeleteFuncType func(eventCtx *EventContext, roleID discord.Snowflake) error

// OnIntegrationCreate.
func OnIntegrationCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var integrationCreatePayload discord.IntegrationCreate
	if err := eventCtx.DecodeContent(payload, &integrationCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if integrationCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*integrationCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Integration(integrationCreatePayload)))
		}
	}

	return nil
}

type OnIntegrationCreateFuncType func(eventCtx *EventContext, integration discord.Integration) error

// OnIntegrationUpdate.
func OnIntegrationUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var integrationUpdatePayload discord.IntegrationUpdate
	if err := eventCtx.DecodeContent(payload, &integrationUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if integrationUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*integrationUpdatePayload.GuildID)
	}

	var beforeIntegration discord.Integration
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeIntegration); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeIntegration, discord.Integration(integrationUpdatePayload)))
		}
	}

	return nil
}

type OnIntegrationUpdateFuncType func(eventCtx *EventContext, before discord.Integration, after discord.Integration) error

// OnIntegrationDelete.
func OnIntegrationDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var integrationDeletePayload discord.IntegrationDelete
	if err := eventCtx.DecodeContent(payload, &integrationDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(integrationDeletePayload.GuildID)

	var applicationID discord.Snowflake
	if !integrationDeletePayload.ApplicationID.IsNil() {
		applicationID = integrationDeletePayload.ApplicationID
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, integrationDeletePayload.ID, applicationID))
		}
	}

	return nil
}

type OnIntegrationDeleteFuncType func(eventCtx *EventContext, integrationID discord.Snowflake, applicationID discord.Snowflake) error

// OnInteractionCreate.
func OnInteractionCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var interactionCreatePayload discord.InteractionCreate
	if err := eventCtx.DecodeContent(payload, &interactionCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if interactionCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*interactionCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInteractionCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Interaction(interactionCreatePayload)))
		}
	}

	return nil
}

type OnInteractionCreateFuncType func(eventCtx *EventContext, interaction discord.Interaction) error

// OnInviteCreate.
func OnInviteCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var inviteCreatePayload discord.InviteCreate
	if err := eventCtx.DecodeContent(payload, &inviteCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if inviteCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*inviteCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInviteCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Invite(inviteCreatePayload)))
		}
	}

	return nil
}

type OnInviteCreateFuncType func(eventCtx *EventContext, invite discord.Invite) error

// OnInviteDelete.
func OnInviteDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var inviteDeletePayload discord.InviteDelete
	if err := eventCtx.DecodeContent(payload, &inviteDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if inviteDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*inviteDeletePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInviteDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Invite(inviteDeletePayload)))
		}
	}

	return nil
}

type OnInviteDeleteFuncType func(eventCtx *EventContext, invite discord.Invite) error

// OnMessageCreate.
func OnMessageCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageCreatePayload discord.MessageCreate
	if err := eventCtx.DecodeContent(payload, &messageCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageCreatePayload.GuildID)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.Message(messageCreatePayload)))
		}
	}

	return nil
}

type OnMessageCreateFuncType func(eventCtx *EventContext, message discord.Message) error

// OnMessageUpdate.
func OnMessageUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageUpdatePayload discord.MessageUpdate
	if err := eventCtx.DecodeContent(payload, &messageUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageUpdatePayload.GuildID)
	}

	var beforeMessage discord.Message
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeMessage); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeMessage, discord.Message(messageUpdatePayload)))
		}
	}

	return nil
}

type OnMessageUpdateFuncType func(eventCtx *EventContext, before discord.Message, after discord.Message) error

// OnMessageDelete.
func OnMessageDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageDeletePayload discord.MessageDelete
	if err := eventCtx.DecodeContent(payload, &messageDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageDeletePayload.GuildID)
	}

	channel := NewChannel(messageDeletePayload.GuildID, messageDeletePayload.ChannelID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageDeletePayload.ID))
		}
	}

	return nil
}

type OnMessageDeleteFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake) error

// OnMessageDeleteBulk.
func OnMessageDeleteBulk(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageDeleteBulkPayload discord.MessageDeleteBulk
	if err := eventCtx.DecodeContent(payload, &messageDeleteBulkPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageDeleteBulkPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageDeleteBulkPayload.GuildID)
	}

	channel := NewChannel(messageDeleteBulkPayload.GuildID, messageDeleteBulkPayload.ChannelID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteBulkFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageDeleteBulkPayload.IDs))
		}
	}

	return nil
}

type OnMessageDeleteBulkFuncType func(eventCtx *EventContext, channel *discord.Channel, messageIDs []discord.Snowflake) error

// OnMessageReactionAdd.
func OnMessageReactionAdd(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionAddPayload discord.MessageReactionAdd
	if err := eventCtx.DecodeContent(payload, &messageReactionAddPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(messageReactionAddPayload.GuildID)

	channel := NewChannel(&messageReactionAddPayload.GuildID, messageReactionAddPayload.ChannelID)

	var guildMember discord.GuildMember
	if messageReactionAddPayload.Member != nil {
		guildMember = *messageReactionAddPayload.Member
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionAddFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionAddPayload.MessageID, messageReactionAddPayload.Emoji, guildMember))
		}
	}

	return nil
}

type OnMessageReactionAddFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake, emoji discord.Emoji, guildMember discord.GuildMember) error

// OnMessageReactionRemove.
func OnMessageReactionRemove(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionRemovePayload discord.MessageReactionRemove
	if err := eventCtx.DecodeContent(payload, &messageReactionRemovePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageReactionRemovePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageReactionRemovePayload.GuildID)
	}

	channel := NewChannel(messageReactionRemovePayload.GuildID, messageReactionRemovePayload.ChannelID)
	user := NewUser(messageReactionRemovePayload.UserID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemovePayload.MessageID, messageReactionRemovePayload.Emoji, user))
		}
	}

	return nil
}

type OnMessageReactionRemoveFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake, emoji discord.Emoji, user *discord.User) error

// OnMessageReactionRemoveAll.
func OnMessageReactionRemoveAll(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionRemoveAllPayload discord.MessageReactionRemoveAll
	if err := eventCtx.DecodeContent(payload, &messageReactionRemoveAllPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(messageReactionRemoveAllPayload.GuildID)

	channel := NewChannel(&messageReactionRemoveAllPayload.GuildID, messageReactionRemoveAllPayload.ChannelID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveAllFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemoveAllPayload.MessageID))
		}
	}

	return nil
}

type OnMessageReactionRemoveAllFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake) error

// OnMessageReactionRemoveEmoji.
func OnMessageReactionRemoveEmoji(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var messageReactionRemoveEmojiPayload discord.MessageReactionRemoveEmoji
	if err := eventCtx.DecodeContent(payload, &messageReactionRemoveEmojiPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if messageReactionRemoveEmojiPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*messageReactionRemoveEmojiPayload.GuildID)
	}

	channel := NewChannel(messageReactionRemoveEmojiPayload.GuildID, messageReactionRemoveEmojiPayload.ChannelID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveEmojiFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemoveEmojiPayload.MessageID, messageReactionRemoveEmojiPayload.Emoji))
		}
	}

	return nil
}

type OnMessageReactionRemoveEmojiFuncType func(eventCtx *EventContext, channel *discord.Channel, messageID discord.Snowflake, emoji discord.Emoji) error

// OnPresenceUpdate.
func OnPresenceUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var presenceUpdatePayload discord.PresenceUpdate
	if err := eventCtx.DecodeContent(payload, &presenceUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(presenceUpdatePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnPresenceUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, presenceUpdatePayload.User, presenceUpdatePayload))
		}
	}

	return nil
}

type OnPresenceUpdateFuncType func(eventCtx *EventContext, user discord.User, payload discord.PresenceUpdate) error

// OnStageInstanceCreate.
func OnStageInstanceCreate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var stageInstanceCreatePayload discord.StageInstanceCreate
	if err := eventCtx.DecodeContent(payload, &stageInstanceCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(stageInstanceCreatePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceCreateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.StageInstance(stageInstanceCreatePayload)))
		}
	}

	return nil
}

type OnStageInstanceCreateFuncType func(eventCtx *EventContext, stage discord.StageInstance) error

// OnStageInstanceUpdate.
func OnStageInstanceUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var stageInstanceUpdatePayload discord.StageInstanceUpdate
	if err := eventCtx.DecodeContent(payload, &stageInstanceUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(stageInstanceUpdatePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.StageInstance(stageInstanceUpdatePayload)))
		}
	}

	return nil
}

type OnStageInstanceUpdateFuncType func(eventCtx *EventContext, stage discord.StageInstance) error

// OnStageInstanceDelete.
func OnStageInstanceDelete(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var stageInstanceDeletePayload discord.StageInstanceDelete
	if err := eventCtx.DecodeContent(payload, &stageInstanceDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceDeleteFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.StageInstance(stageInstanceDeletePayload)))
		}
	}

	return nil
}

type OnStageInstanceDeleteFuncType func(eventCtx *EventContext, stage discord.StageInstance) error

// OnTypingStart.
func OnTypingStart(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var typingStartPayload discord.TypingStart
	if err := eventCtx.DecodeContent(payload, &typingStartPayload); err != nil {
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

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnTypingStartFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, member, user, timestamp))
		}
	}

	return nil
}

type OnTypingStartFuncType func(eventCtx *EventContext, channel *discord.Channel, member *discord.GuildMember, user *discord.User, timestamp time.Time) error

// OnUserUpdate.
func OnUserUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var userUpdatePayload discord.UserUpdate
	if err := eventCtx.DecodeContent(payload, &userUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	var beforeUser discord.User
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeUser); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnUserUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, beforeUser, discord.User(userUpdatePayload)))
		}
	}

	return nil
}

type OnUserUpdateFuncType func(eventCtx *EventContext, before discord.User, after discord.User) error

// OnVoiceStateUpdate.
func OnVoiceStateUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var voiceStateUpdatePayload discord.VoiceStateUpdate
	if err := eventCtx.DecodeContent(payload, &voiceStateUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if voiceStateUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(*voiceStateUpdatePayload.GuildID)
		voiceStateUpdatePayload.Member.GuildID = voiceStateUpdatePayload.GuildID
	}

	var beforeVoiceState discord.VoiceState
	if _, err := eventCtx.DecodeExtra(payload, "before", &beforeVoiceState); err != nil {
		return fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	var guildMember discord.GuildMember
	if voiceStateUpdatePayload.Member != nil {
		guildMember = *voiceStateUpdatePayload.Member
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnVoiceStateUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildMember, beforeVoiceState, discord.VoiceState(voiceStateUpdatePayload)))
		}
	}

	return nil
}

type OnVoiceStateUpdateFuncType func(eventCtx *EventContext, member discord.GuildMember, before discord.VoiceState, after discord.VoiceState) error

// OnVoiceServerUpdate.
func OnVoiceServerUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var voiceServerUpdatePayload discord.VoiceServerUpdate
	if err := eventCtx.DecodeContent(payload, &voiceServerUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(voiceServerUpdatePayload.GuildID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnVoiceServerUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, voiceServerUpdatePayload))
		}
	}

	return nil
}

type OnVoiceServerUpdateFuncType func(eventCtx *EventContext, payload discord.VoiceServerUpdate) error

// OnWebhookUpdate.
func OnWebhookUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var webhookUpdatePayload discord.WebhookUpdate
	if err := eventCtx.DecodeContent(payload, &webhookUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(webhookUpdatePayload.GuildID)

	channel := NewChannel(&webhookUpdatePayload.GuildID, webhookUpdatePayload.ChannelID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnWebhookUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel))
		}
	}

	return nil
}

type OnWebhookUpdateFuncType func(eventCtx *EventContext, channel *discord.Channel) error

// OnGuildJoin.
func OnGuildJoin(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildCreatePayload discord.GuildCreate
	if err := eventCtx.DecodeContent(payload, &guildCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	guild := discord.Guild(guildCreatePayload)
	eventCtx.Guild = &guild

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildJoinFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guild))
		}
	}

	return nil
}

type OnGuildJoinFuncType func(eventCtx *EventContext, guild discord.Guild) error

// OnGuildAvailable.
func OnGuildAvailable(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildCreatePayload discord.GuildCreate
	if err := eventCtx.DecodeContent(payload, &guildCreatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	guild := discord.Guild(guildCreatePayload)
	eventCtx.Guild = &guild

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildJoinFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guild))
		}
	}

	return nil
}

type OnGuildAvailableFuncType func(eventCtx *EventContext, guild discord.Guild) error

// OnGuildLeave.
func OnGuildLeave(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildDeletePayload discord.GuildDelete
	if err := eventCtx.DecodeContent(payload, &guildDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildDeletePayload.ID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildLeaveFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.UnavailableGuild(guildDeletePayload)))
		}
	}

	return nil
}

type OnGuildLeaveFuncType func(eventCtx *EventContext, unavailableGuild discord.UnavailableGuild) error

// OnGuildUnavailable.
func OnGuildUnavailable(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var guildDeletePayload discord.GuildDelete
	if err := eventCtx.DecodeContent(payload, &guildDeletePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.Guild = NewGuild(guildDeletePayload.ID)

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildUnavailableFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, discord.UnavailableGuild(guildDeletePayload)))
		}
	}

	return nil
}

type OnGuildUnavailableFuncType func(eventCtx *EventContext, unavailableGuild discord.UnavailableGuild) error

// Sandwich Events.

// OnSandwichConfigurationReload.
func OnSandwichConfigurationReload(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichConfigurationReloadFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnSandwichConfigurationReloadFuncType func(eventCtx *EventContext) error

// OnSandwichShardStatusUpdate.
func OnSandwichShardStatusUpdate(eventCtx *EventContext, payload sandwich_structs.SandwichPayload) error {
	var shardStatusUpdatePayload sandwich_structs.ShardStatusUpdate
	if err := eventCtx.DecodeContent(payload, &shardStatusUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichShardStatusUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(
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
	if err := eventCtx.DecodeContent(payload, &shardGroupStatusUpdatePayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	eventCtx.EventHandler.EventsMu.RLock()
	defer eventCtx.EventHandler.EventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichShardGroupStatusUpdateFuncType); ok {
			eventCtx.Handlers.WrapFuncType(eventCtx, f(
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
