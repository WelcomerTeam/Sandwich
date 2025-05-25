package internal

import (
	"context"
	"log/slog"

	"github.com/WelcomerTeam/Discord/discord"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
)

type GRPCContext struct {
	context.Context

	Logger *slog.Logger

	SandwichClient sandwich_protobuf.SandwichClient

	Session    *discord.Session
	Identifier *sandwich_protobuf.SandwichApplication
}
