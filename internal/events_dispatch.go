package internal

import (
	"context"
	"sync"
	"time"

	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
	"github.com/WelcomerTeam/Sandwich-Daemon/structs"
	"golang.org/x/xerrors"
)

type Handlers struct {
	eventHandlersMu sync.RWMutex
	EventHandlers   map[string]*EventHandler
}

func NewDiscordHandlers() (h *Handlers) {
	h = &Handlers{
		eventHandlersMu: sync.RWMutex{},
		EventHandlers:   make(map[string]*EventHandler),
	}

	h.RegisterEventHandler(structs.DiscordEventReady, OnReady)
	h.RegisterEventHandler(structs.DiscordEventResumed, OnResumed)
	h.RegisterEventHandler(structs.DiscordEventApplicationCommandCreate, OnApplicationCommandCreate)
	h.RegisterEventHandler(structs.DiscordEventApplicationCommandUpdate, OnApplicationCommandUpdate)
	h.RegisterEventHandler(structs.DiscordEventApplicationCommandDelete, OnApplicationCommandDelete)
	h.RegisterEventHandler(structs.DiscordEventChannelCreate, OnChannelCreate)
	h.RegisterEventHandler(structs.DiscordEventChannelUpdate, OnChannelUpdate)
	h.RegisterEventHandler(structs.DiscordEventChannelDelete, OnChannelDelete)
	h.RegisterEventHandler(structs.DiscordEventChannelPinsUpdate, OnChannelPinsUpdate)
	h.RegisterEventHandler(structs.DiscordEventThreadCreate, OnThreadCreate)
	h.RegisterEventHandler(structs.DiscordEventThreadUpdate, OnThreadUpdate)
	h.RegisterEventHandler(structs.DiscordEventThreadDelete, OnThreadDelete)
	// h.NewEventHandler(structs.DiscordEventThreadListSync, OnThreadListSync)
	h.RegisterEventHandler(structs.DiscordEventThreadMemberUpdate, OnThreadMemberUpdate)
	h.RegisterEventHandler(structs.DiscordEventThreadMembersUpdate, OnThreadMembersUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildCreate, OnGuildCreate)
	h.RegisterEventHandler(structs.DiscordEventGuildUpdate, OnGuildUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildDelete, OnGuildDelete)
	h.RegisterEventHandler(structs.DiscordEventGuildBanAdd, OnGuildBanAdd)
	h.RegisterEventHandler(structs.DiscordEventGuildBanRemove, OnGuildBanRemove)
	h.RegisterEventHandler(structs.DiscordEventGuildEmojisUpdate, OnGuildEmojisUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildStickersUpdate, OnGuildStickersUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildIntegrationsUpdate, OnGuildIntegrationsUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildMemberAdd, OnGuildMemberAdd)
	h.RegisterEventHandler(structs.DiscordEventGuildMemberRemove, OnGuildMemberRemove)
	h.RegisterEventHandler(structs.DiscordEventGuildMemberUpdate, OnGuildMemberUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildRoleCreate, OnGuildRoleCreate)
	h.RegisterEventHandler(structs.DiscordEventGuildRoleUpdate, OnGuildRoleUpdate)
	h.RegisterEventHandler(structs.DiscordEventGuildRoleDelete, OnGuildRoleDelete)
	h.RegisterEventHandler(structs.DiscordEventIntegrationCreate, OnIntegrationCreate)
	h.RegisterEventHandler(structs.DiscordEventIntegrationUpdate, OnIntegrationUpdate)
	h.RegisterEventHandler(structs.DiscordEventIntegrationDelete, OnIntegrationDelete)
	h.RegisterEventHandler(structs.DiscordEventInteractionCreate, OnInteractionCreate)
	h.RegisterEventHandler(structs.DiscordEventInviteCreate, OnInviteCreate)
	h.RegisterEventHandler(structs.DiscordEventInviteDelete, OnInviteDelete)
	h.RegisterEventHandler(structs.DiscordEventMessageCreate, OnMessageCreate)
	h.RegisterEventHandler(structs.DiscordEventMessageUpdate, OnMessageUpdate)
	h.RegisterEventHandler(structs.DiscordEventMessageDelete, OnMessageDelete)
	h.RegisterEventHandler(structs.DiscordEventMessageDeleteBulk, OnMessageDeleteBulk)
	h.RegisterEventHandler(structs.DiscordEventMessageReactionAdd, OnMessageReactionAdd)
	h.RegisterEventHandler(structs.DiscordEventMessageReactionRemove, OnMessageReactionRemove)
	h.RegisterEventHandler(structs.DiscordEventMessageReactionRemoveAll, OnMessageReactionRemoveAll)
	h.RegisterEventHandler(structs.DiscordEventMessageReactionRemoveEmoji, OnMessageReactionRemoveEmoji)
	h.RegisterEventHandler(structs.DiscordEventPresenceUpdate, OnPresenceUpdate)
	h.RegisterEventHandler(structs.DiscordEventStageInstanceCreate, OnStageInstanceCreate)
	h.RegisterEventHandler(structs.DiscordEventStageInstanceUpdate, OnStageInstanceUpdate)
	h.RegisterEventHandler(structs.DiscordEventStageInstanceDelete, OnStageInstanceDelete)
	h.RegisterEventHandler(structs.DiscordEventTypingStart, OnTypingStart)
	h.RegisterEventHandler(structs.DiscordEventUserUpdate, OnUserUpdate)
	h.RegisterEventHandler(structs.DiscordEventVoiceStateUpdate, OnVoiceStateUpdate)
	h.RegisterEventHandler(structs.DiscordEventVoiceServerUpdate, OnVoiceServerUpdate)
	h.RegisterEventHandler(structs.DiscordEventWebhookUpdate, OnWebhookUpdate)

	// Custom Events.
	h.RegisterEventHandler(structs.DiscordEventGuildJoin, OnGuildJoin)
	h.RegisterEventHandler(structs.DiscordEventGuildAvailable, OnGuildAvailable)

	h.RegisterEventHandler(structs.DiscordEventGuildLeave, OnGuildLeave)
	h.RegisterEventHandler(structs.DiscordEventGuildUnavailable, OnGuildUnavailable)
	h.RegisterEventHandler(structs.EventError, nil)

	return h
}

