package internal

import (
	"strings"
	"sync"
)

type ArgumentType uint16

const (
	ArgumentTypeSnowflake ArgumentType = iota + 1
	ArgumentTypeMember
	ArgumentTypeUser
	ArgumentTypeMessage
	ArgumentTypePartialMessage
	ArgumentTypeTextchannel
	ArgumentTypeInvite
	ArgumentTypeGuild
	ArgumentTypeRole
	ArgumentTypeGame
	ArgumentTypeColour
	ArgumentTypeVoiceChannel
	ArgumentTypeStageChannel
	ArgumentTypeEmoji
	ArgumentTypePartialEmoji
	ArgumentTypeCategoryChannel
	ArgumentTypeStoreChannel
	ArgumentTypeThread
	ArgumentTypeGuildChannel
	ArgumentTypeGuildSticker
	ArgumentTypeString
	ArgumentTypeBool
	ArgumentTypeInt
	ArgumentTypeFloat
	ArgumentTypeFill
)

type Commandable struct {
	commandsMu sync.RWMutex
	commands   map[string]*Commandable

	Name    string
	Aliases []string

	Checks []CommandCheckFuncType

	ArgumentParameters []ArgumentParameter
	Command            *CommandHandler

	InvokeWithoutCommand bool

	Parent *Commandable
}

func NewCommandable(handler *CommandHandler, parent *Commandable, invokeWithoutCommand bool, commandName string, commandAliases ...string) (c *Commandable) {
	c = &Commandable{
		commandsMu: sync.RWMutex{},
		commands:   make(map[string]*Commandable),

		Name:    commandName,
		Aliases: commandAliases,

		Checks: make([]CommandCheckFuncType, 0),

		ArgumentParameters: make([]ArgumentParameter, 0),
		Command:            nil,

		InvokeWithoutCommand: invokeWithoutCommand,

		Parent: nil,
	}

	if handler != nil {
		c.Command = handler
	}

	if parent != nil {
		c.Parent = parent
	}

	return c
}

func (c *Commandable) getCommand(name string) (commandable *Commandable, ok bool) {
	commandable, ok = c.commands[strings.ToLower(name)]

	return
}

func (c *Commandable) deleteCommand(name string) {
	delete(c.commands, strings.ToLower(name))
}

func (c *Commandable) setCommand(name string, commandable *Commandable) {
	c.commands[strings.ToLower(name)] = commandable
}

func (c *Commandable) AddGroup(invokeWithoutCommand bool, groupName string, groupAliases ...string) (cc *Commandable, err error) {
	c.commandsMu.Lock()
	defer c.commandsMu.Unlock()

	for _, commandName := range append(groupAliases, groupName) {
		commandName = strings.ToLower(commandName)
		if _, ok := c.getCommand(commandName); ok {
			err = ErrCommandAlreadyRegistered

			return
		}
	}

	cc = NewCommandable(nil, c, invokeWithoutCommand, groupName, groupAliases...)

	for _, commandName := range append(groupAliases, groupName) {
		c.setCommand(commandName, cc)
	}

	return
}

func (c *Commandable) AddCommand(handler CommandHandler, commandName string, commandAliases ...string) (cc *Commandable, err error) {
	c.commandsMu.Lock()
	defer c.commandsMu.Unlock()

	for _, commandName := range append(commandAliases, commandName) {
		commandName = strings.ToLower(commandName)
		if _, ok := c.getCommand(commandName); ok {
			err = ErrCommandAlreadyRegistered

			return
		}
	}

	cc = NewCommandable(&handler, c, false, commandName, commandAliases...)

	for _, commandName := range append(commandAliases, commandName) {
		c.setCommand(commandName, cc)
	}

	return
}

func (c *Commandable) RemoveCommand(name string) (command *Commandable) {
	c.commandsMu.RLock()
	command, ok := c.getCommand(name)
	c.commandsMu.RUnlock()

	if !ok {
		return nil
	}

	c.commandsMu.Lock()
	defer c.commandsMu.Unlock()

	c.deleteCommand(name)

	if contains(command.Aliases, name) {
		// We do not want to remove the original command.
		return command
	}

	for _, alias := range command.Aliases {
		aliasCommand, ok := c.getCommand(alias)

		// Whilst unlikely, an alias may already conflict.
		// Make sure the alias is not the same as the current
		// command we are looking to remove
		if ok && aliasCommand == command {
			c.deleteCommand(alias)
		}
	}

	return command
}

func (c *Commandable) GetCommand(name string) (commandable *Commandable) {
	if !strings.Contains(name, " ") {
		commandable, _ = c.getCommand(name)

		return
	}

	names := strings.Split(name, " ")
	if len(names) == 0 {
		return nil
	}

	commandable = c.GetCommand(names[0])
	if !commandable.IsGroup() {
		return
	}

	var ok bool

	for _, name := range names[1:] {
		commandable, ok = commandable.getCommand(name)
		if !ok {
			return nil
		}
	}

	return commandable
}

func (c *Commandable) RecursivelyRemoveAllCommands() {
	c.commandsMu.RLock()
	defer c.commandsMu.RUnlock()

	for _, command := range c.commands {
		if command.IsGroup() {
			command.RecursivelyRemoveAllCommands()
		}

		c.RemoveCommand(command.Name)
	}
}

