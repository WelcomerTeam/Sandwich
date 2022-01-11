package internal

import (
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

func NewHandlers() *Handlers {
	return &Handlers{
		eventHandlersMu: sync.RWMutex{},
		EventHandlers:   make(map[string]*EventHandler),
	}
}

type EventHandler struct {
	eventName string

	eventsMu sync.RWMutex
	Events   []interface{}

	Parser EventParser

	_handlers *Handlers
}

type EventParser func(ctx *Context, payload structs.SandwichPayload) (err error)

// RegisterParsers creates a new EventHandler. If there is already an event
// registered with the name, it is ignored.
func (h *Handlers) NewEventHandler(eventName string, parser EventParser) (eh *EventHandler) {
	h.eventHandlersMu.Lock()
	defer h.eventHandlersMu.Unlock()

	eh, ok := h.EventHandlers[eventName]

	if !ok {
		eh = &EventHandler{
			eventName: eventName,
			eventsMu:  sync.RWMutex{},
			Events:    make([]interface{}, 0),
			Parser:    parser,
			_handlers: h,
		}

		h.EventHandlers[eventName] = eh
	}

	return eh
}

// Dispatch dispatches a payload.
func (h *Handlers) Dispatch(ctx *Context, payload structs.SandwichPayload) (err error) {
	if ev, ok := h.EventHandlers[payload.Type]; ok {
		ctx.Handlers = h
		ctx.EventHandler = ev

		return ev.Parser(ctx, payload)
	}

	return xerrors.Errorf("No event handler: %s", payload.Type)
}

// DispatchType is similar to Dispatch however a custom event name
// can be passed, preserving the origional payload.
func (h *Handlers) DispatchType(ctx *Context, eventName string, payload structs.SandwichPayload) (err error) {
	if ev, ok := h.EventHandlers[eventName]; ok {
		ctx.Handlers = h
		ctx.EventHandler = ev

		return ev.Parser(ctx, payload)
	}

	return xerrors.Errorf("No event handler: %s", payload.Type)
}

// OnReady
func OnReady(ctx *Context, payload structs.SandwichPayload) (err error) {
	var readyPayload discord.Ready
	if err = ctx.decodeContent(payload, &readyPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnReadyFuncType); ok {
			f(ctx)
		}
	}
	return nil
}

type OnReadyFuncType func(ctx *Context) (err error)

// OnResumed
func OnResumed(ctx *Context, payload structs.SandwichPayload) (err error) {
	var resumePayload discord.Resume
	if err = ctx.decodeContent(payload, &resumePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnResumedFuncType); ok {
			f(ctx)
		}
	}
	return nil
}

type OnResumedFuncType func(ctx *Context) (err error)

// OnApplicationCommandCreate
func OnApplicationCommandCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var applicationCommandCreatePayload discord.ApplicationCommandCreate
	if err = ctx.decodeContent(payload, &applicationCommandCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if applicationCommandCreatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *applicationCommandCreatePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandCreateFuncType); ok {
			f(ctx, ApplicationCommand(*applicationCommandCreatePayload))
		}
	}
	return nil
}

type OnApplicationCommandCreateFuncType func(ctx *Context, command ApplicationCommand) (err error)

// OnApplicationCommandUpdate
func OnApplicationCommandUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var applicationCommandUpdatePayload discord.ApplicationCommandUpdate
	if err = ctx.decodeContent(payload, &applicationCommandUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if applicationCommandUpdatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *applicationCommandUpdatePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandUpdateFuncType); ok {
			f(ctx, ApplicationCommand(*applicationCommandUpdatePayload))
		}
	}
	return nil
}

type OnApplicationCommandUpdateFuncType func(ctx *Context, command ApplicationCommand) (err error)

// OnApplicationCommandDelete
func OnApplicationCommandDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var applicationCommandDeletePayload discord.ApplicationCommandDelete
	if err = ctx.decodeContent(payload, &applicationCommandDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if applicationCommandDeletePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *applicationCommandDeletePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnApplicationCommandDeleteFuncType); ok {
			f(ctx, ApplicationCommand(*applicationCommandDeletePayload))
		}
	}
	return nil
}

type OnApplicationCommandDeleteFuncType func(ctx *Context, command ApplicationCommand) (err error)

// OnChannelCreate
func OnChannelCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var channelCreatePayload discord.ChannelCreate
	if err = ctx.decodeContent(payload, &channelCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if channelCreatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *channelCreatePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnChannelCreateFuncType); ok {
			f(ctx, Channel(*channelCreatePayload))
		}
	}
	return nil
}

type OnChannelCreateFuncType func(ctx *Context, channel Channel) (err error)

