package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func NewEmoji(ctx *EventContext, guildID *discord.Snowflake, emojiID discord.Snowflake) *discord.Emoji {
	return &discord.Emoji{
		ID:      emojiID,
		GuildID: guildID,
	}
}

func FetchEmoji(ctx *EventContext, emoji *discord.Emoji) (*discord.Emoji, error) {
	if emoji.Name != "" {
		return emoji, nil
	}

	if emoji.GuildID == nil {
		return emoji, ErrFetchMissingGuild
	}

	emoji, err := ctx.Sandwich.GRPCInterface.FetchEmojiByID(ctx.ToGRPCContext(), *emoji.GuildID, emoji.ID)
	if err != nil {
		return emoji, errors.Errorf("Failed to fetch emoji: %v", err)
	}

	if emoji == nil {
		return emoji, ErrEmojiNotFound
	}

	return emoji, nil
}
