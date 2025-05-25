package internal

import (
	"fmt"
	"log/slog"
)

type Bot struct {
	Logger *slog.Logger

	Cogs map[string]Cog

	*Handlers
}

func NewBot(logger *slog.Logger) *Bot {
	bot := &Bot{
		Logger:   logger,
		Cogs:     make(map[string]Cog),
		Handlers: newDiscordHandlers(),
	}

	return bot
}

// Cogs

func (bot *Bot) MustRegisterCog(cog Cog) {
	if err := bot.RegisterCog(cog); err != nil {
		panic(fmt.Sprintf(`sandwich: RegisterCog(%v): %v`, cog, err.Error()))
	}
}

func (bot *Bot) RegisterCog(cog Cog) error {
	cogInfo := cog.CogInfo()

	if _, ok := bot.Cogs[cogInfo.Name]; ok {
		return ErrCogAlreadyRegistered
	}

	if err := cog.RegisterCog(bot); err != nil {
		bot.Logger.Error("Failed to register cog", "cog", cogInfo.Name, "error", err)
		panic(fmt.Sprintf(`sandwich: RegisterCog(%v): %v`, cog, err.Error()))
	}

	bot.Cogs[cogInfo.Name] = cog

	bot.Logger.Info("Loaded cog", "cog", cogInfo.Name)

	if cast, ok := cog.(CogWithBotLoad); ok {
		bot.Logger.Info("Cog has BotLoad", "cog", cogInfo.Name)

		cast.BotLoad(bot)
	}

	if cast, ok := cog.(CogWithEvents); ok {
		bot.Logger.Info("Cog has events", "cog", cogInfo.Name)

		bot.RegisterCogEvents(cast.GetEventHandlers())
	}

	return nil
}

func (bot *Bot) RegisterCogEvents(events *Handlers) {
	events.eventHandlersMu.RLock()
	defer events.eventHandlersMu.RUnlock()

	for _, eventHandler := range events.EventHandlers {
		eventHandler.EventsMu.RLock()

		if len(eventHandler.Events) > 0 {
			bot.eventHandlersMu.Lock()

			botEventHandler, ok := bot.EventHandlers[eventHandler.eventName]
			if !ok {
				bot.EventHandlers[eventHandler.eventName] = eventHandler

				bot.Logger.Info("Registered new event handler", "event", eventHandler.eventName)
			} else {
				botEventHandler.EventsMu.Lock()
				eventHandler.EventsMu.RLock()

				botEventHandler.Events = append(botEventHandler.Events, eventHandler.Events...)

				bot.Logger.Info("Registered new events",
					"event", eventHandler.eventName,
					"events", len(eventHandler.Events))

				eventHandler.EventsMu.RUnlock()
				botEventHandler.EventsMu.Unlock()
			}

			bot.eventHandlersMu.Unlock()
		}

		eventHandler.EventsMu.RUnlock()
	}
}