// OnChannelUpdate
func OnChannelUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var channelUpdatePayload discord.ChannelUpdate
	if err = ctx.decodeContent(payload, &channelUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if channelUpdatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *channelUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := ctx.decodeExtra(payload, "before", &beforeChannel); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnChannelUpdateFuncType); ok {
			f(ctx, Channel(beforeChannel), Channel(*channelUpdatePayload))
		}
	}
	return nil
}

type OnChannelUpdateFuncType func(ctx *Context, before Channel, after Channel) (err error)

// OnChannelDelete
func OnChannelDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var channelDeletePayload discord.ChannelDelete
	if err = ctx.decodeContent(payload, &channelDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if channelDeletePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *channelDeletePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnChannelDeleteFuncType); ok {
			f(ctx, Channel(*channelDeletePayload))
		}
	}
	return nil
}

type OnChannelDeleteFuncType func(ctx *Context, channel Channel) (err error)

// OnChannelPinsUpdate
func OnChannelPinsUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var channelPinsUpdatePayload discord.ChannelPinsUpdate
	if err = ctx.decodeContent(payload, &channelPinsUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, channelPinsUpdatePayload.GuildID)

	var timestamp time.Time
	if channelPinsUpdatePayload.LastPinTimestamp != nil {
		timestamp, _ = parseTimeStamp(*channelPinsUpdatePayload.LastPinTimestamp)
	}

	channel := NewChannel(ctx, channelPinsUpdatePayload.ChannelID, &channelPinsUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnChannelPinsUpdateFuncType); ok {
			f(ctx, channel, timestamp)
		}
	}
	return nil
}

type OnChannelPinsUpdateFuncType func(ctx *Context, channel Channel, lastPinTimestamp time.Time) (err error)

// OnThreadCreate
func OnThreadCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var threadCreatePayload discord.ThreadCreate
	if err = ctx.decodeContent(payload, &threadCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadCreatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *threadCreatePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnThreadCreateFuncType); ok {
			f(ctx, Channel(*threadCreatePayload))
		}
	}
	return nil
}

type OnThreadCreateFuncType func(ctx *Context, thread Channel) (err error)

// OnThreadUpdate
func OnThreadUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var threadUpdatePayload discord.ThreadUpdate
	if err = ctx.decodeContent(payload, &threadUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadUpdatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *threadUpdatePayload.GuildID)
	}

	var beforeChannel discord.Channel
	if _, err := ctx.decodeExtra(payload, "before", &beforeChannel); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnThreadUpdateFuncType); ok {
			f(ctx, Channel(beforeChannel), Channel(*threadUpdatePayload))
		}
	}
	return nil
}

type OnThreadUpdateFuncType func(ctx *Context, before Channel, after Channel) (err error)

// OnThreadDelete
func OnThreadDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var threadDeletePayload discord.ThreadDelete
	if err = ctx.decodeContent(payload, &threadDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if threadDeletePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *threadDeletePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnThreadDeleteFuncType); ok {
			f(ctx, Channel(*threadDeletePayload))
		}
	}
	return nil
}

type OnThreadDeleteFuncType func(ctx *Context, thread Channel) (err error)

// // OnThreadListSync
// func OnThreadListSync(ctx *Context, payload structs.SandwichPayload) (err error) {
// 	var threadListSyncPayload discord.ThreadListSync
// 	if err = ctx.decodeContent(payload, &threadListSyncPayload); err != nil {
// 		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
// 	}

// 	for _, event := range ctx.EventHandler.Events {
// 		if f, ok := event.(OnThreadListSyncFuncType); ok {
// 			f(ctx, Channel())
// 		}
// 	}
// 	return nil
// }

// type OnThreadListSyncFuncType func(ctx *Context, thread Channel) (err error)

// OnThreadMemberUpdate
func OnThreadMemberUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var threadMemberUpdatePayload discord.ThreadMemberUpdate
	if err = ctx.decodeContent(payload, &threadMemberUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	channel := NewChannel(ctx, threadMemberUpdatePayload.UserID, threadMemberUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnThreadMemberUpdateFuncType); ok {
			f(ctx, channel, ThreadMember(threadMemberUpdatePayload))
		}
	}
	return nil
}

type OnThreadMemberUpdateFuncType func(ctx *Context, thread Channel, user ThreadMember) (err error)

// OnThreadMembersUpdate
func OnThreadMembersUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var threadMembersUpdatePayload discord.ThreadMembersUpdate
	if err = ctx.decodeContent(payload, &threadMembersUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, threadMembersUpdatePayload.GuildID)

	channel := NewChannel(ctx, threadMembersUpdatePayload.ID, &threadMembersUpdatePayload.GuildID)

	addedUsers := make([]User, 0, len(threadMembersUpdatePayload.AddedMembers))
	for _, addedMember := range threadMembersUpdatePayload.AddedMembers {
		addedUsers = append(addedUsers, NewUser(ctx, *addedMember.UserID))
	}

	removedUsers := make([]User, 0, len(threadMembersUpdatePayload.RemovedMemberIDs))
	for _, removedUser := range threadMembersUpdatePayload.RemovedMemberIDs {
		removedUsers = append(removedUsers, NewUser(ctx, removedUser))
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnThreadMembersUpdateFuncType); ok {
			f(ctx, channel, addedUsers, removedUsers)
		}
	}
	return nil
}