func NewSandwichHandlers() (h *Handlers) {
	h = &Handlers{
		eventHandlersMu: sync.RWMutex{},
		EventHandlers:   make(map[string]*EventHandler),
	}

	h.RegisterEventHandler(structs.SandwichEventConfigurationReload, OnSandwichConfigurationReload)
	h.RegisterEventHandler(structs.SandwichEventShardStatusUpdate, OnSandwichShardStatusUpdate)
	h.RegisterEventHandler(structs.SandwichEventShardGroupStatusUpdate, OnSandwichShardGroupStatusUpdate)

	// Register events that are handled by default.
	h.RegisterOnSandwichConfigurationReload(func(eventCtx *EventContext) (err error) {
		identifiers, err := eventCtx.Sandwich.grpcInterface.FetchConsumerConfiguration(eventCtx, "")
		if err != nil {
			return xerrors.Errorf("Failed to fetch consumer configuration: %v", err)
		}

		eventCtx.Sandwich.identifiersMu.Lock()
		eventCtx.Sandwich.Identifiers = map[string]*structs.ManagerConsumerConfiguration{}

		for k := range identifiers.Identifiers {
			v := identifiers.Identifiers[k]
			eventCtx.Sandwich.Identifiers[k] = &v
		}
		eventCtx.Sandwich.identifiersMu.Unlock()

		return nil
	})

	return h
}

type EventHandler struct {
	eventName string

	eventsMu sync.RWMutex
	Events   []interface{}

	Parser EventParser

	_handlers *Handlers
}

type EventParser func(eventCtx *EventContext, payload structs.SandwichPayload) (err error)

// Discord Events.

// RegisterEventHandler adds a new event handler. If there is already
// an event registered with the name, it is overridden.
func (h *Handlers) RegisterEventHandler(eventName string, parser EventParser) (eh *EventHandler) {
	h.eventHandlersMu.Lock()
	defer h.eventHandlersMu.Unlock()

	eh = &EventHandler{
		eventName: eventName,
		eventsMu:  sync.RWMutex{},
		Events:    make([]interface{}, 0),
		Parser:    parser,
		_handlers: h,
	}

	h.EventHandlers[eventName] = eh

	return eh
}

// Dispatch. dispatches a payload.
func (h *Handlers) Dispatch(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	return h.DispatchType(eventCtx, payload.Type, payload)
}

// DispatchType. is similar to Dispatch however a custom event name
// can. be passed, preserving the original payload.
func (h *Handlers) DispatchType(eventCtx *EventContext, eventName string, payload structs.SandwichPayload) (err error) {
	identifier, ok, err := eventCtx.Sandwich.FetchIdentifier(context.TODO(), payload.Metadata.Application)
	if !ok || err != nil {
		eventCtx.Logger.Warn().Err(err).Msg("Failed to fetch identifier for application")

		return err
	}

	eventCtx.Identifier = identifier

	if ev, ok := h.EventHandlers[eventName]; ok {
		eventCtx.EventHandler = ev

		defer func() {
			errorValue := recover()
			if errorValue != nil {
				eventCtx.Sandwich.RecoverEventPanic(errorValue, eventCtx, payload)
			}
		}()

		return ev.Parser(eventCtx, payload)
	}

	eventCtx.Logger.Info().Str("type", payload.Type).Msg("Unknown event handler")

	return ErrUnknownEvent
}

