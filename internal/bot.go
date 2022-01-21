package internal

import (
	"strconv"
	"strings"
)

type Bot struct {
	*Commandable
	*Handlers

	Prefix PrefixCheckFuncType
}

// Func Type of prefix checking. Returns the prefixes that can be used on a command.
type PrefixCheckFuncType func(eventCtx *EventContext, message Message) (prefixes []string, err error)

// Func Type used for command checks.
type CommandCheckFuncType func(commandCtx *CommandContext) (canRun bool, err error)

func NewBot(prefix PrefixCheckFuncType) (b *Bot) {
	b = &Bot{
		Commandable: NewCommandable(nil, nil, false, ""),
		Handlers:    NewDiscordHandlers(),
		Prefix:      prefix,
	}

	return b
}

// Prefix helpers

func StaticPrefixCheck(passedPrefixes ...string) (fun PrefixCheckFuncType) {
	return func(eventCtx *EventContext, message Message) (prefixes []string, err error) {
		return passedPrefixes, nil
	}
}

func WhenMentionedOr(passedPrefixes ...string) (fun PrefixCheckFuncType) {
	return func(eventCtx *EventContext, message Message) (prefixes []string, err error) {
		prefixes = append(prefixes, passedPrefixes...)
		prefixes = append(prefixes, "<@"+strconv.FormatInt(eventCtx.Identifier.ID, 10)+">")
		prefixes = append(prefixes, "<@!"+strconv.FormatInt(eventCtx.Identifier.ID, 10)+">")

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

// ProcessCommand processes the commands that have been registered to the bot.
// This also checks that the message's author is not a bot.
func (b *Bot) ProcessCommands(eventCtx *EventContext, message Message) (err error) {
	if message.Author == nil {
		return nil
	}

	if message.Author.Bot != nil && *message.Author.Bot {
		return nil
	}

	commandCtx, err := b.GetContext(eventCtx, message)
	if err != nil {
		return err
	}

	return b.Invoke(commandCtx)
}

// GetContext returns the command context from a message.
func (b *Bot) GetContext(eventCtx *EventContext, message Message) (commandContext *CommandContext, err error) {
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

	b.commandsMu.RLock()
	command := b.GetCommand(invoker)
	b.commandsMu.RUnlock()

	commandContext.InvokedWith = invoker
	commandContext.Prefix = invokedPrefix
	commandContext.Command = command

	return commandContext, nil
}

// CanRun checks all global bot checks and returns if the message passes them all.
// If an error occurs, the message will be treated as not being able to run.
func (b *Bot) CanRun(ctx *CommandContext) (canRun bool, err error) {
	for _, check := range b.Checks {
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