type OnThreadMembersUpdateFuncType func(ctx *Context, thread Channel, addedUsers []User, removedUsers []User) (err error)

// OnGuildCreate
func OnGuildCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildCreatePayload discord.GuildCreate
	if err = ctx.decodeContent(payload, &guildCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildCreatePayload.ID)

	var lazy bool
	if _, err := ctx.decodeExtra(payload, "lazy", &lazy); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	// If true, the guild was previously unavailable.
	var unavailable bool
	if _, err := ctx.decodeExtra(payload, "unavailable", &unavailable); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	// Dispatches either a GUILD_JOIN or GUILD_AVAILABLE event.
	// Guilds that are lazy loaded are not handled.

	if unavailable {
		return ctx.Handlers.DispatchType(ctx, "GUILD_AVAILABLE", payload)
	} else if !lazy {
		return ctx.Handlers.DispatchType(ctx, "GUILD_JOIN", payload)
	}
	return nil
}

// OnGuildUpdate
func OnGuildUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildUpdatePayload discord.GuildUpdate
	if err = ctx.decodeContent(payload, &guildUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	guild := Guild(*guildUpdatePayload)
	ctx.Guild = &guild

	var beforeGuild discord.Guild
	if _, err := ctx.decodeExtra(payload, "before", &beforeGuild); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildUpdateFuncType); ok {
			f(ctx, Guild(beforeGuild), guild)
		}
	}
	return nil
}

type OnGuildUpdateFuncType func(ctx *Context, before Guild, after Guild) (err error)

// OnGuildDelete
func OnGuildDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildDeletePayload discord.GuildDelete
	if err = ctx.decodeContent(payload, &guildDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildDeletePayload.ID)

	if guildDeletePayload.Unavailable {
		return ctx.Handlers.DispatchType(ctx, "GUILD_UNAVAILABLE", payload)
	} else {
		return ctx.Handlers.DispatchType(ctx, "GUILD_REMOVE", payload)
	}
}

// OnGuildBanAdd
func OnGuildBanAdd(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildBanAddPayload discord.GuildBanAdd
	if err = ctx.decodeContent(payload, &guildBanAddPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildBanAddPayload.GuildID)

	user := User(*guildBanAddPayload.User)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildBanAddFuncType); ok {
			f(ctx, user)
		}
	}
	return nil
}

type OnGuildBanAddFuncType func(ctx *Context, user User) (err error)

// OnGuildBanRemove
func OnGuildBanRemove(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildBanRemovePayload discord.GuildBanRemove
	if err = ctx.decodeContent(payload, &guildBanRemovePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildBanRemovePayload.GuildID)

	user := User(*guildBanRemovePayload.User)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildBanRemoveFuncType); ok {
			f(ctx, user)
		}
	}
	return nil
}

type OnGuildBanRemoveFuncType func(ctx *Context, user User) (err error)

// OnGuildEmojisUpdate
func OnGuildEmojisUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildEmojisUpdatePayload discord.GuildEmojisUpdate
	if err = ctx.decodeContent(payload, &guildEmojisUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildEmojisUpdatePayload.GuildID)

	var before []Emoji
	if _, err := ctx.decodeExtra(payload, "before", &before); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	after := make([]Emoji, 0, len(guildEmojisUpdatePayload.Emojis))
	for _, emoji := range guildEmojisUpdatePayload.Emojis {
		after = append(after, Emoji(*emoji))
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildEmojisUpdateFuncType); ok {
			f(ctx, before, after)
		}
	}
	return nil
}

type OnGuildEmojisUpdateFuncType func(ctx *Context, before []Emoji, after []Emoji) (err error)

// OnGuildStickersUpdate
func OnGuildStickersUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildStickersUpdatePayload discord.GuildStickersUpdate
	if err = ctx.decodeContent(payload, &guildStickersUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildStickersUpdatePayload.GuildID)

	var before []Sticker
	if _, err := ctx.decodeExtra(payload, "before", &before); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	after := make([]Sticker, 0, len(guildStickersUpdatePayload.Stickers))
	for _, sticker := range guildStickersUpdatePayload.Stickers {
		after = append(after, Sticker(*sticker))
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildStickersUpdateFuncType); ok {
			f(ctx, before, after)
		}
	}
	return nil
}

type OnGuildStickersUpdateFuncType func(ctx *Context, before []Sticker, after []Sticker) (err error)

