package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"golang.org/x/xerrors"
)

func NewEmoji(ctx *EventContext, guildID *discord.Snowflake, emojiID discord.Snowflake) *discord.Emoji {
	return &discord.Emoji{
		ID:      emojiID,
		GuildID: guildID,
	}
}

func FetchEmoji(ctx *EventContext, e *discord.Emoji) (emoji *discord.Emoji, err error) {
	if e.Name != "" {
		return e, nil
	}

	if e.GuildID == nil {
		return e, ErrFetchMissingGuild
	}

	emoji, err = ctx.Sandwich.GRPCInterface.FetchEmojiByID(ctx, *e.GuildID, e.ID)
	if err != nil {
		return e, xerrors.Errorf("Failed to fetch emoji: %v", err)
	}

	if emoji == nil {
		return e, ErrEmojiNotFound
	}

	return
}