// WrapFuncType handles the error of a FuncType if it returns an error.
// It will call any ERROR handlers. Errors that occur in the ERROR handler
// will not trigger the ERROR handler.
func (h *Handlers) WrapFuncType(eventCtx *EventContext, funcTypeErr error) (err error) {
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
func OnReady(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var readyPayload discord.Ready
	if err = eventCtx.decodeContent(payload, &readyPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
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

type OnReadyFuncType func(eventCtx *EventContext) (err error)

// OnResumed.
func OnResumed(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var resumePayload discord.Resume
	if err = eventCtx.decodeContent(payload, &resumePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
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

type OnResumedFuncType func(eventCtx *EventContext) (err error)

// OnApplicationCommandCreate.
func OnApplicationCommandCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var applicationCommandCreatePayload discord.ApplicationCommandCreate
	if err = eventCtx.decodeContent(payload, &applicationCommandCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if applicationCommandCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *applicationCommandCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, ApplicationCommand(*applicationCommandCreatePayload)))
		}
	}

	return nil
}

type OnApplicationCommandCreateFuncType func(eventCtx *EventContext, command ApplicationCommand) (err error)

// OnApplicationCommandUpdate.
func OnApplicationCommandUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var applicationCommandUpdatePayload discord.ApplicationCommandUpdate
	if err = eventCtx.decodeContent(payload, &applicationCommandUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if applicationCommandUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *applicationCommandUpdatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, ApplicationCommand(*applicationCommandUpdatePayload)))
		}
	}

	return nil
}

type OnApplicationCommandUpdateFuncType func(eventCtx *EventContext, command ApplicationCommand) (err error)

// OnApplicationCommandDelete.
func OnApplicationCommandDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var applicationCommandDeletePayload discord.ApplicationCommandDelete
	if err = eventCtx.decodeContent(payload, &applicationCommandDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if applicationCommandDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *applicationCommandDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, ApplicationCommand(*applicationCommandDeletePayload)))
		}
	}

	return nil
}

type OnApplicationCommandDeleteFuncType func(eventCtx *EventContext, command ApplicationCommand) (err error)

// OnChannelCreate.
func OnChannelCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var channelCreatePayload discord.ChannelCreate
	if err = eventCtx.decodeContent(payload, &channelCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if channelCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *channelCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel(*channelCreatePayload)))
		}
	}

	return nil
}

type OnChannelCreateFuncType func(eventCtx *EventContext, channel Channel) (err error)

// OnChannelUpdate.
func OnChannelUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var channelUpdatePayload discord.ChannelUpdate
	if err = eventCtx.decodeContent(payload, &channelUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if channelUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *channelUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeChannel); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel(beforeChannel), Channel(*channelUpdatePayload)))
		}
	}

	return nil
}

type OnChannelUpdateFuncType func(eventCtx *EventContext, before Channel, after Channel) (err error)

// OnChannelDelete.
func OnChannelDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var channelDeletePayload discord.ChannelDelete
	if err = eventCtx.decodeContent(payload, &channelDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if channelDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *channelDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel(*channelDeletePayload)))
		}
	}

	return nil
}

type OnChannelDeleteFuncType func(eventCtx *EventContext, channel Channel) (err error)

// OnChannelPinsUpdate.
func OnChannelPinsUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var channelPinsUpdatePayload discord.ChannelPinsUpdate
	if err = eventCtx.decodeContent(payload, &channelPinsUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, channelPinsUpdatePayload.GuildID)

	var timestamp time.Time
	if channelPinsUpdatePayload.LastPinTimestamp != "" {
		timestamp, _ = parseTimeStamp(channelPinsUpdatePayload.LastPinTimestamp)
	}

	channel := NewChannel(eventCtx, &channelPinsUpdatePayload.GuildID, channelPinsUpdatePayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnChannelPinsUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, timestamp))
		}
	}

	return nil
}

type OnChannelPinsUpdateFuncType func(eventCtx *EventContext, channel *Channel, lastPinTimestamp time.Time) (err error)

// OnThreadCreate.
func OnThreadCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var threadCreatePayload discord.ThreadCreate
	if err = eventCtx.decodeContent(payload, &threadCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *threadCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel(*threadCreatePayload)))
		}
	}

	return nil
}

type OnThreadCreateFuncType func(eventCtx *EventContext, thread Channel) (err error)

// OnThreadUpdate.
func OnThreadUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var threadUpdatePayload discord.ThreadUpdate
	if err = eventCtx.decodeContent(payload, &threadUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *threadUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeChannel); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel(beforeChannel), Channel(*threadUpdatePayload)))
		}
	}

	return nil
}

type OnThreadUpdateFuncType func(eventCtx *EventContext, before Channel, after Channel) (err error)

// OnThreadDelete.
func OnThreadDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var threadDeletePayload discord.ThreadDelete
	if err = eventCtx.decodeContent(payload, &threadDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *threadDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel(*threadDeletePayload)))
		}
	}

	return nil
}

type OnThreadDeleteFuncType func(eventCtx *EventContext, thread Channel) (err error)

