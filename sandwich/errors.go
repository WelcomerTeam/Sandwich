package internal

import "github.com/pkg/errors"

var (
	ErrInvalidIdentifier  = errors.New("payload does not include a valid identifier")
	ErrInvalidApplication = errors.New("could not find identifier matching application")
	ErrInvalidToken       = errors.Errorf("invalid token was passed")
	ErrUnknownEvent       = errors.New("event type does not have a handler")
	ErrUnknownGRPCError   = errors.New("GRPC returned unknown error")

	ErrCommandTreeNotAllowed = errors.New("branching command tree not allowed. You must add commands to a group.")

	ErrInvalidTarget = errors.New("invalid target. Are you doing an operation on a user who is not the running application?")

	ErrFetchMissingGuild     = errors.New("object requires guild ID to fetch")
	ErrFetchMissingSnowflake = errors.New("object requires snowflake to fetch")

	ErrCogAlreadyRegistered     = errors.New("cog with this name already exists")
	ErrCommandAlreadyRegistered = errors.New("command with this name already exists")
	ErrInvalidArgumentType      = errors.New("argument value is not correct type for converter used")

	ErrUnexpectedQuoteError          = errors.New("encountered a quote mark inside a non-quoted string")
	ErrExpectedClosingQuoteError     = errors.New("quote character is expected but not found")
	ErrInvalidEndOfQuotedStringError = errors.New("space is expected after the closing quote")

	ErrCommandNotFound         = errors.New("command with this name was not found")
	ErrCheckFailure            = errors.New("command failed built-in checks")
	ErrMissingRequiredArgument = errors.New("command missing required arguments")
	ErrArgumentNotFound        = errors.New("command argument was not found")
	ErrConverterNotFound       = errors.New("command converter is not setup")

	// Converter errors.

	ErrSnowflakeNotFound = errors.New("id does not follow a valid id or mention format")
	ErrMemberNotFound    = errors.New("member provided was not found")
	ErrUserNotFound      = errors.New("user provided was not found")
	ErrChannelNotFound   = errors.New("channel provided was not found")
	ErrGuildNotFound     = errors.New("guild provided was not found")
	ErrRoleNotFound      = errors.New("role provided was not found")
	ErrEmojiNotFound     = errors.New("emoji provided was not found")

	ErrBadInviteArgument  = errors.New("invite provided was invalid or expired")
	ErrBadColourArgument  = errors.New("colour provided was not in valid format")
	ErrBadBoolArgument    = errors.New("bool provided was not in valid format")
	ErrBadIntArgument     = errors.New("int provided was not in valid format")
	ErrBadFloatArgument   = errors.New("float provided was not in valid format")
	ErrBadWebhookArgument = errors.New("webhook url provided was not in valid format")
)