// OnGuildIntegrationsUpdate
func OnGuildIntegrationsUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildIntegrationsUpdatePayload discord.GuildIntegrationsUpdate
	if err = ctx.decodeContent(payload, &guildIntegrationsUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildIntegrationsUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildIntegrationsUpdateFuncType); ok {
			f(ctx)
		}
	}
	return nil
}

type OnGuildIntegrationsUpdateFuncType func(ctx *Context) (err error)

// OnGuildMemberAdd
func OnGuildMemberAdd(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildMemberAddPayload discord.GuildMemberAdd
	if err = ctx.decodeContent(payload, &guildMemberAddPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildMemberAddPayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberAddFuncType); ok {
			f(ctx, GuildMember(*guildMemberAddPayload.GuildMember))
		}
	}
	return nil
}

type OnGuildMemberAddFuncType func(ctx *Context, member GuildMember) (err error)

// OnGuildMemberRemove
func OnGuildMemberRemove(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildMemberRemovePayload discord.GuildMemberRemove
	if err = ctx.decodeContent(payload, &guildMemberRemovePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildMemberRemovePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberRemoveFuncType); ok {
			f(ctx, User(*guildMemberRemovePayload.User))
		}
	}
	return nil
}

type OnGuildMemberRemoveFuncType func(ctx *Context, member User) (err error)

// OnGuildMemberUpdate
func OnGuildMemberUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildMemberUpdatePayload discord.GuildMemberUpdate
	if err = ctx.decodeContent(payload, &guildMemberUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildMemberUpdatePayload.GuildID)

	var beforeGuildMember discord.GuildMember
	if _, err := ctx.decodeExtra(payload, "before", &beforeGuildMember); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildMemberUpdateFuncType); ok {
			f(ctx, GuildMember(beforeGuildMember), GuildMember(*guildMemberUpdatePayload.GuildMember))
		}
	}
	return nil
}

type OnGuildMemberUpdateFuncType func(ctx *Context, before GuildMember, after GuildMember) (err error)

// OnGuildRoleCreate
func OnGuildRoleCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildRoleCreatePayload discord.GuildRoleCreate
	if err = ctx.decodeContent(payload, &guildRoleCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildRoleCreatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleCreateFuncType); ok {
			f(ctx, Role(*guildRoleCreatePayload.Role))
		}
	}
	return nil
}

type OnGuildRoleCreateFuncType func(ctx *Context, role Role) (err error)

// OnGuildRoleUpdate
func OnGuildRoleUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildRoleUpdatePayload discord.GuildRoleUpdate
	if err = ctx.decodeContent(payload, &guildRoleUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildRoleUpdatePayload.GuildID)

	var beforeRole discord.Role
	if _, err := ctx.decodeExtra(payload, "before", &beforeRole); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleUpdateFuncType); ok {
			f(ctx, Role(beforeRole), Role(*guildRoleUpdatePayload.Role))
		}
	}
	return nil
}

type OnGuildRoleUpdateFuncType func(ctx *Context, before Role, after Role) (err error)

// OnGuildRoleDelete
func OnGuildRoleDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildRoleDeletePayload discord.GuildRoleDelete
	if err = ctx.decodeContent(payload, &guildRoleDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildRoleDeletePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildRoleDeleteFuncType); ok {
			f(ctx, guildRoleDeletePayload.RoleID)
		}
	}
	return nil
}

type OnGuildRoleDeleteFuncType func(ctx *Context, roleID discord.Snowflake) (err error)

// OnIntegrationCreate
func OnIntegrationCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var integrationCreatePayload discord.IntegrationCreate
	if err = ctx.decodeContent(payload, &integrationCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, integrationCreatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnIntegrationCreateFuncType); ok {
			f(ctx, Integration(*integrationCreatePayload.Integration))
		}
	}
	return nil
}

type OnIntegrationCreateFuncType func(ctx *Context, integration Integration) (err error)

// OnIntegrationUpdate
func OnIntegrationUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var integrationUpdatePayload discord.IntegrationUpdate
	if err = ctx.decodeContent(payload, &integrationUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, integrationUpdatePayload.GuildID)

	var beforeIntegration discord.Integration
	if _, err := ctx.decodeExtra(payload, "before", &beforeIntegration); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnIntegrationUpdateFuncType); ok {
			f(ctx, Integration(beforeIntegration), Integration(*integrationUpdatePayload.Integration))
		}
	}
	return nil
}

type OnIntegrationUpdateFuncType func(ctx *Context, before Integration, after Integration) (err error)

// OnIntegrationDelete
func OnIntegrationDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var integrationDeletePayload discord.IntegrationDelete
	if err = ctx.decodeContent(payload, &integrationDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, integrationDeletePayload.GuildID)

	var applicationID discord.Snowflake
	if integrationDeletePayload.ApplicationID != nil {
		applicationID = *integrationDeletePayload.ApplicationID
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnIntegrationDeleteFuncType); ok {
			f(ctx, integrationDeletePayload.ID, applicationID)
		}
	}
	return nil
}