// // OnThreadListSync.
// func. OnThreadListSync(eventCtx *Context, payload structs.SandwichPayload) (err error) {
// 	var threadListSyncPayload discord.ThreadListSync
// 	if err = eventCtx.decodeContent(payload, &threadListSyncPayload); err != nil {
// 		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
// 	}

// 	for _, event := range eventCtx.EventHandler.Events {
// 		if f, ok := event.(OnThreadListSyncFuncType); ok {
// 			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Channel()))
// 		}
// 	}
//
// 	return nil
// }

// type. OnThreadListSyncFuncType func(eventCtx *Context, thread Channel) (err error)

// OnThreadMemberUpdate.
func OnThreadMemberUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var threadMemberUpdatePayload discord.ThreadMemberUpdate
	if err = eventCtx.decodeContent(payload, &threadMemberUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadMemberUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *threadMemberUpdatePayload.GuildID)
	}

	channel := NewChannel(eventCtx, threadMemberUpdatePayload.GuildID, *threadMemberUpdatePayload.UserID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnThreadMemberUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, ThreadMember(*threadMemberUpdatePayload)))
		}
	}

	return nil
}

type OnThreadMemberUpdateFuncType func(eventCtx *EventContext, thread *Channel, user ThreadMember) (err error)

// OnThreadMembersUpdate.
func OnThreadMembersUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var threadMembersUpdatePayload discord.ThreadMembersUpdate
	if err = eventCtx.decodeContent(payload, &threadMembersUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, threadMembersUpdatePayload.GuildID)

	channel := NewChannel(eventCtx, &threadMembersUpdatePayload.GuildID, threadMembersUpdatePayload.ID)

	addedUsers := make([]*User, 0, len(threadMembersUpdatePayload.AddedMembers))
	for _, addedMember := range threadMembersUpdatePayload.AddedMembers {
		addedUsers = append(addedUsers, NewUser(eventCtx, *addedMember.UserID))
	}

	removedUsers := make([]*User, 0, len(threadMembersUpdatePayload.RemovedMemberIDs))
	for _, removedUser := range threadMembersUpdatePayload.RemovedMemberIDs {
		removedUsers = append(removedUsers, NewUser(eventCtx, removedUser))
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

type OnThreadMembersUpdateFuncType func(eventCtx *EventContext, thread *Channel, addedUsers []*User, removedUsers []*User) (err error)

// OnGuildCreate.
func OnGuildCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildCreatePayload discord.GuildCreate
	if err = eventCtx.decodeContent(payload, &guildCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildCreatePayload.ID)

	var lazy bool
	if _, err := eventCtx.decodeExtra(payload, "lazy", &lazy); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	// If. true, the guild was previously unavailable.
	var unavailable bool
	if _, err := eventCtx.decodeExtra(payload, "unavailable", &unavailable); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
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
func OnGuildUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildUpdatePayload discord.GuildUpdate
	if err = eventCtx.decodeContent(payload, &guildUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	guild := Guild(*guildUpdatePayload)
	eventCtx.Guild = &guild

	var beforeGuild discord.Guild
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeGuild); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Guild(beforeGuild), guild))
		}
	}

	return nil
}

type OnGuildUpdateFuncType func(eventCtx *EventContext, before Guild, after Guild) (err error)

// OnGuildDelete.
func OnGuildDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildDeletePayload discord.GuildDelete
	if err = eventCtx.decodeContent(payload, &guildDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildDeletePayload.ID)

	if guildDeletePayload.Unavailable {
		return eventCtx.Handlers.DispatchType(eventCtx, "GUILD_UNAVAILABLE", payload)
	}

	return eventCtx.Handlers.DispatchType(eventCtx, "GUILD_REMOVE", payload)
}

// OnGuildBanAdd.
func OnGuildBanAdd(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildBanAddPayload discord.GuildBanAdd
	if err = eventCtx.decodeContent(payload, &guildBanAddPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildBanAddPayload.GuildID)

	user := User(*guildBanAddPayload.User)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildBanAddFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, user))
		}
	}

	return nil
}

type OnGuildBanAddFuncType func(eventCtx *EventContext, user User) (err error)

// OnGuildBanRemove.
func OnGuildBanRemove(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildBanRemovePayload discord.GuildBanRemove
	if err = eventCtx.decodeContent(payload, &guildBanRemovePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildBanRemovePayload.GuildID)

	user := User(*guildBanRemovePayload.User)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildBanRemoveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, user))
		}
	}

	return nil
}

type OnGuildBanRemoveFuncType func(eventCtx *EventContext, user User) (err error)

