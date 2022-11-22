package internal

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/rs/zerolog"
)

type Bot struct {
	Logger zerolog.Logger

	Commands            *Commandable
	InteractionCommands *InteractionCommandable

	Cogs map[string]Cog

	Converters            *Converters
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

func NewBot(prefix PrefixCheckFuncType, logger zerolog.Logger) *Bot {
	bot := &Bot{
		Logger: logger,

		Commands:            SetupCommandable(&Commandable{}),
		InteractionCommands: SetupInteractionCommandable(&InteractionCommandable{}),

		Cogs: make(map[string]Cog),

		Converters:            NewDefaultConverters(),
		InteractionConverters: NewInteractionConverters(),

		Handlers: NewDiscordHandlers(),
		Prefix:   prefix,
	}

	return bot
}

func (bot *Bot) Close(wg *sync.WaitGroup) {
	for _, cog := range bot.Cogs {
		if cast, ok := cog.(CogWithBotUnload); ok {
			wg.Add(1)

			cogInfo := cog.CogInfo()

			bot.Logger.Info().Str("cog", cogInfo.Name).Msg("Cog has BotUnload")

			cast.BotUnload(bot, wg)
		}
	}

	wg.Done()
}

// Prefix helpers

func StaticPrefixCheck(passedPrefixes ...string) (fun PrefixCheckFuncType) {
	return func(eventCtx *EventContext, message discord.Message) (prefixes []string, err error) {
		return passedPrefixes, nil
	}
}

func WhenMentionedOr(passedPrefixes ...string) (fun PrefixCheckFuncType) {
	return func(eventCtx *EventContext, message discord.Message) ([]string, error) {
		prefixes := []string{}
		prefixes = append(prefixes, passedPrefixes...)
		prefixes = append(prefixes, "<@"+strconv.FormatInt(int64(eventCtx.Identifier.ID), 10)+">")
		prefixes = append(prefixes, "<@!"+strconv.FormatInt(int64(eventCtx.Identifier.ID), 10)+">")

		return prefixes, nil
	}
}

// Cogs

func (bot *Bot) MustRegisterCog(cog Cog) {
	if err := bot.RegisterCog(cog); err != nil {
		panic(fmt.Sprintf(`sandwich: RegisterCog(%v): %v`, cog, err.Error()))
	}
}

func (bot *Bot) RegisterCog(cog Cog) (err error) {
	cogInfo := cog.CogInfo()

	if _, ok := bot.Cogs[cogInfo.Name]; ok {
		return ErrCogAlreadyRegistered
	}

	err = cog.RegisterCog(bot)
	if err != nil {
		bot.Logger.Panic().Str("cog", cogInfo.Name).Err(err).Msg("Failed to register cog")

		return
	}

	bot.Cogs[cogInfo.Name] = cog

	bot.Logger.Info().Str("cog", cogInfo.Name).Msg("Loaded cog")

	if cast, ok := cog.(CogWithBotLoad); ok {
		bot.Logger.Info().Str("cog", cogInfo.Name).Msg("Cog has BotLoad")

		cast.BotLoad(bot)
	}

	if cast, ok := cog.(CogWithCommands); ok {
		commandable := cast.GetCommandable()

		bot.Logger.Info().Str("cog", cogInfo.Name).Int("commands", len(commandable.GetAllCommands())).Msg("Cog has commands")

		bot.RegisterCogCommandable(cog, commandable)
	}

	if cast, ok := cog.(CogWithInteractionCommands); ok {
		interactionCommandable := cast.GetInteractionCommandable()

		bot.Logger.Info().Str("cog", cogInfo.Name).Int("commands", len(interactionCommandable.GetAllCommands())).Msg("Cog has interaction commands")

		bot.RegisterCogInteractionCommandable(cog, interactionCommandable)
	}

	if cast, ok := cog.(CogWithEvents); ok {
		bot.Logger.Info().Str("cog", cogInfo.Name).Msg("Cog has events")

		bot.RegisterCogEvents(cog, cast.GetEventHandlers())
	}

	return nil
}

func (bot *Bot) RegisterCogCommandable(cog Cog, commandable *Commandable) {
	for _, command := range commandable.GetAllCommands() {
		command := command

		// Add cog checks to all commands.
		command.Checks = append(commandable.Checks, command.Checks...)

		bot.Logger.Debug().Str("name", command.Name).Msg("Registering command")

		bot.Commands.MustAddCommand(command)
	}
}

func (bot *Bot) RegisterCogInteractionCommandable(cog Cog, interactionCommandable *InteractionCommandable) {
	for _, command := range interactionCommandable.GetAllCommands() {
		// Add cog checks to all commands.
		command.Checks = append(interactionCommandable.Checks, command.Checks...)

		bot.Logger.Debug().Str("name", command.Name).Msg("Registering interaction command")

		bot.InteractionCommands.MustAddInteractionCommand(command)
	}
}

func (bot *Bot) RegisterCogEvents(cog Cog, events *Handlers) {
	events.eventHandlersMu.RLock()
	defer events.eventHandlersMu.RUnlock()

	for _, eventHandler := range events.EventHandlers {
		eventHandler.eventsMu.RLock()

		if len(eventHandler.Events) > 0 {
			bot.eventHandlersMu.Lock()

			botEventHandler, ok := bot.EventHandlers[eventHandler.eventName]
			if !ok {
				bot.EventHandlers[eventHandler.eventName] = eventHandler

				bot.Logger.Info().Str("event", eventHandler.eventName).Msg("Registered new event handler")
			} else {
				botEventHandler.eventsMu.Lock()
				eventHandler.eventsMu.RLock()

				botEventHandler.Events = append(botEventHandler.Events, eventHandler.Events...)

				bot.Logger.Info().Str("event", eventHandler.eventName).Int("events", len(eventHandler.Events)).Msg("Registered new events")

				eventHandler.eventsMu.RUnlock()
				botEventHandler.eventsMu.Unlock()
			}

			bot.eventHandlersMu.Unlock()
		}

		eventHandler.eventsMu.RUnlock()
	}
}

// Commands

// Invoke invokes the command given under the context and handles any extra internal
// dispatch mechanisms.
func (bot *Bot) Invoke(ctx *CommandContext) (err error) {
	if ctx.Command != nil {
		// dispatch command event
		ok, err := bot.CanRun(ctx)

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
func (bot *Bot) InvokeInteraction(ctx *InteractionContext) (resp *discord.InteractionResponse, err error) {
	return ctx.InteractionCommand.Invoke(ctx)
}

// ProcessCommand processes the commands that have been registered to the bot.
// This also checks that the message's author is not a bot.
func (bot *Bot) ProcessCommands(eventCtx *EventContext, message discord.Message) (err error) {
	if message.Author == nil {
		return nil
	}

	if message.Author.Bot {
		return nil
	}

	commandCtx, err := bot.GetContext(eventCtx, message)
	if err != nil {
		return err
	}

	return bot.Invoke(commandCtx)
}

// ProcessInteraction processes the interaction that has been registered to the bot.
func (bot *Bot) ProcessInteraction(eventCtx *EventContext, interaction discord.Interaction) (resp *discord.InteractionResponse, err error) {
	interactionCtx, err := bot.GetInteractionContext(eventCtx, interaction)
	if err != nil {
		return nil, err
	}

	if interactionCtx.InteractionCommand == nil {
		return nil, ErrCommandNotFound
	}

	return bot.InvokeInteraction(interactionCtx)
}

// GetContext returns the command context from a message.
func (bot *Bot) GetContext(eventCtx *EventContext, message discord.Message) (commandContext *CommandContext, err error) {
	view := NewStringView(message.Content)

	commandContext = NewCommandContext(eventCtx, bot, &message, view)

	if message.Author.ID == eventCtx.Identifier.User.ID {
		return
	}

	prefixes, err := bot.Prefix(eventCtx, message)

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

	view.SkipWS()

	invoker := view.GetWord()

	eventCtx.Logger.Debug().Str("invoker", invoker).Str("prefix", invokedPrefix).Msg("Created context")

	command := bot.Commands.GetCommand(invoker)

	commandContext.InvokedWith = invoker
	commandContext.Prefix = invokedPrefix
	commandContext.Command = command

	return commandContext, nil
}

// GetInteractionContext returns the interaction context from an interaction.
func (bot *Bot) GetInteractionContext(eventCtx *EventContext, interaction discord.Interaction) (interactionContext *InteractionContext, err error) {
	interactionContext = NewInteractionContext(eventCtx, bot, &interaction)

	commandTree := constructCommandTree(interaction.Data.Options, make([]string, 0))

	command := bot.InteractionCommands.GetCommand(interaction.Data.Name)

	interactionContext.InteractionCommand = command
	interactionContext.commandBranch = commandTree
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
func (bot *Bot) CanRun(ctx *CommandContext) (canRun bool, err error) {
	for _, check := range bot.Commands.Checks {
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
