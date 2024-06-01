package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	discord "github.com/WelcomerTeam/Discord/discord"
	protobuf "github.com/WelcomerTeam/Sandwich-Daemon/protobuf"
	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// VERSION follows semantic versioning.
const VERSION = "0.5"

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

	ErrorOnInvalidIdentifier bool
}

func NewSandwich(conn grpc.ClientConnInterface, restInterface discord.RESTInterface, logger io.Writer) *Sandwich {
	sandwich := &Sandwich{
		Logger: zerolog.New(logger).With().Timestamp().Logger(),

		botsMu: sync.RWMutex{},
		Bots:   make(map[string]*Bot),

		SandwichEvents: newSandwichHandlers(),

		identifiersMu: sync.RWMutex{},
		Identifiers:   make(map[string]*sandwich_structs.ManagerConsumerConfiguration),

		lastIdentifierRequestMu: sync.RWMutex{},
		LastIdentifierRequest:   make(map[string]time.Time),

		RESTInterface: restInterface,

		SandwichClient: protobuf.NewSandwichClient(conn),
		GRPCInterface:  NewDefaultGRPCClient(),

		ErrorOnInvalidIdentifier: false,
	}

	return sandwich
}

func (sandwich *Sandwich) SetErrorOnInvalidIdentifier(value bool) {
	sandwich.ErrorOnInvalidIdentifier = value
}

func (sandwich *Sandwich) ListenToChannel(ctx context.Context, channel chan []byte) error {

	// Signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Register message channels
	grpcMessages := make(chan *protobuf.ListenResponse)

	go func() {
		for {
			grpcListener, err := sandwich.SandwichClient.Listen(ctx, &protobuf.ListenRequest{
				Identifier: "",
			})
			if err != nil {
				sandwich.Logger.Warn().Err(err).Msg("Failed to listen to grpc")

				time.Sleep(time.Second)
			} else {
				for {
					var listenResponse protobuf.ListenResponse

					err = grpcListener.RecvMsg(&listenResponse)
					if err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}

						sandwich.Logger.Warn().Err(err).Msg("Failed to receive grpc message")

						break
					} else {
						grpcMessages <- &listenResponse
					}
				}
			}
		}
	}()

	// Event Loop
eventLoop:
	for {
		select {
		case grpcMessage := <-grpcMessages:
			var payload sandwich_structs.SandwichPayload

			err := json.Unmarshal(grpcMessage.Data, &payload)
			if err != nil {
				sandwich.Logger.Warn().Err(err).Msg("Failed to unmarshal grpc message")
			} else {
				sandwich.DispatchGRPCPayload(ctx, payload)
			}
		case stanMessage := <-channel:
			var payload sandwich_structs.SandwichPayload

			err := json.Unmarshal(stanMessage, &payload)
			if err != nil {
				sandwich.Logger.Warn().Err(err).Msg("Failed to unmarshal stan message")
			} else {
				err = sandwich.DispatchSandwichPayload(ctx, payload)
				if err != nil {
					sandwich.Logger.Warn().Err(err).Msg("Failed to dispatch sandwich payload")
				}
			}
		case <-signalCh:
			break eventLoop
		}
	}

	return nil
}

func (sandwich *Sandwich) DispatchGRPCPayload(ctx context.Context, payload sandwich_structs.SandwichPayload) {
	logger := sandwich.Logger.With().Str("application", payload.Metadata.Application).Logger()

	sandwich.SandwichEvents.Dispatch(&EventContext{
		Logger:   logger,
		Sandwich: sandwich,
		Session:  discord.NewSession(ctx, "", sandwich.RESTInterface),
		Handlers: sandwich.SandwichEvents,
		Context:  ctx,
		Payload:  &payload,
	}, payload)
}