// OnGuildEmojisUpdate.
func OnGuildEmojisUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildEmojisUpdatePayload discord.GuildEmojisUpdate
	if err = eventCtx.decodeContent(payload, &guildEmojisUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildEmojisUpdatePayload.GuildID)

	var before []Emoji
	if _, err := eventCtx.decodeExtra(payload, "before", &before); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	after := make([]Emoji, 0, len(guildEmojisUpdatePayload.Emojis))
	for _, emoji := range guildEmojisUpdatePayload.Emojis {
		after = append(after, Emoji(*emoji))
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

type OnGuildEmojisUpdateFuncType func(eventCtx *EventContext, before []Emoji, after []Emoji) (err error)

// OnGuildStickersUpdate.
func OnGuildStickersUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildStickersUpdatePayload discord.GuildStickersUpdate
	if err = eventCtx.decodeContent(payload, &guildStickersUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildStickersUpdatePayload.GuildID)

	var before []Sticker
	if _, err := eventCtx.decodeExtra(payload, "before", &before); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	after := make([]Sticker, 0, len(guildStickersUpdatePayload.Stickers))
	for _, sticker := range guildStickersUpdatePayload.Stickers {
		after = append(after, Sticker(*sticker))
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

type OnGuildStickersUpdateFuncType func(eventCtx *EventContext, before []Sticker, after []Sticker) (err error)

// OnGuildIntegrationsUpdate.
func OnGuildIntegrationsUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildIntegrationsUpdatePayload discord.GuildIntegrationsUpdate
	if err = eventCtx.decodeContent(payload, &guildIntegrationsUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildIntegrationsUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildIntegrationsUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnGuildIntegrationsUpdateFuncType func(eventCtx *EventContext) (err error)

// OnGuildMemberAdd.
func OnGuildMemberAdd(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildMemberAddPayload discord.GuildMemberAdd
	if err = eventCtx.decodeContent(payload, &guildMemberAddPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, *guildMemberAddPayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberAddFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, GuildMember(*guildMemberAddPayload)))
		}
	}

	return nil
}

type OnGuildMemberAddFuncType func(eventCtx *EventContext, member GuildMember) (err error)

// OnGuildMemberRemove.
func OnGuildMemberRemove(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildMemberRemovePayload discord.GuildMemberRemove
	if err = eventCtx.decodeContent(payload, &guildMemberRemovePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildMemberRemovePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberRemoveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, User(*guildMemberRemovePayload.User)))
		}
	}

	return nil
}

type OnGuildMemberRemoveFuncType func(eventCtx *EventContext, member User) (err error)

// OnGuildMemberUpdate.
func OnGuildMemberUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildMemberUpdatePayload discord.GuildMemberUpdate
	if err = eventCtx.decodeContent(payload, &guildMemberUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, *guildMemberUpdatePayload.GuildID)

	var beforeGuildMember discord.GuildMember
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeGuildMember); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, GuildMember(beforeGuildMember), GuildMember(*guildMemberUpdatePayload.GuildMember)))
		}
	}

	return nil
}

type OnGuildMemberUpdateFuncType func(eventCtx *EventContext, before GuildMember, after GuildMember) (err error)

// OnGuildRoleCreate.
func OnGuildRoleCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildRoleCreatePayload discord.GuildRoleCreate
	if err = eventCtx.decodeContent(payload, &guildRoleCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if guildRoleCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *guildRoleCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Role(*guildRoleCreatePayload)))
		}
	}

	return nil
}

type OnGuildRoleCreateFuncType func(eventCtx *EventContext, role Role) (err error)

// OnGuildRoleUpdate.
func OnGuildRoleUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildRoleUpdatePayload discord.GuildRoleUpdate
	if err = eventCtx.decodeContent(payload, &guildRoleUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if guildRoleUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *guildRoleUpdatePayload.GuildID)
	}

	var beforeRole discord.Role
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeRole); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Role(beforeRole), Role(*guildRoleUpdatePayload)))
		}
	}

	return nil
}

type OnGuildRoleUpdateFuncType func(eventCtx *EventContext, before Role, after Role) (err error)

// OnGuildRoleDelete.
func OnGuildRoleDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildRoleDeletePayload discord.GuildRoleDelete
	if err = eventCtx.decodeContent(payload, &guildRoleDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildRoleDeletePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, guildRoleDeletePayload.RoleID))
		}
	}

	return nil
}

type OnGuildRoleDeleteFuncType func(eventCtx *EventContext, roleID discord.Snowflake) (err error)

// OnIntegrationCreate.
func OnIntegrationCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var integrationCreatePayload discord.IntegrationCreate
	if err = eventCtx.decodeContent(payload, &integrationCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if integrationCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *integrationCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Integration(*integrationCreatePayload)))
		}
	}

	return nil
}

type OnIntegrationCreateFuncType func(eventCtx *EventContext, integration Integration) (err error)

// OnIntegrationUpdate.
func OnIntegrationUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var integrationUpdatePayload discord.IntegrationUpdate
	if err = eventCtx.decodeContent(payload, &integrationUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if integrationUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *integrationUpdatePayload.GuildID)
	}

	var beforeIntegration discord.Integration
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeIntegration); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnIntegrationUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Integration(beforeIntegration), Integration(*integrationUpdatePayload)))
		}
	}

	return nil
}

