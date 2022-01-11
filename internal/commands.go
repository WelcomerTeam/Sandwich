package internal

import "sync"

type Commandable struct {
	commandsMu sync.RWMutex
	Commands   map[string]*Commandable

	Checks  []*Check
	Command *Command
}

type Command struct {
	Handler *CommandHandler
}

type Check struct{}

func (c *Commandable) RegisterCommand(commandName string, handler CommandHandler) (cc *Commandable) {
	return
}

type CommandHandler func() (err error)

type CommandContext struct {
	// Details about identifier
}
