package internal

import "github.com/pkg/errors"

var (
	ErrInvalidIdentifier  = errors.New("payload does not include a valid identifier")
	ErrInvalidApplication = errors.New("could not find identifier matching application")
	ErrInvalidToken       = errors.New("invalid token was passed")
	ErrUnknownEvent       = errors.New("event type does not have a handler")
	ErrUnknownGRPCError   = errors.New("grpc returned unknown error")

	ErrCogAlreadyRegistered = errors.New("cog with this name already exists")

	ErrFetchMissingGuild     = errors.New("object requires guild ID to fetch")
	ErrFetchMissingSnowflake = errors.New("object requires snowflake to fetch")

	// Converter errors.

	ErrMemberNotFound     = errors.New("member provided was not found")
	ErrUserNotFound       = errors.New("user provided was not found")
	ErrChannelNotFound    = errors.New("channel provided was not found")
	ErrGuildNotFound      = errors.New("guild provided was not found")
	ErrRoleNotFound       = errors.New("role provided was not found")
	ErrEmojiNotFound      = errors.New("emoji provided was not found")
	ErrBadWebhookArgument = errors.New("webhook url provided was not in valid format")
)