type OnIntegrationUpdateFuncType func(eventCtx *EventContext, before Integration, after Integration) (err error)

// OnIntegrationDelete.
func OnIntegrationDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var integrationDeletePayload discord.IntegrationDelete
	if err = eventCtx.decodeContent(payload, &integrationDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, integrationDeletePayload.GuildID)

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

type OnIntegrationDeleteFuncType func(eventCtx *EventContext, integrationID discord.Snowflake, applicationID discord.Snowflake) (err error)

// OnInteractionCreate.
func OnInteractionCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var interactionCreatePayload discord.InteractionCreate
	if err = eventCtx.decodeContent(payload, &interactionCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if interactionCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *interactionCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInteractionCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Interaction(*interactionCreatePayload)))
		}
	}

	return nil
}

type OnInteractionCreateFuncType func(eventCtx *EventContext, interaction Interaction) (err error)

// OnInviteCreate.
func OnInviteCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var inviteCreatePayload discord.InviteCreate
	if err = eventCtx.decodeContent(payload, &inviteCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if inviteCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *inviteCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInviteCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Invite(*inviteCreatePayload)))
		}
	}

	return nil
}

type OnInviteCreateFuncType func(eventCtx *EventContext, invite Invite) (err error)

// OnInviteDelete.
func OnInviteDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var inviteDeletePayload discord.InviteDelete
	if err = eventCtx.decodeContent(payload, &inviteDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if inviteDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *inviteDeletePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnInviteDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Invite(*inviteDeletePayload)))
		}
	}

	return nil
}

type OnInviteDeleteFuncType func(eventCtx *EventContext, invite Invite) (err error)

// OnMessageCreate.
func OnMessageCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageCreatePayload discord.MessageCreate
	if err = eventCtx.decodeContent(payload, &messageCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageCreatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *messageCreatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Message(*messageCreatePayload)))
		}
	}

	return nil
}

type OnMessageCreateFuncType func(eventCtx *EventContext, message Message) (err error)

// OnMessageUpdate.
func OnMessageUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageUpdatePayload discord.MessageUpdate
	if err = eventCtx.decodeContent(payload, &messageUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *messageUpdatePayload.GuildID)
	}

	var beforeMessage discord.Message
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeMessage); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, Message(beforeMessage), Message(*messageUpdatePayload)))
		}
	}

	return nil
}

type OnMessageUpdateFuncType func(eventCtx *EventContext, before Message, after Message) (err error)

// OnMessageDelete.
func OnMessageDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageDeletePayload discord.MessageDelete
	if err = eventCtx.decodeContent(payload, &messageDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageDeletePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *messageDeletePayload.GuildID)
	}

	channel := NewChannel(eventCtx, messageDeletePayload.GuildID, messageDeletePayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageDeletePayload.ID))
		}
	}

	return nil
}

type OnMessageDeleteFuncType func(eventCtx *EventContext, channel *Channel, messageID discord.Snowflake) (err error)

// OnMessageDeleteBulk.
func OnMessageDeleteBulk(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageDeleteBulkPayload discord.MessageDeleteBulk
	if err = eventCtx.decodeContent(payload, &messageDeleteBulkPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageDeleteBulkPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *messageDeleteBulkPayload.GuildID)
	}

	channel := NewChannel(eventCtx, messageDeleteBulkPayload.GuildID, messageDeleteBulkPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteBulkFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageDeleteBulkPayload.IDs))
		}
	}

	return nil
}

type OnMessageDeleteBulkFuncType func(eventCtx *EventContext, channel *Channel, messageIDs []discord.Snowflake) (err error)

// OnMessageReactionAdd.
func OnMessageReactionAdd(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageReactionAddPayload discord.MessageReactionAdd
	if err = eventCtx.decodeContent(payload, &messageReactionAddPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, messageReactionAddPayload.GuildID)

	channel := NewChannel(eventCtx, &messageReactionAddPayload.GuildID, messageReactionAddPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionAddFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionAddPayload.MessageID, Emoji(*messageReactionAddPayload.Emoji), GuildMember(*messageReactionAddPayload.Member)))
		}
	}

	return nil
}

type OnMessageReactionAddFuncType func(eventCtx *EventContext, channel *Channel, messageID discord.Snowflake, emoji Emoji, guildMember GuildMember) (err error)

// OnMessageReactionRemove.
func OnMessageReactionRemove(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageReactionRemovePayload discord.MessageReactionRemove
	if err = eventCtx.decodeContent(payload, &messageReactionRemovePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageReactionRemovePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *messageReactionRemovePayload.GuildID)
	}

	channel := NewChannel(eventCtx, messageReactionRemovePayload.GuildID, messageReactionRemovePayload.ChannelID)
	user := NewUser(eventCtx, messageReactionRemovePayload.UserID)
	emoji := Emoji(*messageReactionRemovePayload.Emoji)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemovePayload.MessageID, &emoji, user))
		}
	}

	return nil
}