type OnIntegrationDeleteFuncType func(ctx *Context, integrationID discord.Snowflake, applicationID discord.Snowflake) (err error)

// OnInteractionCreate
func OnInteractionCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var interactionCreatePayload discord.InteractionCreate
	if err = ctx.decodeContent(payload, &interactionCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, interactionCreatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnInteractionCreateFuncType); ok {
			f(ctx, Interaction(*interactionCreatePayload))
		}
	}
	return nil
}

type OnInteractionCreateFuncType func(ctx *Context, interaction Interaction) (err error)

// OnInviteCreate
func OnInviteCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var inviteCreatePayload discord.InviteCreate
	if err = ctx.decodeContent(payload, &inviteCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if inviteCreatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *inviteCreatePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnInviteCreateFuncType); ok {
			f(ctx, Invite(*inviteCreatePayload))
		}
	}
	return nil
}

type OnInviteCreateFuncType func(ctx *Context, invite Invite) (err error)

// OnInviteDelete
func OnInviteDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var inviteDeletePayload discord.InviteDelete
	if err = ctx.decodeContent(payload, &inviteDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if inviteDeletePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *inviteDeletePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnInviteDeleteFuncType); ok {
			f(ctx, Invite(*inviteDeletePayload))
		}
	}
	return nil
}

type OnInviteDeleteFuncType func(ctx *Context, invite Invite) (err error)

// OnMessageCreate
func OnMessageCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageCreatePayload discord.MessageCreate
	if err = ctx.decodeContent(payload, &messageCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageCreatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *messageCreatePayload.GuildID)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageCreateFuncType); ok {
			f(ctx, Message(*messageCreatePayload))
		}
	}
	return nil
}

type OnMessageCreateFuncType func(ctx *Context, message Message) (err error)

// OnMessageUpdate
func OnMessageUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageUpdatePayload discord.MessageUpdate
	if err = ctx.decodeContent(payload, &messageUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageUpdatePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *messageUpdatePayload.GuildID)
	}

	var beforeMessage discord.Message
	if _, err := ctx.decodeExtra(payload, "before", &beforeMessage); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageUpdateFuncType); ok {
			f(ctx, Message(beforeMessage), Message(*messageUpdatePayload))
		}
	}
	return nil
}

type OnMessageUpdateFuncType func(ctx *Context, before Message, after Message) (err error)

// OnMessageDelete
func OnMessageDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageDeletePayload discord.MessageDelete
	if err = ctx.decodeContent(payload, &messageDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageDeletePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *messageDeletePayload.GuildID)
	}

	channel := NewChannel(ctx, messageDeletePayload.ChannelID, messageDeletePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteFuncType); ok {
			f(ctx, channel, messageDeletePayload.ID)
		}
	}
	return nil
}

type OnMessageDeleteFuncType func(ctx *Context, channel Channel, messageID discord.Snowflake) (err error)

// OnMessageDeleteBulk
func OnMessageDeleteBulk(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageDeleteBulkPayload discord.MessageDeleteBulk
	if err = ctx.decodeContent(payload, &messageDeleteBulkPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageDeleteBulkPayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *messageDeleteBulkPayload.GuildID)
	}

	channel := NewChannel(ctx, messageDeleteBulkPayload.ChannelID, messageDeleteBulkPayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageDeleteBulkFuncType); ok {
			f(ctx, channel, messageDeleteBulkPayload.IDs)
		}
	}
	return nil
}

type OnMessageDeleteBulkFuncType func(ctx *Context, channel Channel, messageIDs []discord.Snowflake) (err error)

// OnMessageReactionAdd
func OnMessageReactionAdd(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageReactionAddPayload discord.MessageReactionAdd
	if err = ctx.decodeContent(payload, &messageReactionAddPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, messageReactionAddPayload.GuildID)

	channel := NewChannel(ctx, messageReactionAddPayload.ChannelID, &messageReactionAddPayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionAddFuncType); ok {
			f(ctx, channel, messageReactionAddPayload.MessageID, Emoji(*messageReactionAddPayload.Emoji), GuildMember(*messageReactionAddPayload.Member))
		}
	}
	return nil
}

type OnMessageReactionAddFuncType func(ctx *Context, channel Channel, messageID discord.Snowflake, emoji Emoji, guildMember GuildMember) (err error)

