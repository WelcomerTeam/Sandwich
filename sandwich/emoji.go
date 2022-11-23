package internal

import (
	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/pkg/errors"
)

func NewEmoji(guildID *discord.Snowflake, emojiID discord.Snowflake) *discord.Emoji {
	return &discord.Emoji{
		ID:      emojiID,
		GuildID: guildID,
	}
}

func FetchEmoji(ctx *GRPCContext, emoji *discord.Emoji) (*discord.Emoji, error) {
	if emoji.Name != "" {
		return emoji, nil
	}

	if emoji.GuildID == nil {
		return emoji, ErrFetchMissingGuild
	}

	emoji, err := ctx.GRPCInterface.FetchEmojiByID(ctx, *emoji.GuildID, emoji.ID)
	if err != nil {
		return emoji, errors.Errorf("Failed to fetch emoji: %v", err)
	}

	if emoji == nil {
		return emoji, ErrEmojiNotFound
	}

	return emoji, nil
}
