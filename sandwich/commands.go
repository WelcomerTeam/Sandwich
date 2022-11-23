package internal

import (
	"fmt"
	"strings"

	discord "github.com/WelcomerTeam/Discord/discord"
)

type ArgumentType uint16

const (
	ArgumentTypeSnowflake ArgumentType = iota + 1
	ArgumentTypeMember
	ArgumentTypeUser
	_
	_
	ArgumentTypeTextChannel
	ArgumentTypeGuild
	ArgumentTypeRole
	ArgumentTypeColour
	ArgumentTypeVoiceChannel
	ArgumentTypeStageChannel
	ArgumentTypeEmoji
	ArgumentTypePartialEmoji
	ArgumentTypeCategoryChannel
	ArgumentTypeStoreChannel
	ArgumentTypeThread
	ArgumentTypeGuildChannel
	_
	ArgumentTypeString
	ArgumentTypeBool
	ArgumentTypeInt
	ArgumentTypeFloat
	ArgumentTypeFill
)

type Commandable struct {
	Name    string
	Aliases []string

	Description string

	Checks             []CommandCheckFuncType
	ArgumentParameters []ArgumentParameter

	Handler CommandHandler

	InvokeWithoutCommand bool

	commands map[string]*Commandable
	parent   *Commandable
}

type CommandHandler func(ctx *CommandContext) (err error)

func (c *Commandable) MustAddCommand(commandable *Commandable) (cc *Commandable) {
	cc, err := c.AddCommand(commandable)
	if err != nil {
		panic(fmt.Sprintf(`sandwich: AddCommand(%v): %v`, commandable, err.Error()))
	}

	return cc
}

func (c *Commandable) AddCommand(commandable *Commandable) (cc *Commandable, err error) {
	for _, commandName := range append(commandable.Aliases, commandable.Name) {
		commandName = strings.ToLower(commandName)
		if _, ok := c.getCommand(commandName); ok {
			err = ErrCommandAlreadyRegistered

			return nil, err
		}
	}

	commandable = SetupCommandable(commandable)
	commandable.parent = c

	cc = commandable

	for _, commandName := range append(commandable.Aliases, commandable.Name) {
		c.setCommand(commandName, cc)
	}

	return cc, nil
}

func (c *Commandable) RemoveCommand(name string) (command *Commandable) {
	command, ok := c.getCommand(name)

	if !ok {
		return nil
	}

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

func (c *Commandable) RecursivelyRemoveAllCommands() {
	for _, command := range c.commands {
		if command.IsGroup() {
			command.RecursivelyRemoveAllCommands()
		}

		c.RemoveCommand(command.Name)
	}
}

// GetAllCommands returns all commands and removes duplicates due to aliases.
func (c *Commandable) GetAllCommands() (commandables []*Commandable) {
	commandables = make([]*Commandable, 0)

	for key, commandable := range c.commands {
		// If the commandable's name is the same as it's key, it is safe to assume
		// it is not an alias.
		if strings.EqualFold(commandable.Name, key) {
			commandables = append(commandables, commandable)
		}
	}

	return commandables
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

// IsGroup returns true if the command contains other commands.
func (c *Commandable) IsGroup() bool {
	return len(c.commands) > 0
}

// Invoke handles the execution of a command or a group.
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
			if c.Handler != nil {
				err = c.Handler(ctx)
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
		}
	} else {
		err = c.prepare(ctx)
		if err != nil {
			return err
		}

		ctx.InvokedSubcommand = nil
		ctx.SubcommandPassed = nil
	}

	defer func() {
		errorValue := recover()
		if errorValue != nil {
			ctx.EventContext.Sandwich.RecoverEventPanic(errorValue, ctx.EventContext, ctx.EventContext.payload)
		}
	}()

	if c.Handler != nil {
		err = c.Handler(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// CanRun checks commandable checks and returns if the message passes them all.
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

// parseArguments generates the arguments for a command.
func (c *Commandable) parseArguments(ctx *CommandContext) (err error) {
	ctx.Arguments = map[string]*Argument{}

	for _, argumentParameter := range c.ArgumentParameters {
		ctx.currentParameter = &argumentParameter

		transformed, err := c.transform(ctx, argumentParameter)
		if err != nil {
			return err
		}

		ctx.Arguments[argumentParameter.Name] = &Argument{
			ArgumentType: argumentParameter.ArgumentType,
			value:        transformed,
		}
	}

	return nil
}

// transform returns a output value based on the argument parameter passed in.
func (c *Commandable) transform(ctx *CommandContext, argumentParameter ArgumentParameter) (out interface{}, err error) {
	required := argumentParameter.Required
	consumeRestIsSpecial := argumentParameter.ArgumentType == ArgumentTypeFill

	converter := ctx.Bot.Converters.GetConverter(argumentParameter.ArgumentType)
	if converter == nil {
		return nil, ErrConverterNotFound
	}

	view := ctx.View
	view.SkipWS()

	if view.EOF() {
		if required {
			return nil, ErrMissingRequiredArgument
		}

		return converter.data, nil
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

	return converter.converterType(ctx, argument)
}

type CommandContext struct {
	Bot          *Bot
	EventContext *EventContext

	*discord.Message

	Guild   *discord.Guild
	Channel *discord.Channel

	Prefix  string
	Command *Commandable
	View    *StringView

	InvokedWith    string
	InvokedParents []string

	InvokedSubcommand *Commandable
	SubcommandPassed  *string
	CommandFailed     bool

	currentParameter *ArgumentParameter

	Arguments map[string]*Argument
}

// NewCommandContext creates a new command context.
func NewCommandContext(eventContext *EventContext, bot *Bot, message *discord.Message, view *StringView) (commandContext *CommandContext) {
	commandContext = &CommandContext{
		Bot:          bot,
		EventContext: eventContext,

		Message: message,

		Guild:   eventContext.Guild,
		Channel: NewChannel(message.GuildID, message.ChannelID),

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

// MustGetArgument returns an argument based on its name. Panics on error.
func (ctx *CommandContext) MustGetArgument(name string) (a *Argument) {
	arg, err := ctx.GetArgument(name)
	if err != nil {
		panic(fmt.Sprintf(`ctx: GetArgument(%s): %v`, name, err.Error()))
	}

	return arg
}

// GetArgument returns an argument based on its name.
func (ctx *CommandContext) GetArgument(name string) (arg *Argument, err error) {
	arg, ok := ctx.Arguments[name]
	if !ok {
		return nil, ErrArgumentNotFound
	}

	return arg, nil
}

// SetupCommandable ensures all nullable variables are properly constructed.
func SetupCommandable(in *Commandable) (out *Commandable) {
	if in.commands == nil {
		in.commands = make(map[string]*Commandable)
	}

	if in.Aliases == nil {
		in.Aliases = make([]string, 0)
	}

	if in.ArgumentParameters == nil {
		in.ArgumentParameters = make([]ArgumentParameter, 0)
	}

	if in.Checks == nil {
		in.Checks = make([]CommandCheckFuncType, 0)
	}

	return in
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