// OnMessageReactionRemove
func OnMessageReactionRemove(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageReactionRemovePayload discord.MessageReactionRemove
	if err = ctx.decodeContent(payload, &messageReactionRemovePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageReactionRemovePayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *messageReactionRemovePayload.GuildID)
	}

	channel := NewChannel(ctx, messageReactionRemovePayload.ChannelID, messageReactionRemovePayload.GuildID)
	user := NewUser(ctx, messageReactionRemovePayload.UserID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveFuncType); ok {
			f(ctx, channel, messageReactionRemovePayload.MessageID, Emoji(*messageReactionRemovePayload.Emoji), user)
		}
	}
	return nil
}

type OnMessageReactionRemoveFuncType func(ctx *Context, channel Channel, messageID discord.Snowflake, emoji Emoji, user User) (err error)

// OnMessageReactionRemoveAll
func OnMessageReactionRemoveAll(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageReactionRemoveAllPayload discord.MessageReactionRemoveAll
	if err = ctx.decodeContent(payload, &messageReactionRemoveAllPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, messageReactionRemoveAllPayload.GuildID)

	channel := NewChannel(ctx, messageReactionRemoveAllPayload.ChannelID, &messageReactionRemoveAllPayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveAllFuncType); ok {
			f(ctx, channel, messageReactionRemoveAllPayload.MessageID)
		}
	}
	return nil
}

type OnMessageReactionRemoveAllFuncType func(ctx *Context, channel Channel, messageID discord.Snowflake) (err error)

// OnMessageReactionRemoveEmoji
func OnMessageReactionRemoveEmoji(ctx *Context, payload structs.SandwichPayload) (err error) {
	var messageReactionRemoveEmojiPayload discord.MessageReactionRemoveEmoji
	if err = ctx.decodeContent(payload, &messageReactionRemoveEmojiPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if messageReactionRemoveEmojiPayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *messageReactionRemoveEmojiPayload.GuildID)
	}

	channel := NewChannel(ctx, messageReactionRemoveEmojiPayload.ChannelID, messageReactionRemoveEmojiPayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnMessageReactionRemoveEmojiFuncType); ok {
			f(ctx, channel, messageReactionRemoveEmojiPayload.MessageID, Emoji(*messageReactionRemoveEmojiPayload.Emoji))
		}
	}
	return nil
}

type OnMessageReactionRemoveEmojiFuncType func(ctx *Context, channel Channel, messageID discord.Snowflake, emoji Emoji) (err error)

// OnPresenceUpdate
func OnPresenceUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var presenceUpdatePayload discord.PresenceUpdate
	if err = ctx.decodeContent(payload, &presenceUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, presenceUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnPresenceUpdateFuncType); ok {
			f(ctx, User(*presenceUpdatePayload.User), presenceUpdatePayload)
		}
	}
	return nil
}

type OnPresenceUpdateFuncType func(ctx *Context, user User, payload discord.PresenceUpdate) (err error)

// OnStageInstanceCreate
func OnStageInstanceCreate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var stageInstanceCreatePayload discord.StageInstanceCreate
	if err = ctx.decodeContent(payload, &stageInstanceCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, stageInstanceCreatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceCreateFuncType); ok {
			f(ctx, StageInstance(*stageInstanceCreatePayload))
		}
	}
	return nil
}

type OnStageInstanceCreateFuncType func(ctx *Context, stage StageInstance) (err error)

// OnStageInstanceUpdate
func OnStageInstanceUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var stageInstanceUpdatePayload discord.StageInstanceUpdate
	if err = ctx.decodeContent(payload, &stageInstanceUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, stageInstanceUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceUpdateFuncType); ok {
			f(ctx, StageInstance(*stageInstanceUpdatePayload))
		}
	}
	return nil
}

type OnStageInstanceUpdateFuncType func(ctx *Context, stage StageInstance) (err error)

// OnStageInstanceDelete
func OnStageInstanceDelete(ctx *Context, payload structs.SandwichPayload) (err error) {
	var stageInstanceDeletePayload discord.StageInstanceDelete
	if err = ctx.decodeContent(payload, &stageInstanceDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnStageInstanceDeleteFuncType); ok {
			f(ctx, StageInstance(*stageInstanceDeletePayload))
		}
	}
	return nil
}

type OnStageInstanceDeleteFuncType func(ctx *Context, stage StageInstance) (err error)

// OnTypingStart
func OnTypingStart(ctx *Context, payload structs.SandwichPayload) (err error) {
	var typingStartPayload discord.TypingStart
	if err = ctx.decodeContent(payload, &typingStartPayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	if typingStartPayload.GuildID != nil {
		ctx.Guild = NewGuild(ctx, *typingStartPayload.GuildID)
	}

	timestamp := time.Unix(int64(typingStartPayload.Timestamp), 0)

	channel := NewChannel(ctx, typingStartPayload.ChannelID, typingStartPayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnTypingStartFuncType); ok {
			f(ctx, channel, GuildMember(*typingStartPayload.Member), timestamp)
		}
	}
	return nil
}

