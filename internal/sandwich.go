package internal

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"

	protobuf "github.com/WelcomerTeam/Sandwich-Daemon/protobuf"
	"github.com/WelcomerTeam/Sandwich-Daemon/structs"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

var (
	LastRequestTimeout = time.Minute * 60
)

type Sandwich struct {
	Logger zerolog.Logger

	botsMu sync.RWMutex
	Bots   map[string]*Bot

	SandwichEvents *Handlers

	identifiersMu sync.RWMutex
	Identifiers   map[string]*Identifier

	lastIdentifierRequestMu sync.RWMutex
	LastIdentifierRequest   map[string]time.Time

	sandwichClient protobuf.SandwichClient
}

func NewSandwich(conn grpc.ClientConnInterface, logger io.Writer) (s *Sandwich) {
	s = &Sandwich{
		Logger: zerolog.New(logger).With().Timestamp().Logger(),

		botsMu: sync.RWMutex{},
		Bots:   make(map[string]*Bot),

		SandwichEvents: NewSandwichHandlers(),

		identifiersMu: sync.RWMutex{},
		Identifiers:   make(map[string]*Identifier),

		lastIdentifierRequestMu: sync.RWMutex{},
		LastIdentifierRequest:   make(map[string]time.Time),

		sandwichClient: protobuf.NewSandwichClient(conn),
	}

	return
}

func (s *Sandwich) DispatchGRPCPayload(payload structs.SandwichPayload) (err error) {
	return s.SandwichEvents.Dispatch(&Context{
		Logger:   s.Logger.With().Str("application", payload.Metadata.Application).Logger(),
		Sandwich: s,
		Handlers: s.SandwichEvents,
	}, payload)
}

func (s *Sandwich) DispatchSandwichPayload(payload structs.SandwichPayload) (err error) {
	s.botsMu.RLock()
	b, ok := s.Bots[payload.Metadata.Identifier]
	s.botsMu.RUnlock()

	if !ok {
		return ErrInvalidIdentifier
	}

	return b.Dispatch(&Context{
		Logger:   s.Logger.With().Str("application", payload.Metadata.Application).Logger(),
		Sandwich: s,
		Handlers: b.Handlers,
	}, payload)
}

func (s *Sandwich) RegisterBot(identifier string, bot *Bot) {
	s.botsMu.Lock()
	s.Bots[identifier] = bot
	s.botsMu.Unlock()
}

func (s *Sandwich) FetchIdentifier(ctx context.Context, applicationName string) (identifier *Identifier, ok bool, err error) {
	s.identifiersMu.RLock()
	identifier, ok = s.Identifiers[applicationName]
	s.identifiersMu.RUnlock()

	if !ok {
		s.lastIdentifierRequestMu.RLock()
		lastRequest, ok := s.LastIdentifierRequest[applicationName]
		s.lastIdentifierRequestMu.RUnlock()

		if !ok || (ok && time.Now().Add(LastRequestTimeout).Before(lastRequest)) {
			identifiers, err := s.FetchAllIdentifiers(ctx)
			if err != nil {
				return nil, false, err
			}

			s.identifiersMu.Lock()
			s.Identifiers = map[string]*Identifier{}

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
	}

	return identifier, true, nil
}

func (s *Sandwich) FetchAllIdentifiers(ctx context.Context) (identifiers *Identifiers, err error) {
	res, err := s.sandwichClient.FetchConsumerConfiguration(ctx, &protobuf.FetchConsumerConfigurationRequest{})
	if err != nil {
		return identifiers, xerrors.Errorf("Failed to fetch consumer configuration: %v", err)
	}

	identifiers = &Identifiers{}

	err = json.Unmarshal(res.File, &identifiers)
	if err != nil {
		return identifiers, xerrors.Errorf("Failed to unmarshal consumer configuration: %v", err)
	}

	return identifiers, nil
}

type Identifiers struct {
	Version     string                `json:"v"`
	Identifiers map[string]Identifier `json:"identifiers"`
}

type Identifier struct {
	Token string `json:"token"`
	ID    int64  `json:"id"`
	User  User   `json:"user"`
}

// Context is extra data passed to event handlers.
// This is not the same as a command's context.
type Context struct {
	context.Context

	Logger zerolog.Logger

	Sandwich *Sandwich

	// Filled in on dispatch
	EventHandler *EventHandler
	Handlers     *Handlers

	Identifier *Identifier

	Guild *Guild
}

func (ctx *Context) decodeContent(msg structs.SandwichPayload, out interface{}) (err error) {
	err = jsoniter.Unmarshal(msg.Data, &out)
	if err != nil {
		return xerrors.Errorf("Failed to unmarshal gateway payload: %v", err)
	}

	return
}

func (ctx *Context) decodeExtra(msg structs.SandwichPayload, key string, out interface{}) (ok bool, err error) {
	val, ok := msg.Extra[key]
	if ok {
		if len(val) == 0 {
			ok = false

			return
		}

		err = jsoniter.Unmarshal(val, &out)
		if err != nil {
			return ok, xerrors.Errorf("Failed to unmarshal extra: %v", err)
		}
	}

	return
}
