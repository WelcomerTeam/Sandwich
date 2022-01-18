package internal

import "golang.org/x/xerrors"

var (
	ErrInvalidIdentifier  = xerrors.New("Payload does not include a valid identifier")
	ErrInvalidApplication = xerrors.New("Could not find identifier matching application")
	ErrUnknownEvent       = xerrors.New("Event type does not have a handler")
)