type OnTypingStartFuncType func(ctx *Context, channel Channel, member GuildMember, timestamp time.Time) (err error)

// OnUserUpdate
func OnUserUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var userUpdatePayload discord.UserUpdate
	if err = ctx.decodeContent(payload, &userUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	var beforeUser discord.User
	if _, err := ctx.decodeExtra(payload, "before", &beforeUser); err != nil {
		return xerrors.Errorf("Failed to unmarshal extra: %v", err)
	}

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnUserUpdateFuncType); ok {
			f(ctx, User(beforeUser), User(*userUpdatePayload))
		}
	}
	return nil
}

type OnUserUpdateFuncType func(ctx *Context, before User, after User) (err error)

// OnVoiceStateUpdate
func OnVoiceStateUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var voiceStateUpdatePayload discord.VoiceStateUpdate
	if err = ctx.decodeContent(payload, &voiceStateUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, voiceStateUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnVoiceStateUpdateFuncType); ok {
			f(ctx, GuildMember(*voiceStateUpdatePayload.Member))
		}
	}
	return nil
}

type OnVoiceStateUpdateFuncType func(ctx *Context, member GuildMember, voice VoiceState) (err error)

// OnVoiceServerUpdate
func OnVoiceServerUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var voiceServerUpdatePayload discord.VoiceServerUpdate
	if err = ctx.decodeContent(payload, &voiceServerUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, voiceServerUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnVoiceServerUpdateFuncType); ok {
			f(ctx, voiceServerUpdatePayload)
		}
	}
	return nil
}

type OnVoiceServerUpdateFuncType func(ctx *Context, payload discord.VoiceServerUpdate) (err error)

// OnWebhookUpdate
func OnWebhookUpdate(ctx *Context, payload structs.SandwichPayload) (err error) {
	var webhookUpdatePayload discord.WebhookUpdate
	if err = ctx.decodeContent(payload, &webhookUpdatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, webhookUpdatePayload.GuildID)

	channel := NewChannel(ctx, webhookUpdatePayload.ChannelID, &webhookUpdatePayload.GuildID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnWebhookUpdateFuncType); ok {
			f(ctx, channel)
		}
	}
	return nil
}

type OnWebhookUpdateFuncType func(ctx *Context, channel Channel) (err error)

// OnGuildJoin
func OnGuildJoin(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildCreatePayload discord.GuildCreate
	if err = ctx.decodeContent(payload, &guildCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	guild := Guild(*guildCreatePayload.Guild)
	ctx.Guild = &guild

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildJoinFuncType); ok {
			f(ctx, guild)
		}
	}
	return nil
}

type OnGuildJoinFuncType func(ctx *Context, guild Guild) (err error)

// OnGuildAvailable
func OnGuildAvailable(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildCreatePayload discord.GuildCreate
	if err = ctx.decodeContent(payload, &guildCreatePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	guild := Guild(*guildCreatePayload.Guild)
	ctx.Guild = &guild

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildJoinFuncType); ok {
			f(ctx, guild)
		}
	}
	return nil
}

type OnGuildAvailableFuncType func(ctx *Context, guild Guild) (err error)

// OnGuildLeave
func OnGuildLeave(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildDeletePayload discord.GuildDelete
	if err = ctx.decodeContent(payload, &guildDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildDeletePayload.ID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildLeaveFuncType); ok {
			f(ctx, UnavailableGuild(*guildDeletePayload))
		}
	}
	return nil
}

type OnGuildLeaveFuncType func(ctx *Context, unavailableGuild UnavailableGuild) (err error)

// OnGuildUnavailable
func OnGuildUnavailable(ctx *Context, payload structs.SandwichPayload) (err error) {
	var guildDeletePayload discord.GuildDelete
	if err = ctx.decodeContent(payload, &guildDeletePayload); err != nil {
		return xerrors.Errorf("Failed to unmarshal payload: %v", err)
	}

	ctx.Guild = NewGuild(ctx, guildDeletePayload.ID)

	for _, event := range ctx.EventHandler.Events {
		if f, ok := event.(OnGuildUnavailableFuncType); ok {
			f(ctx, UnavailableGuild(*guildDeletePayload))
		}
	}
	return nil
}

type OnGuildUnavailableFuncType func(ctx *Context, unavailableGuild UnavailableGuild) (err error)