func (sandwich *Sandwich) DispatchSandwichPayload(ctx context.Context, payload sandwich_structs.SandwichPayload) error {
	sandwich.botsMu.RLock()
	bot, ok := sandwich.Bots[payload.Metadata.Identifier]
	sandwich.botsMu.RUnlock()

	if !ok {
		if !sandwich.ErrorOnInvalidIdentifier {
			return nil
		} else {
			sandwich.Logger.Debug().
				Str("identifier", payload.Metadata.Identifier).
				Str("application", payload.Metadata.Application).
				Msg(ErrInvalidIdentifier.Error())

			return ErrInvalidIdentifier
		}
	}

	logger := sandwich.Logger.With().Str("application", payload.Metadata.Application).Logger()

	bot.Dispatch(&EventContext{
		Logger:   logger,
		Sandwich: sandwich,
		Session:  discord.NewSession(ctx, "", sandwich.RESTInterface),
		Handlers: bot.Handlers,
		Context:  ctx,
		Payload:  &payload,
	}, payload)

	return nil
}

func (sandwich *Sandwich) RegisterBot(identifier string, bot *Bot) {
	sandwich.botsMu.Lock()
	sandwich.Bots[identifier] = bot
	sandwich.botsMu.Unlock()
}

func (sandwich *Sandwich) RecoverEventPanic(errorValue interface{}, eventCtx *EventContext, payload *sandwich_structs.SandwichPayload) {
	sandwich.Logger.Error().Interface("errorValue", errorValue).Str("type", payload.Type).Msg("Recovered panic on event dispatch")

	fmt.Println(string(debug.Stack()))
}

func (sandwich *Sandwich) FetchIdentifier(ctx context.Context, applicationName string) (identifier *sandwich_structs.ManagerConsumerConfiguration, ok bool, err error) {
	sandwich.identifiersMu.RLock()
	identifier, ok = sandwich.Identifiers[applicationName]
	sandwich.identifiersMu.RUnlock()

	if ok {
		return identifier, true, nil
	}

	sandwich.lastIdentifierRequestMu.RLock()
	lastRequest, ok := sandwich.LastIdentifierRequest[applicationName]
	sandwich.lastIdentifierRequestMu.RUnlock()

	if !ok || (ok && time.Now().Add(LastRequestTimeout).Before(lastRequest)) {
		identifiers, err := sandwich.GRPCInterface.FetchConsumerConfiguration((&EventContext{
			Sandwich: sandwich,
			Session:  discord.NewSession(ctx, "", sandwich.RESTInterface),
			Context:  ctx,
		}).ToGRPCContext(), "")
		if err != nil {
			return nil, false, fmt.Errorf("failed to fetch consumer configuration: %w", err)
		}

		sandwich.identifiersMu.Lock()
		sandwich.Identifiers = map[string]*sandwich_structs.ManagerConsumerConfiguration{}

		for k := range identifiers.Identifiers {
			v := identifiers.Identifiers[k]
			sandwich.Identifiers[k] = &v
		}
		sandwich.identifiersMu.Unlock()

		identifier, ok = sandwich.Identifiers[applicationName]
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

	Payload *sandwich_structs.SandwichPayload
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
	if eventCtx.Payload != nil {
		return eventCtx.Payload.Trace
	}

	return nil
}

func (eventCtx *EventContext) DecodeContent(msg sandwich_structs.SandwichPayload, out interface{}) error {
	err := json.Unmarshal(msg.Data, &out)
	if err != nil {
		return errors.Errorf("failed to unmarshal gateway payload: %v", err)
	}

	return nil
}

func (eventCtx *EventContext) DecodeExtra(msg sandwich_structs.SandwichPayload, key string, out interface{}) (ok bool, err error) {
	val, ok := msg.Extra[key]
	if ok {
		if len(val) == 0 {
			ok = false

			return
		}

		err = json.Unmarshal(val, &out)
		if err != nil {
			return ok, fmt.Errorf("failed to unmarshal extra: %w", err)
		}
	}

	return
}
