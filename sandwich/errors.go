package internal

import "github.com/pkg/errors"

var (
	ErrInvalidIdentifier  = errors.New("Payload does not include a valid identifier")
	ErrInvalidApplication = errors.New("Could not find identifier matching application")
	ErrInvalidToken       = errors.Errorf("Invalid token was passed")
	ErrUnknownEvent       = errors.New("Event type does not have a handler")
	ErrUnknownGRPCError   = errors.New("GRPC returned unknown error")

	ErrCommandTreeNotAllowed = errors.New("Branching command tree not allowed. You must add commands to a group.")

	ErrInvalidTarget = errors.New("Invalid target. Are you doing an operation on a user who is not the running application?")

	ErrFetchMissingGuild     = errors.New("Object requires guild ID to fetch")
	ErrFetchMissingSnowflake = errors.New("Object requires snowflake to fetch")

	ErrCogAlreadyRegistered     = errors.New("Cog with this name already exists")
	ErrCommandAlreadyRegistered = errors.New("Command with this name already exists")
	ErrInvalidArgumentType      = errors.New("Argument value is not correct type for converter used")

	ErrUnexpectedQuoteError          = errors.New("Encountered a quote mark inside a non-quoted string")
	ErrExpectedClosingQuoteError     = errors.New("Quote character is expected but not found")
	ErrInvalidEndOfQuotedStringError = errors.New("Space is expected after the closing quote")

	ErrCommandNotFound         = errors.New("Command with this name was not found")
	ErrCheckFailure            = errors.New("Command failed built-in checks")
	ErrMissingRequiredArgument = errors.New("Command missing required arguments")
	ErrArgumentNotFound        = errors.New("Command argument was not found")
	ErrConverterNotFound       = errors.New("Command converter is not setup")

	// Converter errors.

	ErrSnowflakeNotFound = errors.New("ID does not follow a valid ID or mention format")
	ErrMemberNotFound    = errors.New("Member provided was not found")
	ErrUserNotFound      = errors.New("User provided was not found")
	ErrChannelNotFound   = errors.New("Channel provided was not found")
	ErrGuildNotFound     = errors.New("Guild provided was not found")
	ErrRoleNotFound      = errors.New("Role provided was not found")
	ErrEmojiNotFound     = errors.New("Emoji provided was not found")

	ErrBadInviteArgument  = errors.New("Invite provided was invalid or expired")
	ErrBadColourArgument  = errors.New("Colour provided was not in valid format")
	ErrBadBoolArgument    = errors.New("Bool provided was not in valid format")
	ErrBadIntArgument     = errors.New("Int provided was not in valid format")
	ErrBadFloatArgument   = errors.New("Float provided was not in valid format")
	ErrBadWebhookArgument = errors.New("Webhook url provided was not in valid format")
)
