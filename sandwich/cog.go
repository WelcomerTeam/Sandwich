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
	RegisterCog(b *Bot) error
}

// CogWithCommands is an interface for any cog that implements methods that return cog commands.
type CogWithCommands interface {
	GetCommandable() *Commandable
}

// CogWithInteractionCommands is an interface for any cog that implements methods that return interaction commands.
type CogWithInteractionCommands interface {
	GetInteractionCommandable() *InteractionCommandable
}

// CogWithBotLoad is an interface for any cog that implements methods that run when a bot loads.
type CogWithBotLoad interface {
	BotLoad(b *Bot)
}

// CogWithBotUnload is an interface for any cog that implements methods that run when a bot unloads.
type CogWithBotUnload interface {
	BotUnload(b *Bot, wg *sync.WaitGroup)
}