func (c *Commandable) SetHandler(handler CommandHandler) (cc *Commandable) {
	c.Command = &handler

	return cc
}

func (c *Commandable) AddChecks(checks ...CommandCheckFuncType) (cc *Commandable) {
	c.Checks = append(c.Checks, checks...)

	return c
}

func (c *Commandable) AddArguments(argumentParameters ...ArgumentParameter) (cc *Commandable) {
	c.ArgumentParameters = append(c.ArgumentParameters, argumentParameters...)

	return c
}

func (c *Commandable) IsGroup() bool {
	c.commandsMu.RLock()
	defer c.commandsMu.RUnlock()

	return len(c.commands) > 0
}

func (c *Commandable) Invoke(ctx *CommandContext) (err error) {
	if c.IsGroup() {
		ctx.InvokedSubcommand = nil
		ctx.SubcommandPassed = nil

		earlyInvoke := !c.InvokeWithoutCommand
		if earlyInvoke {
			err = c.prepare(ctx)
			if err != nil {
				return err
			}
		}

		view := ctx.View
		previous := view.index
		view.SkipWS()
		trigger := view.GetWord()

		if trigger != "" {
			ctx.SubcommandPassed = &trigger
			ctx.InvokedSubcommand = c.GetCommand(trigger)
		}

		if earlyInvoke {
			if c.Command != nil {
				callback := *c.Command

				err = callback(ctx)
				if err != nil {
					return err
				}
			}
		}

		ctx.InvokedParents = append(ctx.InvokedParents, ctx.InvokedWith)

		if trigger != "" && ctx.InvokedSubcommand != nil {
			ctx.InvokedWith = trigger

			err = ctx.InvokedSubcommand.Invoke(ctx)
			if err != nil {
				return err
			}
		} else if !earlyInvoke {
			view.index = previous
			view.previous = previous

			if c.Command != nil {
				callback := *c.Command

				println(c.Command, callback)

				err = callback(ctx)
				if err != nil {
					return err
				}
			}
		}
	} else {
		err = c.prepare(ctx)
		if err != nil {
			return err
		}

		ctx.InvokedSubcommand = nil
		ctx.SubcommandPassed = nil

		if c.Command != nil {
			callback := *c.Command

			err = callback(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CanRun checks commandable bot checks and returns if the message passes them all.
// If an error occurs, the message will be treated as not being able to run.
func (c *Commandable) CanRun(ctx *CommandContext) (canRun bool, err error) {
	for _, check := range c.Checks {
		canRun, err := check(ctx)
		if err != nil {
			return false, err
		}

		if !canRun {
			return false, nil
		}
	}

	return true, nil
}

func (c *Commandable) prepare(ctx *CommandContext) (err error) {
	ctx.Command = c

	ok, err := c.CanRun(ctx)

	switch {
	case !ok:
		return ErrCheckFailure
	case err != nil:
		return err
	}

	err = c.parseArguments(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commandable) parseArguments(ctx *CommandContext) (err error) {
	ctx.Arguments = map[string]Argument{}

	for _, argumentParameter := range c.ArgumentParameters {
		ctx.currentParameter = &argumentParameter

		transformed, err := c.transform(ctx, argumentParameter)
		if err != nil {
			return err
		}

		ctx.Arguments[argumentParameter.Name] = Argument{
			ArgumentType: argumentParameter.ArgumentType,
			value:        transformed,
		}
	}

	return nil
}

func (c *Commandable) transform(ctx *CommandContext, argumentParameter ArgumentParameter) (out interface{}, err error) {
	required := argumentParameter.Required
	consumeRestIsSpecial := argumentParameter.ArgumentType == ArgumentTypeFill

	converter := getConverter(argumentParameter.ArgumentType)

	view := ctx.View
	view.SkipWS()

	if view.EOF() {
		if required {
			return nil, ErrMissingRequiredArgument
		}
	}

	previous := view.index

	var argument string

	if consumeRestIsSpecial {
		argument = strings.TrimSpace(view.ReadRest())
	} else {
		argument, _, err = view.GetQuotedWord()
		if err != nil {
			return nil, err
		}
	}

	view.previous = previous

	return converter(ctx, argument)
}

type CommandHandler func(ctx *CommandContext) (err error)

type CommandContext struct {
	Bot          *Bot
	EventContext *EventContext

	*Message

	Guild *Guild

	Prefix  string
	Command *Commandable
	View    *StringView

	InvokedWith    string
	InvokedParents []string

	InvokedSubcommand *Commandable
	SubcommandPassed  *string
	CommandFailed     bool

	currentParameter *ArgumentParameter

	Arguments map[string]Argument
}

func getConverter(converterType ArgumentType) ArgumentConverterType {
	return func(ctx *CommandContext, argument string) (out interface{}, err error) {
		return argument, nil
	}
}

func NewCommandContext(eventContext *EventContext, bot *Bot, message *Message, view *StringView) (commandContext *CommandContext) {
	commandContext = &CommandContext{
		Bot:          bot,
		EventContext: eventContext,

		Message: message,
		Guild:   eventContext.Guild,

		Prefix:  "",
		Command: nil,
		View:    view,

		InvokedWith:    "",
		InvokedParents: make([]string, 0),

		InvokedSubcommand: nil,
		SubcommandPassed:  nil,
		CommandFailed:     false,
	}

	return commandContext
}
