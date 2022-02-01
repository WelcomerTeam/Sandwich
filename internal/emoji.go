package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

func NewEmoji(ctx *EventContext, guildID *discord.Snowflake, emojiID discord.Snowflake) *discord_structs.Emoji {
	return &discord_structs.Emoji{
		ID:      emojiID,
		GuildID: guildID,
	}
}

func FetchEmoji(ctx *EventContext, e *discord_structs.Emoji) (emoji *discord_structs.Emoji, err error) {
	if e.Name != "" {
		return e, nil
	}

	if e.GuildID == nil {
		return e, ErrFetchMissingGuild
	}

	emoji, err = ctx.Sandwich.grpcInterface.FetchEmojiByID(ctx, *e.GuildID, e.ID)
	if err != nil {
		return e, xerrors.Errorf("Failed to fetch emoji: %v", err)
	}

	if emoji == nil {
		return e, ErrEmojiNotFound
	}

	return
}