type OnMessageReactionRemoveFuncType func(eventCtx *EventContext, channel *Channel, messageID discord.Snowflake, emoji *Emoji, user *User) (err error)

// OnMessageReactionRemoveAll.
func OnMessageReactionRemoveAll(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageReactionRemoveAllPayload discord.MessageReactionRemoveAll
	if err = eventCtx.decodeContent(payload, &messageReactionRemoveAllPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, messageReactionRemoveAllPayload.GuildID)

	channel := NewChannel(eventCtx, &messageReactionRemoveAllPayload.GuildID, messageReactionRemoveAllPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveAllFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemoveAllPayload.MessageID))
		}
	}

	return nil
}

type OnMessageReactionRemoveAllFuncType func(eventCtx *EventContext, channel *Channel, messageID discord.Snowflake) (err error)

// OnMessageReactionRemoveEmoji.
func OnMessageReactionRemoveEmoji(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var messageReactionRemoveEmojiPayload discord.MessageReactionRemoveEmoji
	if err = eventCtx.decodeContent(payload, &messageReactionRemoveEmojiPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageReactionRemoveEmojiPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *messageReactionRemoveEmojiPayload.GuildID)
	}

	channel := NewChannel(eventCtx, messageReactionRemoveEmojiPayload.GuildID, messageReactionRemoveEmojiPayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveEmojiFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel, messageReactionRemoveEmojiPayload.MessageID, Emoji(*messageReactionRemoveEmojiPayload.Emoji)))
		}
	}

	return nil
}

type OnMessageReactionRemoveEmojiFuncType func(eventCtx *EventContext, channel *Channel, messageID discord.Snowflake, emoji Emoji) (err error)

// OnPresenceUpdate.
func OnPresenceUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var presenceUpdatePayload discord.PresenceUpdate
	if err = eventCtx.decodeContent(payload, &presenceUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, presenceUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnPresenceUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, User(*presenceUpdatePayload.User), presenceUpdatePayload))
		}
	}

	return nil
}

type OnPresenceUpdateFuncType func(eventCtx *EventContext, user User, payload discord.PresenceUpdate) (err error)

// OnStageInstanceCreate.
func OnStageInstanceCreate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var stageInstanceCreatePayload discord.StageInstanceCreate
	if err = eventCtx.decodeContent(payload, &stageInstanceCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, stageInstanceCreatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceCreateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, StageInstance(*stageInstanceCreatePayload)))
		}
	}

	return nil
}

type OnStageInstanceCreateFuncType func(eventCtx *EventContext, stage StageInstance) (err error)

// OnStageInstanceUpdate.
func OnStageInstanceUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var stageInstanceUpdatePayload discord.StageInstanceUpdate
	if err = eventCtx.decodeContent(payload, &stageInstanceUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, stageInstanceUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, StageInstance(*stageInstanceUpdatePayload)))
		}
	}

	return nil
}

type OnStageInstanceUpdateFuncType func(eventCtx *EventContext, stage StageInstance) (err error)

// OnStageInstanceDelete.
func OnStageInstanceDelete(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var stageInstanceDeletePayload discord.StageInstanceDelete
	if err = eventCtx.decodeContent(payload, &stageInstanceDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceDeleteFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, StageInstance(*stageInstanceDeletePayload)))
		}
	}

	return nil
}

type OnStageInstanceDeleteFuncType func(eventCtx *EventContext, stage StageInstance) (err error)

// OnTypingStart.
func OnTypingStart(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var typingStartPayload discord.TypingStart
	if err = eventCtx.decodeContent(payload, &typingStartPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if typingStartPayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *typingStartPayload.GuildID)
	}

	timestamp := time.Unix(int64(typingStartPayload.Timestamp), 0)
	channel := NewChannel(eventCtx, typingStartPayload.GuildID, typingStartPayload.ChannelID)

	var user *User

	var member *GuildMember

	if typingStartPayload.Member != nil {
		member = (*GuildMember)(typingStartPayload.Member)
		user = (*User)(typingStartPayload.Member.User)
	} else {
		user = NewUser(eventCtx, typingStartPayload.UserID)
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

type OnTypingStartFuncType func(eventCtx *EventContext, channel *Channel, member *GuildMember, user *User, timestamp time.Time) (err error)

// OnUserUpdate.
func OnUserUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var userUpdatePayload discord.UserUpdate
	if err = eventCtx.decodeContent(payload, &userUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	var beforeUser discord.User
	if _, err := eventCtx.decodeExtra(payload, "before", &beforeUser); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnUserUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, User(beforeUser), User(*userUpdatePayload)))
		}
	}

	return nil
}

