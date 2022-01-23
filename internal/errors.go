package internal

import "golang.org/x/xerrors"

var (
	ErrInvalidIdentifier  = xerrors.New("Payload does not include a valid identifier")
	ErrInvalidApplication = xerrors.New("Could not find identifier matching application")
	ErrUnknownEvent       = xerrors.New("Event type does not have a handler")
	ErrUnknownGRPCError   = xerrors.New("GRPC returned unknown error")

	ErrFetchMissingGuild = xerrors.New("Object requires guild ID to fetch")

	ErrCommandAlreadyRegistered = xerrors.New("Command with this name already exists")
	ErrInvalidArgumentType      = xerrors.New("Argument value is not correct type for converter used")

	ErrUnexpectedQuoteError          = xerrors.New("Encountered a quote mark inside a non-quoted string")
	ErrExpectedClosingQuoteError     = xerrors.New("Quote character is expected but not found")
	ErrInvalidEndOfQuotedStringError = xerrors.New("Space is expected after the closing quote")

	ErrCommandNotFound         = xerrors.New("Command with this name was not found")
	ErrCheckFailure            = xerrors.New("Command failed built-in checks")
	ErrMissingRequiredArgument = xerrors.New("Command missing required arguments")
	ErrArgumentNotFound        = xerrors.New("Command argument was not found")
	ErrConverterNotFound       = xerrors.New("Command converter is not setup")

	// Converter errors.

	ErrSnowflakeNotFound = xerrors.New("ID does not follow a valid ID or mention format")
	ErrMemberNotFound    = xerrors.New("Member provided was not found")
	ErrUserNotFound      = xerrors.New("User provided was not found")
	ErrChannelNotFound   = xerrors.New("Channel provided was not found")
	ErrGuildNotFound     = xerrors.New("Guild provided was not found")
	ErrRoleNotFound      = xerrors.New("Role provided was not found")
	ErrEmojiNotFound     = xerrors.New("Emoji provided was not found")

	ErrBadInviteArgument = xerrors.New("Invite provided was invalid or expired")
	ErrBadColourArgument = xerrors.New("Colour provided was not in valid format")
	ErrBadBoolArgument   = xerrors.New("Bool provided was not in valid format")
	ErrBadIntArgument    = xerrors.New("Int provided was not in valid format")
	ErrBadFloatArgument  = xerrors.New("Float provided was not in valid format")
)
