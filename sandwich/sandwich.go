package internal

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"

	discord "github.com/WelcomerTeam/Discord/discord"
	protobuf "github.com/WelcomerTeam/Sandwich-Daemon/protobuf"
	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// VERSION follows semantic versioning.
const VERSION = "0.2.1"

var LastRequestTimeout = time.Minute * 60

type Sandwich struct {
	Logger zerolog.Logger

	botsMu sync.RWMutex
	Bots   map[string]*Bot

	SandwichEvents *Handlers

	identifiersMu sync.RWMutex
	Identifiers   map[string]*sandwich_structs.ManagerConsumerConfiguration

	lastIdentifierRequestMu sync.RWMutex
	LastIdentifierRequest   map[string]time.Time

	RESTInterface discord.RESTInterface

	SandwichClient protobuf.SandwichClient
	GRPCInterface  GRPC
}

func NewSandwich(conn grpc.ClientConnInterface, restInterface discord.RESTInterface, logger io.Writer) (s *Sandwich) {
	s = &Sandwich{
		Logger: zerolog.New(logger).With().Timestamp().Logger(),

		botsMu: sync.RWMutex{},
		Bots:   make(map[string]*Bot),

		SandwichEvents: NewSandwichHandlers(),

		identifiersMu: sync.RWMutex{},
		Identifiers:   make(map[string]*sandwich_structs.ManagerConsumerConfiguration),

		lastIdentifierRequestMu: sync.RWMutex{},
		LastIdentifierRequest:   make(map[string]time.Time),

		RESTInterface: restInterface,

		SandwichClient: protobuf.NewSandwichClient(conn),
		GRPCInterface:  NewDefaultGRPCClient(),
	}

	return
}

func (s *Sandwich) Close() (err error) {
	wg := &sync.WaitGroup{}

	s.botsMu.RLock()
	for _, bot := range s.Bots {
		wg.Add(1)

		bot.Close(wg)
	}
	s.botsMu.RUnlock()

	return
}

func (s *Sandwich) DispatchGRPCPayload(context context.Context, payload sandwich_structs.SandwichPayload) (err error) {
	logger := s.Logger.With().Str("application", payload.Metadata.Application).Logger()

	return s.SandwichEvents.Dispatch(&EventContext{
		Logger:   logger,
		Sandwich: s,
		Session:  discord.NewSession(context, "", s.RESTInterface, logger),
		Handlers: s.SandwichEvents,
		Context:  context,
		payload:  &payload,
	}, payload)
}

func (s *Sandwich) DispatchSandwichPayload(context context.Context, payload sandwich_structs.SandwichPayload) (err error) {
	s.botsMu.RLock()
	bot, ok := s.Bots[payload.Metadata.Identifier]
	s.botsMu.RUnlock()

	if !ok {
		s.Logger.Debug().
			Str("identifier", payload.Metadata.Identifier).
			Str("application", payload.Metadata.Application).
			Msg(ErrInvalidIdentifier.Error())

		return ErrInvalidIdentifier
	}

	logger := s.Logger.With().Str("application", payload.Metadata.Application).Logger()

	return bot.Dispatch(&EventContext{
		Logger:   logger,
		Sandwich: s,
		Session:  discord.NewSession(context, "", s.RESTInterface, logger),
		Handlers: bot.Handlers,
		Context:  context,
		payload:  &payload,
	}, payload)
}

func (s *Sandwich) RegisterBot(identifier string, bot *Bot) {
	s.botsMu.Lock()
	s.Bots[identifier] = bot
	s.botsMu.Unlock()
}

func (s *Sandwich) RecoverEventPanic(errorValue interface{}, eventCtx *EventContext, payload *sandwich_structs.SandwichPayload) {
	s.Logger.Error().Interface("errorValue", errorValue).Str("type", payload.Type).Msg("Recovered panic on event dispatch")

	fmt.Println(string(debug.Stack()))
}

func (s *Sandwich) FetchIdentifier(context context.Context, applicationName string) (identifier *sandwich_structs.ManagerConsumerConfiguration, ok bool, err error) {
	s.identifiersMu.RLock()
	identifier, ok = s.Identifiers[applicationName]
	s.identifiersMu.RUnlock()

	if ok {
		return identifier, true, nil
	}

	s.lastIdentifierRequestMu.RLock()
	lastRequest, ok := s.LastIdentifierRequest[applicationName]
	s.lastIdentifierRequestMu.RUnlock()

	if !ok || (ok && time.Now().Add(LastRequestTimeout).Before(lastRequest)) {
		identifiers, err := s.GRPCInterface.FetchConsumerConfiguration((&EventContext{
			Sandwich: s,
			Session:  discord.NewSession(context, "", s.RESTInterface, s.Logger),
			Context:  context,
		}).ToGRPCContext(), "")
		if err != nil {
			return nil, false, errors.Errorf("Failed to fetch consumer configuration: %v", err)
		}

		s.identifiersMu.Lock()
		s.Identifiers = map[string]*sandwich_structs.ManagerConsumerConfiguration{}

		for k := range identifiers.Identifiers {
			v := identifiers.Identifiers[k]
			s.Identifiers[k] = &v
		}
		s.identifiersMu.Unlock()

		identifier, ok = s.Identifiers[applicationName]
		if !ok {
			return nil, false, ErrInvalidApplication
		}
	} else {
		return nil, false, ErrInvalidApplication
	}

	return identifier, true, nil
}

// EventContext is extra data passed to event handlers.
// This is not the same as a command's context.
type EventContext struct {
	context.Context

	Logger zerolog.Logger

	Sandwich *Sandwich

	Session *discord.Session

	// Filled in on dispatch
	EventHandler *EventHandler
	Handlers     *Handlers

	Identifier *sandwich_structs.ManagerConsumerConfiguration

	Guild *discord.Guild

	payload *sandwich_structs.SandwichPayload
}

func (eventCtx *EventContext) ToGRPCContext() *GRPCContext {
	return &GRPCContext{
		Context:        eventCtx.Context,
		Logger:         eventCtx.Logger,
		SandwichClient: eventCtx.Sandwich.SandwichClient,
		GRPCInterface:  eventCtx.Sandwich.GRPCInterface,
		Session:        eventCtx.Session,
		Identifier:     eventCtx.Identifier,
	}
}

func (eventCtx *EventContext) Trace() sandwich_structs.SandwichTrace {
	if eventCtx.payload != nil {
		return eventCtx.payload.Trace
	}

	return nil
}

func (eventCtx *EventContext) decodeContent(msg sandwich_structs.SandwichPayload, out interface{}) (err error) {
	err = jsoniter.Unmarshal(msg.Data, &out)
	if err != nil {
		return errors.Errorf("Failed to unmarshal gateway payload: %v", err)
	}

	return
}

func (eventCtx *EventContext) decodeExtra(msg sandwich_structs.SandwichPayload, key string, out interface{}) (ok bool, err error) {
	val, ok := msg.Extra[key]
	if ok {
		if len(val) == 0 {
			ok = false

			return
		}

		err = jsoniter.Unmarshal(val, &out)
		if err != nil {
			return ok, errors.Errorf("Failed to unmarshal extra: %v", err)
		}
	}

	return
}
