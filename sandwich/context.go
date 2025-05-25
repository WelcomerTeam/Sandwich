package internal

import (
	"context"

	"github.com/WelcomerTeam/Discord/discord"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
	"github.com/rs/zerolog"
)

type GRPCContext struct {
	context.Context

	Logger zerolog.Logger

	SandwichClient sandwich_protobuf.SandwichClient

	Session    *discord.Session
	Identifier *sandwich_protobuf.SandwichApplication
}
