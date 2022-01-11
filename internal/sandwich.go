package internal

import (
	"context"
	"encoding/json"
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

func NewSandwich(ctx context.Context, conn *grpc.ClientConn) (s *Sandwich) {
	s = &Sandwich{
		botsMu: sync.RWMutex{},
		Bots:   make(map[string]*Bot),

		identifiersMu: sync.RWMutex{},
		Identifiers:   make(map[string]*Identifier),

		sandwichClient: protobuf.NewSandwichClient(conn),
	}

	return
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
		err = jsoniter.Unmarshal(val, &out)
	}

	return
}
