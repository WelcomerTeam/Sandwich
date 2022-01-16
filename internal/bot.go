package internal

type Bot struct {
	*Commandable
	*Handlers
}

func NewBot() (b *Bot) {
	b = &Bot{
		Commandable: NewCommandable(nil),
		Handlers:    NewHandlers(),
	}

	return b
}
