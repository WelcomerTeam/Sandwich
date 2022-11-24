package internal

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Bot struct {
	Logger zerolog.Logger

	Cogs map[string]Cog

	*Handlers
}

func NewBot(logger zerolog.Logger) *Bot {
	bot := &Bot{
		Logger:   logger,
		Cogs:     make(map[string]Cog),
		Handlers: NewDiscordHandlers(),
	}

	return bot
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

	if cast, ok := cog.(CogWithEvents); ok {
		bot.Logger.Info().Str("cog", cogInfo.Name).Msg("Cog has events")

		bot.RegisterCogEvents(cast.GetEventHandlers())
	}

	return nil
}

func (bot *Bot) RegisterCogEvents(events *Handlers) {
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

				bot.Logger.Info().
					Str("event", eventHandler.eventName).
					Int("events", len(eventHandler.Events)).
					Msg("Registered new events")

				eventHandler.eventsMu.RUnlock()
				botEventHandler.eventsMu.Unlock()
			}

			bot.eventHandlersMu.Unlock()
		}

		eventHandler.eventsMu.RUnlock()
	}
}
