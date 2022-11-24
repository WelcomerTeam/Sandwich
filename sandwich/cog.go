package internal

import "sync"

type CogInfo struct {
	Name        string
	Description string

	Meta interface{}
}

// Cog is the basic interface for any cog. This must provide information about the cog
// such as its name and description.
type Cog interface {
	CogInfo() *CogInfo
	RegisterCog(bot *Bot) error
}

// CogWithEvents is an interface for any cog that implements custom event listeners.
type CogWithEvents interface {
	GetEventHandlers() *Handlers
}

// CogWithBotLoad is an interface for any cog that implements methods that run when a bot loads.
type CogWithBotLoad interface {
	BotLoad(bot *Bot)
}

// CogWithBotUnload is an interface for any cog that implements methods that run when a bot unloads.
type CogWithBotUnload interface {
	BotUnload(bot *Bot, wg *sync.WaitGroup)
}
