package internal

import (
	"github.com/rs/zerolog"
)

type Bot struct {
	Logger zerolog.Logger

	*Handlers
}

func NewBot(logger zerolog.Logger) *Bot {
	bot := &Bot{
		Logger:   logger,
		Handlers: NewDiscordHandlers(),
	}

	return bot
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
