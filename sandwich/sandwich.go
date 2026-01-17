package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	discord "github.com/WelcomerTeam/Discord/discord"
	sandwich_daemon "github.com/WelcomerTeam/Sandwich-Daemon"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// VERSION follows semantic versioning.
const VERSION = "1.1.0"

var LastRequestTimeout = time.Minute * 60

type Sandwich struct {
	Logger *slog.Logger

	botsMu sync.RWMutex
	Bots   map[string]*Bot

	SandwichEvents *Handlers

	identifiersMu sync.RWMutex
	Identifiers   map[string]*sandwich_protobuf.SandwichApplication

	lastIdentifierRequestMu sync.RWMutex
	LastIdentifierRequest   map[string]time.Time

	RESTInterface discord.RESTInterface

	SandwichClient sandwich_protobuf.SandwichClient

	ErrorOnInvalidIdentifier bool
}

func NewSandwich(conn grpc.ClientConnInterface, restInterface discord.RESTInterface, logger io.Writer) *Sandwich {
	sandwich := &Sandwich{
		Logger: slog.New(slog.NewTextHandler(logger, nil)),

		botsMu: sync.RWMutex{},
		Bots:   make(map[string]*Bot),

		SandwichEvents: newSandwichHandlers(),

		identifiersMu: sync.RWMutex{},
		Identifiers:   make(map[string]*sandwich_protobuf.SandwichApplication),

		lastIdentifierRequestMu: sync.RWMutex{},
		LastIdentifierRequest:   make(map[string]time.Time),

		RESTInterface: restInterface,

		SandwichClient: sandwich_protobuf.NewSandwichClient(conn),

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
	grpcMessages := make(chan *sandwich_protobuf.ListenResponse)

	go func() {
		for {
			grpcListener, err := sandwich.SandwichClient.Listen(ctx, &sandwich_protobuf.ListenRequest{
				Identifier: "",
			})
			if err != nil {
				sandwich.Logger.Warn("Failed to listen to grpc", "error", err)

				time.Sleep(time.Second)
			} else {
				for {
					var listenResponse sandwich_protobuf.ListenResponse

					err = grpcListener.RecvMsg(&listenResponse)
					if err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}

						sandwich.Logger.Warn("Failed to receive grpc message", "error", err)

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
			var payload sandwich_daemon.ProducedPayload

			err := json.Unmarshal(grpcMessage.Data, &payload)
			if err != nil {
				sandwich.Logger.Warn("Failed to unmarshal grpc message", "error", err)
			} else {
				sandwich.DispatchGRPCPayload(ctx, payload)
			}
		case stanMessage := <-channel:
			var payload sandwich_daemon.ProducedPayload

			err := json.Unmarshal(stanMessage, &payload)
			if err != nil {
				sandwich.Logger.Warn("Failed to unmarshal stan message", "error", err)
			} else {
				err = sandwich.DispatchProducedPayload(ctx, payload)
				if err != nil {
					sandwich.Logger.Warn("Failed to dispatch sandwich payload", "error", err)
				}
			}
		case <-signalCh:
			break eventLoop
		}
	}

	return nil
}

func (sandwich *Sandwich) DispatchGRPCPayload(ctx context.Context, payload sandwich_daemon.ProducedPayload) {
	logger := sandwich.Logger.With("application", payload.Metadata.Application)

	sandwich.SandwichEvents.Dispatch(&EventContext{
		Logger:   logger,
		Sandwich: sandwich,
		Session:  discord.NewSession("", sandwich.RESTInterface),
		Handlers: sandwich.SandwichEvents,
		Context:  ctx,
		Payload:  &payload,
	}, payload)
}

func (sandwich *Sandwich) DispatchProducedPayload(ctx context.Context, payload sandwich_daemon.ProducedPayload) error {
	sandwich.botsMu.RLock()
	bot, ok := sandwich.Bots[payload.Metadata.Identifier]
	sandwich.botsMu.RUnlock()

	if !ok {
		if !sandwich.ErrorOnInvalidIdentifier {
			return nil
		} else {
			sandwich.Logger.Debug("Invalid identifier",
				"identifier", payload.Metadata.Identifier,
				"application", payload.Metadata.Application,
				"error", ErrInvalidIdentifier)

			return ErrInvalidIdentifier
		}
	}

	logger := sandwich.Logger.With("application", payload.Metadata.Application)

	bot.Dispatch(&EventContext{
		Logger:   logger,
		Sandwich: sandwich,
		Session:  discord.NewSession("", sandwich.RESTInterface),
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

func (sandwich *Sandwich) RecoverEventPanic(errorValue interface{}, eventCtx *EventContext, payload *sandwich_daemon.ProducedPayload) {
	sandwich.Logger.Error("Recovered panic on event dispatch",
		"errorValue", errorValue,
		"type", payload.Type)

	stackTrace := debug.Stack()
	println(string(stackTrace))
}

func (sandwich *Sandwich) FetchIdentifier(ctx context.Context, applicationName string) (identifier *sandwich_protobuf.SandwichApplication, ok bool, err error) {
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
		identifiers, err := sandwich.SandwichClient.FetchApplication((&EventContext{
			Sandwich: sandwich,
			Session:  discord.NewSession("", sandwich.RESTInterface),
			Context:  ctx,
		}).ToGRPCContext(), &sandwich_protobuf.ApplicationIdentifier{})
		if err != nil {
			return nil, false, fmt.Errorf("failed to fetch consumer configuration: %w", err)
		}

		sandwich.identifiersMu.Lock()

		sandwich.Identifiers = map[string]*sandwich_protobuf.SandwichApplication{}

		for k := range identifiers.Applications {
			v := identifiers.Applications[k]
			sandwich.Identifiers[k] = v
		}

		identifier, ok = sandwich.Identifiers[applicationName]

		sandwich.identifiersMu.Unlock()

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

	Logger *slog.Logger

	Sandwich *Sandwich

	Session *discord.Session

	// Filled in on dispatch
	EventHandler *EventHandler
	Handlers     *Handlers

	Identifier *sandwich_protobuf.SandwichApplication

	Guild *discord.Guild

	Payload *sandwich_daemon.ProducedPayload
}

func (eventCtx *EventContext) ToGRPCContext() *GRPCContext {
	return &GRPCContext{
		Context:        eventCtx.Context,
		Logger:         eventCtx.Logger,
		SandwichClient: eventCtx.Sandwich.SandwichClient,
		Session:        eventCtx.Session,
		Identifier:     eventCtx.Identifier,
	}
}

func (eventCtx *EventContext) Trace() sandwich_daemon.Trace {
	if eventCtx.Payload != nil {
		return eventCtx.Payload.Trace
	}

	return nil
}

func (eventCtx *EventContext) DecodeContent(msg sandwich_daemon.ProducedPayload, out interface{}) error {
	err := json.Unmarshal(msg.Data, &out)
	if err != nil {
		return errors.Errorf("failed to unmarshal gateway payload: %v", err)
	}

	return nil
}

func (eventCtx *EventContext) DecodeExtra(msg sandwich_daemon.ProducedPayload, key string, out interface{}) (ok bool, err error) {
	valBytes, ok := msg.Extra[key]
	if !ok {
		return false, nil
	}

	if len(valBytes) == 0 {
		return false, nil
	}

	err = json.Unmarshal(valBytes, &out)
	if err != nil {
		return true, fmt.Errorf("failed to unmarshal extra: %w", err)
	}

	return true, nil
}
