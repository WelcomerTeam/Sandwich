package internal

import (
	"fmt"
	"strings"

	"github.com/WelcomerTeam/Discord/discord"
)

type InteractionCommandableType uint8

const (
	InteractionCommandableTypeCommand InteractionCommandableType = iota
	InteractionCommandableTypeSubcommandGroup
	InteractionCommandableTypeSubcommand
)

type InteractionCommandable struct {
	Name string
	Type InteractionCommandableType

	Checks            []InteractionCheckFuncType
	ArgumentParameter []ArgumentParameter

	Handler InteractionHandler

	commands map[string]*InteractionCommandable
	parent   *InteractionCommandable
}

func (ic *InteractionCommandable) AddInteractionCommand(interactionCommandable *InteractionCommandable) (icc *InteractionCommandable, err error) {
	// Convert interactionCommandable parent to SubcommandGroup if it is a subcommand.
	// Convert interactionCommandable to SubcommandGroup if it is not a Command.
	if ic.Type != InteractionCommandableTypeCommand {
		if ic.parent != nil {
			if ic.parent.Type == InteractionCommandableTypeSubcommand {
				ic.parent.Type = InteractionCommandableTypeSubcommandGroup
			}

			ic.Type = InteractionCommandableTypeSubcommandGroup
		}
	}

	commandName := strings.ToLower(interactionCommandable.Name)
	if _, ok := ic.getCommand(commandName); ok {
		err = ErrCommandAlreadyRegistered

		return nil, err
	}

	interactionCommandable = setupInteractionCommandable(interactionCommandable)
	icc.parent = ic

	icc = interactionCommandable

	ic.setCommand(commandName, icc)

	return
}

func (ic *InteractionCommandable) RemoveCommand(name string) (command *InteractionCommandable) {
	command, ok := ic.getCommand(name)

	if !ok {
		return nil
	}

	ic.deleteCommand(name)

	return command
}

func (ic *InteractionCommandable) RecursivelyRemoveAllCommands() {
	for _, command := range ic.commands {
		if command.IsGroup() {
			command.RecursivelyRemoveAllCommands()
		}

		ic.RemoveCommand(command.Name)
	}
}

// GetAllCommands returns all commands.
func (ic *InteractionCommandable) GetAllCommands() (interactionCommandables []*InteractionCommandable) {
	interactionCommandables = make([]*InteractionCommandable, 0)

	for _, commandable := range ic.commands {
		interactionCommandables = append(interactionCommandables, commandable)
	}

	return interactionCommandables
}

func (ic *InteractionCommandable) GetCommand(name string) (interactionCommandable *InteractionCommandable) {
	if !strings.Contains(name, " ") {
		interactionCommandable, _ = ic.getCommand(name)

		return
	}

	names := strings.Split(name, " ")
	if len(names) == 0 {
		return nil
	}

	interactionCommandable = ic.GetCommand(names[0])
	if !interactionCommandable.IsGroup() {
		return
	}

	var ok bool

	for _, name := range names[1:] {
		interactionCommandable, ok = interactionCommandable.getCommand(name)
		if !ok {
			return nil
		}
	}

	return interactionCommandable
}

// IsGroup returns true if the command contains other commands.
func (ic *InteractionCommandable) IsGroup() bool {
	return ic.Type == InteractionCommandableTypeCommand || ic.Type == InteractionCommandableTypeSubcommandGroup
}

// Invoke handles the execution of a command or a group.
func (ic *InteractionCommandable) Invoke(ctx *InteractionContext) (resp *InteractionResponse, err error) {
	if len(ctx.CommandTree) > 0 {
		if ic.IsGroup() {
			branch := ctx.CommandTree[0]
			ctx.CommandTree = ctx.CommandTree[1:]

			commandable := ic.GetCommand(branch)

			if commandable == nil {
				return nil, ErrCommandNotFound
			}

			return commandable.Invoke(ctx)
		} else {
			ctx.EventContext.Logger.Warn().
				Str("command", ic.Name).
				Str("branch", ctx.CommandTree[0]).
				Msg("Encountered non-group whilst traversing command tree.")
		}
	}

	err = ic.prepare(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		errorValue := recover()
		if errorValue != nil {
			ctx.EventContext.Sandwich.RecoverEventPanic(errorValue, ctx.EventContext, ctx.EventContext.payload)
		}
	}()

	if ic.Handler != nil {
		resp, err = ic.Handler(ctx)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

// CanRun checks interactionCommandable checks and returns if the interaction passes them all.
// If an error occurs, the message will be treated as not being able to run.
func (ic *InteractionCommandable) CanRun(ctx *InteractionContext) (canRun bool, err error) {
	for _, check := range ic.Checks {
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

func (ic *InteractionCommandable) prepare(ctx *InteractionContext) (err error) {
	ctx.InteractionCommand = ic

	ok, err := ic.CanRun(ctx)

	switch {
	case !ok:
		return ErrCheckFailure
	case err != nil:
		return err
	}

	// TODO: Parse arguments

	return nil
}

type InteractionContext struct {
	Bot          *Bot
	EventContext *EventContext

	*discord.Interaction

	CommandTree        []string
	InteractionCommand *InteractionCommandable

	Arguments map[string]*Argument
}

// NewInteractionContext creates a new interaction context.
func NewInteractionContext(eventContext *EventContext, bot *Bot, interaction *discord.Interaction) (interactionContext *InteractionContext) {
	return &InteractionContext{
		Bot:          bot,
		EventContext: eventContext,

		Interaction: interaction,

		InteractionCommand: nil,

		Arguments: make(map[string]*Argument),
	}
}

type InteractionHandler func(ctx *InteractionContext) (resp *InteractionResponse, err error)

type InteractionResponse struct {
	Type *discord.InteractionCallbackType
	Data *discord.InteractionCallbackData
}

// MustGetArgument returns an argument based on its name. Panics on error.
func (ctx *InteractionContext) MustGetArgument(name string) (a *Argument) {
	arg, err := ctx.GetArgument(name)
	if err != nil {
		panic(fmt.Sprintf(`ctx: GetArgument(%s): %v`, name, err.Error()))
	}

	return arg
}

// GetArgument returns an argument based on its name.
func (ctx *InteractionContext) GetArgument(name string) (arg *Argument, err error) {
	arg, ok := ctx.Arguments[name]
	if !ok {
		return nil, ErrArgumentNotFound
	}

	return arg, nil
}

// setupInteractionCommandable ensures all nullable variables are properly constructed.
func setupInteractionCommandable(in *InteractionCommandable) (out *InteractionCommandable) {
	if in.commands == nil {
		in.commands = make(map[string]*InteractionCommandable)
	}

	if in.Checks == nil {
		in.Checks = make([]InteractionCheckFuncType, 0)
	}

	return in
}

func (ic *InteractionCommandable) getCommand(name string) (commandable *InteractionCommandable, ok bool) {
	commandable, ok = ic.commands[strings.ToLower(name)]

	return
}

func (ic *InteractionCommandable) deleteCommand(name string) {
	delete(ic.commands, strings.ToLower(name))
}

func (ic *InteractionCommandable) setCommand(name string, commandable *InteractionCommandable) {
	ic.commands[strings.ToLower(name)] = commandable
}
