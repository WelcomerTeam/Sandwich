package internal

import (
	discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"
	"golang.org/x/xerrors"
)

type Emoji discord.Emoji

func NewEmoji(ctx *EventContext, guildID *discord.Snowflake, emojiID discord.Snowflake) *Emoji {
	return &Emoji{
		ID:      emojiID,
		GuildID: guildID,
	}
}

func (e *Emoji) Fetch(ctx *EventContext) (err error) {
	if e.Name != "" {
		return nil
	}

	if e.GuildID == nil {
		return ErrFetchMissingGuild
	}

	emoji, err := ctx.Sandwich.grpcInterface.FetchEmojiByID(ctx, *e.GuildID, e.ID)
	if err != nil {
		return xerrors.Errorf("Failed to fetch emoji: %v", err)
	}

	if emoji != nil {
		*e = *emoji
	} else {
		return ErrEmojiNotFound
	}

	return
}
