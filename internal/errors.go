package internal

import "golang.org/x/xerrors"

var (
	ErrInvalidIdentifier  = xerrors.New("Payload does not include a valid identifier")
	ErrInvalidApplication = xerrors.New("Could not find identifier matching application")
	ErrUnknownEvent       = xerrors.New("Event type does not have a handler")

	ErrCommandAlreadyRegistered = xerrors.New("Command with this name already exists")

	ErrUnexpectedQuoteError          = xerrors.New("Encountered a quote mark inside a non-quoted string")
	ErrExpectedClosingQuoteError     = xerrors.New("Quote character is expected but not found")
	ErrInvalidEndOfQuotedStringError = xerrors.New("Space is expected after the closing quote")

	ErrCommandNotFound         = xerrors.New("Command with this name was not found")
	ErrCheckFailure            = xerrors.New("Command failed built-in checks")
	ErrMissingRequiredArgument = xerrors.New("Command missing required arguments")

	// Converter errors.

	ErrSnowflakeNotFound = xerrors.New("ID does not follow a valid ID or mention format")
	ErrMemberNotFound    = xerrors.New("Member provided was not found")
	ErrUserNotFound      = xerrors.New("User provided was not found")
	ErrChannelNotFound   = xerrors.New("Channel provided was not found")
	ErrGuildNotFound     = xerrors.New("Guild provided was not found")
	ErrRoleNotFound      = xerrors.New("Role provided was not found")

	ErrBadInviteArgument = xerrors.New("Invite provided was invalid or expired")
	ErrBadColourArgument = xerrors.New("Colour provided was not in valid format")
)
