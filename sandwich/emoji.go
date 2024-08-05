package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
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

	gEmoji, err := ctx.GRPCInterface.FetchEmojiByID(ctx, *emoji.GuildID, emoji.ID)
	if err != nil {
		return emoji, fmt.Errorf("failed to fetch emoji: %w", err)
	}

	emoji = &gEmoji

	if emoji.ID.IsNil() {
		return emoji, ErrEmojiNotFound
	}

	return emoji, nil
}
