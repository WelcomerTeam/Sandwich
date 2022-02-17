package internal

import (
	"strconv"
	"strings"

	discord "github.com/WelcomerTeam/Discord/discord"
)

type Bot struct {
	Commands   *Commandable
	Converters *Converters

	InteractionCommands   *InteractionCommandable
	InteractionConverters *InteractionConverters

	*Handlers
	Prefix PrefixCheckFuncType
}

// Func Type of prefix checking. Returns the prefixes that can be used on a command.
type PrefixCheckFuncType func(eventCtx *EventContext, message discord.Message) (prefixes []string, err error)

// Func Type used for command checks.
type CommandCheckFuncType func(commandCtx *CommandContext) (canRun bool, err error)

// Func Type used for command checks.
type InteractionCheckFuncType func(interactionCtx *InteractionContext) (canRun bool, err error)

func NewBot(prefix PrefixCheckFuncType) (b *Bot) {
	b = &Bot{
		Commands:   setupCommandable(&Commandable{}),
		Converters: NewDefaultConverters(),

		InteractionCommands:   setupInteractionCommandable(&InteractionCommandable{}),
		InteractionConverters: NewInteractionConverters(),

		Handlers: NewDiscordHandlers(),
		Prefix:   prefix,
	}

	return b
}

// Prefix helpers

func StaticPrefixCheck(passedPrefixes ...string) (fun PrefixCheckFuncType) {
	return func(eventCtx *EventContext, message discord.Message) (prefixes []string, err error) {
		return passedPrefixes, nil
	}
}

func WhenMentionedOr(passedPrefixes ...string) (fun PrefixCheckFuncType) {
	return func(eventCtx *EventContext, message discord.Message) (prefixes []string, err error) {
		prefixes = append(prefixes, passedPrefixes...)
		prefixes = append(prefixes, "<@"+strconv.FormatInt(int64(eventCtx.Identifier.ID), 10)+">")
		prefixes = append(prefixes, "<@!"+strconv.FormatInt(int64(eventCtx.Identifier.ID), 10)+">")

		return prefixes, nil
	}
}

// Commands

// Invoke invokes the command given under the context and handles any extra internal
// dispatch mechanisms.
func (b *Bot) Invoke(ctx *CommandContext) (err error) {
	if ctx.Command != nil {
		// dispatch command event
		ok, err := b.CanRun(ctx)

		switch {
		case ok:
			if err = ctx.Command.Invoke(ctx); err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			return ErrCheckFailure
		}
	} else if ctx.InvokedWith != "" {
		return ErrCommandNotFound
	}

	return
}

// Invoke invokes the command given under the context and handles any extra internal
// dispatch mechanisms.
func (b *Bot) InvokeInteraction(ctx *InteractionContext) (resp *InteractionResponse, err error) {
	return ctx.InteractionCommand.Invoke(ctx)
}

// ProcessCommand processes the commands that have been registered to the bot.
// This also checks that the message's author is not a bot.
func (b *Bot) ProcessCommands(eventCtx *EventContext, message discord.Message) (err error) {
	if message.Author == nil {
		return nil
	}

	if message.Author.Bot {
		return nil
	}

	commandCtx, err := b.GetContext(eventCtx, message)
	if err != nil {
		return err
	}

	return b.Invoke(commandCtx)
}

// ProcessInteraction processes the interaction that has been registered to the bot.
func (b *Bot) ProcessInteraction(eventCtx *EventContext, interaction discord.Interaction) (resp *InteractionResponse, err error) {
	interactionCtx, err := b.GetInteractionContext(eventCtx, interaction)
	if err != nil {
		return nil, err
	}

	return b.InvokeInteraction(interactionCtx)
}

// GetContext returns the command context from a message.
func (b *Bot) GetContext(eventCtx *EventContext, message discord.Message) (commandContext *CommandContext, err error) {
	view := NewStringView(message.Content)

	commandContext = NewCommandContext(eventCtx, b, &message, view)

	if message.Author.ID == eventCtx.Identifier.User.ID {
		return
	}

	prefixes, err := b.Prefix(eventCtx, message)

	var invokedPrefix string

	foundPrefix := false

	for _, prefix := range prefixes {
		if strings.HasPrefix(message.Content, prefix) {
			if view.SkipString(prefix) {
				invokedPrefix = prefix
				foundPrefix = true

				break
			}
		}
	}

	eventCtx.Logger.Debug().Bool("foundPrefix", foundPrefix).Msg("GetContext start")

	if !foundPrefix {
		return
	}

	// view.skip_ws()

	invoker := view.GetWord()

	eventCtx.Logger.Debug().Str("invoker", invoker).Str("prefix", invokedPrefix).Msg("Created context")

	command := b.Commands.GetCommand(invoker)

	commandContext.InvokedWith = invoker
	commandContext.Prefix = invokedPrefix
	commandContext.Command = command

	return commandContext, nil
}

// GetInteractionContext returns the interaction context from an interaction.
func (b *Bot) GetInteractionContext(eventCtx *EventContext, interaction discord.Interaction) (interactionContext *InteractionContext, err error) {
	interactionContext = NewInteractionContext(eventCtx, b, &interaction)

	commandTree := constructCommandTree(interaction.Data.Options, make([]string, 0))

	command := b.InteractionCommands.GetCommand(interaction.Data.Name)

	interactionContext.InteractionCommand = command
	interactionContext.CommandTree = commandTree

	return interactionContext, nil
}

func constructCommandTree(options []*discord.InteractionDataOption, tree []string) (newTree []string) {
	newTree = tree

	for _, option := range options {
		switch option.Type {
		case discord.ApplicationCommandOptionTypeSubCommandGroup:
		case discord.ApplicationCommandOptionTypeSubCommand:
			newTree = append(newTree, option.Name)
			newTree = constructCommandTree(option.Options, newTree)
		default:
		}
	}

	return
}

// CanRun checks all global bot checks and returns if the message passes them all.
// If an error occurs, the message will be treated as not being able to run.
func (b *Bot) CanRun(ctx *CommandContext) (canRun bool, err error) {
	for _, check := range b.Commands.Checks {
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
