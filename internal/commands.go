package internal

import "sync"

type Commandable struct {
	commandsMu sync.RWMutex
	Commands   map[string]*Commandable

	Checks  []*Check
	Command *Command
}

func NewCommandable(handler *CommandHandler) (c *Commandable) {
	c = &Commandable{
		commandsMu: sync.RWMutex{},
		Commands:   make(map[string]*Commandable),

		Checks:  make([]*Check, 0),
		Command: nil,
	}

	if handler != nil {
		c.Command = &Command{
			Handler: handler,
		}
	}

	return c
}

type Command struct {
	Handler *CommandHandler
}

type Check struct{}

func (c *Commandable) RegisterCommand(commandName string, handler *CommandHandler) (cc *Commandable) {
	c.commandsMu.Lock()
	defer c.commandsMu.Unlock()

	c.Commands[commandName] = NewCommandable(handler)

	return cc
}

type CommandHandler func(ctx *CommandContext) (err error)

type CommandContext struct {
	// Details about identifier
}