func init() {
	h := NewHandlers()

	h.NewEventHandler("READY", OnReady)
	h.NewEventHandler("RESUMED", OnResumed)
	h.NewEventHandler("APPLICATION_COMMAND_CREATE", OnApplicationCommandCreate)
	h.NewEventHandler("APPLICATION_COMMAND_UPDATE", OnApplicationCommandUpdate)
	h.NewEventHandler("APPLICATION_COMMAND_DELETE", OnApplicationCommandDelete)
	h.NewEventHandler("CHANNEL_CREATE", OnChannelCreate)
	h.NewEventHandler("CHANNEL_UPDATE", OnChannelUpdate)
	h.NewEventHandler("CHANNEL_DELETE", OnChannelDelete)
	h.NewEventHandler("CHANNEL_PINS_UPDATE", OnChannelPinsUpdate)
	h.NewEventHandler("THREAD_CREATE", OnThreadCreate)
	h.NewEventHandler("THREAD_UPDATE", OnThreadUpdate)
	h.NewEventHandler("THREAD_DELETE", OnThreadDelete)
	// h.NewEventHandler("THREAD_LIST_SYNC", OnThreadListSync)
	h.NewEventHandler("THREAD_MEMBER_UPDATE", OnThreadMemberUpdate)
	h.NewEventHandler("THREAD_MEMBERS_UPDATE", OnThreadMembersUpdate)
	h.NewEventHandler("GUILD_CREATE", OnGuildCreate)
	h.NewEventHandler("GUILD_UPDATE", OnGuildUpdate)
	h.NewEventHandler("GUILD_DELETE", OnGuildDelete)
	h.NewEventHandler("GUILD_BAN_ADD", OnGuildBanAdd)
	h.NewEventHandler("GUILD_BAN_REMOVE", OnGuildBanRemove)
	h.NewEventHandler("GUILD_EMOJIS_UPDATE", OnGuildEmojisUpdate)
	h.NewEventHandler("GUILD_STICKERS_UPDATE", OnGuildStickersUpdate)
	h.NewEventHandler("GUILD_INTEGRATIONS_UPDATE", OnGuildIntegrationsUpdate)
	h.NewEventHandler("GUILD_MEMBER_ADD", OnGuildMemberAdd)
	h.NewEventHandler("GUILD_MEMBER_REMOVE", OnGuildMemberRemove)
	h.NewEventHandler("GUILD_MEMBER_UPDATE", OnGuildMemberUpdate)
	h.NewEventHandler("GUILD_ROLE_CREATE", OnGuildRoleCreate)
	h.NewEventHandler("GUILD_ROLE_UPDATE", OnGuildRoleUpdate)
	h.NewEventHandler("GUILD_ROLE_DELETE", OnGuildRoleDelete)
	h.NewEventHandler("INTEGRATION_CREATE", OnIntegrationCreate)
	h.NewEventHandler("INTEGRATION_UPDATE", OnIntegrationUpdate)
	h.NewEventHandler("INTEGRATION_DELETE", OnIntegrationDelete)
	h.NewEventHandler("INTERACTION_CREATE", OnInteractionCreate)
	h.NewEventHandler("INVITE_CREATE", OnInviteCreate)
	h.NewEventHandler("INVITE_DELETE", OnInviteDelete)
	h.NewEventHandler("MESSAGE_CREATE", OnMessageCreate)
	h.NewEventHandler("MESSAGE_UPDATE", OnMessageUpdate)
	h.NewEventHandler("MESSAGE_DELETE", OnMessageDelete)
	h.NewEventHandler("MESSAGE_DELETE_BULK", OnMessageDeleteBulk)
	h.NewEventHandler("MESSAGE_REACTION_ADD", OnMessageReactionAdd)
	h.NewEventHandler("MESSAGE_REACTION_REMOVE", OnMessageReactionRemove)
	h.NewEventHandler("MESSAGE_REACTION_REMOVE_ALL", OnMessageReactionRemoveAll)
	h.NewEventHandler("MESSAGE_REACTION_REMOVE_EMOJI", OnMessageReactionRemoveEmoji)
	h.NewEventHandler("PRESENCE_UPDATE", OnPresenceUpdate)
	h.NewEventHandler("STAGE_INSTANCE_CREATE", OnStageInstanceCreate)
	h.NewEventHandler("STAGE_INSTANCE_UPDATE", OnStageInstanceUpdate)
	h.NewEventHandler("STAGE_INSTANCE_DELETE", OnStageInstanceDelete)
	h.NewEventHandler("TYPING_START", OnTypingStart)
	h.NewEventHandler("USER_UPDATE", OnUserUpdate)
	h.NewEventHandler("VOICE_STATE_UPDATE", OnVoiceStateUpdate)
	h.NewEventHandler("VOICE_SERVER_UPDATE", OnVoiceServerUpdate)
	h.NewEventHandler("WEBHOOKS_UPDATE", OnWebhookUpdate)

	// Custom Events
	h.NewEventHandler("GUILD_JOIN", OnGuildJoin)
	h.NewEventHandler("GUILD_AVAILABLE", OnGuildAvailable)

	h.NewEventHandler("GUILD_LEAVE", OnGuildLeave)
	h.NewEventHandler("GUILD_UNAVAILABLE", OnGuildUnavailable)

}
