package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	protobuf "github.com/WelcomerTeam/Sandwich-Daemon/protobuf"
	"github.com/WelcomerTeam/Sandwich-Daemon/structs"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

type Sandwich struct {
	botsMu sync.RWMutex
	Bots   map[string]*Bot

	identifiersMu sync.RWMutex
	Identifiers   map[string]*Identifier

	sandwichClient protobuf.SandwichClient
}

func NewSandwich(ctx context.Context, conn grpc.ClientConnInterface) (s *Sandwich) {
	s = &Sandwich{
		botsMu: sync.RWMutex{},
		Bots:   make(map[string]*Bot),

		identifiersMu: sync.RWMutex{},
		Identifiers:   make(map[string]*Identifier),

		sandwichClient: protobuf.NewSandwichClient(conn),
	}

	return
}

func (s *Sandwich) DispatchSandwichPayload(payload structs.SandwichPayload) (err error) {
	s.botsMu.RLock()
	b, ok := s.Bots[payload.Metadata.Identifier]
	s.botsMu.RUnlock()

	if !ok {
		return xerrors.New("No identifier")
	}

	return b.Dispatch(&Context{
		Sandwich: s,
	}, payload)
}

func (s *Sandwich) RegisterBot(identifier string, bot *Bot) {
	s.botsMu.Lock()
	s.Bots[identifier] = bot
	s.botsMu.Unlock()
}

func (s *Sandwich) FetchIdentifiers(ctx context.Context) (identifiers *Identifiers, err error) {
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
	Sandwich *Sandwich

	// Filled in on dispatch
	EventHandler *EventHandler
	Handlers     *Handlers

	Guild *Guild
}

func (ctx *Context) wrapFuncType(err error) {
	if err != nil {
		fmt.Printf("We errored: %v", err)
	}
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