type OnUserUpdateFuncType func(eventCtx *EventContext, before User, after User) (err error)

// OnVoiceStateUpdate.
func OnVoiceStateUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var voiceStateUpdatePayload discord.VoiceStateUpdate
	if err = eventCtx.decodeContent(payload, &voiceStateUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if voiceStateUpdatePayload.GuildID != nil {
		eventCtx.Guild = NewGuild(eventCtx, *voiceStateUpdatePayload.GuildID)
	}

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnVoiceStateUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, GuildMember(*voiceStateUpdatePayload.Member), VoiceState(*voiceStateUpdatePayload)))
		}
	}

	return nil
}

type OnVoiceStateUpdateFuncType func(eventCtx *EventContext, member GuildMember, voice VoiceState) (err error)

// OnVoiceServerUpdate.
func OnVoiceServerUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var voiceServerUpdatePayload discord.VoiceServerUpdate
	if err = eventCtx.decodeContent(payload, &voiceServerUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, voiceServerUpdatePayload.GuildID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnVoiceServerUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, voiceServerUpdatePayload))
		}
	}

	return nil
}

type OnVoiceServerUpdateFuncType func(eventCtx *EventContext, payload discord.VoiceServerUpdate) (err error)

// OnWebhookUpdate.
func OnWebhookUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var webhookUpdatePayload discord.WebhookUpdate
	if err = eventCtx.decodeContent(payload, &webhookUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, webhookUpdatePayload.GuildID)

	channel := NewChannel(eventCtx, &webhookUpdatePayload.GuildID, webhookUpdatePayload.ChannelID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnWebhookUpdateFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, channel))
		}
	}

	return nil
}

type OnWebhookUpdateFuncType func(eventCtx *EventContext, channel *Channel) (err error)

// OnGuildJoin.
func OnGuildJoin(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildCreatePayload discord.GuildCreate
	if err = eventCtx.decodeContent(payload, &guildCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	guild := Guild(*guildCreatePayload)
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

type OnGuildJoinFuncType func(eventCtx *EventContext, guild Guild) (err error)

// OnGuildAvailable.
func OnGuildAvailable(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildCreatePayload discord.GuildCreate
	if err = eventCtx.decodeContent(payload, &guildCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	guild := Guild(*guildCreatePayload)
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

type OnGuildAvailableFuncType func(eventCtx *EventContext, guild Guild) (err error)

// OnGuildLeave.
func OnGuildLeave(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildDeletePayload discord.GuildDelete
	if err = eventCtx.decodeContent(payload, &guildDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildDeletePayload.ID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildLeaveFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, UnavailableGuild(*guildDeletePayload)))
		}
	}

	return nil
}

type OnGuildLeaveFuncType func(eventCtx *EventContext, unavailableGuild UnavailableGuild) (err error)

// OnGuildUnavailable.
func OnGuildUnavailable(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var guildDeletePayload discord.GuildDelete
	if err = eventCtx.decodeContent(payload, &guildDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	eventCtx.Guild = NewGuild(eventCtx, guildDeletePayload.ID)

	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnGuildUnavailableFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx, UnavailableGuild(*guildDeletePayload)))
		}
	}

	return nil
}

type OnGuildUnavailableFuncType func(eventCtx *EventContext, unavailableGuild UnavailableGuild) (err error)

// Sandwich Events.

// OnSandwichConfigurationReload.
func OnSandwichConfigurationReload(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	eventCtx.EventHandler.eventsMu.RLock()
	defer eventCtx.EventHandler.eventsMu.RUnlock()

	for _, event := range eventCtx.EventHandler.Events {
		if f, ok := event.(OnSandwichConfigurationReloadFuncType); ok {
			return eventCtx.Handlers.WrapFuncType(eventCtx, f(eventCtx))
		}
	}

	return nil
}

type OnSandwichConfigurationReloadFuncType func(eventCtx *EventContext) (err error)

// OnSandwichShardStatusUpdate.
func OnSandwichShardStatusUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var shardStatusUpdatePayload structs.ShardStatusUpdate
	if err = eventCtx.decodeContent(payload, &shardStatusUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
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

type OnSandwichShardStatusUpdateFuncType func(eventCtx *EventContext, manager string, shardGroup int32, shard int32, status structs.ShardStatus) (err error)

// OnSandwichShardGroupStatusUpdate.
func OnSandwichShardGroupStatusUpdate(eventCtx *EventContext, payload structs.SandwichPayload) (err error) {
	var shardGroupStatusUpdatePayload structs.ShardGroupStatusUpdate
	if err = eventCtx.decodeContent(payload, &shardGroupStatusUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
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

type OnSandwichShardGroupStatusUpdateFuncType func(eventCtx *EventContext, manager string, shardGroup int32, status structs.ShardGroupStatus) (err error)

// Generic Events.

type OnErrorFuncType func(eventCtx *EventContext, eventErr error) (err error)
